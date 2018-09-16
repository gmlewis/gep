// -*- compile-command: "go test"; -*-
// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"log"
	"strconv"

	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
)

// argOrder generates a slice of argument indices (1-based) for every function
// within the list of symbols. It takes into account the arity of each function.
//
// argOrder is used to build up the actual evaluatable expression tree.
//
// For example:
//   '+.*.-./' => [[1, 2], [3, 4], [5, 6], [7, 8]]
//   '+.d0.c0./' => [[1, 2], nil, nil, [3, 4]]
func (g *Gene) getMathArgOrder() [][]int {
	argOrder := make([][]int, len(g.Symbols))
	argCount := 0
	for i := 0; i < len(g.Symbols); i++ {
		sym := g.Symbols[i]
		s, ok := mn.Math[sym]
		if !ok {
			continue
		}
		n := s.Terminals()
		if n <= 0 {
			continue
		}
		args := make([]int, n)
		for j := 0; j < n; j++ {
			argCount++
			args[j] = argCount
		}
		argOrder[i] = args
	}
	return argOrder
}

// SymbolCount returns the count of the number of times the symbol
// is actually used in the Gene.
// Note that this count is typically different from the number
// of times the symbol appears in the Karva expression.  This can be
// a handy metric to assist in the fitness evaluation of a Gene.
// Note also that this currently only works for Math expressions.
// Hopefully this restriction will be lifted in the future.
// A workaround for using it with other types is to evaluate the
// Gene, and then g.symbolCount will already be populated.
func (g *Gene) SymbolCount(sym string) int {
	if g.SymbolMap == nil {
		argOrder := g.getMathArgOrder()
		g.SymbolMap = make(map[string]int)
		g.mf = g.buildMathTree(0, argOrder)
	}
	return g.SymbolMap[sym]
}

// EvalMath evaluates the gene as a floating-point expression and returns the result.
// in represents the float64 inputs available to the gene.
func (g *Gene) EvalMath(in []float64) float64 {
	if g.mf == nil {
		argOrder := g.getMathArgOrder()
		g.SymbolMap = make(map[string]int)
		g.mf = g.buildMathTree(0, argOrder)
	}
	return g.mf(in)
}

func (g *Gene) buildMathTree(symbolIndex int, argOrder [][]int) func([]float64) float64 {
	// count := make(map[string]int)
	// log.Infof("buildMathTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex > len(g.Symbols) {
		log.Printf("bad symbolIndex %v for symbols: %v", symbolIndex, g.Symbols)
		return func(a []float64) float64 { return 0.0 }
	}
	sym := g.Symbols[symbolIndex]
	g.SymbolMap[sym]++
	if s, ok := mn.Math[sym]; ok {
		args := argOrder[symbolIndex]
		var funcs []func([]float64) float64
		for _, arg := range args {
			f := g.buildMathTree(arg, argOrder)
			funcs = append(funcs, f)
		}
		return func(in []float64) float64 {
			var values []float64
			for _, f := range funcs {
				values = append(values, f(in))
			}
			return s.Float64Function(values)
		}
	} else { // No named symbol found - look for d0, d1, ...
		if sym[0:1] == "d" {
			if index, err := strconv.Atoi(sym[1:]); err != nil {
				log.Printf("unable to parse variable index: sym=%q", sym)
			} else {
				return func(in []float64) float64 {
					if index >= len(in) {
						log.Printf("error evaluating gene %q: index %v >= d length (%v)", sym, index, len(in))
						return 0.0
					}
					return in[index]
				}
			}
		} else if sym[0:1] == "c" {
			if index, err := strconv.Atoi(sym[1:]); err != nil {
				log.Printf("unable to parse constant index: sym=%v", sym)
			} else {
				return func(in []float64) float64 {
					if index >= len(g.Constants) {
						log.Printf("error evaluating gene %q: index %v >= c length (%v)", sym, index, len(g.Constants))
						return 0.0
					}
					return g.Constants[index]
				}
			}
		}
	}
	log.Printf("unable to return function: unknown gene symbol %q", sym)
	return func(in []float64) float64 { return 0.0 }
}
