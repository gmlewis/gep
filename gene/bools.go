// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"log"
	"strconv"

	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
)

func (g *Gene) generateBoolFunc() {
	argOrder := g.getArgOrder()
	g.SymbolMap = make(map[string]int)
	g.bf = g.buildBoolTree(0, argOrder)
}

// EvalBool evaluates the gene as a boolean expression and returns the result.
// "in" represents the boolean inputs available to the gene.
func (g *Gene) EvalBool(in []bool) bool {
	if g.bf == nil {
		g.generateBoolFunc()
	}
	return g.bf(in)
}

func (g *Gene) buildBoolTree(symbolIndex int, argOrder [][]int) func([]bool) bool {
	// count := make(map[string]int)
	// log.Infof("buildBoolTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex >= len(g.Symbols) {
		log.Printf("bad symbolIndex %v for symbols: %v", symbolIndex, g.Symbols)
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
				log.Printf("unable to parse variable index: sym=%q", sym)
			} else {
				return func(in []bool) bool {
					if index >= len(in) {
						log.Printf("error evaluating gene symbol %q: index %v >= d length (%v)", sym, index, len(in))
						return false
					}
					return in[index]
				}
			}
		}
		// Note that constants c0, c1, ... don't make sense for bool expressions
	}
	log.Printf("unable to return function: unknown gene symbol %q", sym)
	return func(in []bool) bool { return false }
}
