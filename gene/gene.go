// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package gene provides the basis for a single gene in GEP.
package gene

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gmlewis/gep/v2/functions"
	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
	in "github.com/gmlewis/gep/v2/functions/int_nodes"
	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
	vin "github.com/gmlewis/gep/v2/functions/vector_int_nodes"
)

const (
	constRange = 100
)

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
	HeadSize    int
	choiceSlice []string
	// numTerminals is the number of inputs to the genetic program.
	// It is important to retain this information in order to correctly
	// distinguish between terminals (inputs and constants) and
	// functions in the choiceSlice.  The first numTerminals entries
	// are entirely inputs ("d*") and constants ("c*") whereas all
	// choices following that are strictly function symbols.
	numTerminals int
	// rng is an optional seeded random source for reproducible evolution.
	// When nil, the global math/rand source is used.
	rng *rand.Rand
}

// randIntn returns a random int in [0, n) using g.rng when set, else the global source.
func (g *Gene) randIntn(n int) int {
	if g.rng != nil {
		return g.rng.Intn(n)
	}
	return rand.Intn(n)
}

// randFloat64 returns a random float64 in [0.0, 1.0) using g.rng when set, else the global source.
func (g *Gene) randFloat64() float64 {
	if g.rng != nil {
		return g.rng.Float64()
	}
	return rand.Float64()
}

// randPerm returns a random permutation of [0, n) using g.rng when set, else the global source.
func (g *Gene) randPerm(n int) []int {
	if g.rng != nil {
		return g.rng.Perm(n)
	}
	return rand.Perm(n)
}

// New creates a new gene based on the Karva string representation
// and returns an error when the symbol indexes are malformed.
func New(x string, funcType functions.FuncType) (*Gene, error) {
	parts := strings.Split(x, ".")
	numConstants, numTerminals := 0, 0
	var errs []error
	for _, sym := range parts {
		if len(sym) == 0 {
			errs = append(errs, fmt.Errorf("empty symbol"))
			continue
		}
		if sym[0:1] == "d" {
			index, err := strconv.Atoi(sym[1:])
			if err != nil {
				errs = append(errs, fmt.Errorf("unable to parse variable index %q: %w", sym, err))
				continue
			}
			if index >= numTerminals {
				numTerminals = index + 1
			}
		} else if sym[0:1] == "c" {
			index, err := strconv.Atoi(sym[1:])
			if err != nil {
				errs = append(errs, fmt.Errorf("unable to parse constant index %q: %w", sym, err))
				continue
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
	}, errors.Join(errs...)
}

// RandomNew generates a new, random gene for further manipulation by the GEP
// algorithm. The headSize, tailSize, numTerminals, and numConstants determine the respective
// properties of the gene, and functions provide the available functions and
// their respective weights to be used in the creation of the gene.
// An optional *rand.Rand may be passed as the last argument to enable deterministic
// gene generation; if omitted or nil, the global math/rand source is used.
func RandomNew(headSize, tailSize, numTerminals, numConstants int, functions []FuncWeight, funcType functions.FuncType, rngs ...*rand.Rand) *Gene {
	var rng *rand.Rand
	if len(rngs) > 0 {
		rng = rngs[0]
	}
	totalWeight := numTerminals + numConstants
	for _, f := range functions {
		totalWeight += f.Weight
	}
	choiceSlice := make([]string, 0, totalWeight)
	for i := 0; i < numTerminals; i++ {
		choiceSlice = append(choiceSlice, fmt.Sprintf("d%v", i))
	}
	constants := make([]float64, 0, numConstants)
	r := &Gene{
		Symbols:      make([]string, 0, headSize+tailSize),
		Constants:    constants,
		funcType:     funcType,
		HeadSize:     headSize,
		numTerminals: numTerminals + numConstants,
		rng:          rng,
	}
	for i := 0; i < numConstants; i++ {
		choiceSlice = append(choiceSlice, fmt.Sprintf("c%v", i))
		r.Constants = append(r.Constants, math.Round(constRange*r.randFloat64()))
	}
	for _, f := range functions {
		for i := 0; i < f.Weight; i++ {
			choiceSlice = append(choiceSlice, f.Symbol)
		}
	}
	r.choiceSlice = choiceSlice
	choices := r.randPerm(totalWeight)
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

// String returns the Karva representation of the gene and any
// symbol-format errors encountered while rendering constants.
func (g Gene) String() (string, error) {
	var syms []string
	var errs []error
	for _, s := range g.Symbols {
		if strings.HasPrefix(s, "c") {
			i, err := strconv.Atoi(s[1:])
			if err != nil || i < 0 || i >= len(g.Constants) {
				if err != nil {
					errs = append(errs, fmt.Errorf("bad constant name %q: %w", s, err))
				} else {
					errs = append(errs, fmt.Errorf("constant index out of range: symbol=%q index=%v len(constants)=%v", s, i, len(g.Constants)))
				}
				syms = append(syms, s)
				continue
			}
			syms = append(syms, fmt.Sprintf("%v(%v)", s, g.Constants[i]))
		} else {
			syms = append(syms, s)
		}
	}
	return strings.Join(syms, "."), errors.Join(errs...)
}

// DotGraph returns a graphviz "dot" language representation of the gene.
func (g Gene) DotGraph() string {
	var lines []string
	return strings.Join(lines, "\n")
}

// SymbolCount returns the count of the number of times the symbol
// is actually used in the Gene.
// Note that this count is typically different from the number
// of times the symbol appears in the Karva expression.  This can be
// a handy metric to assist in the fitness evaluation of a Gene.
func (g *Gene) SymbolCount(sym string) (int, error) {
	if g.SymbolMap == nil {
		switch g.funcType {
		case functions.Bool:
			if err := g.generateBoolFunc(); err != nil {
				return 0, err
			}
		case functions.Int:
			if err := g.generateIntFunc(); err != nil {
				return 0, err
			}
		case functions.Float64:
			if err := g.generateMathFunc(); err != nil {
				return 0, err
			}
		case functions.VectorInts:
			if err := g.generateVectorIntFunc(); err != nil {
				return 0, err
			}
		default:
			return 0, fmt.Errorf("unknown funcType: %v", g.funcType)
		}
	}
	return g.SymbolMap[sym], nil
}

// Mutate mutates a gene by performing a single random symbol exchange
// within the gene and surfaces invalid mutation preconditions.
func (g *Gene) Mutate() error {
	position := g.randIntn(len(g.Symbols))
	if g.numTerminals < 2 {
		position %= g.HeadSize // Force choice to be within the head
	}
	if position < g.HeadSize {
		if len(g.choiceSlice) < 2 {
			return fmt.Errorf("must have choice of more than one function")
		}
		symbol := g.Symbols[position]
		for symbol == g.Symbols[position] { // Force new symbol to be different from old one
			n := g.randIntn(len(g.choiceSlice))
			symbol = g.choiceSlice[n]
		}
		// fmt.Printf("\nChanging symbol #%v from %q to %q\n", position, g.Symbols[position], symbol)
		g.Symbols[position] = symbol
	} else { // Must choose strictly from terminals
		terminal := g.Symbols[position]
		for terminal == g.Symbols[position] { // Force new terminal to be different from old one
			n := g.randIntn(g.numTerminals)
			terminal = g.choiceSlice[n]
		}
		// fmt.Printf("\nChanging terminal #%v from %q to %q\n", position, g.Symbols[position], terminal)
		g.Symbols[position] = terminal
	}
	g.InvalidateCache()
	return nil
}

// Dup duplicates the gene into the provided destination gene.
func (g *Gene) Dup() (*Gene, error) {
	if g == nil {
		return nil, fmt.Errorf("gene.Dup error: src and dst must be non-nil")
	}
	r := &Gene{
		Symbols:      make([]string, len(g.Symbols)),
		Constants:    make([]float64, len(g.Constants)),
		funcType:     g.funcType,
		HeadSize:     g.HeadSize,
		choiceSlice:  make([]string, len(g.choiceSlice)),
		numTerminals: g.numTerminals,
		rng:          g.rng,
	}
	copy(r.Symbols, g.Symbols)
	copy(r.Constants, g.Constants)
	copy(r.choiceSlice, g.choiceSlice)
	return r, nil
}

// InvalidateCache clears all cached generated functions and symbol counts.
func (g *Gene) InvalidateCache() {
	g.SymbolMap = nil
	g.bf = nil
	g.intF = nil
	g.mf = nil
	g.vif = nil
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
	if g1.HeadSize != g2.HeadSize {
		return fmt.Errorf("g1.HeadSize=%v != g2.HeadSize=%v", g1.HeadSize, g2.HeadSize)
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
func (g *Gene) getArgOrder() ([][]int, error) {
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
		return nil, fmt.Errorf("unknown funcType: %v", g.funcType)
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
	return argOrder, nil
}
