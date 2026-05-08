// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package vectorIntNodes

import (
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/functions"
)

func TestCatalogFrom_RegistersAllSymbols(t *testing.T) {
	cat, err := CatalogFrom(VectorIntFuncs)
	if err != nil {
		t.Fatalf("CatalogFrom(VectorIntFuncs): %v", err)
	}
	for sym := range VectorIntFuncs {
		if _, ok := cat.Lookup(sym); !ok {
			t.Errorf("Lookup(%q) returned false, want true", sym)
		}
	}
}

func TestCatalogFrom_NodeArityAndSymbol(t *testing.T) {
	cat, err := CatalogFrom(VectorIntFuncs)
	if err != nil {
		t.Fatalf("CatalogFrom(VectorIntFuncs): %v", err)
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
	cat, err := CatalogFrom(VectorIntFuncs)
	if err != nil {
		t.Fatalf("CatalogFrom(VectorIntFuncs): %v", err)
	}
	node, _ := cat.Lookup("+")
	// [1,2] + [3,4] = [4,6]
	a := functions.VectorInt{1, 2}
	b := functions.VectorInt{3, 4}
	got := node.Eval([]functions.VectorInt{a, b})
	want := functions.VectorInt{4, 6}
	if len(got) != len(want) {
		t.Fatalf("+.Eval([%v %v]) len=%d, want %d", a, b, len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("+.Eval([%v %v])[%d]=%v, want %v", a, b, i, got[i], want[i])
		}
	}
}

func TestCatalogFrom_RoundTrip_ParseAndEval(t *testing.T) {
	// "+.d0.d1" with d0=[1,2], d1=[3,4] → [4,6]
	cat, err := CatalogFrom(VectorIntFuncs)
	if err != nil {
		t.Fatalf("CatalogFrom(VectorIntFuncs): %v", err)
	}
	syms, err := core.ParseSymbols(strings.Split("+.d0.d1", "."), cat)
	if err != nil {
		t.Fatalf("ParseSymbols: %v", err)
	}
	gene := core.Gene[functions.VectorInt]{Symbols: syms}
	got, err := gene.Eval([]functions.VectorInt{{1, 2}, {3, 4}})
	if err != nil {
		t.Fatalf("gene.Eval: %v", err)
	}
	want := functions.VectorInt{4, 6}
	if len(got) != len(want) {
		t.Fatalf("+.d0.d1 len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("+.d0.d1[%d]=%v, want %v", i, got[i], want[i])
		}
	}
}

func TestCatalogFrom_RoundTrip_MultiGeneGenome(t *testing.T) {
	// gene0: +.d0.d1 with d0=[10,20], d1=[1,2] → [11,22]
	// gene1: *.d0.d1 with d0=[10,20], d1=[2,3] — but we use d0=[10,20], d1=[1,2] → [10,40]
	// For simplicity: gene0: +.d0.d1 → [11,22]; gene1: -.d0.d1 → [9,18]
	// link: element-wise + of the two gene outputs → [20,40]
	cat, err := CatalogFrom(VectorIntFuncs)
	if err != nil {
		t.Fatalf("CatalogFrom(VectorIntFuncs): %v", err)
	}

	parse := func(karva string) core.Gene[functions.VectorInt] {
		t.Helper()
		syms, err := core.ParseSymbols(strings.Split(karva, "."), cat)
		if err != nil {
			t.Fatalf("ParseSymbols(%q): %v", karva, err)
		}
		return core.Gene[functions.VectorInt]{Symbols: syms}
	}

	link, err := core.NewLinkFunc[functions.VectorInt]("+", func(v []functions.VectorInt) functions.VectorInt {
		if len(v) == 0 {
			return nil
		}
		result := make(functions.VectorInt, len(v[0]))
		for _, vec := range v {
			for i, val := range vec {
				result[i] += val
			}
		}
		return result
	})
	if err != nil {
		t.Fatalf("NewLinkFunc: %v", err)
	}

	genome := core.Genome[functions.VectorInt]{
		Genes: []core.Gene[functions.VectorInt]{parse("+.d0.d1"), parse("-.d0.d1")},
		Link:  link,
	}
	// d0=[10,20], d1=[1,2]
	// gene0: [10+1, 20+2] = [11,22]
	// gene1: [10-1, 20-2] = [9,18]
	// link:  [11+9, 22+18] = [20,40]
	got, err := genome.Eval([]functions.VectorInt{{10, 20}, {1, 2}})
	if err != nil {
		t.Fatalf("genome.Eval: %v", err)
	}
	want := functions.VectorInt{20, 40}
	if len(got) != len(want) {
		t.Fatalf("genome.Eval len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("genome.Eval[%d]=%v, want %v", i, got[i], want[i])
		}
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
