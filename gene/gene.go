// -*- compile-command: "go test"; -*-
// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package gene provides the basis for a single gene in GEP.
package gene

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gmlewis/gep/v2/functions"
	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
)

// FuncWeight contains the symbol name and its weight to be used during
// a run of the GEP algorithm. A symbol with weight 5, for example, will
// be five times more likely to be used than a symbol with weight 1.
type FuncWeight struct {
	Symbol string
	Weight int
}

// Gene contains all the information needed to represent a single gene
// in a GEP expression.
type Gene struct {
	// Symbols is the slice of strings being used in this gene's expression.
	Symbols []string
	// Constants is the slice of floats available for use by this gene.
	Constants []float64

	SymbolMap   map[string]int          // do not use directly.  Use SymbolCount() instead.
	bf          func([]bool) bool       // boolean generated function
	mf          func([]float64) float64 // math generated function
	headSize    int
	choiceSlice []string
	// numTerminals is the number of inputs to the genetic program.
	// It is important to retain this information in order to correctly
	// distinguish between terminals (inputs and constants) and
	// functions in the choiceSlice.  The first numTerminals entries
	// are entirely inputs ("d*") and constants ("c*") whereas all
	// choices following that are strictly function symbols.
	numTerminals int
}

// New creates a new gene based on the Karva string representation.
func New(x string) *Gene {
	parts := strings.Split(x, ".")
	numConstants, numTerminals := 0, 0
	for _, sym := range parts {
		if sym[0:1] == "d" {
			index, err := strconv.Atoi(sym[1:])
			if err != nil {
				log.Fatalf("unable to parse variable index %q: %v", sym, err)
			}
			if index >= numTerminals {
				numTerminals = index + 1
			}
		} else if sym[0:1] == "c" {
			index, err := strconv.Atoi(sym[1:])
			if err != nil {
				log.Fatalf("unable to parse constant index %q: %v", sym, err)
			}
			if index >= numConstants {
				numConstants = index + 1
			}
		}
	}
	return &Gene{
		Symbols:      parts,
		Constants:    make([]float64, numConstants),
		numTerminals: numTerminals + numConstants,
	}
}

// RandomNew generates a new, random gene for further manipulation by the GEP
// algorithm. The headSize, tailSize, numTerminals, and numConstants determine the respective
// properties of the gene, and functions provide the available functions and
// their respective weights to be used in the creation of the gene.
func RandomNew(headSize, tailSize, numTerminals, numConstants int, functions []FuncWeight) *Gene {
	totalWeight := numTerminals + numConstants
	for _, f := range functions {
		totalWeight += f.Weight
	}
	choiceSlice := make([]string, 0, totalWeight)
	for i := 0; i < numTerminals; i++ {
		choiceSlice = append(choiceSlice, fmt.Sprintf("d%v", i))
	}
	var constants []float64
	for i := 0; i < numConstants; i++ {
		choiceSlice = append(choiceSlice, fmt.Sprintf("c%v", i))
		constants = append(constants, rand.Float64())
	}
	for _, f := range functions {
		for i := 0; i < f.Weight; i++ {
			choiceSlice = append(choiceSlice, f.Symbol)
		}
	}
	choices := rand.Perm(totalWeight)
	r := &Gene{
		Symbols:      make([]string, 0, headSize+tailSize),
		Constants:    constants,
		headSize:     headSize,
		choiceSlice:  choiceSlice,
		numTerminals: numTerminals + numConstants,
	}
	for i := 0; i < headSize; i++ { // head is made up of any symbol (function, input, or constant)
		choice := choices[i%len(choices)]
		r.Symbols = append(r.Symbols, choiceSlice[choice])
	}
	for i := 0; i < tailSize; i++ { // tail is strictly made up of terminals (input or constant)
		choice := choices[i%len(choices)]
		r.Symbols = append(r.Symbols, choiceSlice[choice%r.numTerminals])
	}
	return r
}

// String returns the Karva representation of the gene.
func (g Gene) String() string {
	return strings.Join(g.Symbols, ".")
}

// argOrder generates a slice of argument indices (1-based) for every function
// within the list of symbols. It takes into account the arity of each function.
//
// argOrder is used to build up the actual evaluatable expression tree.
//
// For example:
//   'Or.And.Not.Nor' => [[1, 2], [3, 4], [5, 6], [7, 8]]
//   'Or.d0.c0.Nor' => [[1, 2], nil, nil, [3, 4]]
func (g *Gene) getBoolArgOrder(nodes functions.FuncMap) [][]int {
	argOrder := make([][]int, len(g.Symbols))
	argCount := 0
	for i := 0; i < len(g.Symbols); i++ {
		sym := g.Symbols[i]
		s, ok := nodes[sym]
		if !ok {
			continue
		}
		n := s.Terminals()
		if n <= 0 {
			continue
		}
		args := make([]int, n)
		for j := 0; j < n; j++ {
			argCount++
			args[j] = argCount
		}
		argOrder[i] = args
	}
	return argOrder
}

// EvalBool evaluates the gene as a boolean expression and returns the result.
// "in" represents the boolean inputs available to the gene.
// "nodes" is the map of available boolean functions to the gene.
func (g *Gene) EvalBool(in []bool, nodes functions.FuncMap) bool {
	if g.bf == nil {
		argOrder := g.getBoolArgOrder(nodes)
		g.SymbolMap = make(map[string]int)
		g.bf = g.buildBoolTree(0, argOrder, nodes)
	}
	return g.bf(in)
}

func merge(dst *map[string]int, src map[string]int) {
	for k, v := range src {
		(*dst)[k] += v
	}
}

func (g *Gene) buildBoolTree(symbolIndex int, argOrder [][]int, nodes functions.FuncMap) func([]bool) bool {
	// count := make(map[string]int)
	// log.Infof("buildBoolTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex >= len(g.Symbols) {
		log.Printf("bad symbolIndex %v for symbols: %v", symbolIndex, g.Symbols)
		return func(a []bool) bool { return false } // , count
	}
	sym := g.Symbols[symbolIndex]
	g.SymbolMap[sym]++
	if s, ok := nodes[sym]; ok {
		/*
			switch s.Terminals() {
			case 0:
				return func(in []bool) bool {
					return s.BoolFunction(false, false, false, false)
				}, count
			case 1:
				args := argOrder[symbolIndex]
				f, c := g.buildBoolTree(args[0], argOrder, nodes)
				merge(&count, c)
				return func(in []bool) bool {
					return s.BoolFunction(f(in), false, false, false)
				}, count
			case 2:
				args := argOrder[symbolIndex]
				left, lc := g.buildBoolTree(args[0], argOrder, nodes)
				merge(&count, lc)
				right, rc := g.buildBoolTree(args[1], argOrder, nodes)
				merge(&count, rc)
				return func(in []bool) bool {
					return s.BoolFunction(left(in), right(in), false, false)
				}, count
			case 3:
				args := argOrder[symbolIndex]
				f1, c1 := g.buildBoolTree(args[0], argOrder, nodes)
				merge(&count, c1)
				f2, c2 := g.buildBoolTree(args[1], argOrder, nodes)
				merge(&count, c2)
				f3, c3 := g.buildBoolTree(args[2], argOrder, nodes)
				merge(&count, c3)
				return func(in []bool) bool {
					return s.BoolFunction(f1(in), f2(in), f3(in), false)
				}, count
			case 4:
				args := argOrder[symbolIndex]
				f1, c1 := g.buildBoolTree(args[0], argOrder, nodes)
				merge(&count, c1)
				f2, c2 := g.buildBoolTree(args[1], argOrder, nodes)
				merge(&count, c2)
				f3, c3 := g.buildBoolTree(args[2], argOrder, nodes)
				merge(&count, c3)
				f4, c4 := g.buildBoolTree(args[3], argOrder, nodes)
				merge(&count, c4)
				return func(in []bool) bool {
					return s.BoolFunction(f1(in), f2(in), f3(in), f4(in))
				}, count
			}
		*/
		args := argOrder[symbolIndex]
		var funcs []func([]bool) bool
		for _, arg := range args {
			f := g.buildBoolTree(arg, argOrder, nodes)
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
				} // , count
			}
		}
		// Note that constants c0, c1, ... don't make sense for bool expressions
	}
	log.Printf("unable to return function: unknown gene symbol %q", sym)
	return func(in []bool) bool { return false } // , count
}

// argOrder generates a slice of argument indices (1-based) for every function
// within the list of symbols. It takes into account the arity of each function.
//
// argOrder is used to build up the actual evaluatable expression tree.
//
// For example:
//   '+.*.-./' => [[1, 2], [3, 4], [5, 6], [7, 8]]
//   '+.d0.c0./' => [[1, 2], nil, nil, [3, 4]]
func (g *Gene) getMathArgOrder() [][]int {
	argOrder := make([][]int, len(g.Symbols))
	argCount := 0
	for i := 0; i < len(g.Symbols); i++ {
		sym := g.Symbols[i]
		s, ok := mn.Math[sym]
		if !ok {
			continue
		}
		n := s.Terminals()
		if n <= 0 {
			continue
		}
		args := make([]int, n)
		for j := 0; j < n; j++ {
			argCount++
			args[j] = argCount
		}
		argOrder[i] = args
	}
	return argOrder
}

// SymbolCount returns the count of the number of times the symbol
// is actually used in the Gene.
// Note that this count is typically different from the number
// of times the symbol appears in the Karva expression.  This can be
// a handy metric to assist in the fitness evaluation of a Gene.
// Note also that this currently only works for Math expressions.
// Hopefully this restriction will be lifted in the future.
// A workaround for using it with other types is to evaluate the
// Gene, and then g.symbolCount will already be populated.
func (g *Gene) SymbolCount(sym string) int {
	if g.SymbolMap == nil {
		argOrder := g.getMathArgOrder()
		g.SymbolMap = make(map[string]int)
		g.mf = g.buildMathTree(0, argOrder)
	}
	return g.SymbolMap[sym]
}

// EvalMath evaluates the gene as a floating-point expression and returns the result.
// in represents the float64 inputs available to the gene.
func (g *Gene) EvalMath(in []float64) float64 {
	if g.mf == nil {
		argOrder := g.getMathArgOrder()
		g.SymbolMap = make(map[string]int)
		g.mf = g.buildMathTree(0, argOrder)
	}
	return g.mf(in)
}

func (g *Gene) buildMathTree(symbolIndex int, argOrder [][]int) func([]float64) float64 {
	// count := make(map[string]int)
	// log.Infof("buildMathTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex > len(g.Symbols) {
		log.Printf("bad symbolIndex %v for symbols: %v", symbolIndex, g.Symbols)
		return func(a []float64) float64 { return 0.0 } // , count
	}
	sym := g.Symbols[symbolIndex]
	g.SymbolMap[sym]++
	if s, ok := mn.Math[sym]; ok {
		/*
			switch s.Terminals() {
			case 0:
				return func(in []float64) float64 {
					return s.Float64Function(0.0, 0.0, 0.0, 0.0)
				}, count
			case 1:
				args := argOrder[symbolIndex]
				f, c := g.buildMathTree(args[0], argOrder)
				merge(&count, c)
				return func(in []float64) float64 {
					return s.Float64Function(f(in), 0.0, 0.0, 0.0)
				}, count
			case 2:
				args := argOrder[symbolIndex]
				left, lc := g.buildMathTree(args[0], argOrder)
				merge(&count, lc)
				right, rc := g.buildMathTree(args[1], argOrder)
				merge(&count, rc)
				return func(in []float64) float64 {
					return s.Float64Function(left(in), right(in), 0.0, 0.0)
				}, count
			case 3:
				args := argOrder[symbolIndex]
				f1, c1 := g.buildMathTree(args[0], argOrder)
				merge(&count, c1)
				f2, c2 := g.buildMathTree(args[1], argOrder)
				merge(&count, c2)
				f3, c3 := g.buildMathTree(args[2], argOrder)
				merge(&count, c3)
				return func(in []float64) float64 {
					return s.Float64Function(f1(in), f2(in), f3(in), 0.0)
				}, count
			case 4:
				args := argOrder[symbolIndex]
				f1, c1 := g.buildMathTree(args[0], argOrder)
				merge(&count, c1)
				f2, c2 := g.buildMathTree(args[1], argOrder)
				merge(&count, c2)
				f3, c3 := g.buildMathTree(args[2], argOrder)
				merge(&count, c3)
				f4, c4 := g.buildMathTree(args[3], argOrder)
				merge(&count, c4)
				return func(in []float64) float64 {
					return s.Float64Function(f1(in), f2(in), f3(in), f4(in))
				}, count
			}
		*/
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
				log.Printf("unable to parse variable index: sym=%q", sym)
			} else {
				return func(in []float64) float64 {
					if index >= len(in) {
						log.Printf("error evaluating gene %q: index %v >= d length (%v)", sym, index, len(in))
						return 0.0
					}
					return in[index]
				} // , count
			}
		} else if sym[0:1] == "c" {
			if index, err := strconv.Atoi(sym[1:]); err != nil {
				log.Printf("unable to parse constant index: sym=%v", sym)
			} else {
				return func(in []float64) float64 {
					if index >= len(g.Constants) {
						log.Printf("error evaluating gene %q: index %v >= c length (%v)", sym, index, len(g.Constants))
						return 0.0
					}
					return g.Constants[index]
				} // , count
			}
		}
	}
	log.Printf("unable to return function: unknown gene symbol %q", sym)
	return func(in []float64) float64 { return 0.0 } // , count
}

// Mutate mutates a gene by performing a single random symbol exchange within the gene.
func (g *Gene) Mutate() {
	position := rand.Intn(len(g.Symbols))
	if g.numTerminals < 2 {
		position %= g.headSize // Force choice to be within the head
	}
	if position < g.headSize {
		if len(g.choiceSlice) < 2 {
			log.Printf("error: must have choice of more than one function")
			return
		}
		symbol := g.Symbols[position]
		for symbol == g.Symbols[position] { // Force new symbol to be different from old one
			n := rand.Intn(len(g.choiceSlice))
			symbol = g.choiceSlice[n]
		}
		// fmt.Printf("\nChanging symbol #%v from %q to %q\n", position, g.Symbols[position], symbol)
		g.Symbols[position] = symbol
	} else { // Must choose strictly from terminals
		terminal := g.Symbols[position]
		for terminal == g.Symbols[position] { // Force new terminal to be different from old one
			n := rand.Intn(g.numTerminals)
			terminal = g.choiceSlice[n]
		}
		// fmt.Printf("\nChanging terminal #%v from %q to %q\n", position, g.Symbols[position], terminal)
		g.Symbols[position] = terminal
	}
	// Invalidate the cached function
	g.bf = nil
	g.mf = nil
}

// Dup duplicates the gene into the provided destination gene.
func (g *Gene) Dup() *Gene {
	if g == nil {
		log.Printf("gene.Dup error: src and dst must be non-nil")
		return nil
	}
	r := &Gene{
		Symbols:      make([]string, len(g.Symbols)),
		Constants:    make([]float64, len(g.Constants)),
		bf:           g.bf,
		mf:           g.mf,
		headSize:     g.headSize,
		choiceSlice:  make([]string, len(g.choiceSlice)),
		numTerminals: g.numTerminals,
	}
	for i := range g.Symbols {
		r.Symbols[i] = g.Symbols[i]
	}
	for i := range g.Constants {
		r.Constants[i] = g.Constants[i]
	}
	for i := range g.choiceSlice {
		r.choiceSlice[i] = g.choiceSlice[i]
	}
	return r
}

// CheckEqual is used for testing purposes only (exported to use in genome_test.go).
func CheckEqual(g1 *Gene, g2 *Gene) error {
	if g1 == nil || g2 == nil {
		return fmt.Errorf("gene.CheckEqual error: g1 and g2 must be non-nil")
	}
	if len(g1.Symbols) != len(g2.Symbols) {
		return fmt.Errorf("len(g1.Symbols)=%v != len(g2.Symbols)=%v", len(g1.Symbols), len(g2.Symbols))
	}
	for i, v1 := range g1.Symbols {
		if v1 != g2.Symbols[i] {
			return fmt.Errorf("g1.Symbols[%v]=%v != g2.Symbols[%v]=%v", i, v1, i, g2.Symbols[i])
		}
	}
	if len(g1.Constants) != len(g2.Constants) {
		return fmt.Errorf("len(g1.Constants)=%v != len(g2.Constants)=%v", len(g1.Constants), len(g2.Constants))
	}
	for i, v1 := range g1.Constants {
		if v1 != g2.Constants[i] {
			return fmt.Errorf("g1.Constants[%v]=%v != g2.Constants[%v]=%v", i, v1, i, g2.Constants[i])
		}
	}
	if len(g1.choiceSlice) != len(g2.choiceSlice) {
		return fmt.Errorf("len(g1.choiceSlice)=%v != len(g2.choiceSlice)=%v", len(g1.choiceSlice), len(g2.choiceSlice))
	}
	for i, v1 := range g1.choiceSlice {
		if v1 != g2.choiceSlice[i] {
			return fmt.Errorf("g1.choiceSlice[%v]=%v != g2.choiceSlice[%v]=%v", i, v1, i, g2.choiceSlice[i])
		}
	}
	if g1.headSize != g2.headSize {
		return fmt.Errorf("g1.headSize=%v != g2.headSize=%v", g1.headSize, g2.headSize)
	}
	if g1.numTerminals != g2.numTerminals {
		return fmt.Errorf("g1.numTerminals=%v != g2.numTerminals=%v", g1.numTerminals, g2.numTerminals)
	}
	return nil
}
