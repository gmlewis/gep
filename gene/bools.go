// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"errors"
	"fmt"
	"strconv"

	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
)

func (g *Gene) generateBoolFunc() error {
	argOrder, err := g.getArgOrder()
	if err != nil {
		return err
	}
	g.SymbolMap = make(map[string]int)
	g.bf, err = g.buildBoolTree(0, argOrder)
	if err != nil {
		return err
	}
	return nil
}

// EvalBool evaluates the gene as a boolean expression and returns the result.
// "in" represents the boolean inputs available to the gene.
func (g *Gene) EvalBool(in []bool) (bool, error) {
	if err := g.validateBoolSymbols(in); err != nil {
		return false, err
	}
	if g.bf == nil {
		if err := g.generateBoolFunc(); err != nil {
			return false, err
		}
	}
	if g.bf == nil {
		return false, errors.New("unable to generate bool evaluator")
	}
	return g.bf(in), nil
}

func (g *Gene) buildBoolTree(symbolIndex int, argOrder [][]int) (func([]bool) bool, error) {
	// count := make(map[string]int)
	// log.Infof("buildBoolTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex >= len(g.Symbols) {
		return nil, fmt.Errorf("gene.buildBoolTree error: symbolIndex %d out of bounds [0,%d)", symbolIndex, len(g.Symbols))
	}
	sym := g.Symbols[symbolIndex]
	g.SymbolMap[sym]++
	if s, ok := bn.BoolAllGates[sym]; ok {
		args := argOrder[symbolIndex]
		var funcs []func([]bool) bool
		for _, arg := range args {
			f, err := g.buildBoolTree(arg, argOrder)
			if err != nil {
				return nil, err
			}
			funcs = append(funcs, f)
		}
		return func(in []bool) bool {
			var values []bool
			for _, f := range funcs {
				values = append(values, f(in))
			}
			return s.BoolFunction(values)
		}, nil
	}
	if sym == "" {
		return nil, errors.New("gene.buildBoolTree error: empty symbol")
	}
	if sym[0:1] == "d" { // No named symbol found - look for d0, d1, ...
		index, err := strconv.Atoi(sym[1:])
		if err != nil {
			return nil, fmt.Errorf("gene.buildBoolTree error: unable to parse terminal index for symbol %q: %w", sym, err)
		}
		return func(in []bool) bool {
			if index >= len(in) {
				return false
			}
			return in[index]
		}, nil
	}
	// Note that constants c0, c1, ... don't make sense for bool expressions.
	return nil, fmt.Errorf("gene.buildBoolTree error: unknown symbol %q", sym)
}
