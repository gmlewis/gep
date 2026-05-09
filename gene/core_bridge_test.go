// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"reflect"
	"testing"

	"github.com/gmlewis/gep/v2/functions"
)

func mustNewGene(t *testing.T, karva string, funcType functions.FuncType) *Gene {
	t.Helper()
	g, err := New(karva, funcType)
	if err != nil {
		t.Fatalf("New(%q) error: %v", karva, err)
	}
	return g
}

func TestCoreBool_RoundTrip(t *testing.T) {
	legacy := mustNewGene(t, "Nand.d0.d1", functions.Bool)

	typed, err := legacy.CoreBool()
	if err != nil {
		t.Fatalf("CoreBool(): %v", err)
	}
	got, err := typed.Eval([]bool{true, false})
	if err != nil {
		t.Fatalf("typed.Eval(): %v", err)
	}
	want, err := legacy.EvalBool([]bool{true, false})
	if err != nil {
		t.Fatalf("legacy.EvalBool(): %v", err)
	}
	if got != want {
		t.Fatalf("typed.Eval()=%v, want %v", got, want)
	}

	roundTrip, err := NewFromCoreBool(typed)
	if err != nil {
		t.Fatalf("NewFromCoreBool(): %v", err)
	}
	if !reflect.DeepEqual(roundTrip.Symbols, legacy.Symbols) {
		t.Fatalf("round-trip symbols=%v, want %v", roundTrip.Symbols, legacy.Symbols)
	}
}

func TestCoreInt_RoundTrip(t *testing.T) {
	legacy := mustNewGene(t, "+.c0.d0", functions.Int)
	legacy.Constants[0] = 7

	typed, err := legacy.CoreInt()
	if err != nil {
		t.Fatalf("CoreInt(): %v", err)
	}
	got, err := typed.Eval([]int{5})
	if err != nil {
		t.Fatalf("typed.Eval(): %v", err)
	}
	want, err := legacy.EvalInt([]int{5})
	if err != nil {
		t.Fatalf("legacy.EvalInt(): %v", err)
	}
	if got != want {
		t.Fatalf("typed.Eval()=%v, want %v", got, want)
	}

	roundTrip, err := NewFromCoreInt(typed)
	if err != nil {
		t.Fatalf("NewFromCoreInt(): %v", err)
	}
	if !reflect.DeepEqual(roundTrip.Symbols, legacy.Symbols) {
		t.Fatalf("round-trip symbols=%v, want %v", roundTrip.Symbols, legacy.Symbols)
	}
	if !reflect.DeepEqual(roundTrip.Constants, legacy.Constants) {
		t.Fatalf("round-trip constants=%v, want %v", roundTrip.Constants, legacy.Constants)
	}
}

func TestCoreFloat64_RoundTrip(t *testing.T) {
	legacy := mustNewGene(t, "+.c0.d0", functions.Float64)
	legacy.Constants[0] = 2.5

	typed, err := legacy.CoreFloat64()
	if err != nil {
		t.Fatalf("CoreFloat64(): %v", err)
	}
	got, err := typed.Eval([]float64{4})
	if err != nil {
		t.Fatalf("typed.Eval(): %v", err)
	}
	want, err := legacy.EvalMath([]float64{4})
	if err != nil {
		t.Fatalf("legacy.EvalMath(): %v", err)
	}
	if got != want {
		t.Fatalf("typed.Eval()=%v, want %v", got, want)
	}

	roundTrip, err := NewFromCoreFloat64(typed)
	if err != nil {
		t.Fatalf("NewFromCoreFloat64(): %v", err)
	}
	if !reflect.DeepEqual(roundTrip.Symbols, legacy.Symbols) {
		t.Fatalf("round-trip symbols=%v, want %v", roundTrip.Symbols, legacy.Symbols)
	}
	if !reflect.DeepEqual(roundTrip.Constants, legacy.Constants) {
		t.Fatalf("round-trip constants=%v, want %v", roundTrip.Constants, legacy.Constants)
	}
}

func TestCoreVectorInt_RoundTrip(t *testing.T) {
	legacy := mustNewGene(t, "+.d0.d1", functions.VectorInts)

	typed, err := legacy.CoreVectorInt()
	if err != nil {
		t.Fatalf("CoreVectorInt(): %v", err)
	}
	got, err := typed.Eval([]functions.VectorInt{{1, 2}, {3, 4}})
	if err != nil {
		t.Fatalf("typed.Eval(): %v", err)
	}
	want, err := legacy.EvalVectorInt([]functions.VectorInt{{1, 2}, {3, 4}})
	if err != nil {
		t.Fatalf("legacy.EvalVectorInt(): %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("typed.Eval()=%v, want %v", got, want)
	}

	roundTrip, err := NewFromCoreVectorInt(typed)
	if err != nil {
		t.Fatalf("NewFromCoreVectorInt(): %v", err)
	}
	if !reflect.DeepEqual(roundTrip.Symbols, legacy.Symbols) {
		t.Fatalf("round-trip symbols=%v, want %v", roundTrip.Symbols, legacy.Symbols)
	}
}

func TestCoreBridge_UnsupportedConstants(t *testing.T) {
	boolGene := mustNewGene(t, "c0", functions.Bool)
	if _, err := boolGene.CoreBool(); err == nil {
		t.Fatal("CoreBool() with constants: got nil error, want non-nil")
	}

	vectorGene := mustNewGene(t, "c0", functions.VectorInts)
	if _, err := vectorGene.CoreVectorInt(); err == nil {
		t.Fatal("CoreVectorInt() with constants: got nil error, want non-nil")
	}
}
