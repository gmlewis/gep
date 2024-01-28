// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/genome"
	"github.com/gmlewis/gep/v2/model"
)

var nandTests = []struct {
	in  []bool
	out bool
}{
	{[]bool{false, false}, true},
	{[]bool{false, true}, true},
	{[]bool{true, false}, true},
	{[]bool{true, true}, false},
}

func validateNand(g *genome.Genome) float64 {
	correct := 0
	for _, n := range nandTests {
		r := g.EvalBool(n.in)
		if r == n.out {
			correct++
		}
	}
	return 1000.0 * float64(correct) / float64(len(nandTests))
}

func main() {
	funcs := []gene.FuncWeight{
		{Symbol: "Not", Weight: 1},
		{Symbol: "And", Weight: 5},
		{Symbol: "Or", Weight: 5},
	}
	e := model.New(funcs, functions.Bool, 30, 7, 1, 2, 0, "Or", validateNand, false)
	s := e.Evolve(1000)
	fmt.Printf("nand solution: %v, score=%v\n", s, validateNand(s))
}
