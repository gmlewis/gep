// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package core

import (
	"strings"
	"testing"
)

type intNode struct {
	symbol string
	arity  int
	fn     func([]int) int
}

func (n intNode) Symbol() string   { return n.symbol }
func (n intNode) Arity() int       { return n.arity }
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

// --- Catalog tests ---

func newIntCatalog(t *testing.T) *Catalog[int] {
	t.Helper()
	cat := NewCatalog[int]()
	for _, n := range []intNode{
		{symbol: "+", arity: 2, fn: func(v []int) int { return v[0] + v[1] }},
		{symbol: "-", arity: 2, fn: func(v []int) int { return v[0] - v[1] }},
		{symbol: "*", arity: 2, fn: func(v []int) int { return v[0] * v[1] }},
	} {
		if err := cat.Register(n); err != nil {
			t.Fatalf("Register(%q): %v", n.symbol, err)
		}
	}
	return cat
}

func TestCatalog_RegisterAndLookup(t *testing.T) {
	cat := NewCatalog[int]()
	n := intNode{symbol: "+", arity: 2, fn: func(v []int) int { return v[0] + v[1] }}
	if err := cat.Register(n); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	got, ok := cat.Lookup("+")
	if !ok {
		t.Fatalf("Lookup(%q) returned ok=false, want true", "+")
	}
	if got.Symbol() != "+" {
		t.Fatalf("Lookup(%q).Symbol()=%q, want %q", "+", got.Symbol(), "+")
	}
}

func TestCatalog_LookupMissing(t *testing.T) {
	cat := NewCatalog[int]()
	if _, ok := cat.Lookup("missing"); ok {
		t.Fatalf("Lookup(missing): got ok=true, want false")
	}
}

func TestCatalog_Register_Errors(t *testing.T) {
	cat := NewCatalog[int]()

	// nil node
	if err := cat.Register(nil); err == nil {
		t.Fatal("Register(nil): got nil error, want non-nil")
	}
	// empty symbol
	if err := cat.Register(intNode{symbol: "", arity: 1, fn: func(v []int) int { return 0 }}); err == nil {
		t.Fatal("Register(empty symbol): got nil error, want non-nil")
	}
	// duplicate symbol
	n := intNode{symbol: "+", arity: 2, fn: func(v []int) int { return v[0] + v[1] }}
	if err := cat.Register(n); err != nil {
		t.Fatalf("first Register: %v", err)
	}
	if err := cat.Register(n); err == nil {
		t.Fatal("duplicate Register: got nil error, want non-nil")
	}
}

// --- ParseSymbols tests ---

func TestParseSymbols_BasicExpression(t *testing.T) {
	cat := newIntCatalog(t)
	// "+.d0.d1"
	syms, err := ParseSymbols([]string{"+", "d0", "d1"}, cat)
	if err != nil {
		t.Fatalf("ParseSymbols returned error: %v", err)
	}
	if len(syms) != 3 {
		t.Fatalf("len(syms)=%d, want 3", len(syms))
	}
	if syms[0].Kind != SymbolKindFunction || syms[0].Name != "+" {
		t.Fatalf("syms[0]=%+v, want function '+'", syms[0])
	}
	if syms[1].Kind != SymbolKindTerminal || syms[1].TerminalIndex != 0 {
		t.Fatalf("syms[1]=%+v, want terminal d0", syms[1])
	}
	if syms[2].Kind != SymbolKindTerminal || syms[2].TerminalIndex != 1 {
		t.Fatalf("syms[2]=%+v, want terminal d1", syms[2])
	}
}

func TestParseSymbols_WithConstants(t *testing.T) {
	cat := newIntCatalog(t)
	syms, err := ParseSymbols([]string{"+", "c0", "d0"}, cat)
	if err != nil {
		t.Fatalf("ParseSymbols returned error: %v", err)
	}
	if syms[1].Kind != SymbolKindConstant || syms[1].ConstantIndex != 0 {
		t.Fatalf("syms[1]=%+v, want constant c0", syms[1])
	}
}

func TestParseSymbols_Errors(t *testing.T) {
	cat := newIntCatalog(t)

	// nil catalog
	if _, err := ParseSymbols[int]([]string{"+"}, nil); err == nil {
		t.Fatal("ParseSymbols(nil catalog): got nil error, want non-nil")
	}
	// empty symbol
	if _, err := ParseSymbols([]string{""}, cat); err == nil {
		t.Fatal("ParseSymbols(empty symbol): got nil error, want non-nil")
	}
	// unknown symbol
	if _, err := ParseSymbols([]string{"unknown"}, cat); err == nil {
		t.Fatal("ParseSymbols(unknown symbol): got nil error, want non-nil")
	}
}

// --- Gene.Eval tests ---

// makeGene builds a typed Gene from a dot-separated Karva expression string and
// optional integer constants.
func makeGene(t *testing.T, karva string, cat *Catalog[int], constants []int) Gene[int] {
	t.Helper()
	parts := strings.Split(karva, ".")
	syms, err := ParseSymbols(parts, cat)
	if err != nil {
		t.Fatalf("ParseSymbols(%q): %v", karva, err)
	}
	return Gene[int]{Symbols: syms, Constants: constants}
}

func TestGeneEval_SimpleAdd(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0.d1", cat, nil)
	got, err := g.Eval([]int{1, 2})
	if err != nil {
		t.Fatalf("Eval returned error: %v", err)
	}
	if got != 3 {
		t.Fatalf("Eval([1 2])=%v, want 3", got)
	}
}

func TestGeneEval_WithConstant(t *testing.T) {
	cat := newIntCatalog(t)
	// "+.c0.d0" with constant[0]=10, terminal[0]=5 → 15
	g := makeGene(t, "+.c0.d0", cat, []int{10})
	got, err := g.Eval([]int{5})
	if err != nil {
		t.Fatalf("Eval returned error: %v", err)
	}
	if got != 15 {
		t.Fatalf("Eval([5])=%v, want 15", got)
	}
}

func TestGeneEval_NestedExpression(t *testing.T) {
	cat := newIntCatalog(t)
	// "-.+.+.-.-.*.d0.d0.d0.d0.d0.d0.d0" = -(x² + x), x=8 → -72
	g := makeGene(t, "-.+.+.-.-.*.d0.d0.d0.d0.d0.d0.d0", cat, nil)
	got, err := g.Eval([]int{8})
	if err != nil {
		t.Fatalf("Eval returned error: %v", err)
	}
	if got != -72 {
		t.Fatalf("Eval([8])=%v, want -72", got)
	}
}

func TestGeneEval_SingleTerminal(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "d0", cat, nil)
	got, err := g.Eval([]int{42})
	if err != nil {
		t.Fatalf("Eval returned error: %v", err)
	}
	if got != 42 {
		t.Fatalf("Eval([42])=%v, want 42", got)
	}
}

func TestGeneEval_TerminalOutOfRange(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "d5", cat, nil)
	if _, err := g.Eval([]int{1, 2}); err == nil {
		t.Fatal("Eval with out-of-range terminal: got nil error, want non-nil")
	}
}

func TestGeneEval_ConstantOutOfRange(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "c2", cat, []int{0, 1}) // c2 requires 3 constants
	if _, err := g.Eval(nil); err == nil {
		t.Fatal("Eval with out-of-range constant: got nil error, want non-nil")
	}
}

func TestGeneEval_InvalidGeneErrors(t *testing.T) {
	g := Gene[int]{} // empty symbols → invalid
	if _, err := g.Eval(nil); err == nil {
		t.Fatal("Eval on invalid gene: got nil error, want non-nil")
	}
}

// --- Genome.Eval tests ---

func TestGenomeEval_SingleGene(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0.d1", cat, nil)
	link, err := NewLinkFunc[int]("id", func(v []int) int { return v[0] })
	if err != nil {
		t.Fatalf("NewLinkFunc: %v", err)
	}
	genome := Genome[int]{Genes: []Gene[int]{g}, Link: link}
	got, err := genome.Eval([]int{3, 4})
	if err != nil {
		t.Fatalf("Eval returned error: %v", err)
	}
	if got != 7 {
		t.Fatalf("Eval([3 4])=%v, want 7", got)
	}
}

func TestGenomeEval_MultipleGenes(t *testing.T) {
	cat := newIntCatalog(t)
	// gene0: d0+d1, gene1: d0*d1
	// link: sum → (d0+d1) + (d0*d1)
	// with d0=3, d1=4: (3+4) + (3*4) = 7 + 12 = 19
	gene0 := makeGene(t, "+.d0.d1", cat, nil)
	gene1 := makeGene(t, "*.d0.d1", cat, nil)
	link, err := NewLinkFunc[int]("+", func(v []int) int {
		sum := 0
		for _, x := range v {
			sum += x
		}
		return sum
	})
	if err != nil {
		t.Fatalf("NewLinkFunc: %v", err)
	}
	genome := Genome[int]{Genes: []Gene[int]{gene0, gene1}, Link: link}
	got, err := genome.Eval([]int{3, 4})
	if err != nil {
		t.Fatalf("Eval returned error: %v", err)
	}
	if got != 19 {
		t.Fatalf("Eval([3 4])=%v, want 19", got)
	}
}

func TestGenomeEval_InvalidGenomeErrors(t *testing.T) {
	genome := Genome[int]{} // empty genes → invalid
	if _, err := genome.Eval(nil); err == nil {
		t.Fatal("Eval on invalid genome: got nil error, want non-nil")
	}
}
