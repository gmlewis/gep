// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"errors"
	"strconv"

	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
)

func (g *Gene) generateBoolFunc() error {
	argOrder, err := g.getArgOrder()
	if err != nil {
		return err
	}
	g.SymbolMap = make(map[string]int)
	g.bf = g.buildBoolTree(0, argOrder)
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

func (g *Gene) buildBoolTree(symbolIndex int, argOrder [][]int) func([]bool) bool {
	// count := make(map[string]int)
	// log.Infof("buildBoolTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex >= len(g.Symbols) {
		return func(a []bool) bool { return false }
	}
	sym := g.Symbols[symbolIndex]
	g.SymbolMap[sym]++
	if s, ok := bn.BoolAllGates[sym]; ok {
		args := argOrder[symbolIndex]
		var funcs []func([]bool) bool
		for _, arg := range args {
			f := g.buildBoolTree(arg, argOrder)
			funcs = append(funcs, f)
		}
		return func(in []bool) bool {
			var values []bool
			for _, f := range funcs {
				values = append(values, f(in))
			}
			return s.BoolFunction(values)
		}
	} else { // No named symbol found - look for d0, d1, ...
		if sym[0:1] == "d" {
			if index, err := strconv.Atoi(sym[1:]); err != nil {
			} else {
				return func(in []bool) bool {
					if index >= len(in) {
						return false
					}
					return in[index]
				}
			}
		}
		// Note that constants c0, c1, ... don't make sense for bool expressions
	}
	return func(in []bool) bool { return false }
}
