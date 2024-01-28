// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package model provides the complete representation of the model for a given GEP problem.
package model

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

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

	debug bool
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
	debug bool) *Generation {
	r := &Generation{
		Individuals: make([]*genome.Genome, numIndividuals),
		Funcs:       fs,
		ScoringFunc: sf,
		debug:       debug,
	}
	n := maxArity(fs, funcType)
	tailSize := headSize*(n-1) + 1
	for i := range r.Individuals {
		genes := make([]*gene.Gene, numGenesPerGenome)
		for j := range genes {
			genes[j] = gene.RandomNew(headSize, tailSize, numTerminals, numConstants, fs, funcType)
		}
		r.Individuals[i] = genome.New(genes, linkFunc)
	}
	return r
}

// Evolve runs the GEP algorithm for the given number of iterations, or until a score of 1000 (or more) is reached.
func (g *Generation) Evolve(iterations int) *genome.Genome {
	// Algorithm flow diagram, figure 3.1, book page 56
	for i := 0; i < iterations; i++ {
		// fmt.Printf("Iteration #%v...\n", i)
		bestGenome := g.getBest() // Preserve the best genome
		if bestGenome.Score >= 1000.0 {
			fmt.Printf("Stopping after generation #%v\n", i)
			return bestGenome
		}
		// fmt.Printf("Best genome (score %v): %v\n", bestGenome.Score, *bestGenome)
		saveCopy := bestGenome.Dup()
		g.replication() // Section 3.3.1, book page 75
		g.mutation()    // Section 3.3.2, book page 77
		// g.isTransposition()
		// g.risTransposition()
		// g.geneTransposition()
		// g.onePointRecombination()
		// g.twoPointRecombination()
		// g.geneRecombination()
		// Now that replication is done, restore the best genome (aka "elitism")
		g.Individuals[0] = saveCopy
	}
	fmt.Printf("Stopping after generation #%v\n", iterations)
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
func (g *Generation) replication() {
	// roulette wheel selection - see www.youtube.com/watch?v=aHLslaWO-AQ
	minWeight, maxWeight := 0.0, 0.0
	for i, v := range g.Individuals {
		if i == 0 || v.Score > maxWeight {
			maxWeight = v.Score
		}
		if i == 0 || v.Score < minWeight {
			minWeight = v.Score
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
	index := rand.Intn(len(g.Individuals))
	beta := 0.0
	for i := 0; i < len(g.Individuals); i++ {
		beta += rand.Float64() * 2.0
		scaledScore := f(g.Individuals[index].Score)
		for beta > scaledScore {
			beta -= scaledScore
			index = (index + 1) % len(g.Individuals)
		}
		result = append(result, g.Individuals[index].Dup())
	}
	g.Individuals = result
}

func (g *Generation) singleMutation(index int) {
	gen := g.Individuals[index]
	// Determine the total number of mutations to perform within the genome
	numMutations := 1 + rand.Intn(2)
	// fmt.Printf("\nMutating genome #%v %v times, before:\n%v\n", genomeNum, numMutations, genome)
	gen.Mutate(numMutations)
	// fmt.Printf("after:\n%v\n", gen)
}

func (g *Generation) mutation() {
	// Determine the total number of individuals to mutate
	numMutations := 1 + rand.Intn(len(g.Individuals)-1)
	for i := 0; i < numMutations; i++ {
		// Pick a random genome
		genomeNum := rand.Intn(len(g.Individuals))
		g.singleMutation(genomeNum)
	}
}

func (g *Generation) singleCrossover(idx1, idx2 int) {
	genome1 := g.Individuals[idx1]
	genome2 := g.Individuals[idx2]

	// Pick a random gene from genome1 and genome2
	geneIdx1 := rand.Intn(len(genome1.Genes))
	gene1 := genome1.Genes[geneIdx1]
	geneIdx2 := rand.Intn(len(genome2.Genes))
	gene2 := genome2.Genes[geneIdx2]

	if len(gene1.Symbols) != len(gene2.Symbols) || gene1.HeadSize != gene2.HeadSize {
		log.Fatalf("programming error: gene1: %v symbols (headSize=%v), gene2: %v symbols (headSize=%v)", len(gene1.Symbols), gene1.HeadSize, len(gene2.Symbols), gene2.HeadSize)
	}

	// Pick a random location within the head of both gene's symbols to crossover.
	// The length of both symbols slices will stay the same.
	symbolIdx := rand.Intn(gene1.HeadSize)
	head1, tail1 := gene1.Symbols[:gene1.HeadSize], gene1.Symbols[gene1.HeadSize:]
	head2, tail2 := gene2.Symbols[:gene2.HeadSize], gene2.Symbols[gene2.HeadSize:]
	newSyms1 := append([]string{}, head2[symbolIdx:]...)
	newSyms1 = append(newSyms1, head1[:symbolIdx]...)
	newSyms1 = append(newSyms1, tail1...)
	newSyms2 := append([]string{}, head1[symbolIdx:]...)
	newSyms2 = append(newSyms2, head2[:symbolIdx]...)
	newSyms2 = append(newSyms2, tail2...)

	if len(newSyms1) != len(newSyms2) || len(newSyms1) != len(gene1.Symbols) {
		log.Fatalf("programming error: newSyms1: %v symbols, newSyms2: %v symbols, gene1: %v symbols", len(newSyms1), len(newSyms2), len(gene1.Symbols))
	}

	if g.debug {
		log.Printf("singleCrossover:\nbefore genome[%v].gene[%v]=%v\nbefore genome[%v].gene[%v]=%v\nafter genome[%v].gene[%v]=%v\nafter genome[%v].gene[%v]=%v",
			idx1, geneIdx1, strings.Join(gene1.Symbols, "."),
			idx2, geneIdx2, strings.Join(gene2.Symbols, "."),
			idx1, geneIdx1, strings.Join(newSyms1, "."),
			idx2, geneIdx2, strings.Join(newSyms2, "."))
	}

	gene1.Symbols = newSyms1
	gene2.Symbols = newSyms2
}

func (g *Generation) crossover() {
	if len(g.Individuals) < 2 {
		return
	}

	// Determine the total number of individuals pairs to crossover
	numCrossovers := 1 + rand.Intn(len(g.Individuals)-1)
	for i := 0; i < numCrossovers; i++ {
		// Pick two random genomes
		genomeNum1 := rand.Intn(len(g.Individuals))
		var genomeNum2 int
		for {
			genomeNum2 = rand.Intn(len(g.Individuals))
			if genomeNum2 != genomeNum1 {
				break
			}
		}
		g.singleCrossover(genomeNum1, genomeNum2)
	}
}

// getBest evaluates all individuals and returns a pointer to the best one.
func (g *Generation) getBest() *genome.Genome {
	bestScore := 0.0
	bestGenome := g.Individuals[0]
	c := make(chan *genome.Genome)
	for i := 0; i < len(g.Individuals); i++ { // Evaluate individuals concurrently
		go g.Individuals[i].EvaluateWithScore(g.ScoringFunc, c)
	}
	for i := 0; i < len(g.Individuals); i++ { // Collect and return the highest scoring Genome
		gn := <-c
		if gn.Score > bestScore {
			bestGenome = gn
			bestScore = gn.Score
		}
	}
	return bestGenome
}

// maxArity determines the maximum number of input terminals for the given set of symbols.
func maxArity(fs []gene.FuncWeight, funcType functions.FuncType) int {
	var lookup functions.FuncMap
	switch funcType {
	case functions.Bool:
		lookup = bn.BoolAllGates
	case functions.Int:
		lookup = in.Int
	case functions.Float64:
		lookup = mn.Math
	default:
		log.Fatalf("unknown funcType: %v", funcType)
	}

	r := 0
	for _, f := range fs {
		if fn, ok := lookup[f.Symbol]; ok {
			if fn.Terminals() > r {
				r = fn.Terminals()
			}
		} else {
			log.Printf("unable to find symbol %v in function map", f.Symbol)
		}
	}
	return r
}
