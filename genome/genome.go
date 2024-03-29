// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package genome provides the basis for a single GEP genome
// which represents a single individual in a population.
// A genome consists of one or more genes.
// Each gene within an individual is used to generate
// the outputs (or actions) of that individual.
// A link function determines how the genes are combined
// or mapped to the output or actions of the genome.
package genome

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/grammars"
)

// Genome contains the genes that make up the genome.
// It also provides the linking function and the score that results
// from evaluating the genome against the fitness function.
type Genome struct {
	Genes    []*gene.Gene
	LinkFunc string
	Score    float64

	SymbolMap map[string]int // do not use directly.  Use SymbolCount() instead.
}

// New creates a new genome from the given genes and linking function.
func New(genes []*gene.Gene, linkFunc string) *Genome {
	return &Genome{Genes: genes, LinkFunc: linkFunc}
}

func merge(dst *map[string]int, src map[string]int) {
	for k, v := range src {
		(*dst)[k] += v
	}
}

// SymbolCount returns the count of the number of times the symbol
// is actually used in the Genome.
// Note that this count is typically different from the number
// of times the symbol appears in the Karva expression.  This can be
// a handy metric to assist in the fitness evaluation of a Genome.
func (g *Genome) SymbolCount(sym string) int {
	if g.SymbolMap == nil {
		g.SymbolMap = make(map[string]int)
		g.SymbolMap[g.LinkFunc] = len(g.Genes) - 1
		for i := 0; i < len(g.Genes); i++ {
			g.Genes[i].SymbolCount(sym) // force evaluation
			m := g.Genes[i].SymbolMap
			merge(&(g.SymbolMap), m)
		}
	}
	return g.SymbolMap[sym]
}

// String returns the Karva representation of the genome.
func (g Genome) String() string {
	var result []string
	for _, gene := range g.Genes {
		result = append(result, gene.String())
	}
	return fmt.Sprintf("%v, score=%v", strings.Join(result, "|"+g.LinkFunc+"|"), g.Score)
}

// Expression returns the expression of the genome.
func (g Genome) Expression(grammar *grammars.Grammar, helpers grammars.HelperMap) (string, error) {
	var result []string
	for _, gene := range g.Genes {
		s, err := gene.Expression(grammar, helpers)
		if err != nil {
			return "", err
		}
		result = append(result, s)
	}

	return fmt.Sprintf("%v, score=%v", strings.Join(result, " "+g.LinkFunc+" "), g.Score), nil
}

// DotGraph returns a graphviz "dot" language representation of the genome.
func (g Genome) DotGraph() string {
	var lines []string
	// for i, gene := range g.Genes {
	// 	karva := gene.DotGraph()
	// }
	return strings.Join(lines, "\n")
}

// Mutate mutates a genome by performing numMutations random symbol exchanges within the genome.
func (g *Genome) Mutate(numMutations int) {
	for i := 0; i < numMutations; i++ {
		n := rand.Intn(len(g.Genes))
		// fmt.Printf("\nMutating gene #%v, before:\n%v\n", n, g.Genes[n])
		g.Genes[n].Mutate()
		// fmt.Printf("after:\n%v\n", g.Genes[n])
	}
}

// Dup duplicates the genome into the provided destination genome.
func (g *Genome) Dup() *Genome {
	if g == nil {
		log.Printf("denome.Dup error: src and dst must be non-nil")
		return nil
	}
	dst := &Genome{
		Genes:    make([]*gene.Gene, len(g.Genes)),
		LinkFunc: g.LinkFunc,
		Score:    g.Score,
	}
	for i := range g.Genes {
		dst.Genes[i] = g.Genes[i].Dup()
	}
	return dst
}

// ScoringFunc is the function that is used to evaluate the fitness of the model.
// Typically, a return value of 0 means that the function is nowhere close to being
// a valid solution and a return value of 1000 (or higher) means a perfect solution.
type ScoringFunc func(g *Genome) float64

// EvaluateWithScore scores a genome and sends the result to a channel.
func (g *Genome) EvaluateWithScore(sf ScoringFunc, c chan<- *Genome) {
	if sf == nil {
		log.Fatalf("genome.EvaluateWithScore: ScoringFunc must not be nil")
	}
	g.Score = sf(g)
	c <- g
}

// Evaluate runs the model with the observations and populates the provided action
// based on the link function.
func (g *Genome) Evaluate(observations []int, action any) error {
	result := g.EvalIntTuple(observations)

	switch v := action.(type) {
	case *[]int:
		*v = result
	case *int:
		*v = result[0]
	default:
		return fmt.Errorf("genome.Evaluate: action type '%T' not yet supported", v)
	}
	return nil
}
