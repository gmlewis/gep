// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package genome

import (
	"fmt"

	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
)

// EvalBool evaluates the genome as a boolean expression and returns the result.
// in represents the boolean inputs available to the genome.
// fm is the map of available boolean functions to the genome.
func (g *Genome) EvalBool(in []bool) (bool, error) {
	lf, ok := bn.BoolAllGates[g.LinkFunc]
	if !ok {
		return false, fmt.Errorf("unable to find linking function: %v", g.LinkFunc)
	}
	result, err := g.Genes[0].EvalBool(in)
	if err != nil {
		return false, err
	}
	for _, gene := range g.Genes[1:] {
		next, err := gene.EvalBool(in)
		if err != nil {
			return false, err
		}
		x := []bool{result, next}
		result = lf.BoolFunction(x)
	}
	return result, nil
}
