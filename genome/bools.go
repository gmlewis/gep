// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package genome

import (
	"log"

	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
)

// EvalBool evaluates the genome as a boolean expression and returns the result.
// in represents the boolean inputs available to the genome.
// fm is the map of available boolean functions to the genome.
func (g *Genome) EvalBool(in []bool) bool {
	lf, ok := bn.BoolAllGates[g.LinkFunc]
	if !ok {
		log.Printf("Unable to find linking function: %v", g.LinkFunc)
		return false
	}
	result := g.Genes[0].EvalBool(in)
	for i := 1; i < len(g.Genes); i++ {
		x := []bool{result, g.Genes[i].EvalBool(in)}
		result = lf.BoolFunction(x)
	}
	return result
}
