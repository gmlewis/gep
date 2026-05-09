// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package genome

import (
	"fmt"

	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
)

// EvalMath evaluates the genome as a floating-point expression and returns the result.
// in represents the float64 inputs available to the genome.
func (g *Genome) EvalMath(in []float64) (float64, error) {
	lf, ok := mn.Math[g.LinkFunc]
	if !ok {
		return 0.0, fmt.Errorf("unable to find linking function: %v", g.LinkFunc)
	}
	result, err := g.Genes[0].EvalMath(in)
	if err != nil {
		return 0.0, err
	}
	for i := 1; i < len(g.Genes); i++ {
		next, err := g.Genes[i].EvalMath(in)
		if err != nil {
			return 0.0, err
		}
		x := []float64{result, next}
		result = lf.Float64Function(x)
	}
	return result, nil
}
