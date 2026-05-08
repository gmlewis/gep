// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package evolution

import (
	"math/rand"
	"sync/atomic"
	"testing"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/evolution/mutation"
)

// --- helpers ---

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

// newTestGeneration creates a small deterministic generation for testing.
func newTestGeneration(t *testing.T, seed int64, sf ScoringFunc[int]) *Generation[int] {
	t.Helper()
	cat := newIntCatalog(t)
	link := newSumLink(t)
	g, err := NewWithSeed(seed, cat, 10, 4, 2, 2, 0, link, sf)
	if err != nil {
		t.Fatalf("NewWithSeed: %v", err)
	}
	return g
}

// --- New / NewWithSeed tests ---

func TestNew_ValidPopulation(t *testing.T) {
	cat := newIntCatalog(t)
	link := newSumLink(t)
	g, err := New(cat, 5, 3, 2, 2, 0, link, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if len(g.Individuals) != 5 {
		t.Fatalf("len(Individuals)=%d, want 5", len(g.Individuals))
	}
	for i, ind := range g.Individuals {
		if err := ind.Genome.Validate(); err != nil {
			t.Errorf("individual[%d].Genome.Validate: %v", i, err)
		}
	}
}

func TestNew_Errors(t *testing.T) {
	cat := newIntCatalog(t)
	link := newSumLink(t)

	if _, err := New[int](nil, 5, 3, 2, 2, 0, link, nil); err == nil {
		t.Error("New(nil cat): got nil error, want non-nil")
	}
	if _, err := New(cat, 0, 3, 2, 2, 0, link, nil); err == nil {
		t.Error("New(numIndividuals=0): got nil error, want non-nil")
	}
	if _, err := New(cat, -1, 3, 2, 2, 0, link, nil); err == nil {
		t.Error("New(numIndividuals=-1): got nil error, want non-nil")
	}
	if _, err := New[int](cat, 5, 3, 2, 2, 0, nil, nil); err == nil {
		t.Error("New(nil link): got nil error, want non-nil")
	}
}

func TestNewWithSeed_DeterministicPopulation(t *testing.T) {
	cat := newIntCatalog(t)
	link := newSumLink(t)
	const seed = int64(42)

	g1, err := NewWithSeed(seed, cat, 8, 4, 2, 2, 0, link, nil)
	if err != nil {
		t.Fatalf("NewWithSeed (first): %v", err)
	}
	g2, err := NewWithSeed(seed, cat, 8, 4, 2, 2, 0, link, nil)
	if err != nil {
		t.Fatalf("NewWithSeed (second): %v", err)
	}

	if len(g1.Individuals) != len(g2.Individuals) {
		t.Fatalf("population size mismatch: %d vs %d", len(g1.Individuals), len(g2.Individuals))
	}
	for i, ind1 := range g1.Individuals {
		ind2 := g2.Individuals[i]
		for j, gene1 := range ind1.Genome.Genes {
			gene2 := ind2.Genome.Genes[j]
			if len(gene1.Symbols) != len(gene2.Symbols) {
				t.Fatalf("individual[%d].gene[%d]: symbol count %d vs %d", i, j, len(gene1.Symbols), len(gene2.Symbols))
			}
			for k, sym1 := range gene1.Symbols {
				sym2 := gene2.Symbols[k]
				if sym1.Kind != sym2.Kind || sym1.Name != sym2.Name {
					t.Fatalf("individual[%d].gene[%d].symbol[%d]: {%v,%q} vs {%v,%q}",
						i, j, k, sym1.Kind, sym1.Name, sym2.Kind, sym2.Name)
				}
			}
		}
	}
}

func TestNewWithSeed_DifferentSeedsProduceDifferentPopulations(t *testing.T) {
	cat := newIntCatalog(t)
	link := newSumLink(t)

	g1, err := NewWithSeed(1, cat, 8, 4, 2, 2, 0, link, nil)
	if err != nil {
		t.Fatalf("NewWithSeed(1): %v", err)
	}
	g2, err := NewWithSeed(2, cat, 8, 4, 2, 2, 0, link, nil)
	if err != nil {
		t.Fatalf("NewWithSeed(2): %v", err)
	}

	different := false
outer:
	for i, ind1 := range g1.Individuals {
		for j, gene1 := range ind1.Genome.Genes {
			gene2 := g2.Individuals[i].Genome.Genes[j]
			for k, sym1 := range gene1.Symbols {
				if sym1.Name != gene2.Symbols[k].Name {
					different = true
					break outer
				}
			}
		}
	}
	if !different {
		t.Fatal("populations from different seeds are identical; expected differences")
	}
}

// --- Individual.Dup tests ---

func TestIndividualDup(t *testing.T) {
	cat := newIntCatalog(t)
	g, err := newRandomGeneHelper(cat)
	if err != nil {
		t.Fatalf("newRandomGeneHelper: %v", err)
	}
	link := newSumLink(t)
	ind := Individual[int]{
		Genome: core.Genome[int]{Genes: []core.Gene[int]{g}, Link: link},
		Score:  3.14,
	}
	dup := ind.Dup()
	if dup.Score != ind.Score {
		t.Fatalf("Dup().Score=%v, want %v", dup.Score, ind.Score)
	}
	// Mutate the dup's symbol; original must be unchanged.
	sym, _ := core.NewTerminalSymbol[int](99)
	dup.Genome.Genes[0].Symbols[0] = sym
	if ind.Genome.Genes[0].Symbols[0].Name == dup.Genome.Genes[0].Symbols[0].Name {
		t.Fatal("Dup is not a deep copy: modifying dup affected original")
	}
}

// newRandomGeneHelper is a small helper used in tests.
func newRandomGeneHelper(cat *core.Catalog[int]) (core.Gene[int], error) {
	return core.NewRandomGene(cat, 3, 2, 0, rand.New(rand.NewSource(1)))
}

// --- Evaluate tests ---

func TestEvaluate_ScoresAllIndividuals(t *testing.T) {
	var counter atomic.Int64
	sf := func(core.Genome[int]) float64 {
		return float64(counter.Add(1))
	}
	g := newTestGeneration(t, 1, sf)
	// Reset counter so we can verify all individuals are scored.
	counter.Store(0)
	g.Evaluate()
	got := int(counter.Load())
	if got != len(g.Individuals) {
		t.Fatalf("ScoringFunc called %d times, want %d", got, len(g.Individuals))
	}
	for i, ind := range g.Individuals {
		if ind.Score == 0 {
			t.Errorf("individual[%d].Score is still zero after Evaluate", i)
		}
	}
}

func TestEvaluate_NilScoringFuncIsNoop(t *testing.T) {
	g := newTestGeneration(t, 1, nil)
	for i := range g.Individuals {
		g.Individuals[i].Score = 7.0
	}
	g.Evaluate()
	for i, ind := range g.Individuals {
		if ind.Score != 7.0 {
			t.Errorf("individual[%d].Score changed by nil-func Evaluate: got %v, want 7.0", i, ind.Score)
		}
	}
}

// --- BestIndividual tests ---

func TestBestIndividual_MaximizeScore(t *testing.T) {
	g := newTestGeneration(t, 1, nil)
	for i := range g.Individuals {
		g.Individuals[i].Score = float64(i) // index 9 is best
	}
	best := g.BestIndividual()
	want := float64(len(g.Individuals) - 1)
	if best.Score != want {
		t.Fatalf("BestIndividual().Score=%v, want %v", best.Score, want)
	}
}

func TestBestIndividual_MinimizeScore(t *testing.T) {
	g := newTestGeneration(t, 1, nil)
	g.MinimizeScore = true
	for i := range g.Individuals {
		g.Individuals[i].Score = float64(i + 1) // index 0 (score=1) is best when minimizing
	}
	best := g.BestIndividual()
	if best.Score != 1 {
		t.Fatalf("BestIndividual(minimize).Score=%v, want 1", best.Score)
	}
}

func TestBestIndividual_AllNegativeScores(t *testing.T) {
	g := newTestGeneration(t, 1, nil)
	for i := range g.Individuals {
		g.Individuals[i].Score = -float64(i + 1) // index 0 (score=-1) is best
	}
	best := g.BestIndividual()
	if best.Score != -1 {
		t.Fatalf("BestIndividual(all-negative).Score=%v, want -1", best.Score)
	}
}

// --- Select tests ---

func TestSelect_PreservesPopulationSize(t *testing.T) {
	sf := func(core.Genome[int]) float64 { return 1.0 }
	g := newTestGeneration(t, 5, sf)
	before := len(g.Individuals)
	g.Evaluate()
	g.Select()
	if got := len(g.Individuals); got != before {
		t.Fatalf("Select changed population size: got %d, want %d", got, before)
	}
}

func TestSelect_AllIndividualsAreDeepCopies(t *testing.T) {
	sf := func(core.Genome[int]) float64 { return 1.0 }
	g := newTestGeneration(t, 7, sf)
	g.Evaluate()

	// Capture original pointers.
	origPtrs := make([]interface{}, len(g.Individuals))
	for i := range g.Individuals {
		origPtrs[i] = &g.Individuals[i].Genome.Genes[0].Symbols[0]
	}

	g.Select()

	// After selection, all individuals should be fresh copies (no aliasing with
	// the original population).
	for i := range g.Individuals {
		newPtr := &g.Individuals[i].Genome.Genes[0].Symbols[0]
		for _, origPtr := range origPtrs {
			if origPtr == newPtr {
				t.Errorf("individual[%d] shares symbol slice with original population (not a deep copy)", i)
			}
		}
	}
}

func TestSelect_FavorsBetterIndividuals(t *testing.T) {
	// Build a generation where one individual has a vastly higher score.
	// After many selections it should appear far more often than the rest.
	const seed = int64(99)
	cat := newIntCatalog(t)
	link := newSumLink(t)
	g, err := NewWithSeed(seed, cat, 20, 4, 2, 2, 0, link, nil)
	if err != nil {
		t.Fatalf("NewWithSeed: %v", err)
	}

	// Assign scores: index 0 gets 1000, all others get 1.
	g.Individuals[0].Score = 1000
	for i := 1; i < len(g.Individuals); i++ {
		g.Individuals[i].Score = 1
	}
	// Save the genome of the high-score individual so we can identify it later.
	highScoreKarva := g.Individuals[0].Genome.KarvaString()

	g.Select()

	// Count how many times the high-score individual was selected.
	selected := 0
	for _, ind := range g.Individuals {
		if ind.Genome.KarvaString() == highScoreKarva {
			selected++
		}
	}
	// With such a score disparity, the high-score individual should dominate.
	if selected < 5 {
		t.Fatalf("high-score individual selected only %d/%d times; expected dominance", selected, len(g.Individuals))
	}
}

// --- Evolve tests ---

func TestEvolve_StopFuncHaltsEarly(t *testing.T) {
	callCount := 0
	stopAfter := 3
	sf := func(core.Genome[int]) float64 { return 0.5 } // never reaches 1000
	g := newTestGeneration(t, 11, sf)
	g.StopFunc = func(Individual[int]) bool {
		callCount++
		return callCount >= stopAfter
	}
	g.Evolve(1000)
	if callCount != stopAfter {
		t.Fatalf("StopFunc called %d times, want %d", callCount, stopAfter)
	}
}

func TestEvolve_DefaultStopThreshold(t *testing.T) {
	var callCount atomic.Int64
	sf := func(core.Genome[int]) float64 {
		n := callCount.Add(1)
		if n >= int64(3*10) { // 3 full evaluations (10 individuals each)
			return 1000.0
		}
		return 0.0
	}
	g := newTestGeneration(t, 13, sf)
	best := g.Evolve(100)
	if best.Score < 1000.0 {
		t.Fatalf("Evolve (default threshold): best.Score=%v, want >= 1000", best.Score)
	}
}

func TestEvolve_ElitismPreservesBest(t *testing.T) {
	// Assign a score via ScoringFunc: the best individual at each generation
	// must always appear in the next generation's population.
	const nIter = 5
	const seed = int64(7)

	maxSeen := 0.0
	gen := 0
	sf := func(gn core.Genome[int]) float64 {
		// Stable score based on Karva representation length so we can track
		// the best without external state.
		return float64(len(gn.KarvaString()))
	}
	g := newTestGeneration(t, seed, sf)
	g.StopFunc = func(best Individual[int]) bool {
		gen++
		if best.Score > maxSeen {
			maxSeen = best.Score
		}
		// Verify the previous best still has score <= current best.
		if best.Score < maxSeen {
			return true // triggers a failure below
		}
		return gen >= nIter
	}
	g.Evolve(100)
	if gen < nIter {
		t.Fatalf("elitism violated: best score decreased (maxSeen=%v, last=%v)", maxSeen, g.BestIndividual().Score)
	}
}

func TestEvolve_DeterministicWithSeed(t *testing.T) {
	sf := func(core.Genome[int]) float64 { return 0.0 }

	g1 := newTestGeneration(t, 42, sf)
	g2 := newTestGeneration(t, 42, sf)

	g1.StopFunc = func(Individual[int]) bool { return false }
	g2.StopFunc = func(Individual[int]) bool { return false }

	// Run a fixed number of steps on each.
	for i := 0; i < 3; i++ {
		g1.Evaluate()
		g1.Select()
		g2.Evaluate()
		g2.Select()
	}

	for i, ind1 := range g1.Individuals {
		ind2 := g2.Individuals[i]
		if ind1.Genome.KarvaString() != ind2.Genome.KarvaString() {
			t.Fatalf("determinism violated at individual[%d]: %q vs %q",
				i, ind1.Genome.KarvaString(), ind2.Genome.KarvaString())
		}
	}
}

func TestEvolve_MaxIterationsReturnsCurrentBest(t *testing.T) {
	sf := func(core.Genome[int]) float64 { return 0.5 } // never triggers default stop
	g := newTestGeneration(t, 3, sf)
	best := g.Evolve(2)
	if best.Score == 0 && g.ScoringFunc != nil {
		t.Fatal("Evolve did not evaluate genomes before returning")
	}
	_ = best
}

// --- Mutate tests ---

func TestMutate_PreservesPopulationSize(t *testing.T) {
	g := newTestGeneration(t, 5, nil)
	before := len(g.Individuals)
	g.Mutate()
	if got := len(g.Individuals); got != before {
		t.Fatalf("Mutate changed population size: got %d, want %d", got, before)
	}
}

func TestMutate_ZeroRates_GenomesUnchanged(t *testing.T) {
	// With all zero rates, Mutate must not alter any genome.
	g := newTestGeneration(t, 6, nil)
	origKarvas := make([]string, len(g.Individuals))
	for i, ind := range g.Individuals {
		origKarvas[i] = ind.Genome.KarvaString()
	}
	g.Mutate()
	for i, ind := range g.Individuals {
		if ind.Genome.KarvaString() != origKarvas[i] {
			t.Errorf("individual[%d] changed with zero mutation rates", i)
		}
	}
}

func TestMutate_PointMutation_ChangesAtLeastOneGenome(t *testing.T) {
	// With PointMutationRate=1.0 every gene is mutated; some genomes must differ.
	g := newTestGeneration(t, 10, nil)
	origKarvas := make([]string, len(g.Individuals))
	for i, ind := range g.Individuals {
		origKarvas[i] = ind.Genome.KarvaString()
	}
	g.MutationConfig = mutation.Config{PointMutationRate: 1.0}
	g.Mutate()
	changed := 0
	for i, ind := range g.Individuals {
		if ind.Genome.KarvaString() != origKarvas[i] {
			changed++
		}
	}
	if changed == 0 {
		t.Fatal("PointMutationRate=1.0: no genomes were changed by Mutate")
	}
}

func TestMutate_MutatedGenomesAreValid(t *testing.T) {
	g := newTestGeneration(t, 10, nil)
	g.MutationConfig = mutation.Config{
		PointMutationRate:     1.0,
		InversionRate:         1.0,
		ISTranspositionRate:   1.0,
		MaxISLen:              2,
		RISTranspositionRate:  1.0,
		MaxRISLen:             2,
		GeneTranspositionRate: 1.0,
	}
	g.Mutate()
	for i, ind := range g.Individuals {
		if err := ind.Genome.Validate(); err != nil {
			t.Errorf("individual[%d].Genome invalid after Mutate: %v", i, err)
		}
	}
}

func TestMutate_HeadSizeOverriddenFromGeneration(t *testing.T) {
	// Mutate must use the headSize stored at construction time, not any value
	// left in MutationConfig. We verify this by setting a deliberately wrong
	// HeadSize in MutationConfig; the mutation should still produce valid genomes.
	g := newTestGeneration(t, 8, nil)
	g.MutationConfig = mutation.Config{
		PointMutationRate: 1.0,
		HeadSize:          9999, // will be overridden internally
	}
	g.Mutate()
	for i, ind := range g.Individuals {
		if err := ind.Genome.Validate(); err != nil {
			t.Errorf("individual[%d].Genome invalid: %v", i, err)
		}
	}
}

func TestEvolve_WithMutation_StillConverges(t *testing.T) {
	// Evolve must still converge when mutation operators are active.
	var callCount atomic.Int64
	sf := func(core.Genome[int]) float64 {
		n := callCount.Add(1)
		if n >= int64(5*10) {
			return 1000.0
		}
		return 0.0
	}
	g := newTestGeneration(t, 13, sf)
	g.MutationConfig = mutation.Config{
		PointMutationRate:     0.1,
		InversionRate:         0.1,
		ISTranspositionRate:   0.1,
		MaxISLen:              2,
		RISTranspositionRate:  0.1,
		MaxRISLen:             2,
		GeneTranspositionRate: 0.1,
	}
	best := g.Evolve(200)
	if best.Score < 1000.0 {
		t.Fatalf("Evolve with mutation: best.Score=%v, want >= 1000", best.Score)
	}
}
