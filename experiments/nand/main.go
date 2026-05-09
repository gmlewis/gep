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
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/evolution"
	evolutionMutation "github.com/gmlewis/gep/v2/evolution/mutation"
	"github.com/gmlewis/gep/v2/functions"
	boolNodes "github.com/gmlewis/gep/v2/functions/bool_nodes"
	"github.com/gmlewis/gep/v2/genome"
	"github.com/gmlewis/gep/v2/grammars"
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

func validateNand(g core.Genome[bool]) float64 {
	correct := 0
	for _, n := range nandTests {
		r, err := g.Eval(n.in)
		if err != nil {
			return 0
		}
		if r == n.out {
			correct++
		}
	}
	return 1000.0 * float64(correct) / float64(len(nandTests))
}

func main() {
	funcs := []string{
		"Not",
		"And",
		"Or",
	}
	fm := make(functions.FuncMap, len(funcs))
	for _, sym := range funcs {
		fn, ok := boolNodes.BoolAllGates[sym]
		if !ok {
			log.Fatalf("unsupported boolean function %q", sym)
		}
		fm[sym] = fn
	}
	cat, err := boolNodes.CatalogFrom(fm)
	if err != nil {
		log.Fatalf("CatalogFrom failed: %v", err)
	}

	linkNode, ok := boolNodes.BoolAllGates["Or"]
	if !ok {
		log.Fatal(errors.New(`link function "Or" not found`))
	}
	link, err := core.NewLinkFunc[bool]("Or", linkNode.BoolFunction)
	if err != nil {
		log.Fatalf("NewLinkFunc failed: %v", err)
	}

	numIn := len(nandTests[0].in)
	population, err := evolution.New(cat, 30, 7, 1, numIn, 0, link, validateNand)
	if err != nil {
		log.Fatalf("New failed: %v", err)
	}
	population.MutationConfig = evolutionMutation.Config{
		PointMutationRate: 0.044,
		InversionRate:     0.1,
	}
	solution := population.Evolve(1000)
	legacySolution, err := genome.NewFromCoreBool(solution.Genome)
	if err != nil {
		log.Fatalf("NewFromCoreBool failed: %v", err)
	}

	// Write out the Go source code for the solution.
	gr, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		log.Printf("unable to load Boolean grammar: %v", err)
	}

	fmt.Printf("\n// gepModel is auto-generated Go source code for the\n")
	fmt.Printf("// nand solution karva expression:\n// %v\n", legacySolution)
	if err := legacySolution.Write(os.Stdout, gr); err != nil {
		log.Fatalf("Write failed: %v", err)
	}
}
