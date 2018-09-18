// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Symbolic_regression is a simple experiment to run the GEP algorithm using the floating-point math package.
// Given a set of input functions (+, -, *, and /), this solves the equation "a^4 + a^3 + a^2 + a"
// from those basic building blocks. This experiment usually converges to a solution within
// the first 10000 generations of evolution, but not always.
package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/genome"
	"github.com/gmlewis/gep/v2/grammars"
	"github.com/gmlewis/gep/v2/model"
)

// srTests is a random sample of inputs and outputs for the function "a^4 + a^3 + a^2 + a"
var srTests = []struct {
	in  []float64
	out float64
}{
	// {[]float64{0}, 0},
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
	{[]float64{100}, 101010100},
	{[]float64{-100}, 99009900},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func validateFunc(g *genome.Genome) float64 {
	result := 0.0
	for _, n := range srTests {
		r := g.EvalMath(n.in)
		// fmt.Printf("r=%v, n.in=%v, n.out=%v, g=%v\n", r, n.in, n.out, g)
		if math.IsInf(r, 0) {
			return 0.0
		}
		fitness := math.Abs(r - n.out)
		fitness = 1000.0 / (1.0 + fitness) // fitness is normalized and max value is 1000
		// fmt.Printf("r=%v, n.in=%v, n.out=%v, fitness=%v, g=%v\n", r, n.in, n.out, fitness, g)
		result += fitness
	}
	return result / float64(len(srTests))
}

func main() {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
	}
	numIn := len(srTests[0].in)
	e := model.New(funcs, functions.Float64, 30, 8, 4, numIn, 0, "+", validateFunc)
	s := e.Evolve(10000)

	// Write out the Go source code for the solution.
	gr, err := grammars.LoadGoMathGrammar()
	if err != nil {
		log.Printf("unable to load grammar: %v", err)
	}
	fmt.Printf("\n// gepModel is auto-generated Go source code for the\n")
	fmt.Printf("// (a^4 + a^3 + a^2 + a) solution karva expression:\n// %q, score=%v\n", s, validateFunc(s))
	s.Write(os.Stdout, gr)
}
