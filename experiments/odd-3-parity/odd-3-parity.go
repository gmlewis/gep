// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// odd-3-parity is a simple experiment to run the GEP algorithm using the Boolean logic package.
// Given a set of input functions (Not, And, and Or), this solves how to create an
// odd 3 parity function from those basic building blocks.
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

var parityTests = []struct {
	in  []bool
	out bool
}{
	{[]bool{false, false, false}, false},
	{[]bool{false, false, true}, true},
	{[]bool{false, true, false}, true},
	{[]bool{false, true, true}, false},
	{[]bool{true, false, false}, true},
	{[]bool{true, false, true}, false},
	{[]bool{true, true, false}, false},
	{[]bool{true, true, true}, true},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func validateParity(g *genome.Genome) float64 {
	correct := 0
	for _, n := range parityTests {
		r := g.EvalBool(n.in)
		if r == n.out {
			correct++
		}
	}
	return 1000.0 * float64(correct) / float64(len(parityTests))
}

func main() {
	funcs := []gene.FuncWeight{
		{"Not", 10},
		{"And", 20},
		{"Or", 20},
	}
	numIn := len(parityTests[0].in)
	e := model.New(funcs, functions.Bool, 30, 7, 3, numIn, 0, "And", validateParity)
	s := e.Evolve(10000)

	// Write out the Go source code for the solution.
	gr, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		log.Printf("unable to load Boolean grammar: %v", err)
	}
	fmt.Printf("\n// gepModel is auto-generated Go source code for the\n")
	fmt.Printf("// odd-3-parity solution karva expression:\n// %q, score=%v\n", s, validateParity(s))
	s.Write(os.Stdout, gr)
}
