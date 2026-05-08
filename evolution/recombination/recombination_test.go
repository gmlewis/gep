// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package recombination

import (
	"math/rand"
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

func TestApply_EmptyInput(t *testing.T) {
	got := Apply[int](nil, Config{}, nil)
	if got != nil {
		t.Fatalf("Apply(nil): got %v, want nil", got)
	}
	got = Apply([]core.Genome[int]{}, Config{}, nil)
	if got != nil {
		t.Fatalf("Apply(empty): got %v, want nil", got)
	}
}

func TestApply_ZeroRatesReturnsDeepCopies(t *testing.T) {
	genomes := newGenomes(t, 5)
	origKarvas := make([]string, len(genomes))
	for i, g := range genomes {
		origKarvas[i] = g.KarvaString()
	}

	got := Apply(genomes, Config{}, rand.New(rand.NewSource(1)))
	if len(got) != len(genomes) {
		t.Fatalf("Apply returned %d genomes, want %d", len(got), len(genomes))
	}
	for i := range got {
		if got[i].KarvaString() != origKarvas[i] {
			t.Errorf("genome[%d] changed with zero rates: got %q, want %q", i, got[i].KarvaString(), origKarvas[i])
		}
		if &got[i].Genes[0].Symbols[0] == &genomes[i].Genes[0].Symbols[0] {
			t.Errorf("genome[%d] aliases input storage (not a deep copy)", i)
		}
	}
}

func TestApply_PreservesPopulationSize(t *testing.T) {
	genomes := newGenomes(t, 8)
	got := Apply(genomes, Config{
		OnePointRate: 0.5,
		TwoPointRate: 0.5,
	}, rand.New(rand.NewSource(3)))
	if len(got) != len(genomes) {
		t.Fatalf("Apply returned %d genomes, want %d", len(got), len(genomes))
	}
}

func TestApply_OnePoint_ChangesSomeGenomes(t *testing.T) {
	genomes := newGenomes(t, 10)
	origKarvas := make([]string, len(genomes))
	for i, g := range genomes {
		origKarvas[i] = g.KarvaString()
	}

	got := Apply(genomes, Config{OnePointRate: 1.0}, rand.New(rand.NewSource(5)))
	changed := 0
	for i := range got {
		if got[i].KarvaString() != origKarvas[i] {
			changed++
		}
		if err := got[i].Validate(); err != nil {
			t.Errorf("genome[%d] invalid after one-point recombination: %v", i, err)
		}
	}
	if changed == 0 {
		t.Fatal("OnePointRate=1.0: no genomes were changed")
	}
}

func TestApply_TwoPoint_MutatedGenomeValid(t *testing.T) {
	genomes := newGenomes(t, 10)
	got := Apply(genomes, Config{TwoPointRate: 1.0}, rand.New(rand.NewSource(7)))
	for i, g := range got {
		if err := g.Validate(); err != nil {
			t.Errorf("genome[%d] invalid after two-point recombination: %v", i, err)
		}
	}
}

func TestApply_Deterministic(t *testing.T) {
	genomes := newGenomes(t, 6)
	cfg := Config{OnePointRate: 0.5, TwoPointRate: 0.3}
	got1 := Apply(genomes, cfg, rand.New(rand.NewSource(42)))
	got2 := Apply(genomes, cfg, rand.New(rand.NewSource(42)))
	for i := range got1 {
		if got1[i].KarvaString() != got2[i].KarvaString() {
			t.Errorf("Apply not deterministic at genome[%d]: %q vs %q",
				i, got1[i].KarvaString(), got2[i].KarvaString())
		}
	}
}
