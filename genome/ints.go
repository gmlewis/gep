// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package genome

import (
	"log"
	"sync"

	intN "github.com/gmlewis/gep/v2/functions/int_nodes"
)

// EvalInt evaluates the genome as an integer expression and returns the result.
// in represents the int inputs available to the genome.
func (g *Genome) EvalInt(in []int) int {
	lf, ok := intN.Int[g.LinkFunc]
	if !ok {
		log.Printf("Unable to find linking function: %v", g.LinkFunc)
		return 0.0
	}
	result := g.Genes[0].EvalInt(in)
	for i := 1; i < len(g.Genes); i++ {
		x := []int{result, g.Genes[i].EvalInt(in)}
		result = lf.IntFunction(x)
	}
	return result
}

// EvalIntTuple evaluates the genome by evaluating each gene and assigning
// its output to each element of the tuple.
func (g *Genome) EvalIntTuple(in []int) []int {
	result := make([]int, len(g.Genes))

	var wg sync.WaitGroup
	for i := 0; i < len(g.Genes); i++ {
		wg.Add(1)
		go func(i int) {
			result[i] = g.Genes[i].EvalInt(in)
			wg.Done()
		}(i)
	}
	wg.Wait()

	return result
}
