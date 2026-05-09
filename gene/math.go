// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"errors"
	"strconv"

	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
)

func (g *Gene) generateMathFunc() error {
	argOrder, err := g.getArgOrder()
	if err != nil {
		return err
	}
	g.SymbolMap = make(map[string]int)
	g.mf = g.buildMathTree(0, argOrder)
	return nil
}

// EvalMath evaluates the gene as a floating-point expression and returns the result.
// in represents the float64 inputs available to the gene.
func (g *Gene) EvalMath(in []float64) (float64, error) {
	if err := g.validateMathSymbols(in); err != nil {
		return 0.0, err
	}
	if g.mf == nil {
		if err := g.generateMathFunc(); err != nil {
			return 0.0, err
		}
	}
	if g.mf == nil {
		return 0.0, errors.New("unable to generate math evaluator")
	}
	return g.mf(in), nil
}

func (g *Gene) buildMathTree(symbolIndex int, argOrder [][]int) func([]float64) float64 {
	// count := make(map[string]int)
	// log.Infof("buildMathTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex >= len(g.Symbols) {
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
			} else {
				return func(in []float64) float64 {
					if index >= len(in) {
						return 0.0
					}
					return in[index]
				}
			}
		} else if sym[0:1] == "c" {
			if index, err := strconv.Atoi(sym[1:]); err != nil {
			} else {
				return func(in []float64) float64 {
					if index >= len(g.Constants) {
						return 0.0
					}
					return g.Constants[index]
				}
			}
		}
	}
	return func(in []float64) float64 { return 0.0 }
}
