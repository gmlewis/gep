// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package genome

import (
	"fmt"

	intN "github.com/gmlewis/gep/v2/functions/int_nodes"
)

// EvalInt evaluates the genome as an integer expression and returns the result.
// in represents the int inputs available to the genome.
func (g *Genome) EvalInt(in []int) (int, error) {
	lf, ok := intN.Int[g.LinkFunc]
	if !ok {
		return 0, fmt.Errorf("unable to find linking function: %v", g.LinkFunc)
	}
	result, err := g.Genes[0].EvalInt(in)
	if err != nil {
		return 0, err
	}
	for i := 1; i < len(g.Genes); i++ {
		next, err := g.Genes[i].EvalInt(in)
		if err != nil {
			return 0, err
		}
		x := []int{result, next}
		result = lf.IntFunction(x)
	}
	return result, nil
}

// EvalIntTuple evaluates the genome by evaluating each gene and assigning
// its output to each element of the tuple.
func (g *Genome) EvalIntTuple(in []int) ([]int, error) {
	result := make([]int, len(g.Genes))

	// var wg sync.WaitGroup
	for i := 0; i < len(g.Genes); i++ {
		// wg.Add(1)
		// go func(i int) {
		v, err := g.Genes[i].EvalInt(in)
		if err != nil {
			return nil, err
		}
		result[i] = v
		// wg.Done()
		// }(i)
	}
	// wg.Wait()

	return result, nil
}
