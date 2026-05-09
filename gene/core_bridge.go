// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/functions"
	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
	in "github.com/gmlewis/gep/v2/functions/int_nodes"
	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
	vin "github.com/gmlewis/gep/v2/functions/vector_int_nodes"
)

func hasLegacyConstantSymbol(symbols []string) bool {
	for _, sym := range symbols {
		if strings.HasPrefix(sym, "c") {
			return true
		}
	}
	return false
}

func hasCoreConstantSymbol[T any](symbols []core.Symbol[T]) bool {
	for _, sym := range symbols {
		if sym.Kind == core.SymbolKindConstant {
			return true
		}
	}
	return false
}

// CoreBool converts a legacy boolean gene into the typed core representation.
func (g *Gene) CoreBool() (core.Gene[bool], error) {
	if g == nil {
		return core.Gene[bool]{}, errors.New("gene.CoreBool: gene cannot be nil")
	}
	if len(g.Constants) > 0 || hasLegacyConstantSymbol(g.Symbols) {
		return core.Gene[bool]{}, errors.New("gene.CoreBool: bool genes do not support constants")
	}
	cat, err := bn.CatalogFrom(bn.BoolAllGates)
	if err != nil {
		return core.Gene[bool]{}, fmt.Errorf("gene.CoreBool: %w", err)
	}
	syms, err := core.ParseSymbols(g.Symbols, cat)
	if err != nil {
		return core.Gene[bool]{}, fmt.Errorf("gene.CoreBool: %w", err)
	}
	return core.Gene[bool]{Symbols: syms}, nil
}

// CoreInt converts a legacy integer gene into the typed core representation.
func (g *Gene) CoreInt() (core.Gene[int], error) {
	if g == nil {
		return core.Gene[int]{}, errors.New("gene.CoreInt: gene cannot be nil")
	}
	cat, err := in.CatalogFrom(in.Int)
	if err != nil {
		return core.Gene[int]{}, fmt.Errorf("gene.CoreInt: %w", err)
	}
	syms, err := core.ParseSymbols(g.Symbols, cat)
	if err != nil {
		return core.Gene[int]{}, fmt.Errorf("gene.CoreInt: %w", err)
	}
	constants := make([]int, len(g.Constants))
	for i, v := range g.Constants {
		constants[i] = int(v)
	}
	return core.Gene[int]{Symbols: syms, Constants: constants}, nil
}

// CoreFloat64 converts a legacy floating-point gene into the typed core representation.
func (g *Gene) CoreFloat64() (core.Gene[float64], error) {
	if g == nil {
		return core.Gene[float64]{}, errors.New("gene.CoreFloat64: gene cannot be nil")
	}
	cat, err := mn.CatalogFrom(mn.Math)
	if err != nil {
		return core.Gene[float64]{}, fmt.Errorf("gene.CoreFloat64: %w", err)
	}
	syms, err := core.ParseSymbols(g.Symbols, cat)
	if err != nil {
		return core.Gene[float64]{}, fmt.Errorf("gene.CoreFloat64: %w", err)
	}
	constants := make([]float64, len(g.Constants))
	copy(constants, g.Constants)
	return core.Gene[float64]{Symbols: syms, Constants: constants}, nil
}

// CoreVectorInt converts a legacy vector-int gene into the typed core representation.
func (g *Gene) CoreVectorInt() (core.Gene[functions.VectorInt], error) {
	if g == nil {
		return core.Gene[functions.VectorInt]{}, errors.New("gene.CoreVectorInt: gene cannot be nil")
	}
	if len(g.Constants) > 0 || hasLegacyConstantSymbol(g.Symbols) {
		return core.Gene[functions.VectorInt]{}, errors.New("gene.CoreVectorInt: vector-int genes with scalar constants are not supported by core.Gene")
	}
	cat, err := vin.CatalogFrom(vin.VectorIntFuncs)
	if err != nil {
		return core.Gene[functions.VectorInt]{}, fmt.Errorf("gene.CoreVectorInt: %w", err)
	}
	syms, err := core.ParseSymbols(g.Symbols, cat)
	if err != nil {
		return core.Gene[functions.VectorInt]{}, fmt.Errorf("gene.CoreVectorInt: %w", err)
	}
	return core.Gene[functions.VectorInt]{Symbols: syms}, nil
}

// NewFromCoreBool converts a typed boolean gene into the legacy representation.
func NewFromCoreBool(g core.Gene[bool]) (*Gene, error) {
	if len(g.Constants) > 0 || hasCoreConstantSymbol(g.Symbols) {
		return nil, errors.New("gene.NewFromCoreBool: bool genes do not support constants")
	}
	return New(g.KarvaString(), functions.Bool)
}

// NewFromCoreInt converts a typed integer gene into the legacy representation.
func NewFromCoreInt(g core.Gene[int]) (*Gene, error) {
	r, err := New(g.KarvaString(), functions.Int)
	if err != nil {
		return nil, err
	}
	r.Constants = make([]float64, len(g.Constants))
	for i, v := range g.Constants {
		r.Constants[i] = float64(v)
	}
	return r, nil
}

// NewFromCoreFloat64 converts a typed floating-point gene into the legacy representation.
func NewFromCoreFloat64(g core.Gene[float64]) (*Gene, error) {
	r, err := New(g.KarvaString(), functions.Float64)
	if err != nil {
		return nil, err
	}
	r.Constants = make([]float64, len(g.Constants))
	copy(r.Constants, g.Constants)
	return r, nil
}

// NewFromCoreVectorInt converts a typed vector-int gene into the legacy representation.
func NewFromCoreVectorInt(g core.Gene[functions.VectorInt]) (*Gene, error) {
	if len(g.Constants) > 0 || hasCoreConstantSymbol(g.Symbols) {
		return nil, errors.New("gene.NewFromCoreVectorInt: vector-int genes with typed constants are not supported by legacy gene.Gene")
	}
	return New(g.KarvaString(), functions.VectorInts)
}
