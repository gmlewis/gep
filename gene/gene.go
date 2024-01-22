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
	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
	in "github.com/gmlewis/gep/v2/functions/int_nodes"
	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
	vin "github.com/gmlewis/gep/v2/functions/vector_int_nodes"
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

	// funcType keep track of the underlying function types (no generics).
	funcType functions.FuncType
	// Instead of generics, we list all the possibilities:
	bf   func([]bool) bool                               // boolean generated function
	intF func([]int) int                                 // integer generated function
	mf   func([]float64) float64                         // math generated function
	vif  func([]functions.VectorInt) functions.VectorInt // vector of integers generated function

	SymbolMap   map[string]int // do not use directly.  Use SymbolCount() instead.
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
func New(x string, funcType functions.FuncType) *Gene {
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
		funcType:     funcType,
		numTerminals: numTerminals + numConstants,
	}
}

// RandomNew generates a new, random gene for further manipulation by the GEP
// algorithm. The headSize, tailSize, numTerminals, and numConstants determine the respective
// properties of the gene, and functions provide the available functions and
// their respective weights to be used in the creation of the gene.
func RandomNew(headSize, tailSize, numTerminals, numConstants int, functions []FuncWeight, funcType functions.FuncType) *Gene {
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
		funcType:     funcType,
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
	var syms []string
	for _, s := range g.Symbols {
		if strings.HasPrefix(s, "c") {
			i, err := strconv.Atoi(s[1:])
			if err != nil {
				log.Fatalf("bad constant name: %v", s)
			}
			syms = append(syms, fmt.Sprintf("%v(%.2f)", s, g.Constants[i]))
		} else {
			syms = append(syms, s)
		}
	}
	return strings.Join(syms, ".")
}

// SymbolCount returns the count of the number of times the symbol
// is actually used in the Gene.
// Note that this count is typically different from the number
// of times the symbol appears in the Karva expression.  This can be
// a handy metric to assist in the fitness evaluation of a Gene.
func (g *Gene) SymbolCount(sym string) int {
	if g.SymbolMap == nil {
		switch g.funcType {
		case functions.Bool:
			g.generateBoolFunc()
		case functions.Int:
			g.generateIntFunc()
		case functions.Float64:
			g.generateMathFunc()
		case functions.VectorInts:
			g.generateVectorIntFunc()
		default:
			log.Fatalf("unknown funcType: %v", g.funcType)
		}
	}
	return g.SymbolMap[sym]
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
	g.intF = nil
	g.mf = nil
	g.vif = nil
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
		funcType:     g.funcType,
		bf:           g.bf,
		intF:         g.intF,
		mf:           g.mf,
		headSize:     g.headSize,
		choiceSlice:  make([]string, len(g.choiceSlice)),
		numTerminals: g.numTerminals,
	}
	copy(r.Symbols, g.Symbols)
	copy(r.Constants, g.Constants)
	copy(r.choiceSlice, g.choiceSlice)
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

// getArgOrder generates a slice of argument indices (1-based) for every function
// within the list of symbols. It takes into account the arity of each function.
//
// argOrder is used to build up the actual evaluatable expression tree.
//
// For example:
//
//	'+.*.-./' => [[1, 2], [3, 4], [5, 6], [7, 8]]
//	'+.d0.c0./' => [[1, 2], nil, nil, [3, 4]]
func (g *Gene) getArgOrder() [][]int {
	var lookup functions.FuncMap
	switch g.funcType {
	case functions.Bool:
		lookup = bn.BoolAllGates
	case functions.Int:
		lookup = in.Int
	case functions.Float64:
		lookup = mn.Math
	case functions.VectorInts:
		lookup = vin.VectorIntFuncs
	default:
		log.Fatalf("unknown funcType: %v", g.funcType)
	}

	argOrder := make([][]int, len(g.Symbols))
	argCount := 0
	for i := 0; i < len(g.Symbols); i++ {
		sym := g.Symbols[i]
		s, ok := lookup[sym]
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
