// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	mn "github.com/gmlewis/gep/functions/math_nodes"
	"github.com/gmlewis/gep/gene"
	"github.com/gmlewis/gep/genome"
	"github.com/gmlewis/gep/model"
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

func validateFunc(g *genome.Genome) float64 {
	result := 0.0
	for _, n := range srTests {
		r := g.EvalMath(n.in)
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
	funcs := []gene.FuncWeight{
		{"+", 1},
		{"-", 1},
		{"*", 1},
		{"/", 1},
	}
	e := model.New(funcs, mn.Math, 30, 6, 1, 1, 0, "+", validateFunc)
	s := e.Evolve(10000)
	fmt.Printf("(a^4 + a^3 + a^2 + a) solution: %v, score=%v\n", s, validateFunc(s))
}
