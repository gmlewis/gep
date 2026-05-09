// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model

import (
	"io"
	"os"
	"sync/atomic"
	"testing"

	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/genome"
)

func mustGene(t *testing.T, karva string, funcType functions.FuncType) *gene.Gene {
	t.Helper()
	g, err := gene.New(karva, funcType)
	if err != nil {
		t.Fatalf("gene.New(%q) error: %v", karva, err)
	}
	return g
}

func TestMaxArity(t *testing.T) {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 2},
		{Symbol: "*", Weight: 3},
		{Symbol: "/", Weight: 4},
	}
	if g, err := maxArity(funcs, functions.Float64); err != nil || g != 2 {
		w := 2
		t.Errorf("maxArity(%v, functions.Float64) = %v, want %v", funcs, g, w)
	}
	funcs = append(funcs, gene.FuncWeight{
		Symbol: "LT3A",
		Weight: 1,
	})
	if g, err := maxArity(funcs, functions.Float64); err != nil || g != 3 {
		w := 3
		t.Errorf("maxArity(%v, functions.Float64) = %v, want %v", funcs, g, w)
	}
}

func TestMaxArity_UnknownFuncTypeReturnsError(t *testing.T) {
	if _, err := maxArity(nil, functions.FuncType(999)); err == nil {
		t.Fatal("maxArity(nil, unknown) error = nil, want non-nil")
	}
}

func BenchmarkReplication(b *testing.B) {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
	}
	e := New(funcs, functions.Float64, 30, 8, 4, 1, 0, "+", nil, false)
	for i := 0; i < b.N; i++ {
		if err := e.replication(); err != nil {
			b.Fatalf("replication() error: %v", err)
		}
	}
}

func BenchmarkMutation(b *testing.B) {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
	}
	e := New(funcs, functions.Float64, 30, 8, 4, 1, 0, "+", nil, false)
	for i := 0; i < b.N; i++ {
		if err := e.mutation(); err != nil {
			b.Fatalf("mutation() error: %v", err)
		}
	}
}

func TestReplication(t *testing.T) {
	g := &Generation{
		Individuals: []*genome.Genome{
			{Score: -1000},
			{Score: -500},
			{Score: -100},
			{Score: -50},
			{Score: -10},
			{Score: -5},
			{Score: -1},
			{Score: 1},
			{Score: 5},
		},
	}

	before := len(g.Individuals)
	if err := g.replication(); err != nil {
		t.Fatalf("replication() error: %v", err)
	}
	got := len(g.Individuals)
	if want := before; got != want {
		t.Errorf("replication = %v individuals, want %v", got, want)
	}
}

func TestGetBestHandlesAllNegativeScores(t *testing.T) {
	g1 := &genome.Genome{}
	g2 := &genome.Genome{}
	g3 := &genome.Genome{}
	scores := map[*genome.Genome]float64{
		g1: -100,
		g2: -1,
		g3: -50,
	}
	g := &Generation{
		Individuals: []*genome.Genome{g1, g2, g3},
		ScoringFunc: func(gn *genome.Genome) float64 {
			return scores[gn]
		},
	}

	got, err := g.getBest()
	if err != nil {
		t.Fatalf("getBest() error: %v", err)
	}
	if got != g2 {
		t.Fatalf("getBest() = %p (score=%v), want %p (score=%v)", got, scores[got], g2, scores[g2])
	}
	if scores[got] != -1 {
		t.Fatalf("getBest() score = %v, want -1", scores[got])
	}
}

func TestGetBest_NilScoringFuncReturnsError(t *testing.T) {
	g := &Generation{
		Individuals: []*genome.Genome{{}},
		ScoringFunc: nil,
	}
	if _, err := g.getBest(); err == nil {
		t.Fatal("getBest() error = nil, want non-nil")
	}
}

func TestGetBest_MinimizeScore(t *testing.T) {
	g1 := &genome.Genome{}
	g2 := &genome.Genome{}
	g3 := &genome.Genome{}
	scores := map[*genome.Genome]float64{
		g1: 100,
		g2: 10,
		g3: 50,
	}
	g := &Generation{
		Individuals:   []*genome.Genome{g1, g2, g3},
		ScoringFunc:   func(gn *genome.Genome) float64 { return scores[gn] },
		MinimizeScore: true,
	}

	got, err := g.getBest()
	if err != nil {
		t.Fatalf("getBest() error: %v", err)
	}
	if got != g2 {
		t.Fatalf("getBest() = %p (score=%v), want %p (score=%v)", got, scores[got], g2, scores[g2])
	}
	if scores[got] != 10 {
		t.Fatalf("getBest() score = %v, want 10", scores[got])
	}
}

func TestSingleCrossover_MismatchedGenesDoesNotMutate(t *testing.T) {
	gene1 := &gene.Gene{Symbols: []string{"d0"}, HeadSize: 1}
	gene2 := &gene.Gene{Symbols: []string{"d0", "d1"}, HeadSize: 2}
	g := &Generation{
		Individuals: []*genome.Genome{
			{Genes: []*gene.Gene{gene1}},
			{Genes: []*gene.Gene{gene2}},
		},
	}
	if err := g.singleCrossover(0, 1); err == nil {
		t.Fatal("singleCrossover() error = nil, want non-nil")
	}
	if got, want := g.Individuals[0].Genes[0].Symbols, []string{"d0"}; len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("genome[0] symbols changed: got %v, want %v", got, want)
	}
	if got, want := g.Individuals[1].Genes[0].Symbols, []string{"d0", "d1"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("genome[1] symbols changed: got %v, want %v", got, want)
	}
}

func TestSingleCrossover_MismatchedGenesReturnsError(t *testing.T) {
	gene1 := &gene.Gene{Symbols: []string{"d0"}, HeadSize: 1}
	gene2 := &gene.Gene{Symbols: []string{"d0", "d1"}, HeadSize: 2}
	g := &Generation{
		Individuals: []*genome.Genome{
			{Genes: []*gene.Gene{gene1}},
			{Genes: []*gene.Gene{gene2}},
		},
	}
	if err := g.singleCrossover(0, 1); err == nil {
		t.Fatal("singleCrossover() error = nil, want non-nil")
	}
}

func TestSingleCrossover_InvalidatesGeneCaches(t *testing.T) {
	gene1 := mustGene(t, "d0", functions.Float64)
	gene1.HeadSize = 1
	gene2 := mustGene(t, "d1", functions.Float64)
	gene2.HeadSize = 1

	if got, err := gene1.EvalMath([]float64{10, 20}); err != nil || got != 10 {
		t.Fatalf("gene1.EvalMath(before) = %v, want 10", got)
	}
	if got, err := gene2.EvalMath([]float64{10, 20}); err != nil || got != 20 {
		t.Fatalf("gene2.EvalMath(before) = %v, want 20", got)
	}

	g := &Generation{
		Individuals: []*genome.Genome{
			{Genes: []*gene.Gene{gene1}},
			{Genes: []*gene.Gene{gene2}},
		},
	}
	if err := g.singleCrossover(0, 1); err != nil {
		t.Fatalf("singleCrossover() error: %v", err)
	}

	if got, err := gene1.EvalMath([]float64{10, 20}); err != nil || got != 20 {
		t.Fatalf("gene1.EvalMath(after) = %v, want 20", got)
	}
	if got, err := gene2.EvalMath([]float64{10, 20}); err != nil || got != 10 {
		t.Fatalf("gene2.EvalMath(after) = %v, want 10", got)
	}
}

// TestEvolve_StopFuncHaltsBeforeMaxIterations verifies that a custom StopFunc
// terminates evolution as soon as it returns true, without requiring Score >= 1000.
func TestEvolve_StopFuncHaltsBeforeMaxIterations(t *testing.T) {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
	}
	stopAfter := 3
	callCount := 0
	g := New(funcs, functions.Float64, 10, 4, 1, 1, 0, "+", func(gn *genome.Genome) float64 {
		return 0.5 // never reaches 1000
	}, false)
	g.StopFunc = func(best *genome.Genome) bool {
		callCount++
		return callCount >= stopAfter
	}
	if _, err := g.Evolve(1000); err != nil {
		t.Fatalf("Evolve() error: %v", err)
	}
	if callCount != stopAfter {
		t.Fatalf("StopFunc called %v times, want %v", callCount, stopAfter)
	}
}

// TestEvolve_NilStopFuncUsesDefaultThreshold verifies that when StopFunc is nil,
// evolution stops when the best genome's score reaches 1000.
func TestEvolve_NilStopFuncUsesDefaultThreshold(t *testing.T) {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
	}
	var callCount atomic.Int64
	g := New(funcs, functions.Float64, 5, 2, 1, 1, 0, "+", func(gn *genome.Genome) float64 {
		n := callCount.Add(1)
		if n >= 3 {
			return 1000.0 // triggers default stop
		}
		return 0.0
	}, false)
	// StopFunc is nil — default threshold of 1000 applies
	best, err := g.Evolve(100)
	if err != nil {
		t.Fatalf("Evolve() error: %v", err)
	}
	if best.Score < 1000.0 {
		t.Fatalf("expected best.Score >= 1000, got %v", best.Score)
	}
}

func TestEvolve_DoesNotWriteStdoutOnStop(t *testing.T) {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
	}
	g := New(funcs, functions.Float64, 5, 2, 1, 1, 0, "+", func(*genome.Genome) float64 {
		return 1000.0
	}, false)

	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	if _, err := g.Evolve(10); err != nil {
		t.Fatalf("Evolve() error: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Fatalf("stdout close error: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("stdout read error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("Evolve() wrote to stdout: %q", string(out))
	}
}

// TestNewWithSeed_DeterministicPopulation verifies that two calls to NewWithSeed
// with identical parameters and the same seed produce identical initial populations.
func TestNewWithSeed_DeterministicPopulation(t *testing.T) {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
	}
	const seed = int64(42)
	g1 := NewWithSeed(seed, funcs, functions.Float64, 10, 4, 2, 2, 1, "+", nil, false)
	g2 := NewWithSeed(seed, funcs, functions.Float64, 10, 4, 2, 2, 1, "+", nil, false)

	if len(g1.Individuals) != len(g2.Individuals) {
		t.Fatalf("population sizes differ: %v vs %v", len(g1.Individuals), len(g2.Individuals))
	}
	for i, ind1 := range g1.Individuals {
		ind2 := g2.Individuals[i]
		for j, g := range ind1.Genes {
			g2g := ind2.Genes[j]
			if len(g.Symbols) != len(g2g.Symbols) {
				t.Fatalf("individual[%v].gene[%v]: symbol count %v != %v", i, j, len(g.Symbols), len(g2g.Symbols))
			}
			for k, sym := range g.Symbols {
				if sym != g2g.Symbols[k] {
					t.Fatalf("individual[%v].gene[%v].symbol[%v]: %q != %q", i, j, k, sym, g2g.Symbols[k])
				}
			}
		}
	}
}

// TestNewWithSeed_DifferentSeedsProduceDifferentPopulations verifies that different
// seeds lead to different populations (with overwhelming probability).
func TestNewWithSeed_DifferentSeedsProduceDifferentPopulations(t *testing.T) {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
	}
	g1 := NewWithSeed(1, funcs, functions.Float64, 10, 4, 2, 2, 1, "+", nil, false)
	g2 := NewWithSeed(2, funcs, functions.Float64, 10, 4, 2, 2, 1, "+", nil, false)

	different := false
	for i, ind1 := range g1.Individuals {
		ind2 := g2.Individuals[i]
		for j, g := range ind1.Genes {
			g2g := ind2.Genes[j]
			for k, sym := range g.Symbols {
				if sym != g2g.Symbols[k] {
					different = true
					break
				}
			}
			if different {
				break
			}
		}
		if different {
			break
		}
	}
	if !different {
		t.Fatal("populations from different seeds are identical; expected differences")
	}
}

// TestNewWithSeed_DeterministicEvolution verifies that two identical seeded
// generations produce the same evolution trajectory after one Evolve step.
func TestNewWithSeed_DeterministicEvolution(t *testing.T) {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
	}
	const seed = int64(77)
	scorer := func(gn *genome.Genome) float64 { return 0.0 }

	g1 := NewWithSeed(seed, funcs, functions.Float64, 6, 3, 1, 1, 0, "+", scorer, false)
	g2 := NewWithSeed(seed, funcs, functions.Float64, 6, 3, 1, 1, 0, "+", scorer, false)

	// Run one step of evolution on each
	if _, err := g1.getBest(); err != nil {
		t.Fatalf("g1.getBest() error: %v", err)
	}
	if err := g1.replication(); err != nil {
		t.Fatalf("g1.replication() error: %v", err)
	}
	if err := g1.mutation(); err != nil {
		t.Fatalf("g1.mutation() error: %v", err)
	}

	if _, err := g2.getBest(); err != nil {
		t.Fatalf("g2.getBest() error: %v", err)
	}
	if err := g2.replication(); err != nil {
		t.Fatalf("g2.replication() error: %v", err)
	}
	if err := g2.mutation(); err != nil {
		t.Fatalf("g2.mutation() error: %v", err)
	}

	// Populations should be identical after deterministic operations
	for i, ind1 := range g1.Individuals {
		ind2 := g2.Individuals[i]
		for j, g := range ind1.Genes {
			g2g := ind2.Genes[j]
			for k, sym := range g.Symbols {
				if sym != g2g.Symbols[k] {
					t.Fatalf("after evolution: individual[%v].gene[%v].symbol[%v]: %q != %q", i, j, k, sym, g2g.Symbols[k])
				}
			}
		}
	}
}
