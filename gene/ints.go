// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"errors"
	"fmt"
	"strconv"

	in "github.com/gmlewis/gep/v2/functions/int_nodes"
)

func (g *Gene) generateIntFunc() error {
	argOrder, err := g.getArgOrder()
	if err != nil {
		return err
	}
	g.SymbolMap = make(map[string]int)
	g.intF, err = g.buildIntTree(0, argOrder)
	if err != nil {
		return err
	}
	return nil
}

// EvalInt evaluates the gene as a floating-point expression and returns the result.
// in represents the int inputs available to the gene.
func (g *Gene) EvalInt(in []int) (int, error) {
	if err := g.validateIntSymbols(in); err != nil {
		return 0, err
	}
	if g.intF == nil {
		if err := g.generateIntFunc(); err != nil {
			return 0, err
		}
	}
	if g.intF == nil {
		return 0, errors.New("unable to generate int evaluator")
	}
	return g.intF(in), nil
}

func (g *Gene) buildIntTree(symbolIndex int, argOrder [][]int) (func([]int) int, error) {
	// count := make(map[string]int)
	// log.Infof("buildIntTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex >= len(g.Symbols) {
		return nil, fmt.Errorf("gene.buildIntTree error: symbolIndex %d out of bounds [0,%d)", symbolIndex, len(g.Symbols))
	}
	sym := g.Symbols[symbolIndex]
	g.SymbolMap[sym]++
	if s, ok := in.Int[sym]; ok {
		args := argOrder[symbolIndex]
		var funcs []func([]int) int
		for _, arg := range args {
			f, err := g.buildIntTree(arg, argOrder)
			if err != nil {
				return nil, err
			}
			funcs = append(funcs, f)
		}
		return func(in []int) int {
			var values []int
			for _, f := range funcs {
				values = append(values, f(in))
			}
			return s.IntFunction(values)
		}, nil
	}
	if sym == "" {
		return nil, errors.New("gene.buildIntTree error: empty symbol")
	}
	if sym[0:1] == "d" { // No named symbol found - look for d0, d1, ...
		index, err := strconv.Atoi(sym[1:])
		if err != nil {
			return nil, fmt.Errorf("gene.buildIntTree error: unable to parse terminal index for symbol %q: %w", sym, err)
		}
		return func(in []int) int {
			if index >= len(in) {
				return 0
			}
			return in[index]
		}, nil
	}
	if sym[0:1] == "c" {
		index, err := strconv.Atoi(sym[1:])
		if err != nil {
			return nil, fmt.Errorf("gene.buildIntTree error: unable to parse constant index for symbol %q: %w", sym, err)
		}
		return func(in []int) int {
			if index >= len(g.Constants) {
				return 0
			}
			return int(g.Constants[index])
		}, nil
	}
	return nil, fmt.Errorf("gene.buildIntTree error: unknown symbol %q", sym)
}
