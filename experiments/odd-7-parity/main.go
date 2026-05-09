// -*- compile-command: "go run main.go"; -*-
// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// odd-7-parity is a simple experiment to run the GEP algorithm using the Boolean logic package.
// Given a set of input functions (Not, And, and Or), this solves how to create an
// odd 7 parity function from those basic building blocks.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gmlewis/gep/v2/codegen"
	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/evolution"
	evolutionMutation "github.com/gmlewis/gep/v2/evolution/mutation"
	boolNodes "github.com/gmlewis/gep/v2/functions/bool_nodes"
	"github.com/gmlewis/gep/v2/grammars"
)

var parityTests = []struct {
	in  []bool
	out bool
}{
	{[]bool{false, false, false, false, false, false, false}, false}, // 1
	{[]bool{false, false, false, false, false, false, true}, true},
	{[]bool{false, false, false, false, false, true, false}, true},
	{[]bool{false, false, false, false, false, true, true}, false},
	{[]bool{false, false, false, false, true, false, false}, true},
	{[]bool{false, false, false, false, true, false, true}, false},
	{[]bool{false, false, false, false, true, true, false}, false},
	{[]bool{false, false, false, false, true, true, true}, true},
	{[]bool{false, false, false, true, false, false, false}, true},
	{[]bool{false, false, false, true, false, false, true}, false},
	{[]bool{false, false, false, true, false, true, false}, false},
	{[]bool{false, false, false, true, false, true, true}, true},
	{[]bool{false, false, false, true, true, false, false}, false},
	{[]bool{false, false, false, true, true, false, true}, true},
	{[]bool{false, false, false, true, true, true, false}, true},
	{[]bool{false, false, false, true, true, true, true}, false},
	{[]bool{false, false, true, false, false, false, false}, true},
	{[]bool{false, false, true, false, false, false, true}, false},
	{[]bool{false, false, true, false, false, true, false}, false},
	{[]bool{false, false, true, false, false, true, true}, true},
	{[]bool{false, false, true, false, true, false, false}, false},
	{[]bool{false, false, true, false, true, false, true}, true},
	{[]bool{false, false, true, false, true, true, false}, true},
	{[]bool{false, false, true, false, true, true, true}, false},
	{[]bool{false, false, true, true, false, false, false}, false},
	{[]bool{false, false, true, true, false, false, true}, true}, // 26
	{[]bool{false, false, true, true, false, true, false}, true},
	{[]bool{false, false, true, true, false, true, true}, false},
	{[]bool{false, false, true, true, true, false, false}, true},
	{[]bool{false, false, true, true, true, false, true}, false},
	{[]bool{false, false, true, true, true, true, false}, false},
	{[]bool{false, false, true, true, true, true, true}, true},
	{[]bool{false, true, false, false, false, false, false}, true},
	{[]bool{false, true, false, false, false, false, true}, false},
	{[]bool{false, true, false, false, false, true, false}, false},
	{[]bool{false, true, false, false, false, true, true}, true},
	{[]bool{false, true, false, false, true, false, false}, false},
	{[]bool{false, true, false, false, true, false, true}, true},
	{[]bool{false, true, false, false, true, true, false}, true}, // 39
	{[]bool{false, true, false, false, true, true, true}, false},
	{[]bool{false, true, false, true, false, false, false}, false},
	{[]bool{false, true, false, true, false, false, true}, true},
	{[]bool{false, true, false, true, false, true, false}, true},
	{[]bool{false, true, false, true, false, true, true}, false},
	{[]bool{false, true, false, true, true, false, false}, true},
	{[]bool{false, true, false, true, true, false, true}, false},
	{[]bool{false, true, false, true, true, true, false}, false},
	{[]bool{false, true, false, true, true, true, true}, true},
	{[]bool{false, true, true, false, false, false, false}, false},
	{[]bool{false, true, true, false, false, false, true}, true},
	{[]bool{false, true, true, false, false, true, false}, true},
	{[]bool{false, true, true, false, false, true, true}, false}, // 52
	{[]bool{false, true, true, false, true, false, false}, true},
	{[]bool{false, true, true, false, true, false, true}, false},
	{[]bool{false, true, true, false, true, true, false}, false},
	{[]bool{false, true, true, false, true, true, true}, true},
	{[]bool{false, true, true, true, false, false, false}, true},
	{[]bool{false, true, true, true, false, false, true}, false},
	{[]bool{false, true, true, true, false, true, false}, false},
	{[]bool{false, true, true, true, false, true, true}, true},
	{[]bool{false, true, true, true, true, false, false}, false},
	{[]bool{false, true, true, true, true, false, true}, true},
	{[]bool{false, true, true, true, true, true, false}, true},
	{[]bool{false, true, true, true, true, true, true}, false}, // 64
	{[]bool{true, false, false, false, false, false, false}, true},
	{[]bool{true, false, false, false, false, false, true}, false},
	{[]bool{true, false, false, false, false, true, false}, false},
	{[]bool{true, false, false, false, false, true, true}, true},
	{[]bool{true, false, false, false, true, false, false}, false},
	{[]bool{true, false, false, false, true, false, true}, true},
	{[]bool{true, false, false, false, true, true, false}, true},
	{[]bool{true, false, false, false, true, true, true}, false},
	{[]bool{true, false, false, true, false, false, false}, false},
	{[]bool{true, false, false, true, false, false, true}, true},
	{[]bool{true, false, false, true, false, true, false}, true},
	{[]bool{true, false, false, true, false, true, true}, false}, // 76
	{[]bool{true, false, false, true, true, false, false}, true},
	{[]bool{true, false, false, true, true, false, true}, false},
	{[]bool{true, false, false, true, true, true, false}, false},
	{[]bool{true, false, false, true, true, true, true}, true},
	{[]bool{true, false, true, false, false, false, false}, false},
	{[]bool{true, false, true, false, false, false, true}, true},
	{[]bool{true, false, true, false, false, true, false}, true},
	{[]bool{true, false, true, false, false, true, true}, false},
	{[]bool{true, false, true, false, true, false, false}, true},
	{[]bool{true, false, true, false, true, false, true}, false},
	{[]bool{true, false, true, false, true, true, false}, false},
	{[]bool{true, false, true, false, true, true, true}, true}, // 88
	{[]bool{true, false, true, true, false, false, false}, true},
	{[]bool{true, false, true, true, false, false, true}, false},
	{[]bool{true, false, true, true, false, true, false}, false},
	{[]bool{true, false, true, true, false, true, true}, true},
	{[]bool{true, false, true, true, true, false, false}, false},
	{[]bool{true, false, true, true, true, false, true}, true},
	{[]bool{true, false, true, true, true, true, false}, true},
	{[]bool{true, false, true, true, true, true, true}, false},
	{[]bool{true, true, false, false, false, false, false}, false},
	{[]bool{true, true, false, false, false, false, true}, true},
	{[]bool{true, true, false, false, false, true, false}, true},
	{[]bool{true, true, false, false, false, true, true}, false}, // 100
	{[]bool{true, true, false, false, true, false, false}, true},
	{[]bool{true, true, false, false, true, false, true}, false},
	{[]bool{true, true, false, false, true, true, false}, false},
	{[]bool{true, true, false, false, true, true, true}, true},
	{[]bool{true, true, false, true, false, false, false}, true},
	{[]bool{true, true, false, true, false, false, true}, false},
	{[]bool{true, true, false, true, false, true, false}, false},
	{[]bool{true, true, false, true, false, true, true}, true},
	{[]bool{true, true, false, true, true, false, false}, false},
	{[]bool{true, true, false, true, true, false, true}, true},
	{[]bool{true, true, false, true, true, true, false}, true},
	{[]bool{true, true, false, true, true, true, true}, false}, // 112
	{[]bool{true, true, true, false, false, false, false}, true},
	{[]bool{true, true, true, false, false, false, true}, false},
	{[]bool{true, true, true, false, false, true, false}, false},
	{[]bool{true, true, true, false, false, true, true}, true},
	{[]bool{true, true, true, false, true, false, false}, false},
	{[]bool{true, true, true, false, true, false, true}, true},
	{[]bool{true, true, true, false, true, true, false}, true},
	{[]bool{true, true, true, false, true, true, true}, false},
	{[]bool{true, true, true, true, false, false, false}, false},
	{[]bool{true, true, true, true, false, false, true}, true},
	{[]bool{true, true, true, true, false, true, false}, true},
	{[]bool{true, true, true, true, false, true, true}, false}, // 124
	{[]bool{true, true, true, true, true, false, false}, true},
	{[]bool{true, true, true, true, true, false, true}, false},
	{[]bool{true, true, true, true, true, true, false}, false},
	{[]bool{true, true, true, true, true, true, true}, true},
}

func validateParity(g core.Genome[bool]) float64 {
	correct := 0
	for _, n := range parityTests {
		r, err := g.Eval(n.in)
		if err != nil {
			return 0
		}
		if r == n.out {
			correct++
		}
	}
	return 1000.0 * float64(correct) / float64(len(parityTests))
}

func main() {
	funcs := []string{
		"Not",
		"And",
		"Or",
	}
	cat, err := boolNodes.CatalogFromNames(funcs)
	if err != nil {
		log.Fatalf("CatalogFromNames failed: %v", err)
	}

	link, err := boolNodes.LinkFuncFrom("Xor")
	if err != nil {
		log.Fatalf("LinkFuncFrom failed: %v", err)
	}

	numIn := len(parityTests[0].in)
	population, err := evolution.New(cat, 30, 8, 4, numIn, 0, link, validateParity)
	if err != nil {
		log.Fatalf("New failed: %v", err)
	}
	population.MutationConfig = evolutionMutation.Config{
		PointMutationRate: 0.044,
		InversionRate:     0.1,
	}
	solution := population.Evolve(10000)

	// Write out the Go source code for the solution.
	gr, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		log.Printf("unable to load Boolean grammar: %v", err)
	}

	symsPerGene := make([][]string, len(solution.Genome.Genes))
	for i, g := range solution.Genome.Genes {
		symsPerGene[i] = g.SymbolNames()
	}
	prog := codegen.ProgramFromSymbols(symsPerGene, nil, solution.Genome.Link.Symbol())

	fmt.Printf("\n// gepModel is auto-generated Go source code for the\n")
	fmt.Printf("// odd-7-parity solution karva expression:\n// %v\n", solution.Genome.KarvaString())
	if err := codegen.Write(os.Stdout, prog, gr); err != nil {
		log.Fatalf("Write failed: %v", err)
	}
}
