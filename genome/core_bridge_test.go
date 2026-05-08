// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package genome

import (
	"reflect"
	"testing"

	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
)

func TestCoreBool_RoundTrip(t *testing.T) {
	legacy := New([]*gene.Gene{
		gene.New("Nand.d0.d1", functions.Bool),
		gene.New("Or.d0.d1", functions.Bool),
	}, "And")

	typed, err := legacy.CoreBool()
	if err != nil {
		t.Fatalf("CoreBool(): %v", err)
	}
	got, err := typed.Eval([]bool{true, false})
	if err != nil {
		t.Fatalf("typed.Eval(): %v", err)
	}
	if want := legacy.EvalBool([]bool{true, false}); got != want {
		t.Fatalf("typed.Eval()=%v, want %v", got, want)
	}

	roundTrip, err := NewFromCoreBool(typed)
	if err != nil {
		t.Fatalf("NewFromCoreBool(): %v", err)
	}
	if roundTrip.LinkFunc != legacy.LinkFunc {
		t.Fatalf("round-trip link=%q, want %q", roundTrip.LinkFunc, legacy.LinkFunc)
	}
	for i := range legacy.Genes {
		if !reflect.DeepEqual(roundTrip.Genes[i].Symbols, legacy.Genes[i].Symbols) {
			t.Fatalf("gene[%d] symbols=%v, want %v", i, roundTrip.Genes[i].Symbols, legacy.Genes[i].Symbols)
		}
	}
}

func TestCoreInt_RoundTrip(t *testing.T) {
	g0 := gene.New("+.c0.d0", functions.Int)
	g0.Constants[0] = 3
	g1 := gene.New("-.d0.c0", functions.Int)
	g1.Constants[0] = 1
	legacy := New([]*gene.Gene{g0, g1}, "+")

	typed, err := legacy.CoreInt()
	if err != nil {
		t.Fatalf("CoreInt(): %v", err)
	}
	got, err := typed.Eval([]int{5})
	if err != nil {
		t.Fatalf("typed.Eval(): %v", err)
	}
	if want := legacy.EvalInt([]int{5}); got != want {
		t.Fatalf("typed.Eval()=%v, want %v", got, want)
	}

	roundTrip, err := NewFromCoreInt(typed)
	if err != nil {
		t.Fatalf("NewFromCoreInt(): %v", err)
	}
	if roundTrip.LinkFunc != legacy.LinkFunc {
		t.Fatalf("round-trip link=%q, want %q", roundTrip.LinkFunc, legacy.LinkFunc)
	}
	if !reflect.DeepEqual(roundTrip.Genes[0].Constants, legacy.Genes[0].Constants) {
		t.Fatalf("gene[0] constants=%v, want %v", roundTrip.Genes[0].Constants, legacy.Genes[0].Constants)
	}
	if !reflect.DeepEqual(roundTrip.Genes[1].Constants, legacy.Genes[1].Constants) {
		t.Fatalf("gene[1] constants=%v, want %v", roundTrip.Genes[1].Constants, legacy.Genes[1].Constants)
	}
}

func TestCoreFloat64_RoundTrip(t *testing.T) {
	g0 := gene.New("+.c0.d0", functions.Float64)
	g0.Constants[0] = 2.5
	g1 := gene.New("-.d0.c0", functions.Float64)
	g1.Constants[0] = 1.5
	legacy := New([]*gene.Gene{g0, g1}, "+")

	typed, err := legacy.CoreFloat64()
	if err != nil {
		t.Fatalf("CoreFloat64(): %v", err)
	}
	got, err := typed.Eval([]float64{4})
	if err != nil {
		t.Fatalf("typed.Eval(): %v", err)
	}
	if want := legacy.EvalMath([]float64{4}); got != want {
		t.Fatalf("typed.Eval()=%v, want %v", got, want)
	}

	roundTrip, err := NewFromCoreFloat64(typed)
	if err != nil {
		t.Fatalf("NewFromCoreFloat64(): %v", err)
	}
	if roundTrip.LinkFunc != legacy.LinkFunc {
		t.Fatalf("round-trip link=%q, want %q", roundTrip.LinkFunc, legacy.LinkFunc)
	}
	if !reflect.DeepEqual(roundTrip.Genes[0].Constants, legacy.Genes[0].Constants) {
		t.Fatalf("gene[0] constants=%v, want %v", roundTrip.Genes[0].Constants, legacy.Genes[0].Constants)
	}
	if !reflect.DeepEqual(roundTrip.Genes[1].Constants, legacy.Genes[1].Constants) {
		t.Fatalf("gene[1] constants=%v, want %v", roundTrip.Genes[1].Constants, legacy.Genes[1].Constants)
	}
}
