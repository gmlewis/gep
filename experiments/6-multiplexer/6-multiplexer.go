// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// 6-multiplexer is a simple experiment to run the GEP algorithm using the Boolean logic package.
// Given a set of input functions (Not, And, Or, Nand, and Nor), this solves how to create a
// 6-way multiplexer from those basic building blocks.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/genome"
	"github.com/gmlewis/gep/v2/grammars"
	"github.com/gmlewis/gep/v2/model"
)

var multiTests = []struct {
	in  []bool
	out bool
}{
	{[]bool{false, false, false, false, false, false}, false},
	{[]bool{false, false, false, false, false, true}, false},
	{[]bool{false, false, false, false, true, false}, false},
	{[]bool{false, false, false, false, true, true}, false},
	{[]bool{false, false, false, true, false, false}, false},
	{[]bool{false, false, false, true, false, true}, false},
	{[]bool{false, false, false, true, true, false}, false},
	{[]bool{false, false, false, true, true, true}, false},
	{[]bool{false, false, true, false, false, false}, true},
	{[]bool{false, false, true, false, false, true}, true},
	{[]bool{false, false, true, false, true, false}, true},
	{[]bool{false, false, true, false, true, true}, true},
	{[]bool{false, false, true, true, false, false}, true},
	{[]bool{false, false, true, true, false, true}, true},
	{[]bool{false, false, true, true, true, false}, true},
	{[]bool{false, false, true, true, true, true}, true},
	{[]bool{false, true, false, false, false, false}, false},
	{[]bool{false, true, false, false, false, true}, false},
	{[]bool{false, true, false, false, true, false}, false},
	{[]bool{false, true, false, false, true, true}, false},
	{[]bool{false, true, false, true, false, false}, true},
	{[]bool{false, true, false, true, false, true}, true},
	{[]bool{false, true, false, true, true, false}, true},
	{[]bool{false, true, false, true, true, true}, true},
	{[]bool{false, true, true, false, false, false}, false},
	{[]bool{false, true, true, false, false, true}, false},
	{[]bool{false, true, true, false, true, false}, false},
	{[]bool{false, true, true, false, true, true}, false},
	{[]bool{false, true, true, true, false, false}, true},
	{[]bool{false, true, true, true, false, true}, true},
	{[]bool{false, true, true, true, true, false}, true},
	{[]bool{false, true, true, true, true, true}, true},
	{[]bool{true, false, false, false, false, false}, false},
	{[]bool{true, false, false, false, false, true}, false},
	{[]bool{true, false, false, false, true, false}, true},
	{[]bool{true, false, false, false, true, true}, true},
	{[]bool{true, false, false, true, false, false}, false},
	{[]bool{true, false, false, true, false, true}, false},
	{[]bool{true, false, false, true, true, false}, true},
	{[]bool{true, false, false, true, true, true}, true},
	{[]bool{true, false, true, false, false, false}, false},
	{[]bool{true, false, true, false, false, true}, false},
	{[]bool{true, false, true, false, true, false}, true},
	{[]bool{true, false, true, false, true, true}, true},
	{[]bool{true, false, true, true, false, false}, false},
	{[]bool{true, false, true, true, false, true}, false},
	{[]bool{true, false, true, true, true, false}, true},
	{[]bool{true, false, true, true, true, true}, true},
	{[]bool{true, true, false, false, false, false}, false},
	{[]bool{true, true, false, false, false, true}, true},
	{[]bool{true, true, false, false, true, false}, false},
	{[]bool{true, true, false, false, true, true}, true},
	{[]bool{true, true, false, true, false, false}, false},
	{[]bool{true, true, false, true, false, true}, true},
	{[]bool{true, true, false, true, true, false}, false},
	{[]bool{true, true, false, true, true, true}, true},
	{[]bool{true, true, true, false, false, false}, false},
	{[]bool{true, true, true, false, false, true}, true},
	{[]bool{true, true, true, false, true, false}, false},
	{[]bool{true, true, true, false, true, true}, true},
	{[]bool{true, true, true, true, false, false}, false},
	{[]bool{true, true, true, true, false, true}, true},
	{[]bool{true, true, true, true, true, false}, false},
	{[]bool{true, true, true, true, true, true}, true},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func validateMulti(g *genome.Genome) float64 {
	correct := 0
	for _, n := range multiTests {
		r := g.EvalBool(n.in)
		if r == n.out {
			correct++
		}
	}
	return 1000.0 * float64(correct) / float64(len(multiTests))
}

func main() {
	funcs := []gene.FuncWeight{
		{Symbol: "Not", Weight: 10},
		{Symbol: "And", Weight: 20},
		{Symbol: "Or", Weight: 20},
		{Symbol: "Nand", Weight: 20},
		{Symbol: "Nor", Weight: 20},
	}
	numIn := len(multiTests[0].in)
	e := model.New(funcs, functions.Bool, 30, 8, 4, numIn, 0, "And", validateMulti)
	s := e.Evolve(20000)

	// Write out the Go source code for the solution.
	gr, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		log.Printf("unable to load Boolean grammar: %v", err)
	}
	fmt.Printf("\n// gepModel is auto-generated Go source code for the\n")
	fmt.Printf("// 6-multiplexer solution karva expression:\n// %q, score=%v\n", s, validateMulti(s))
	s.Write(os.Stdout, gr)
}
