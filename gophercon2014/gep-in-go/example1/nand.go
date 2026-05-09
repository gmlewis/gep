// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/evolution"
	"github.com/gmlewis/gep/v2/functions"
	boolNodes "github.com/gmlewis/gep/v2/functions/bool_nodes"
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
	funcs := []string{"Not", "And", "Or"}
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
		log.Fatal(`link function "Or" not found`)
	}
	link, err := core.NewLinkFunc[bool]("Or", linkNode.BoolFunction)
	if err != nil {
		log.Fatalf("NewLinkFunc failed: %v", err)
	}

	population, err := evolution.New(cat, 30, 7, 1, 2, 0, link, validateNand)
	if err != nil {
		log.Fatalf("New failed: %v", err)
	}
	solution := population.Evolve(1000)
	fmt.Printf("nand solution: %v, score=%v\n", solution.Genome.KarvaString(), solution.Score)
}
