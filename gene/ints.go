// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"errors"
	"strconv"

	in "github.com/gmlewis/gep/v2/functions/int_nodes"
)

func (g *Gene) generateIntFunc() error {
	argOrder, err := g.getArgOrder()
	if err != nil {
		return err
	}
	g.SymbolMap = make(map[string]int)
	g.intF = g.buildIntTree(0, argOrder)
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

func (g *Gene) buildIntTree(symbolIndex int, argOrder [][]int) func([]int) int {
	// count := make(map[string]int)
	// log.Infof("buildIntTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex >= len(g.Symbols) {
		return func(a []int) int { return 0 }
	}
	sym := g.Symbols[symbolIndex]
	g.SymbolMap[sym]++
	if s, ok := in.Int[sym]; ok {
		args := argOrder[symbolIndex]
		var funcs []func([]int) int
		for _, arg := range args {
			f := g.buildIntTree(arg, argOrder)
			funcs = append(funcs, f)
		}
		return func(in []int) int {
			var values []int
			for _, f := range funcs {
				values = append(values, f(in))
			}
			return s.IntFunction(values)
		}
	} else { // No named symbol found - look for d0, d1, ...
		if sym[0:1] == "d" {
			if index, err := strconv.Atoi(sym[1:]); err != nil {
			} else {
				return func(in []int) int {
					if index >= len(in) {
						return 0
					}
					return in[index]
				}
			}
		} else if sym[0:1] == "c" {
			if index, err := strconv.Atoi(sym[1:]); err != nil {
			} else {
				return func(in []int) int {
					if index >= len(g.Constants) {
						return 0
					}
					return int(g.Constants[index])
				}
			}
		}
	}
	return func(in []int) int { return 0 }
}
