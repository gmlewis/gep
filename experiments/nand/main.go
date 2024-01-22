// -*- compile-command: "go run main.go"; -*-
// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// nand is a simple experiment to run the GEP algorithm using the Boolean logic package.
// Given a set of input functions (Not, And, and Or), this solves how to create a NAND gate
// from those basic building blocks. This experiment usually converges to a solution within
// the first generation of evolution.
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

var nandTests = []struct {
	in  []bool
	out bool
}{
	{[]bool{false, false}, true},
	{[]bool{false, true}, true},
	{[]bool{true, false}, true},
	{[]bool{true, true}, false},
}

func init() {
	rand.Seed(time.Now().UnixNano())
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
	numIn := len(nandTests[0].in)
	population := model.New(funcs, functions.Bool, 30, 7, 1, numIn, 0, "Or", validateNand)
	solution := population.Evolve(1000)

	// Write out the Go source code for the solution.
	gr, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		log.Printf("unable to load Boolean grammar: %v", err)
	}

	fmt.Printf("\n// gepModel is auto-generated Go source code for the\n")
	fmt.Printf("// nand solution karva expression:\n// %q\n", solution)
	solution.Write(os.Stdout, gr)
}
