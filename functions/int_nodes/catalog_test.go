// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package intNodes

import (
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/core"
)

func TestCatalogFrom_RegistersAllSymbols(t *testing.T) {
	cat, err := CatalogFrom(Int)
	if err != nil {
		t.Fatalf("CatalogFrom(Int): %v", err)
	}
	for sym := range Int {
		if _, ok := cat.Lookup(sym); !ok {
			t.Errorf("Lookup(%q) returned false, want true", sym)
		}
	}
}

func TestCatalogFrom_NodeArityAndSymbol(t *testing.T) {
	cat, err := CatalogFrom(Int)
	if err != nil {
		t.Fatalf("CatalogFrom(Int): %v", err)
	}
	node, ok := cat.Lookup("+")
	if !ok {
		t.Fatal("Lookup(\"+\") returned false, want true")
	}
	if node.Symbol() != "+" {
		t.Errorf("node.Symbol()=%q, want %q", node.Symbol(), "+")
	}
	if node.Arity() != 2 {
		t.Errorf("node.Arity()=%d, want 2", node.Arity())
	}
}

func TestCatalogFrom_NodeEval_Add(t *testing.T) {
	cat, err := CatalogFrom(Int)
	if err != nil {
		t.Fatalf("CatalogFrom(Int): %v", err)
	}
	node, _ := cat.Lookup("+")
	if got := node.Eval([]int{3, 4}); got != 7 {
		t.Errorf("+.Eval([3 4])=%v, want 7", got)
	}
}

func TestCatalogFrom_RoundTrip_ParseAndEval(t *testing.T) {
	// "+.d0.d1" with d0=5, d1=6 → 11
	cat, err := CatalogFrom(Int)
	if err != nil {
		t.Fatalf("CatalogFrom(Int): %v", err)
	}
	syms, err := core.ParseSymbols(strings.Split("+.d0.d1", "."), cat)
	if err != nil {
		t.Fatalf("ParseSymbols: %v", err)
	}
	gene := core.Gene[int]{Symbols: syms}
	got, err := gene.Eval([]int{5, 6})
	if err != nil {
		t.Fatalf("gene.Eval: %v", err)
	}
	if got != 11 {
		t.Errorf("+.d0.d1 with d0=5,d1=6 = %v, want 11", got)
	}
}

func TestCatalogFrom_RoundTrip_MultiGeneGenome(t *testing.T) {
	// gene0: *.d0.d1 with d0=3, d1=4 → 12
	// gene1: +.d0.d1 with d0=3, d1=4 → 7
	// link: + → 19
	cat, err := CatalogFrom(Int)
	if err != nil {
		t.Fatalf("CatalogFrom(Int): %v", err)
	}

	parse := func(karva string) core.Gene[int] {
		t.Helper()
		syms, err := core.ParseSymbols(strings.Split(karva, "."), cat)
		if err != nil {
			t.Fatalf("ParseSymbols(%q): %v", karva, err)
		}
		return core.Gene[int]{Symbols: syms}
	}

	link, err := core.NewLinkFunc[int]("+", func(v []int) int {
		sum := 0
		for _, x := range v {
			sum += x
		}
		return sum
	})
	if err != nil {
		t.Fatalf("NewLinkFunc: %v", err)
	}

	genome := core.Genome[int]{
		Genes: []core.Gene[int]{parse("*.d0.d1"), parse("+.d0.d1")},
		Link:  link,
	}
	got, err := genome.Eval([]int{3, 4})
	if err != nil {
		t.Fatalf("genome.Eval: %v", err)
	}
	if got != 19 {
		t.Errorf("genome.Eval([3 4])=%v, want 19", got)
	}
}

func TestCatalogFrom_EmptyFuncMap(t *testing.T) {
	cat, err := CatalogFrom(nil)
	if err != nil {
		t.Fatalf("CatalogFrom(nil): unexpected error: %v", err)
	}
	if _, ok := cat.Lookup("any"); ok {
		t.Error("Lookup on empty catalog returned true, want false")
	}
}

func TestLinkFuncFrom_FoldsSingleAndMultiGeneOutputs(t *testing.T) {
	link, err := LinkFuncFrom("+")
	if err != nil {
		t.Fatalf("LinkFuncFrom(\"+\"): %v", err)
	}

	if got := link.Eval([]int{7}); got != 7 {
		t.Errorf("link.Eval([7])=%v, want 7", got)
	}

	if got := link.Eval([]int{3, 4, 5}); got != 12 {
		t.Errorf("link.Eval([3 4 5])=%v, want 12", got)
	}
}
