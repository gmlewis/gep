// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"log"
	"strconv"

	in "github.com/gmlewis/gep/v2/functions/int_nodes"
)

func (g *Gene) generateIntFunc() {
	argOrder := g.getArgOrder()
	g.SymbolMap = make(map[string]int)
	g.intF = g.buildIntTree(0, argOrder)
}

// EvalInt evaluates the gene as a floating-point expression and returns the result.
// in represents the int inputs available to the gene.
func (g *Gene) EvalInt(in []int) int {
	if g.intF == nil {
		g.generateIntFunc()
	}
	return g.intF(in)
}

func (g *Gene) buildIntTree(symbolIndex int, argOrder [][]int) func([]int) int {
	// count := make(map[string]int)
	// log.Infof("buildIntTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex > len(g.Symbols) {
		log.Printf("bad symbolIndex %v for symbols: %v", symbolIndex, g.Symbols)
		return func(a []int) int { return 0.0 }
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
				log.Printf("unable to parse variable index: sym=%q", sym)
			} else {
				return func(in []int) int {
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
				return func(in []int) int {
					if index >= len(g.Constants) {
						log.Printf("error evaluating gene %q: index %v >= c length (%v)", sym, index, len(g.Constants))
						return 0.0
					}
					return int(g.Constants[index])
				}
			}
		}
	}
	log.Printf("unable to return function: unknown gene symbol %q", sym)
	return func(in []int) int { return 0.0 }
}
