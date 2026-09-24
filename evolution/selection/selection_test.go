// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package selection

import (
	"math/rand"
	"testing"

	"github.com/gmlewis/gep/v3/core"
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

func newPopulation(t *testing.T, n int) []Candidate[int] {
	t.Helper()
	cat := newIntCatalog(t)
	link := newSumLink(t)
	pop := make([]Candidate[int], n)
	for i := range pop {
		genome, err := core.NewRandomGenome(cat, 2, 4, 2, 0, link, rand.New(rand.NewSource(int64(i+1))))
		if err != nil {
			t.Fatalf("NewRandomGenome(%d): %v", i, err)
		}
		pop[i] = Candidate[int]{Genome: genome}
	}
	return pop
}

func TestRoulette_PreservesPopulationSize(t *testing.T) {
	pop := newPopulation(t, 10)
	for i := range pop {
		pop[i].Score = 1
	}
	got := Roulette(pop, Config{}, rand.New(rand.NewSource(5)))
	if len(got) != len(pop) {
		t.Fatalf("len(result)=%d, want %d", len(got), len(pop))
	}
}

func TestRoulette_ReturnsDeepCopies(t *testing.T) {
	pop := newPopulation(t, 10)
	for i := range pop {
		pop[i].Score = float64(i + 1)
	}
	origPtrs := make([]any, len(pop))
	for i := range pop {
		origPtrs[i] = &pop[i].Genome.Genes[0].Symbols[0]
	}

	got := Roulette(pop, Config{}, rand.New(rand.NewSource(7)))
	for i := range got {
		newPtr := &got[i].Genome.Genes[0].Symbols[0]
		for _, origPtr := range origPtrs {
			if origPtr == newPtr {
				t.Fatalf("result[%d] aliases input genome symbol storage", i)
			}
		}
	}
}

func TestRoulette_FavorsHigherScores(t *testing.T) {
	pop := newPopulation(t, 20)
	pop[0].Score = 1000
	for i := 1; i < len(pop); i++ {
		pop[i].Score = 1
	}
	high := pop[0].Genome.KarvaString()

	got := Roulette(pop, Config{}, rand.New(rand.NewSource(99)))
	selected := 0
	for _, ind := range got {
		if ind.Genome.KarvaString() == high {
			selected++
		}
	}
	if selected < 5 {
		t.Fatalf("high-score candidate selected only %d/%d times; expected dominance", selected, len(got))
	}
}

func TestRoulette_FavorsLowerScoresWhenMinimizing(t *testing.T) {
	pop := newPopulation(t, 20)
	pop[0].Score = 1
	for i := 1; i < len(pop); i++ {
		pop[i].Score = 1000
	}
	low := pop[0].Genome.KarvaString()

	got := Roulette(pop, Config{MinimizeScore: true}, rand.New(rand.NewSource(99)))
	selected := 0
	for _, ind := range got {
		if ind.Genome.KarvaString() == low {
			selected++
		}
	}
	if selected < 5 {
		t.Fatalf("low-score candidate selected only %d/%d times with MinimizeScore; expected dominance", selected, len(got))
	}
}

// --- Tournament tests ---

func TestTournament_Empty(t *testing.T) {
	got := Tournament[int](nil, 2, Config{}, nil)
	if got != nil {
		t.Fatalf("Tournament(nil): got %v, want nil", got)
	}
	got = Tournament([]Candidate[int]{}, 2, Config{}, nil)
	if got != nil {
		t.Fatalf("Tournament(empty): got %v, want nil", got)
	}
}

func TestTournament_PreservesPopulationSize(t *testing.T) {
	pop := newPopulation(t, 10)
	for i := range pop {
		pop[i].Score = float64(i)
	}
	got := Tournament(pop, 3, Config{}, rand.New(rand.NewSource(5)))
	if len(got) != len(pop) {
		t.Fatalf("len(result)=%d, want %d", len(got), len(pop))
	}
}

func TestTournament_ReturnsDeepCopies(t *testing.T) {
	pop := newPopulation(t, 10)
	for i := range pop {
		pop[i].Score = float64(i + 1)
	}
	origPtrs := make([]any, len(pop))
	for i := range pop {
		origPtrs[i] = &pop[i].Genome.Genes[0].Symbols[0]
	}

	got := Tournament(pop, 2, Config{}, rand.New(rand.NewSource(7)))
	for i := range got {
		newPtr := &got[i].Genome.Genes[0].Symbols[0]
		for _, origPtr := range origPtrs {
			if origPtr == newPtr {
				t.Fatalf("result[%d] aliases input genome symbol storage", i)
			}
		}
	}
}

func TestTournament_FavorsHigherScores(t *testing.T) {
	pop := newPopulation(t, 20)
	for i := range 5 {
		pop[i].Score = 1000
	}
	for i := 5; i < len(pop); i++ {
		pop[i].Score = 1
	}

	got := Tournament(pop, 5, Config{}, rand.New(rand.NewSource(99)))
	selected := 0
	for _, ind := range got {
		if ind.Score == 1000 {
			selected++
		}
	}
	if selected < 10 {
		t.Fatalf("top candidates selected only %d/%d times; expected >= 10 with k=5", selected, len(got))
	}
}

func TestTournament_FavorsLowerScoresWhenMinimizing(t *testing.T) {
	pop := newPopulation(t, 20)
	for i := range 5 {
		pop[i].Score = 1
	}
	for i := 5; i < len(pop); i++ {
		pop[i].Score = 1000
	}

	got := Tournament(pop, 5, Config{MinimizeScore: true}, rand.New(rand.NewSource(99)))
	selected := 0
	for _, ind := range got {
		if ind.Score == 1 {
			selected++
		}
	}
	if selected < 10 {
		t.Fatalf("top candidates selected only %d/%d times with MinimizeScore; expected >= 10 with k=5", selected, len(got))
	}
}

func TestTournament_Deterministic(t *testing.T) {
	pop := newPopulation(t, 15)
	for i := range pop {
		pop[i].Score = float64(i)
	}
	got1 := Tournament(pop, 3, Config{}, rand.New(rand.NewSource(42)))
	got2 := Tournament(pop, 3, Config{}, rand.New(rand.NewSource(42)))
	for i := range got1 {
		if got1[i].Genome.KarvaString() != got2[i].Genome.KarvaString() {
			t.Errorf("not deterministic at index %d: %q vs %q", i, got1[i].Genome.KarvaString(), got2[i].Genome.KarvaString())
		}
	}
}

func TestTournament_ClampsTournamentSize(t *testing.T) {
	pop := newPopulation(t, 5)
	// tournamentSize <= 0 defaults to 2
	got := Tournament(pop, 0, Config{}, rand.New(rand.NewSource(1)))
	if len(got) != 5 {
		t.Fatalf("got %d candidates, want 5", len(got))
	}
	// tournamentSize > len(pop) clamps to len(pop)
	got = Tournament(pop, 999, Config{}, rand.New(rand.NewSource(1)))
	if len(got) != 5 {
		t.Fatalf("got %d candidates, want 5", len(got))
	}
}
