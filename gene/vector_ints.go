// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"strconv"

	"github.com/gmlewis/gep/v2/functions"
	vin "github.com/gmlewis/gep/v2/functions/vector_int_nodes"
)

type VectorInt = functions.VectorInt

func (g *Gene) generateVectorIntFunc() {
	argOrder := g.getArgOrder()
	g.SymbolMap = make(map[string]int)
	g.vif = g.buildVectorIntTree(0, argOrder)
}

// EvalVectorInt evaluates the gene as a floating-point expression and returns the result.
// in represents the int inputs available to the gene.
func (g *Gene) EvalVectorInt(in []VectorInt) VectorInt {
	if g.vif == nil {
		g.generateVectorIntFunc()
	}
	return g.vif(in)
}

func (g *Gene) buildVectorIntTree(symbolIndex int, argOrder [][]int) func([]VectorInt) VectorInt {
	// count := make(map[string]int)
	// log.Infof("buildVectorIntTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex >= len(g.Symbols) {
		return func(a []VectorInt) VectorInt { return VectorInt{} }
	}
	sym := g.Symbols[symbolIndex]
	g.SymbolMap[sym]++
	if s, ok := vin.VectorIntFuncs[sym]; ok {
		args := argOrder[symbolIndex]
		var funcs []func([]VectorInt) VectorInt
		for _, arg := range args {
			f := g.buildVectorIntTree(arg, argOrder)
			funcs = append(funcs, f)
		}
		return func(in []VectorInt) VectorInt {
			var values []VectorInt
			for _, f := range funcs {
				values = append(values, f(in))
			}
			return s.VectorIntFunction(values)
		}
	} else { // No named symbol found - look for d0, d1, ...
		if sym[0:1] == "d" {
			if index, err := strconv.Atoi(sym[1:]); err != nil {
			} else {
				return func(in []VectorInt) VectorInt {
					if index >= len(in) {
						return VectorInt{}
					}
					return in[index]
				}
			}
		} else if sym[0:1] == "c" {
			if index, err := strconv.Atoi(sym[1:]); err != nil {
			} else {
				return func(in []VectorInt) VectorInt {
					if index >= len(g.Constants) {
						return VectorInt{}
					}
					op := func(in []int) int {
						return int(g.Constants[index])
					}
					return vin.ProcessVector(in, op)
				}
			}
		}
	}
	return func(in []VectorInt) VectorInt { return VectorInt{} }
}
