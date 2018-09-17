// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package genome

import (
	"log"

	intN "github.com/gmlewis/gep/v2/functions/int_nodes"
)

// EvalInt evaluates the genome as a floating-point expression and returns the result.
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
