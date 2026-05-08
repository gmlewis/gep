// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package evolution provides a typed generic GEP (Gene Expression Programming)
// evolution engine built on top of the core package.
//
// The package is parameterised by the numeric type T used inside genomes (e.g.
// float64, int) and exposes the following building blocks:
//
//   - Individual[T] – a genome together with its most-recently evaluated fitness score.
//   - Generation[T] – a population of individuals with pluggable scoring, selection,
//     and stopping criteria.
//   - New / NewWithSeed – constructors that create random populations; the seeded
//     variant produces fully deterministic populations and evolution trajectories.
//   - Evaluate – concurrent per-individual fitness evaluation.
//   - BestIndividual – returns the best individual, respecting both maximization
//     and minimization orientations.
//   - Select – roulette-wheel (fitness-proportionate) selection.
//   - Recombine – applies configurable crossover operators (one-point and
//     two-point recombination) to the population.
//   - Mutate – applies configurable genetic operators (point mutation, inversion,
//     IS/RIS transposition, gene transposition) to the population.
//   - Evolve – runs the complete evolution loop (evaluate → stop-check → select →
//     recombine → mutate → elitism) for a given number of iterations.
//
// Typical usage:
//
//	cat := core.NewCatalog[float64]()
//	// ... register function nodes ...
//	link, _ := core.NewLinkFunc[float64]("+", func(v []float64) float64 { ... })
//
//	gen, _ := evolution.NewWithSeed(42, cat, 50, 7, 3, 2, 0, link,
//	    func(g core.Genome[float64]) float64 { return fitness(g) })
//	gen.StopFunc = func(best evolution.Individual[float64]) bool {
//	    return best.Score >= 1000.0
//	}
//	gen.MutationConfig = mutation.Config{
//	    HeadSize:          7,
//	    NumTerminals:      2,
//	    PointMutationRate: 0.044,
//	    InversionRate:     0.1,
//	}
//	result := gen.Evolve(500)
package evolution

import (
	"fmt"
	"math/rand"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/evolution/evaluation"
	"github.com/gmlewis/gep/v2/evolution/mutation"
	"github.com/gmlewis/gep/v2/evolution/recombination"
	"github.com/gmlewis/gep/v2/evolution/selection"
)

// ScoringFunc computes a fitness score for a genome.
// When MinimizeScore is false (the default), higher scores are better.
// When MinimizeScore is true, lower scores are treated as better.
type ScoringFunc[T any] func(core.Genome[T]) float64

// Individual wraps a typed genome together with its most-recently evaluated
// fitness score.
type Individual[T any] struct {
	Genome core.Genome[T]
	Score  float64
}

// Dup returns a deep copy of the individual. The score is copied by value.
func (ind Individual[T]) Dup() Individual[T] {
	return Individual[T]{Genome: ind.Genome.Dup(), Score: ind.Score}
}

// Generation holds a population of typed individuals and the configuration
// needed to evolve them.
type Generation[T any] struct {
	// Individuals is the current population.
	Individuals []Individual[T]

	// ScoringFunc evaluates the fitness of a genome. If nil, Evaluate is a no-op.
	ScoringFunc ScoringFunc[T]

	// MinimizeScore controls score orientation. When false (default), higher
	// scores are better. When true, lower scores are treated as better.
	MinimizeScore bool

	// StopFunc is an optional callback invoked with the current best individual
	// after each generation's evaluation. When StopFunc returns true, evolution
	// halts early and that individual is returned. When nil, the default
	// stopping criterion of Score >= 1000 is used.
	StopFunc func(Individual[T]) bool

	// MutationConfig controls which genetic operators are applied after
	// selection and at what rates. Zero values disable the corresponding
	// operators. See evolution/mutation.Config for details.
	MutationConfig mutation.Config

	// RecombinationConfig controls one-point and two-point crossover rates
	// applied after selection. Zero values disable recombination operators.
	RecombinationConfig recombination.Config

	// cat is the typed function catalog used for point mutation.
	cat *core.Catalog[T]

	// headSize is the gene head size passed at construction time; forwarded to
	// mutation operators.
	headSize int

	// numTerminals is forwarded to mutation operators.
	numTerminals int

	// numConstants is forwarded to mutation operators.
	numConstants int

	// rng is an optional seeded random source for deterministic evolution.
	// When nil (as created by New), the global math/rand source is used.
	// Use NewWithSeed to obtain a fully deterministic Generation.
	rng *rand.Rand
}

// effectiveScore maps a raw score to an internal "higher is better" scale so
// that both maximization and minimization can share the same selection logic.
func (g *Generation[T]) effectiveScore(score float64) float64 {
	if g.MinimizeScore {
		return -score
	}
	return score
}

// New creates a new random Generation.
//
// cat is the typed function catalog used to build random genomes.
// numIndividuals is the population size (must be > 0).
// headSize, numGenesPerGenome, numTerminals, numConstants are forwarded to
// core.NewRandomGenome.
// link is the link operator used to combine gene outputs within each genome.
// sf is the fitness scoring function (may be nil).
func New[T any](
	cat *core.Catalog[T],
	numIndividuals, headSize, numGenesPerGenome, numTerminals, numConstants int,
	link core.LinkOperator[T],
	sf ScoringFunc[T],
) (*Generation[T], error) {
	return newGeneration(nil, cat, numIndividuals, headSize, numGenesPerGenome, numTerminals, numConstants, link, sf)
}

// NewWithSeed creates a new random Generation with a seeded RNG for fully
// deterministic reproducibility. Two calls with the same seed and identical
// parameters produce identical populations and identical evolution trajectories.
func NewWithSeed[T any](
	seed int64,
	cat *core.Catalog[T],
	numIndividuals, headSize, numGenesPerGenome, numTerminals, numConstants int,
	link core.LinkOperator[T],
	sf ScoringFunc[T],
) (*Generation[T], error) {
	rng := rand.New(rand.NewSource(seed))
	return newGeneration(rng, cat, numIndividuals, headSize, numGenesPerGenome, numTerminals, numConstants, link, sf)
}

// newGeneration is the shared internal constructor.
func newGeneration[T any](
	rng *rand.Rand,
	cat *core.Catalog[T],
	numIndividuals, headSize, numGenesPerGenome, numTerminals, numConstants int,
	link core.LinkOperator[T],
	sf ScoringFunc[T],
) (*Generation[T], error) {
	if numIndividuals <= 0 {
		return nil, fmt.Errorf("evolution.New: numIndividuals must be > 0")
	}
	if cat == nil {
		return nil, fmt.Errorf("evolution.New: catalog cannot be nil")
	}
	if link == nil {
		return nil, fmt.Errorf("evolution.New: link operator cannot be nil")
	}

	individuals := make([]Individual[T], numIndividuals)
	for i := range individuals {
		genome, err := core.NewRandomGenome(cat, numGenesPerGenome, headSize, numTerminals, numConstants, link, rng)
		if err != nil {
			return nil, fmt.Errorf("evolution.New: individual[%d]: %w", i, err)
		}
		individuals[i] = Individual[T]{Genome: genome}
	}

	return &Generation[T]{
		Individuals:  individuals,
		ScoringFunc:  sf,
		cat:          cat,
		headSize:     headSize,
		numTerminals: numTerminals,
		numConstants: numConstants,
		rng:          rng,
	}, nil
}

// Evaluate updates the Score of every individual in the population by calling
// ScoringFunc concurrently. It is a no-op when ScoringFunc is nil.
func (g *Generation[T]) Evaluate() {
	if g.ScoringFunc == nil {
		return
	}
	genomes := make([]core.Genome[T], len(g.Individuals))
	for i, ind := range g.Individuals {
		genomes[i] = ind.Genome
	}
	scores := evaluation.ScoreAll(genomes, evaluation.ScoringFunc[T](g.ScoringFunc), evaluation.Config{})
	for i, score := range scores {
		g.Individuals[i].Score = score
	}
}

// BestIndividual returns the individual with the best effective score in the
// current population. When MinimizeScore is false, the individual with the
// highest Score wins; when MinimizeScore is true, the individual with the
// lowest Score wins.
//
// BestIndividual does NOT call Evaluate; call Evaluate first if scores may be
// stale.
func (g *Generation[T]) BestIndividual() Individual[T] {
	best := g.Individuals[0]
	bestEff := g.effectiveScore(best.Score)
	for _, ind := range g.Individuals[1:] {
		if eff := g.effectiveScore(ind.Score); eff > bestEff {
			best = ind
			bestEff = eff
		}
	}
	return best
}

// Select replaces the current population with a new population drawn using the
// roulette-wheel (fitness-proportionate) selection algorithm.
//
// All individuals' Scores must have been set by a prior call to Evaluate.
// The method is safe to call even when all individuals have equal scores
// (every individual has an equal probability of being selected in that case).
func (g *Generation[T]) Select() {
	candidates := make([]selection.Candidate[T], len(g.Individuals))
	for i, ind := range g.Individuals {
		candidates[i] = selection.Candidate[T]{Genome: ind.Genome, Score: ind.Score}
	}

	selected := selection.Roulette(candidates, selection.Config{
		MinimizeScore: g.MinimizeScore,
	}, g.rng)

	g.Individuals = make([]Individual[T], len(selected))
	for i, ind := range selected {
		g.Individuals[i] = Individual[T]{Genome: ind.Genome, Score: ind.Score}
	}
}

// Recombine applies configured crossover operators to the current population
// in place, replacing each individual with a (possibly recombined) offspring.
//
// When both rates are zero (the default), Recombine is effectively a no-op that
// deep-copies the current population.
func (g *Generation[T]) Recombine() {
	genomes := make([]core.Genome[T], len(g.Individuals))
	for i, ind := range g.Individuals {
		genomes[i] = ind.Genome
	}

	recombined := recombination.Apply(genomes, g.RecombinationConfig, g.rng)
	for i, genome := range recombined {
		g.Individuals[i].Genome = genome
	}
}

// Mutate applies the configured genetic operators to the current population
// in place, replacing each individual with a (possibly mutated) offspring.
//
// The operators applied and their rates are governed by MutationConfig.
// When all rates are zero (the default), Mutate is effectively a no-op that
// deep-copies the current population.
//
// Mutate uses the catalog, head size, numTerminals, and numConstants stored at
// construction time to drive point mutation. MutationConfig.HeadSize,
// MutationConfig.NumTerminals, and MutationConfig.NumConstants are overridden
// by the values stored at construction time so callers do not need to duplicate
// those settings.
func (g *Generation[T]) Mutate() {
	cfg := g.MutationConfig
	cfg.HeadSize = g.headSize
	cfg.NumTerminals = g.numTerminals
	cfg.NumConstants = g.numConstants

	genomes := make([]core.Genome[T], len(g.Individuals))
	for i, ind := range g.Individuals {
		genomes[i] = ind.Genome
	}

	mutated := mutation.Apply(genomes, g.cat, cfg, g.rng)
	for i, genome := range mutated {
		g.Individuals[i].Genome = genome
	}
}

// Evolve runs the GEP algorithm for up to iterations generations.
//
// Each iteration:
//  1. Evaluates all individuals (calls Evaluate).
//  2. Identifies the best individual (calls BestIndividual).
//  3. Checks the stopping criterion (StopFunc or default Score >= 1000).
//     If met, returns the best individual immediately.
//  4. Saves a deep copy of the best individual (elitism).
//  5. Replaces the population via roulette-wheel selection (calls Select).
//  6. Applies configured crossover operators to the population (calls Recombine).
//  7. Applies configured genetic operators to the population (calls Mutate).
//  8. Restores the saved best individual to position 0 (elitism guarantee).
//
// If the loop completes without triggering the stopping criterion, Evolve
// performs a final Evaluate and returns BestIndividual.
func (g *Generation[T]) Evolve(iterations int) Individual[T] {
	for i := 0; i < iterations; i++ {
		g.Evaluate()
		best := g.BestIndividual()

		stop := false
		if g.StopFunc != nil {
			stop = g.StopFunc(best)
		} else {
			stop = best.Score >= 1000.0
		}
		if stop {
			return best
		}

		saveBest := best.Dup()
		g.Select()
		g.Recombine()
		g.Mutate()
		// Elitism: ensure the best individual survives to the next generation.
		g.Individuals[0] = saveBest
	}
	g.Evaluate()
	return g.BestIndividual()
}
