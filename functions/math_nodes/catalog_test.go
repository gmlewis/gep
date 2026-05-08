// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package mathNodes

import (
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/core"
)

func TestCatalogFrom_RegistersAllSymbols(t *testing.T) {
	cat, err := CatalogFrom(Math)
	if err != nil {
		t.Fatalf("CatalogFrom(Math): %v", err)
	}
	for sym := range Math {
		if _, ok := cat.Lookup(sym); !ok {
			t.Errorf("Lookup(%q) returned false, want true", sym)
		}
	}
}

func TestCatalogFrom_NodeArityAndSymbol(t *testing.T) {
	cat, err := CatalogFrom(Math)
	if err != nil {
		t.Fatalf("CatalogFrom(Math): %v", err)
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
	cat, err := CatalogFrom(Math)
	if err != nil {
		t.Fatalf("CatalogFrom(Math): %v", err)
	}
	node, _ := cat.Lookup("+")
	if got := node.Eval([]float64{2.0, 3.0}); got != 5.0 {
		t.Errorf("+.Eval([2 3])=%v, want 5", got)
	}
}

func TestCatalogFrom_RoundTrip_ParseAndEval(t *testing.T) {
	// "+.*.d0.d1.d0.d1" = (d0*d1) + d0*d1 — but that's too deep for a simple gene.
	// Use "+.d0.d1" with d0=3, d1=4 → 7
	cat, err := CatalogFrom(Math)
	if err != nil {
		t.Fatalf("CatalogFrom(Math): %v", err)
	}
	syms, err := core.ParseSymbols(strings.Split("+.d0.d1", "."), cat)
	if err != nil {
		t.Fatalf("ParseSymbols: %v", err)
	}
	gene := core.Gene[float64]{Symbols: syms}
	got, err := gene.Eval([]float64{3.0, 4.0})
	if err != nil {
		t.Fatalf("gene.Eval: %v", err)
	}
	if got != 7.0 {
		t.Errorf("+.d0.d1 with d0=3,d1=4 = %v, want 7", got)
	}
}

func TestCatalogFrom_RoundTrip_MultiGeneGenome(t *testing.T) {
	// gene0: *.d0.d1 with d0=3, d1=4 → 12
	// gene1: +.d0.d1 with d0=3, d1=4 → 7
	// link: + → 12+7 = 19
	cat, err := CatalogFrom(Math)
	if err != nil {
		t.Fatalf("CatalogFrom(Math): %v", err)
	}

	parse := func(karva string) core.Gene[float64] {
		t.Helper()
		syms, err := core.ParseSymbols(strings.Split(karva, "."), cat)
		if err != nil {
			t.Fatalf("ParseSymbols(%q): %v", karva, err)
		}
		return core.Gene[float64]{Symbols: syms}
	}

	link, err := core.NewLinkFunc[float64]("+", func(v []float64) float64 {
		sum := 0.0
		for _, x := range v {
			sum += x
		}
		return sum
	})
	if err != nil {
		t.Fatalf("NewLinkFunc: %v", err)
	}

	genome := core.Genome[float64]{
		Genes: []core.Gene[float64]{parse("*.d0.d1"), parse("+.d0.d1")},
		Link:  link,
	}
	got, err := genome.Eval([]float64{3.0, 4.0})
	if err != nil {
		t.Fatalf("genome.Eval: %v", err)
	}
	if got != 19.0 {
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
