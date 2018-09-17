// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package genome

import (
	"log"

	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
)

// EvalMath evaluates the genome as a floating-point expression and returns the result.
// in represents the float64 inputs available to the genome.
func (g *Genome) EvalMath(in []float64) float64 {
	lf, ok := mn.Math[g.LinkFunc]
	if !ok {
		log.Printf("Unable to find linking function: %v", g.LinkFunc)
		return 0.0
	}
	result := g.Genes[0].EvalMath(in)
	for i := 1; i < len(g.Genes); i++ {
		x := []float64{result, g.Genes[i].EvalMath(in)}
		result = lf.Float64Function(x)
	}
	return result
}
