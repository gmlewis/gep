// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gmlewis/gep/v2/functions"
	vin "github.com/gmlewis/gep/v2/functions/vector_int_nodes"
)

type VectorInt = functions.VectorInt

func (g *Gene) generateVectorIntFunc() error {
	argOrder, err := g.getArgOrder()
	if err != nil {
		return err
	}
	g.SymbolMap = make(map[string]int)
	g.vif, err = g.buildVectorIntTree(0, argOrder)
	if err != nil {
		return err
	}
	return nil
}

// EvalVectorInt evaluates the gene as a floating-point expression and returns the result.
// in represents the int inputs available to the gene.
func (g *Gene) EvalVectorInt(in []VectorInt) (VectorInt, error) {
	if err := g.validateVectorIntSymbols(in); err != nil {
		return VectorInt{}, err
	}
	if g.vif == nil {
		if err := g.generateVectorIntFunc(); err != nil {
			return VectorInt{}, err
		}
	}
	if g.vif == nil {
		return VectorInt{}, errors.New("unable to generate vector-int evaluator")
	}
	return g.vif(in), nil
}

func (g *Gene) buildVectorIntTree(symbolIndex int, argOrder [][]int) (func([]VectorInt) VectorInt, error) {
	// count := make(map[string]int)
	// log.Infof("buildVectorIntTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex >= len(g.Symbols) {
		return nil, fmt.Errorf("gene.buildVectorIntTree error: symbolIndex %d out of bounds [0,%d)", symbolIndex, len(g.Symbols))
	}
	sym := g.Symbols[symbolIndex]
	g.SymbolMap[sym]++
	if s, ok := vin.VectorIntFuncs[sym]; ok {
		args := argOrder[symbolIndex]
		var funcs []func([]VectorInt) VectorInt
		for _, arg := range args {
			f, err := g.buildVectorIntTree(arg, argOrder)
			if err != nil {
				return nil, err
			}
			funcs = append(funcs, f)
		}
		return func(in []VectorInt) VectorInt {
			var values []VectorInt
			for _, f := range funcs {
				values = append(values, f(in))
			}
			return s.VectorIntFunction(values)
		}, nil
	}
	if sym == "" {
		return nil, errors.New("gene.buildVectorIntTree error: empty symbol")
	}
	if sym[0:1] == "d" { // No named symbol found - look for d0, d1, ...
		index, err := strconv.Atoi(sym[1:])
		if err != nil {
			return nil, fmt.Errorf("gene.buildVectorIntTree error: unable to parse terminal index for symbol %q: %w", sym, err)
		}
		return func(in []VectorInt) VectorInt {
			if index >= len(in) {
				return VectorInt{}
			}
			return in[index]
		}, nil
	}
	if sym[0:1] == "c" {
		index, err := strconv.Atoi(sym[1:])
		if err != nil {
			return nil, fmt.Errorf("gene.buildVectorIntTree error: unable to parse constant index for symbol %q: %w", sym, err)
		}
		return func(in []VectorInt) VectorInt {
			if index >= len(g.Constants) {
				return VectorInt{}
			}
			op := func(in []int) int {
				return int(g.Constants[index])
			}
			return vin.ProcessVector(in, op)
		}, nil
	}
	return nil, fmt.Errorf("gene.buildVectorIntTree error: unknown symbol %q", sym)
}
