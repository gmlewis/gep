// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"errors"
	"fmt"
	"strconv"

	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
)

func (g *Gene) generateMathFunc() error {
	argOrder, err := g.getArgOrder()
	if err != nil {
		return err
	}
	g.SymbolMap = make(map[string]int)
	g.mf, err = g.buildMathTree(0, argOrder)
	if err != nil {
		return err
	}
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

func (g *Gene) buildMathTree(symbolIndex int, argOrder [][]int) (func([]float64) float64, error) {
	// count := make(map[string]int)
	// log.Infof("buildMathTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex >= len(g.Symbols) {
		return nil, fmt.Errorf("gene.buildMathTree error: symbolIndex %d out of bounds [0,%d)", symbolIndex, len(g.Symbols))
	}
	sym := g.Symbols[symbolIndex]
	g.SymbolMap[sym]++
	if s, ok := mn.Math[sym]; ok {
		args := argOrder[symbolIndex]
		var funcs []func([]float64) float64
		for _, arg := range args {
			f, err := g.buildMathTree(arg, argOrder)
			if err != nil {
				return nil, err
			}
			funcs = append(funcs, f)
		}
		return func(in []float64) float64 {
			var values []float64
			for _, f := range funcs {
				values = append(values, f(in))
			}
			return s.Float64Function(values)
		}, nil
	}
	if sym == "" {
		return nil, errors.New("gene.buildMathTree error: empty symbol")
	}
	if sym[0:1] == "d" { // No named symbol found - look for d0, d1, ...
		index, err := strconv.Atoi(sym[1:])
		if err != nil {
			return nil, fmt.Errorf("gene.buildMathTree error: unable to parse terminal index for symbol %q: %w", sym, err)
		}
		return func(in []float64) float64 {
			if index >= len(in) {
				return 0.0
			}
			return in[index]
		}, nil
	}
	if sym[0:1] == "c" {
		index, err := strconv.Atoi(sym[1:])
		if err != nil {
			return nil, fmt.Errorf("gene.buildMathTree error: unable to parse constant index for symbol %q: %w", sym, err)
		}
		return func(in []float64) float64 {
			if index >= len(g.Constants) {
				return 0.0
			}
			return g.Constants[index]
		}, nil
	}
	return nil, fmt.Errorf("gene.buildMathTree error: unknown symbol %q", sym)
}
