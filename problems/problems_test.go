// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package problems

import (
	"math"
	"testing"

	"github.com/gmlewis/gep/v2/core"
)

// --- helpers ---

// boolLinkOr returns a LinkOperator[bool] that OR-combines gene outputs.
func boolLinkOr(t *testing.T) core.LinkOperator[bool] {
	t.Helper()
	link, err := core.NewLinkFunc[bool]("Or", func(v []bool) bool {
		r := false
		for _, b := range v {
			r = r || b
		}
		return r
	})
	if err != nil {
		t.Fatalf("NewLinkFunc[bool]: %v", err)
	}
	return link
}

// floatLinkAdd returns a LinkOperator[float64] that sums gene outputs.
func floatLinkAdd(t *testing.T) core.LinkOperator[float64] {
	t.Helper()
	link, err := core.NewLinkFunc[float64]("+", func(v []float64) float64 {
		sum := 0.0
		for _, x := range v {
			sum += x
		}
		return sum
	})
	if err != nil {
		t.Fatalf("NewLinkFunc[float64]: %v", err)
	}
	return link
}

// terminalGene[T] creates a single-gene genome that returns its first terminal (d0).
func terminalBoolGenome(t *testing.T) core.Genome[bool] {
	t.Helper()
	return core.Genome[bool]{
		Genes: []core.Gene[bool]{
			{Symbols: []core.Symbol[bool]{{Kind: core.SymbolKindTerminal, Name: "d0", TerminalIndex: 0}}},
		},
		Link: boolLinkOr(t),
	}
}

// terminalFloatGenome creates a single-gene genome that returns its first terminal (d0).
func terminalFloatGenome(t *testing.T) core.Genome[float64] {
	t.Helper()
	return core.Genome[float64]{
		Genes: []core.Gene[float64]{
			{Symbols: []core.Symbol[float64]{{Kind: core.SymbolKindTerminal, Name: "d0", TerminalIndex: 0}}},
		},
		Link: floatLinkAdd(t),
	}
}

// --- BoolProblem tests ---

func TestNewBoolProblem_EmptyCasesReturnsError(t *testing.T) {
	_, err := NewBoolProblem(nil)
	if err == nil {
		t.Fatal("NewBoolProblem(nil): got nil error, want non-nil")
	}
	_, err = NewBoolProblem([]Case[bool]{})
	if err == nil {
		t.Fatal("NewBoolProblem([]): got nil error, want non-nil")
	}
}

func TestNewBoolProblem_NonEmpty(t *testing.T) {
	cases := []Case[bool]{
		{In: []bool{true}, Out: true},
	}
	p, err := NewBoolProblem(cases)
	if err != nil {
		t.Fatalf("NewBoolProblem: unexpected error: %v", err)
	}
	if len(p.Cases) != 1 {
		t.Fatalf("len(p.Cases) = %d, want 1", len(p.Cases))
	}
}

func TestBoolProblem_NumHitsScoringFunc_AllCorrect(t *testing.T) {
	// Genome returns d0; cases where d0 is always the correct answer.
	cases := []Case[bool]{
		{In: []bool{true}, Out: true},
		{In: []bool{false}, Out: false},
		{In: []bool{true}, Out: true},
		{In: []bool{false}, Out: false},
	}
	p, err := NewBoolProblem(cases)
	if err != nil {
		t.Fatalf("NewBoolProblem: %v", err)
	}

	sf, err := p.NumHitsScoringFunc(1.0)
	if err != nil {
		t.Fatalf("NumHitsScoringFunc: %v", err)
	}

	g := terminalBoolGenome(t)
	got := sf(g)
	want := 4.0 // scaleFactor=1 * 4 correct
	if got != want {
		t.Errorf("NumHitsScoringFunc score = %v, want %v", got, want)
	}
}

func TestBoolProblem_NumHitsScoringFunc_PartialCorrect(t *testing.T) {
	// Genome returns d0; NAND cases where d0 is only sometimes correct.
	cases := []Case[bool]{
		{In: []bool{false, false}, Out: true},  // d0=false, want true  → miss
		{In: []bool{false, true}, Out: true},   // d0=false, want true  → miss
		{In: []bool{true, false}, Out: true},   // d0=true,  want true  → hit
		{In: []bool{true, true}, Out: false},   // d0=true,  want false → miss
	}
	p, err := NewBoolProblem(cases)
	if err != nil {
		t.Fatalf("NewBoolProblem: %v", err)
	}

	sf, err := p.NumHitsScoringFunc(1000.0)
	if err != nil {
		t.Fatalf("NumHitsScoringFunc: %v", err)
	}

	g := terminalBoolGenome(t) // returns d0
	got := sf(g)
	want := 1000.0 // 1 hit * 1000
	if got != want {
		t.Errorf("NumHitsScoringFunc score = %v, want %v", got, want)
	}
}

// --- FloatProblem tests ---

func TestNewFloatProblem_EmptyCasesReturnsError(t *testing.T) {
	_, err := NewFloatProblem(nil)
	if err == nil {
		t.Fatal("NewFloatProblem(nil): got nil error, want non-nil")
	}
	_, err = NewFloatProblem([]Case[float64]{})
	if err == nil {
		t.Fatal("NewFloatProblem([]): got nil error, want non-nil")
	}
}

func TestNewFloatProblem_NonEmpty(t *testing.T) {
	cases := []Case[float64]{
		{In: []float64{1.0}, Out: 1.0},
	}
	p, err := NewFloatProblem(cases)
	if err != nil {
		t.Fatalf("NewFloatProblem: unexpected error: %v", err)
	}
	if len(p.Cases) != 1 {
		t.Fatalf("len(p.Cases) = %d, want 1", len(p.Cases))
	}
}

func TestFloatProblem_NumHitsAbsScoringFunc_AllCorrect(t *testing.T) {
	// Perfect regression: genome returns d0, cases are (x, x).
	cases := []Case[float64]{
		{In: []float64{1.0}, Out: 1.0},
		{In: []float64{2.0}, Out: 2.0},
		{In: []float64{3.0}, Out: 3.0},
	}
	p, err := NewFloatProblem(cases)
	if err != nil {
		t.Fatalf("NewFloatProblem: %v", err)
	}

	sf, err := p.NumHitsAbsScoringFunc(0.0, 1.0) // precision=0, scaleFactor=1
	if err != nil {
		t.Fatalf("NumHitsAbsScoringFunc: %v", err)
	}

	g := terminalFloatGenome(t)
	got := sf(g)
	want := 3.0 // 3 hits * scaleFactor=1
	if got != want {
		t.Errorf("NumHitsAbsScoringFunc score = %v, want %v", got, want)
	}
}

func TestFloatProblem_NumHitsAbsScoringFunc_InvalidPrecisionReturnsError(t *testing.T) {
	cases := []Case[float64]{{In: []float64{1.0}, Out: 1.0}}
	p, err := NewFloatProblem(cases)
	if err != nil {
		t.Fatalf("NewFloatProblem: %v", err)
	}
	_, err = p.NumHitsAbsScoringFunc(-0.1, 1.0) // invalid precision
	if err == nil {
		t.Fatal("NumHitsAbsScoringFunc(-0.1, 1): got nil error, want non-nil")
	}
}

func TestFloatProblem_MeanSquaredErrorAbsScoringFunc_PerfectFit(t *testing.T) {
	// Genome returns d0, cases are (x, x) → MSE = 0 → score = scaleFactor/(1+0) = scaleFactor.
	cases := []Case[float64]{
		{In: []float64{1.0}, Out: 1.0},
		{In: []float64{2.0}, Out: 2.0},
		{In: []float64{4.0}, Out: 4.0},
	}
	p, err := NewFloatProblem(cases)
	if err != nil {
		t.Fatalf("NewFloatProblem: %v", err)
	}

	sf, err := p.MeanSquaredErrorAbsScoringFunc(1000.0)
	if err != nil {
		t.Fatalf("MeanSquaredErrorAbsScoringFunc: %v", err)
	}

	g := terminalFloatGenome(t)
	got := sf(g)
	want := 1000.0 // perfect fit → score = 1000
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("MeanSquaredErrorAbsScoringFunc score = %v, want %v", got, want)
	}
}

func TestFloatProblem_RSquareScoringFunc_PerfectFit(t *testing.T) {
	// Genome returns d0, cases are (x, x) → R²=1 → score = scaleFactor.
	cases := []Case[float64]{
		{In: []float64{1.0}, Out: 1.0},
		{In: []float64{2.0}, Out: 2.0},
		{In: []float64{3.0}, Out: 3.0},
	}
	p, err := NewFloatProblem(cases)
	if err != nil {
		t.Fatalf("NewFloatProblem: %v", err)
	}

	sf, err := p.RSquareScoringFunc(1000.0)
	if err != nil {
		t.Fatalf("RSquareScoringFunc: %v", err)
	}

	g := terminalFloatGenome(t)
	got := sf(g)
	want := 1000.0 // R²=1 → score = 1000
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("RSquareScoringFunc score = %v, want %v", got, want)
	}
}
