// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package evolution provides a typed generic GEP evolution engine built on the
// core package.
//
// This is Phase 3 of the GEP modernization roadmap: the evolution engine is
// rebuilt as an explicit, configurable, and independently testable subsystem.
// The package currently implements the foundation layer:
//
//   - Individual[T] – a scored typed genome.
//   - Generation[T] – a population of individuals with configuration.
//   - New / NewWithSeed – constructors that create random populations.
//   - Evaluate – concurrent fitness evaluation for all individuals.
//   - BestIndividual – best-fitness lookup (supports both max and min goals).
//   - Select – roulette-wheel (fitness-proportionate) selection.
//   - Evolve – a basic evolution loop (evaluate → stop-check → select → elitism).
//
// Genetic operators (mutation, recombination, transposition) are wired in
// subsequent milestones.
package evolution

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/gmlewis/gep/v2/core"
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

	// rng is an optional seeded random source for deterministic evolution.
	// When nil (as created by New), the global math/rand source is used.
	// Use NewWithSeed to obtain a fully deterministic Generation.
	rng *rand.Rand
}

// randIntn returns a random int in [0, n) using g.rng when non-nil.
func (g *Generation[T]) randIntn(n int) int {
	if g.rng != nil {
		return g.rng.Intn(n)
	}
	return rand.Intn(n) //nolint:gosec
}

// randFloat64 returns a random float64 in [0.0, 1.0) using g.rng when non-nil.
func (g *Generation[T]) randFloat64() float64 {
	if g.rng != nil {
		return g.rng.Float64()
	}
	return rand.Float64() //nolint:gosec
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
		Individuals: individuals,
		ScoringFunc: sf,
		rng:         rng,
	}, nil
}

// Evaluate updates the Score of every individual in the population by calling
// ScoringFunc concurrently. It is a no-op when ScoringFunc is nil.
func (g *Generation[T]) Evaluate() {
	if g.ScoringFunc == nil {
		return
	}
	type result struct {
		idx   int
		score float64
	}
	ch := make(chan result, len(g.Individuals))
	var wg sync.WaitGroup
	for i, ind := range g.Individuals {
		wg.Add(1)
		go func(idx int, genome core.Genome[T]) {
			defer wg.Done()
			ch <- result{idx: idx, score: g.ScoringFunc(genome)}
		}(i, ind.Genome)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for r := range ch {
		g.Individuals[r.idx].Score = r.score
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
	// Map effective scores into the [0.1, 1.0] range so all weights are
	// strictly positive.  A weight of exactly 0 would stall the beta loop.
	minEff := math.Inf(1)
	maxEff := math.Inf(-1)
	for _, ind := range g.Individuals {
		e := g.effectiveScore(ind.Score)
		if e < minEff {
			minEff = e
		}
		if e > maxEff {
			maxEff = e
		}
	}
	scale := maxEff - minEff
	if scale <= 0 {
		scale = 1.0
	}
	scaledScore := func(score float64) float64 {
		return 0.1 + (g.effectiveScore(score)-minEff)/scale
	}

	result := make([]Individual[T], 0, len(g.Individuals))
	idx := g.randIntn(len(g.Individuals))
	beta := 0.0
	for range g.Individuals {
		beta += g.randFloat64() * 2.0
		s := scaledScore(g.Individuals[idx].Score)
		for beta > s {
			beta -= s
			idx = (idx + 1) % len(g.Individuals)
			s = scaledScore(g.Individuals[idx].Score)
		}
		result = append(result, g.Individuals[idx].Dup())
	}
	g.Individuals = result
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
//  6. Restores the saved best individual to position 0 (elitism guarantee).
//
// If the loop completes without triggering the stopping criterion, Evolve
// performs a final Evaluate and returns BestIndividual.
//
// Genetic operators (mutation, recombination, transposition) are wired in
// subsequent milestones; only selection and elitism are applied here.
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
		// Elitism: ensure the best individual survives to the next generation.
		g.Individuals[0] = saveBest
	}
	g.Evaluate()
	return g.BestIndividual()
}
