// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package boolNodes

import (
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/core"
)

func TestCatalogFrom_RegistersAllSymbols(t *testing.T) {
	cat, err := CatalogFrom(BoolAllGates)
	if err != nil {
		t.Fatalf("CatalogFrom(BoolAllGates): %v", err)
	}
	for sym := range BoolAllGates {
		if _, ok := cat.Lookup(sym); !ok {
			t.Errorf("Lookup(%q) returned false, want true", sym)
		}
	}
}

func TestCatalogFrom_NodeArityAndSymbol(t *testing.T) {
	cat, err := CatalogFrom(BoolNandOnly)
	if err != nil {
		t.Fatalf("CatalogFrom(BoolNandOnly): %v", err)
	}
	node, ok := cat.Lookup("Nand")
	if !ok {
		t.Fatal("Lookup(\"Nand\") returned false, want true")
	}
	if node.Symbol() != "Nand" {
		t.Errorf("node.Symbol()=%q, want %q", node.Symbol(), "Nand")
	}
	if node.Arity() != 2 {
		t.Errorf("node.Arity()=%d, want 2", node.Arity())
	}
}

func TestCatalogFrom_NodeEval_Nand(t *testing.T) {
	cat, err := CatalogFrom(BoolNandOnly)
	if err != nil {
		t.Fatalf("CatalogFrom(BoolNandOnly): %v", err)
	}
	node, _ := cat.Lookup("Nand")
	tests := []struct {
		in  []bool
		out bool
	}{
		{[]bool{false, false}, true},
		{[]bool{false, true}, true},
		{[]bool{true, false}, true},
		{[]bool{true, true}, false},
	}
	for _, tc := range tests {
		if got := node.Eval(tc.in); got != tc.out {
			t.Errorf("Nand.Eval(%v)=%v, want %v", tc.in, got, tc.out)
		}
	}
}

func TestCatalogFrom_RoundTrip_ParseAndEval(t *testing.T) {
	// "Or.d0.d1" with d0=false, d1=true → true
	cat, err := CatalogFrom(BoolNotAndOrOnly)
	if err != nil {
		t.Fatalf("CatalogFrom(BoolNotAndOrOnly): %v", err)
	}
	syms, err := core.ParseSymbols(strings.Split("Or.d0.d1", "."), cat)
	if err != nil {
		t.Fatalf("ParseSymbols: %v", err)
	}
	gene := core.Gene[bool]{Symbols: syms}
	got, err := gene.Eval([]bool{false, true})
	if err != nil {
		t.Fatalf("gene.Eval: %v", err)
	}
	if !got {
		t.Errorf("Or(false, true) = false, want true")
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
	link, err := LinkFuncFrom("Or")
	if err != nil {
		t.Fatalf("LinkFuncFrom(\"Or\"): %v", err)
	}

	if got := link.Eval([]bool{true}); got != true {
		t.Errorf("link.Eval([true])=%v, want true", got)
	}

	if got := link.Eval([]bool{false, false, true}); got != true {
		t.Errorf("link.Eval([false false true])=%v, want true", got)
	}
}

func TestLinkFuncFrom_RejectsNonBinaryOperator(t *testing.T) {
	if _, err := LinkFuncFrom("And3"); err == nil {
		t.Fatal("LinkFuncFrom(\"And3\") error=nil, want non-binary error")
	}
}
