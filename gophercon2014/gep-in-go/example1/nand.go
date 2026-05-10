// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/evolution"
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
	cat, err := boolNodes.CatalogFromNames([]string{"Not", "And", "Or"})
	if err != nil {
		log.Fatalf("CatalogFromNames failed: %v", err)
	}
	link, err := boolNodes.LinkFuncFrom("Or")
	if err != nil {
		log.Fatalf("LinkFuncFrom failed: %v", err)
	}

	population, err := evolution.New(cat, 30, 7, 1, 2, 0, link, validateNand)
	if err != nil {
		log.Fatalf("New failed: %v", err)
	}
	solution := population.Evolve(1000)
	fmt.Printf("nand solution: %v, score=%v\n", solution.Genome.KarvaString(), solution.Score)
}
