// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"math"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/evolution"
	"github.com/gmlewis/gep/v2/functions"
	mathNodes "github.com/gmlewis/gep/v2/functions/math_nodes"
)

var srTests = []struct {
	in  []float64
	out float64
}{
	{[]float64{2.81}, 95.2425},
	{[]float64{6}, 1554},
	{[]float64{7.043}, 2866.55},
	{[]float64{8}, 4680},
	{[]float64{10}, 11110},
	{[]float64{11.38}, 18386},
	{[]float64{12}, 22620},
	{[]float64{14}, 41370},
	{[]float64{15}, 54240},
	{[]float64{20}, 168420},
}

func validateFunc(g core.Genome[float64]) float64 {
	result := 0.0
	for _, n := range srTests {
		r, err := g.Eval(n.in)
		if err != nil {
			return 0
		}
		if math.IsInf(r, 0) {
			return 0.0
		}
		fitness := math.Abs(r - n.out)
		fitness = 1000.0 / (1.0 + fitness) // fitness is normalized and max value is 1000
		result += fitness
	}
	return result
}

func main() {
	funcs := []string{"+", "-", "*", "/"}
	fm := make(functions.FuncMap, len(funcs))
	for _, sym := range funcs {
		fn, ok := mathNodes.Math[sym]
		if !ok {
			log.Fatalf("unsupported math function %q", sym)
		}
		fm[sym] = fn
	}
	cat, err := mathNodes.CatalogFrom(fm)
	if err != nil {
		log.Fatalf("CatalogFrom failed: %v", err)
	}

	linkNode, ok := mathNodes.Math["+"]
	if !ok {
		log.Fatal(`link function "+" not found`)
	}
	link, err := core.NewLinkFunc[float64]("+", linkNode.Float64Function)
	if err != nil {
		log.Fatalf("NewLinkFunc failed: %v", err)
	}

	population, err := evolution.New(cat, 30, 6, 1, 1, 0, link, validateFunc)
	if err != nil {
		log.Fatalf("New failed: %v", err)
	}
	solution := population.Evolve(10000)
	fmt.Printf("(a^4 + a^3 + a^2 + a) solution: %v, score=%v\n", solution.Genome.KarvaString(), solution.Score)
}
