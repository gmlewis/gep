// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package evaluation

import (
	"math/rand"
	"sync/atomic"
	"testing"

	"github.com/gmlewis/gep/v2/core"
)

type intNode struct {
	symbol string
	arity  int
	fn     func([]int) int
}

func (n intNode) Symbol() string   { return n.symbol }
func (n intNode) Arity() int       { return n.arity }
func (n intNode) Eval(v []int) int { return n.fn(v) }

func newIntCatalog(t *testing.T) *core.Catalog[int] {
	t.Helper()
	cat := core.NewCatalog[int]()
	for _, n := range []intNode{
		{symbol: "+", arity: 2, fn: func(v []int) int { return v[0] + v[1] }},
		{symbol: "-", arity: 2, fn: func(v []int) int { return v[0] - v[1] }},
		{symbol: "*", arity: 2, fn: func(v []int) int { return v[0] * v[1] }},
	} {
		if err := cat.Register(n); err != nil {
			t.Fatalf("catalog.Register(%q): %v", n.symbol, err)
		}
	}
	return cat
}

func newSumLink(t *testing.T) core.LinkOperator[int] {
	t.Helper()
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
	return link
}

func newGenomes(t *testing.T, n int) []core.Genome[int] {
	t.Helper()
	cat := newIntCatalog(t)
	link := newSumLink(t)
	genomes := make([]core.Genome[int], n)
	for i := range genomes {
		genome, err := core.NewRandomGenome(cat, 2, 4, 2, 0, link, rand.New(rand.NewSource(int64(i+1))))
		if err != nil {
			t.Fatalf("NewRandomGenome(%d): %v", i, err)
		}
		genomes[i] = genome
	}
	return genomes
}

func TestScoreAll_NilScoringFunc(t *testing.T) {
	genomes := newGenomes(t, 5)
	got := ScoreAll(genomes, nil, Config{Workers: 3})
	if len(got) != len(genomes) {
		t.Fatalf("len(scores)=%d, want %d", len(got), len(genomes))
	}
	for i, score := range got {
		if score != 0 {
			t.Fatalf("scores[%d]=%v, want 0 with nil scoring func", i, score)
		}
	}
}

func TestScoreAll_ScoresEveryGenome(t *testing.T) {
	genomes := newGenomes(t, 8)
	var calls atomic.Int64
	got := ScoreAll(genomes, func(core.Genome[int]) float64 {
		calls.Add(1)
		return 1.0
	}, Config{Workers: 2})

	if int(calls.Load()) != len(genomes) {
		t.Fatalf("calls=%d, want %d", calls.Load(), len(genomes))
	}
	for i, score := range got {
		if score != 1.0 {
			t.Fatalf("scores[%d]=%v, want 1.0", i, score)
		}
	}
}

func TestScoreAll_IndexAlignedResults(t *testing.T) {
	genomes := newGenomes(t, 6)
	sf := func(g core.Genome[int]) float64 { return float64(len(g.KarvaString())) }
	want := make([]float64, len(genomes))
	for i, g := range genomes {
		want[i] = sf(g)
	}

	got := ScoreAll(genomes, sf, Config{Workers: len(genomes) * 2})
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("scores[%d]=%v, want %v", i, got[i], want[i])
		}
	}
}
