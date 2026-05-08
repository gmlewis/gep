// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package core

import "testing"

type intNode struct {
	symbol string
	arity  int
	fn     func([]int) int
}

func (n intNode) Symbol() string  { return n.symbol }
func (n intNode) Arity() int      { return n.arity }
func (n intNode) Eval(v []int) int { return n.fn(v) }

func TestNewFunctionSymbol(t *testing.T) {
	n := intNode{
		symbol: "+",
		arity:  2,
		fn: func(v []int) int {
			return v[0] + v[1]
		},
	}

	sym, err := NewFunctionSymbol[int](n)
	if err != nil {
		t.Fatalf("NewFunctionSymbol returned error: %v", err)
	}
	if sym.Kind != SymbolKindFunction {
		t.Fatalf("sym.Kind=%v, want %v", sym.Kind, SymbolKindFunction)
	}
	if sym.Name != "+" {
		t.Fatalf("sym.Name=%q, want %q", sym.Name, "+")
	}
}

func TestNewFunctionSymbol_Validation(t *testing.T) {
	tests := []struct {
		name string
		node Node[int]
	}{
		{name: "nil node"},
		{name: "empty symbol", node: intNode{arity: 1, fn: func(v []int) int { return 0 }}},
		{name: "negative arity", node: intNode{symbol: "bad", arity: -1, fn: func(v []int) int { return 0 }}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := NewFunctionSymbol[int](tc.node); err == nil {
				t.Fatalf("NewFunctionSymbol(%q): got nil error, want non-nil", tc.name)
			}
		})
	}
}

func TestTerminalAndConstantSymbols(t *testing.T) {
	terminal, err := NewTerminalSymbol[int](2)
	if err != nil {
		t.Fatalf("NewTerminalSymbol returned error: %v", err)
	}
	if terminal.Kind != SymbolKindTerminal || terminal.Name != "d2" || terminal.TerminalIndex != 2 {
		t.Fatalf("unexpected terminal symbol: %+v", terminal)
	}

	constant, err := NewConstantSymbol[int](3)
	if err != nil {
		t.Fatalf("NewConstantSymbol returned error: %v", err)
	}
	if constant.Kind != SymbolKindConstant || constant.Name != "c3" || constant.ConstantIndex != 3 {
		t.Fatalf("unexpected constant symbol: %+v", constant)
	}
}

func TestTerminalAndConstantSymbols_NegativeIndex(t *testing.T) {
	if _, err := NewTerminalSymbol[int](-1); err == nil {
		t.Fatalf("NewTerminalSymbol(-1): got nil error, want non-nil")
	}
	if _, err := NewConstantSymbol[int](-1); err == nil {
		t.Fatalf("NewConstantSymbol(-1): got nil error, want non-nil")
	}
}

func TestGeneValidate(t *testing.T) {
	fn, err := NewFunctionSymbol[int](intNode{
		symbol: "+",
		arity:  2,
		fn: func(v []int) int {
			return v[0] + v[1]
		},
	})
	if err != nil {
		t.Fatalf("NewFunctionSymbol returned error: %v", err)
	}

	term, err := NewTerminalSymbol[int](0)
	if err != nil {
		t.Fatalf("NewTerminalSymbol returned error: %v", err)
	}
	cst, err := NewConstantSymbol[int](0)
	if err != nil {
		t.Fatalf("NewConstantSymbol returned error: %v", err)
	}

	g := Gene[int]{
		Symbols:   []Symbol[int]{fn, term, cst},
		Constants: []int{42},
	}
	if err := g.Validate(); err != nil {
		t.Fatalf("Gene.Validate returned error: %v", err)
	}
}

func TestGeneValidate_Errors(t *testing.T) {
	tests := []struct {
		name string
		gene Gene[int]
	}{
		{name: "empty symbols", gene: Gene[int]{}},
		{name: "unknown kind", gene: Gene[int]{Symbols: []Symbol[int]{{Kind: SymbolKindUnknown}}}},
		{name: "function symbol missing node", gene: Gene[int]{Symbols: []Symbol[int]{{Kind: SymbolKindFunction, Name: "+"}}}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.gene.Validate(); err == nil {
				t.Fatalf("gene.Validate(): got nil error, want non-nil")
			}
		})
	}
}

func TestGenomeValidate(t *testing.T) {
	fn, err := NewFunctionSymbol[int](intNode{
		symbol: "+",
		arity:  2,
		fn: func(v []int) int {
			return v[0] + v[1]
		},
	})
	if err != nil {
		t.Fatalf("NewFunctionSymbol returned error: %v", err)
	}
	term, err := NewTerminalSymbol[int](0)
	if err != nil {
		t.Fatalf("NewTerminalSymbol returned error: %v", err)
	}
	link, err := NewLinkFunc[int]("sum", func(v []int) int {
		r := 0
		for _, val := range v {
			r += val
		}
		return r
	})
	if err != nil {
		t.Fatalf("NewLinkFunc returned error: %v", err)
	}

	gen := Genome[int]{
		Genes: []Gene[int]{
			{Symbols: []Symbol[int]{fn, term}},
		},
		Link: link,
	}
	if err := gen.Validate(); err != nil {
		t.Fatalf("Genome.Validate returned error: %v", err)
	}
}

func TestGenomeValidate_Errors(t *testing.T) {
	link, err := NewLinkFunc[int]("sum", func(v []int) int { return len(v) })
	if err != nil {
		t.Fatalf("NewLinkFunc returned error: %v", err)
	}

	tests := []struct {
		name   string
		genome Genome[int]
	}{
		{name: "empty genes", genome: Genome[int]{Link: link}},
		{name: "nil link", genome: Genome[int]{Genes: []Gene[int]{{Symbols: []Symbol[int]{{Kind: SymbolKindTerminal}}}}}},
		{name: "invalid gene", genome: Genome[int]{Genes: []Gene[int]{{}}, Link: link}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.genome.Validate(); err == nil {
				t.Fatalf("genome.Validate(): got nil error, want non-nil")
			}
		})
	}
}

func TestNewLinkFunc_ValidationAndEval(t *testing.T) {
	if _, err := NewLinkFunc[int]("", func(v []int) int { return len(v) }); err == nil {
		t.Fatalf("NewLinkFunc with empty name: got nil error, want non-nil")
	}
	if _, err := NewLinkFunc[int]("sum", nil); err == nil {
		t.Fatalf("NewLinkFunc with nil fn: got nil error, want non-nil")
	}

	link, err := NewLinkFunc[int]("sum", func(v []int) int {
		r := 0
		for _, x := range v {
			r += x
		}
		return r
	})
	if err != nil {
		t.Fatalf("NewLinkFunc returned error: %v", err)
	}
	if got := link.Symbol(); got != "sum" {
		t.Fatalf("link.Symbol()=%q, want %q", got, "sum")
	}
	if got := link.Eval([]int{1, 2, 3}); got != 6 {
		t.Fatalf("link.Eval([1 2 3])=%v, want 6", got)
	}
}
