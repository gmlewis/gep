// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package model provides the complete representation of the model for a given GEP problem.
package model

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/gmlewis/gep/v2/functions"
	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
	in "github.com/gmlewis/gep/v2/functions/int_nodes"
	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/genome"
)

// Generation represents one complete generation of the model.
type Generation struct {
	Individuals []*genome.Genome
	Funcs       []gene.FuncWeight
	ScoringFunc genome.ScoringFunc
	// MinimizeScore controls score orientation in selection and best-genome
	// tracking. When false (default), higher scores are better. When true,
	// lower scores are treated as better.
	MinimizeScore bool
	// StopFunc is an optional callback that is called after each generation's
	// best genome is evaluated. When StopFunc returns true, evolution halts
	// early and the best genome is returned. When nil, the default stopping
	// criterion of Score >= 1000 is used.
	StopFunc func(*genome.Genome) bool

	debug bool
	// rng is an optional seeded random source for reproducible evolution.
	// When nil (as produced by New), the global math/rand source is used.
	// Use NewWithSeed to obtain a fully deterministic Generation.
	rng *rand.Rand
}

// randIntn returns a random int in [0, n) using g.rng when set, else the global source.
func (g *Generation) randIntn(n int) int {
	if g.rng != nil {
		return g.rng.Intn(n)
	}
	return rand.Intn(n)
}

// randFloat64 returns a random float64 in [0.0, 1.0) using g.rng when set, else the global source.
func (g *Generation) randFloat64() float64 {
	if g.rng != nil {
		return g.rng.Float64()
	}
	return rand.Float64()
}

// effectiveScore maps raw scores to an internal "higher is better" scale used
// only by selection logic so the same code can support both maximization and
// minimization.
func (g *Generation) effectiveScore(score float64) float64 {
	if g.MinimizeScore {
		return -score
	}
	return score
}

// New creates a new random generation of the model.
// fs is a slice of function weights.
// funcType is the underlying function type (no generics).
// numIndividuals is the number of genomes to use to populate this generation of the model.
// headSize is the number of head symbols to use in a genome.
// numGenesPerGenome is the number of genes to use per genome.
// numTerminals is the number of terminals (inputs) to use within each gene.
// numConstants is the number of constants (inputs) to use within each gene.
// linkFunc is the linking function used to combine the genes within a genome.
// sf is the scoring (or fitness) function.
func New(
	fs []gene.FuncWeight,
	funcType functions.FuncType,
	numIndividuals,
	headSize,
	numGenesPerGenome,
	numTerminals,
	numConstants int,
	linkFunc string,
	sf genome.ScoringFunc,
	debug bool) (*Generation, error) {
	return newGeneration(nil, fs, funcType, numIndividuals, headSize, numGenesPerGenome, numTerminals, numConstants, linkFunc, sf, debug)
}

// NewWithSeed creates a new random generation of the model seeded for deterministic
// reproducibility. Two calls to NewWithSeed with the same seed and identical parameters
// will produce identical populations and identical evolution trajectories.
func NewWithSeed(
	seed int64,
	fs []gene.FuncWeight,
	funcType functions.FuncType,
	numIndividuals,
	headSize,
	numGenesPerGenome,
	numTerminals,
	numConstants int,
	linkFunc string,
	sf genome.ScoringFunc,
	debug bool) (*Generation, error) {
	rng := rand.New(rand.NewSource(seed))
	return newGeneration(rng, fs, funcType, numIndividuals, headSize, numGenesPerGenome, numTerminals, numConstants, linkFunc, sf, debug)
}

// newGeneration is the shared internal constructor used by New and NewWithSeed.
func newGeneration(
	rng *rand.Rand,
	fs []gene.FuncWeight,
	funcType functions.FuncType,
	numIndividuals,
	headSize,
	numGenesPerGenome,
	numTerminals,
	numConstants int,
	linkFunc string,
	sf genome.ScoringFunc,
	debug bool) (*Generation, error) {
	r := &Generation{
		Individuals: make([]*genome.Genome, numIndividuals),
		Funcs:       fs,
		ScoringFunc: sf,
		debug:       debug,
		rng:         rng,
	}
	n, err := maxArity(fs, funcType)
	if err != nil {
		return nil, err
	}
	tailSize := headSize*(n-1) + 1
	for i := range r.Individuals {
		genes := make([]*gene.Gene, numGenesPerGenome)
		for j := range genes {
			genes[j] = gene.RandomNew(headSize, tailSize, numTerminals, numConstants, fs, funcType, rng)
		}
		r.Individuals[i] = genome.New(genes, linkFunc, rng)
	}
	return r, nil
}

// Evolve runs the GEP algorithm for the given number of iterations.
// Evolution stops early when StopFunc (if set) returns true for the current best
// genome. When StopFunc is nil, the default criterion of Score >= 1000 is used.
func (g *Generation) Evolve(iterations int) (*genome.Genome, error) {
	// Algorithm flow diagram, figure 3.1, book page 56
	for i := 0; i < iterations; i++ {
		// fmt.Printf("Iteration #%v...\n", i)
		bestGenome, err := g.getBest() // Preserve the best genome
		if err != nil {
			return nil, err
		}
		stop := false
		if g.StopFunc != nil {
			stop = g.StopFunc(bestGenome)
		} else {
			stop = bestGenome.Score >= 1000.0
		}
		if stop {
			return bestGenome, nil
		}
		// fmt.Printf("Best genome (score %v): %v\n", bestGenome.Score, *bestGenome)
		saveCopy, err := bestGenome.Dup()
		if err != nil {
			return nil, err
		}
		if err := g.replication(); err != nil { // Section 3.3.1, book page 75
			return nil, err
		}
		if err := g.mutation(); err != nil { // Section 3.3.2, book page 77
			return nil, err
		}
		// g.isTransposition()
		// g.risTransposition()
		// g.geneTransposition()
		// g.onePointRecombination()
		// g.twoPointRecombination()
		// g.geneRecombination()
		// Now that replication is done, restore the best genome (aka "elitism")
		g.Individuals[0] = saveCopy
	}
	return g.getBest()
}

// replication replaces all individuals in the population by
// selecting random individuals (weighted by individual
// scores) using the roulette wheel selection algorithm.
// It duplicates those individuals, replacing the population with
// the new collection of individuals.
//
// This algorithm is slightly tricky because the scores can have
// any possible float64 range.
func (g *Generation) replication() error {
	// roulette wheel selection - see www.youtube.com/watch?v=aHLslaWO-AQ
	minWeight, maxWeight := 0.0, 0.0
	for i, v := range g.Individuals {
		score := g.effectiveScore(v.Score)
		if i == 0 || score > maxWeight {
			maxWeight = score
		}
		if i == 0 || score < minWeight {
			minWeight = score
		}
	}
	// Map minWidth->maxWeight to 0.1->1
	// Note that a weight (scaledScore) of 0 would create an
	// infinite loop due to the `beta -= scaledScore`.
	weightScale := maxWeight - minWeight
	if weightScale <= 0 {
		weightScale = 1.0
	}
	f := func(v float64) float64 { return 0.1 + (v-minWeight)/weightScale }

	result := make([]*genome.Genome, 0, len(g.Individuals))
	index := g.randIntn(len(g.Individuals))
	beta := 0.0
	for i := 0; i < len(g.Individuals); i++ {
		beta += g.randFloat64() * 2.0
		scaledScore := f(g.effectiveScore(g.Individuals[index].Score))
		for beta > scaledScore {
			beta -= scaledScore
			index = (index + 1) % len(g.Individuals)
		}
		dup, err := g.Individuals[index].Dup()
		if err != nil {
			return err
		}
		result = append(result, dup)
	}
	g.Individuals = result
	return nil
}

func (g *Generation) singleMutation(index int) error {
	gen := g.Individuals[index]
	// Determine the total number of mutations to perform within the genome
	numMutations := 1 + g.randIntn(2)
	// fmt.Printf("\nMutating genome #%v %v times, before:\n%v\n", genomeNum, numMutations, genome)
	if err := gen.Mutate(numMutations); err != nil {
		return err
	}
	// fmt.Printf("after:\n%v\n", gen)
	return nil
}

func (g *Generation) mutation() error {
	// Determine the total number of individuals to mutate
	numMutations := 1 + g.randIntn(len(g.Individuals)-1)
	for i := 0; i < numMutations; i++ {
		// Pick a random genome
		genomeNum := g.randIntn(len(g.Individuals))
		if err := g.singleMutation(genomeNum); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generation) singleCrossover(idx1, idx2 int) error {
	genome1 := g.Individuals[idx1]
	genome2 := g.Individuals[idx2]

	// Pick a random gene from genome1 and genome2
	geneIdx1 := g.randIntn(len(genome1.Genes))
	gene1 := genome1.Genes[geneIdx1]
	geneIdx2 := g.randIntn(len(genome2.Genes))
	gene2 := genome2.Genes[geneIdx2]

	if len(gene1.Symbols) != len(gene2.Symbols) || gene1.HeadSize != gene2.HeadSize {
		return fmt.Errorf("programming error: gene1: %v symbols (headSize=%v), gene2: %v symbols (headSize=%v)", len(gene1.Symbols), gene1.HeadSize, len(gene2.Symbols), gene2.HeadSize)
	}

	// Pick a random location within the head of both gene's symbols to crossover.
	// The length of both symbols slices will stay the same.
	symbolIdx := g.randIntn(gene1.HeadSize)
	head1, tail1 := gene1.Symbols[:gene1.HeadSize], gene1.Symbols[gene1.HeadSize:]
	head2, tail2 := gene2.Symbols[:gene2.HeadSize], gene2.Symbols[gene2.HeadSize:]
	newSyms1 := append([]string{}, head2[symbolIdx:]...)
	newSyms1 = append(newSyms1, head1[:symbolIdx]...)
	newSyms1 = append(newSyms1, tail1...)
	newSyms2 := append([]string{}, head1[symbolIdx:]...)
	newSyms2 = append(newSyms2, head2[:symbolIdx]...)
	newSyms2 = append(newSyms2, tail2...)

	if len(newSyms1) != len(newSyms2) || len(newSyms1) != len(gene1.Symbols) {
		return fmt.Errorf("programming error: newSyms1: %v symbols, newSyms2: %v symbols, gene1: %v symbols", len(newSyms1), len(newSyms2), len(gene1.Symbols))
	}

	gene1.Symbols = newSyms1
	gene2.Symbols = newSyms2
	gene1.InvalidateCache()
	gene2.InvalidateCache()
	genome1.SymbolMap = nil
	genome2.SymbolMap = nil
	return nil
}

func (g *Generation) crossover() error {
	if len(g.Individuals) < 2 {
		return nil
	}

	// Determine the total number of individuals pairs to crossover
	numCrossovers := 1 + g.randIntn(len(g.Individuals)-1)
	for i := 0; i < numCrossovers; i++ {
		// Pick two random genomes
		genomeNum1 := g.randIntn(len(g.Individuals))
		var genomeNum2 int
		for {
			genomeNum2 = g.randIntn(len(g.Individuals))
			if genomeNum2 != genomeNum1 {
				break
			}
		}
		if err := g.singleCrossover(genomeNum1, genomeNum2); err != nil {
			return err
		}
	}
	return nil
}

// getBest evaluates all individuals and returns a pointer to the best one.
func (g *Generation) getBest() (*genome.Genome, error) {
	if len(g.Individuals) == 0 {
		return nil, fmt.Errorf("no individuals in generation")
	}
	highestEffectiveScore := math.Inf(-1)
	bestGenome := g.Individuals[0]
	for _, gn := range g.Individuals {
		if err := gn.EvaluateWithScore(g.ScoringFunc); err != nil {
			return nil, err
		}
		effectiveScore := g.effectiveScore(gn.Score)
		if effectiveScore > highestEffectiveScore {
			bestGenome = gn
			highestEffectiveScore = effectiveScore
		}
	}
	return bestGenome, nil
}

// maxArity determines the maximum number of input terminals for the given set of symbols.
func maxArity(fs []gene.FuncWeight, funcType functions.FuncType) (int, error) {
	var lookup functions.FuncMap
	switch funcType {
	case functions.Bool:
		lookup = bn.BoolAllGates
	case functions.Int:
		lookup = in.Int
	case functions.Float64:
		lookup = mn.Math
	default:
		return 0, fmt.Errorf("unknown funcType: %v", funcType)
	}

	r := 0
	var errs []error
	for _, f := range fs {
		if fn, ok := lookup[f.Symbol]; ok {
			if fn.Terminals() > r {
				r = fn.Terminals()
			}
		} else {
			errs = append(errs, fmt.Errorf("unable to find symbol %v in function map", f.Symbol))
		}
	}
	return r, errors.Join(errs...)
}
