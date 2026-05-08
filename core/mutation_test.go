// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package core

import (
	"math/rand"
	"strings"
	"testing"
)

// --- KarvaString tests ---

func TestGene_KarvaString_Empty(t *testing.T) {
	g := Gene[int]{}
	if got := g.KarvaString(); got != "" {
		t.Fatalf("KarvaString() on empty gene: got %q, want %q", got, "")
	}
}

func TestGene_KarvaString_SingleTerminal(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "d0", cat, nil)
	if got := g.KarvaString(); got != "d0" {
		t.Fatalf("KarvaString()=%q, want %q", got, "d0")
	}
}

func TestGene_KarvaString_Expression(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0.d1", cat, nil)
	if got := g.KarvaString(); got != "+.d0.d1" {
		t.Fatalf("KarvaString()=%q, want %q", got, "+.d0.d1")
	}
}

func TestGene_KarvaString_WithConstant(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "*.c0.d1", cat, []int{7})
	if got := g.KarvaString(); got != "*.c0.d1" {
		t.Fatalf("KarvaString()=%q, want %q", got, "*.c0.d1")
	}
}

func TestGenome_KarvaString_SingleGene(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0.d1", cat, nil)
	link, _ := NewLinkFunc[int]("id", func(v []int) int { return v[0] })
	genome := Genome[int]{Genes: []Gene[int]{g}, Link: link}
	want := "+.d0.d1"
	if got := genome.KarvaString(); got != want {
		t.Fatalf("Genome.KarvaString()=%q, want %q", got, want)
	}
}

func TestGenome_KarvaString_MultiGene(t *testing.T) {
	cat := newIntCatalog(t)
	g0 := makeGene(t, "+.d0.d1", cat, nil)
	g1 := makeGene(t, "*.d0.d1", cat, nil)
	link, _ := NewLinkFunc[int]("+", func(v []int) int { return v[0] + v[1] })
	genome := Genome[int]{Genes: []Gene[int]{g0, g1}, Link: link}
	want := "+.d0.d1|+|*.d0.d1"
	if got := genome.KarvaString(); got != want {
		t.Fatalf("Genome.KarvaString()=%q, want %q", got, want)
	}
}

func TestGenome_KarvaString_Empty(t *testing.T) {
	g := Genome[int]{}
	if got := g.KarvaString(); got != "" {
		t.Fatalf("KarvaString() on empty genome: got %q, want %q", got, "")
	}
}

// --- Dup tests ---

func TestGene_Dup_IndependentCopy(t *testing.T) {
	cat := newIntCatalog(t)
	orig := makeGene(t, "+.d0.d1", cat, []int{42})
	dup := orig.Dup()

	// Modify dup and verify orig is unaffected.
	dup.Symbols[0].Name = "MODIFIED"
	dup.Constants[0] = 999

	if orig.Symbols[0].Name != "+" {
		t.Errorf("orig.Symbols[0].Name=%q, want %q after dup modification", orig.Symbols[0].Name, "+")
	}
	if orig.Constants[0] != 42 {
		t.Errorf("orig.Constants[0]=%v, want 42 after dup modification", orig.Constants[0])
	}
}

func TestGene_Dup_SameContent(t *testing.T) {
	cat := newIntCatalog(t)
	orig := makeGene(t, "-.*.d0.d1.d0", cat, []int{5, 10})
	dup := orig.Dup()

	if orig.KarvaString() != dup.KarvaString() {
		t.Errorf("Dup() KarvaString mismatch: orig=%q dup=%q", orig.KarvaString(), dup.KarvaString())
	}
	if len(orig.Constants) != len(dup.Constants) {
		t.Errorf("Dup() Constants length mismatch: %d vs %d", len(orig.Constants), len(dup.Constants))
	}
}

func TestGenome_Dup_IndependentCopy(t *testing.T) {
	cat := newIntCatalog(t)
	g0 := makeGene(t, "+.d0.d1", cat, nil)
	g1 := makeGene(t, "*.d0.d1", cat, nil)
	link, _ := NewLinkFunc[int]("+", func(v []int) int { return v[0] + v[1] })
	orig := Genome[int]{Genes: []Gene[int]{g0, g1}, Link: link}

	dup := orig.Dup()
	dup.Genes[0].Symbols[0].Name = "MODIFIED"

	if orig.Genes[0].Symbols[0].Name != "+" {
		t.Errorf("orig.Genes[0].Symbols[0].Name=%q, want '+' after dup modification", orig.Genes[0].Symbols[0].Name)
	}
}

func TestGenome_Dup_SharedLink(t *testing.T) {
	cat := newIntCatalog(t)
	g0 := makeGene(t, "+.d0.d1", cat, nil)
	link, _ := NewLinkFunc[int]("sum", func(v []int) int { return v[0] })
	orig := Genome[int]{Genes: []Gene[int]{g0}, Link: link}
	dup := orig.Dup()

	// Link operator symbol should be preserved.
	if orig.Link.Symbol() != dup.Link.Symbol() {
		t.Errorf("Genome.Dup() link symbol mismatch: orig=%q dup=%q", orig.Link.Symbol(), dup.Link.Symbol())
	}
}

// --- PointMutate tests ---

func TestPointMutate_ValidGene(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0.d1", cat, nil)
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 20; i++ {
		mut, err := PointMutate(g, cat, 1, 2, 0, rng)
		if err != nil {
			t.Fatalf("PointMutate: %v", err)
		}
		if err := mut.Validate(); err != nil {
			t.Fatalf("mutated gene invalid: %v", err)
		}
		// Tail positions must be terminals.
		for j := 1; j < len(mut.Symbols); j++ {
			if mut.Symbols[j].Kind == SymbolKindFunction {
				t.Errorf("tail position %d is a function symbol, want terminal", j)
			}
		}
	}
}

func TestPointMutate_TailRemainsTerminal(t *testing.T) {
	cat := newIntCatalog(t) // maxArity=2, headSize=1 → tail=[1..]
	// Build a larger gene so we can observe tail mutations
	g, err := NewRandomGene(cat, 5, 3, 0, rand.New(rand.NewSource(0)))
	if err != nil {
		t.Fatalf("NewRandomGene: %v", err)
	}
	headSize := 5
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 50; i++ {
		mut, err := PointMutate(g, cat, headSize, 3, 0, rng)
		if err != nil {
			t.Fatalf("PointMutate: %v", err)
		}
		for j := headSize; j < len(mut.Symbols); j++ {
			if mut.Symbols[j].Kind == SymbolKindFunction {
				t.Errorf("tail[%d] is a function after point mutation", j)
			}
		}
	}
}

func TestPointMutate_Errors(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0.d1", cat, nil)

	if _, err := PointMutate[int](Gene[int]{}, cat, 1, 2, 0, nil); err == nil {
		t.Fatal("PointMutate(empty gene): want error, got nil")
	}
	if _, err := PointMutate(g, nil, 1, 2, 0, nil); err == nil {
		t.Fatal("PointMutate(nil cat): want error, got nil")
	}
	if _, err := PointMutate(g, cat, -1, 2, 0, nil); err == nil {
		t.Fatal("PointMutate(headSize=-1): want error, got nil")
	}
	if _, err := PointMutate(g, cat, 1, -1, 0, nil); err == nil {
		t.Fatal("PointMutate(numTerminals=-1): want error, got nil")
	}
	if _, err := PointMutate(g, cat, 1, 0, -1, nil); err == nil {
		t.Fatal("PointMutate(numConstants=-1): want error, got nil")
	}
	if _, err := PointMutate(g, cat, 1, 0, 0, nil); err == nil {
		t.Fatal("PointMutate(no terminals): want error, got nil")
	}
}

// --- Inversion tests ---

func TestInversion_ValidGene(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.*.-.d0.d1.d0.d1", cat, nil)
	rng := rand.New(rand.NewSource(5))
	for i := 0; i < 20; i++ {
		inv, err := Inversion(g, 3, rng)
		if err != nil {
			t.Fatalf("Inversion: %v", err)
		}
		if len(inv.Symbols) != len(g.Symbols) {
			t.Errorf("Inversion changed symbol count: %d vs %d", len(inv.Symbols), len(g.Symbols))
		}
		// Tail must be unchanged (positions 3..end).
		for j := 3; j < len(g.Symbols); j++ {
			if inv.Symbols[j].Name != g.Symbols[j].Name {
				t.Errorf("tail symbol[%d] changed: %q vs %q", j, inv.Symbols[j].Name, g.Symbols[j].Name)
			}
		}
	}
}

func TestInversion_HeadSizeZero(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0.d1", cat, nil)
	inv, err := Inversion(g, 0, nil)
	if err != nil {
		t.Fatalf("Inversion(headSize=0): %v", err)
	}
	if inv.KarvaString() != g.KarvaString() {
		t.Errorf("Inversion(headSize=0) changed symbols: got %q want %q", inv.KarvaString(), g.KarvaString())
	}
}

func TestInversion_HeadSizeOne(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0.d1", cat, nil)
	inv, err := Inversion(g, 1, nil)
	if err != nil {
		t.Fatalf("Inversion(headSize=1): %v", err)
	}
	if inv.KarvaString() != g.KarvaString() {
		t.Errorf("Inversion(headSize=1) changed symbols: got %q want %q", inv.KarvaString(), g.KarvaString())
	}
}

func TestInversion_Error_NegativeHeadSize(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0.d1", cat, nil)
	if _, err := Inversion(g, -1, nil); err == nil {
		t.Fatal("Inversion(headSize=-1): want error, got nil")
	}
}

// --- ISTransposition tests ---

func TestISTransposition_TailUnchanged(t *testing.T) {
	cat := newIntCatalog(t)
	headSize := 4
	// Build a gene large enough to have a meaningful tail.
	g, err := NewRandomGene(cat, headSize, 2, 0, rand.New(rand.NewSource(7)))
	if err != nil {
		t.Fatalf("NewRandomGene: %v", err)
	}
	tailStart := headSize
	rng := rand.New(rand.NewSource(7))
	for i := 0; i < 30; i++ {
		dst, err := ISTransposition(g, headSize, 3, rng)
		if err != nil {
			t.Fatalf("ISTransposition: %v", err)
		}
		for j := tailStart; j < len(g.Symbols); j++ {
			if dst.Symbols[j].Name != g.Symbols[j].Name {
				t.Errorf("tail[%d] changed after ISTransposition: %q vs %q",
					j, dst.Symbols[j].Name, g.Symbols[j].Name)
			}
		}
		if len(dst.Symbols) != len(g.Symbols) {
			t.Errorf("ISTransposition changed gene length: %d vs %d", len(dst.Symbols), len(g.Symbols))
		}
	}
}

func TestISTransposition_SmallHead(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0", cat, nil) // headSize=1, tail=[1]
	dst, err := ISTransposition(g, 1, 1, nil)
	if err != nil {
		t.Fatalf("ISTransposition(headSize=1): %v", err)
	}
	// Should return unchanged copy when headSize <= 1.
	if dst.KarvaString() != g.KarvaString() {
		t.Errorf("ISTransposition(headSize=1) changed gene: %q vs %q", dst.KarvaString(), g.KarvaString())
	}
}

func TestISTransposition_Errors(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0.d1", cat, nil)
	if _, err := ISTransposition(g, -1, 1, nil); err == nil {
		t.Fatal("ISTransposition(headSize=-1): want error, got nil")
	}
	if _, err := ISTransposition(g, 1, 0, nil); err == nil {
		t.Fatal("ISTransposition(maxISLen=0): want error, got nil")
	}
}

// --- RISTransposition tests ---

func TestRISTransposition_StartsWithFunction(t *testing.T) {
	cat := newIntCatalog(t)
	headSize := 5
	g, err := NewRandomGene(cat, headSize, 2, 0, rand.New(rand.NewSource(3)))
	if err != nil {
		t.Fatalf("NewRandomGene: %v", err)
	}
	// Force at least one function in the head by constructing a known gene.
	g2 := makeGene(t, "+.*.d0.d1.d0.d1.d0", cat, nil)
	rng := rand.New(rand.NewSource(99))
	for i := 0; i < 20; i++ {
		dst, err := RISTransposition(g2, headSize, 2, rng)
		if err != nil {
			t.Fatalf("RISTransposition: %v", err)
		}
		// Position 0 of the head must now be a function symbol.
		if len(dst.Symbols) > 0 && dst.Symbols[0].Kind != SymbolKindFunction {
			t.Errorf("RISTransposition: position 0 is not a function symbol: %+v", dst.Symbols[0])
		}
		// Tail unchanged.
		for j := headSize; j < len(g2.Symbols); j++ {
			if dst.Symbols[j].Name != g2.Symbols[j].Name {
				t.Errorf("tail[%d] changed after RISTransposition", j)
			}
		}
		_ = g
	}
}

func TestRISTransposition_NoFunctionInHead(t *testing.T) {
	// A gene whose head contains only terminals should be returned unchanged.
	cat := newIntCatalog(t)
	g := makeGene(t, "d0.d1.d0", cat, nil) // all terminals
	dst, err := RISTransposition(g, 2, 1, nil)
	if err != nil {
		t.Fatalf("RISTransposition: %v", err)
	}
	if dst.KarvaString() != g.KarvaString() {
		t.Errorf("RISTransposition changed all-terminal head: %q vs %q", dst.KarvaString(), g.KarvaString())
	}
}

func TestRISTransposition_Errors(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0.d1", cat, nil)
	if _, err := RISTransposition(g, -1, 1, nil); err == nil {
		t.Fatal("RISTransposition(headSize=-1): want error, got nil")
	}
	if _, err := RISTransposition(g, 1, 0, nil); err == nil {
		t.Fatal("RISTransposition(maxISLen=0): want error, got nil")
	}
}

// --- OnePointRecombine tests ---

func TestOnePointRecombine_ValidGenes(t *testing.T) {
	cat := newIntCatalog(t)
	g1 := makeGene(t, "+.d0.d1", cat, nil)
	g2 := makeGene(t, "*.d0.d1", cat, nil)
	rng := rand.New(rand.NewSource(11))
	for i := 0; i < 20; i++ {
		c1, c2, err := OnePointRecombine(g1, g2, rng)
		if err != nil {
			t.Fatalf("OnePointRecombine: %v", err)
		}
		if len(c1.Symbols) != len(g1.Symbols) {
			t.Errorf("c1 length changed: %d vs %d", len(c1.Symbols), len(g1.Symbols))
		}
		if len(c2.Symbols) != len(g2.Symbols) {
			t.Errorf("c2 length changed: %d vs %d", len(c2.Symbols), len(g2.Symbols))
		}
		// Combined symbol sets should equal the original sets combined.
		origNames := countNames(g1, g2)
		childNames := countNames(c1, c2)
		for k, v := range origNames {
			if childNames[k] != v {
				t.Errorf("symbol %q count changed after recombination: orig=%d child=%d", k, v, childNames[k])
			}
		}
	}
}

func TestOnePointRecombine_Errors(t *testing.T) {
	cat := newIntCatalog(t)
	empty := Gene[int]{}
	g := makeGene(t, "+.d0.d1", cat, nil)
	g3 := makeGene(t, "+.*.-.d0.d1.d0.d1", cat, nil)

	if _, _, err := OnePointRecombine(empty, g, nil); err == nil {
		t.Fatal("OnePointRecombine(empty, g): want error, got nil")
	}
	if _, _, err := OnePointRecombine(g, empty, nil); err == nil {
		t.Fatal("OnePointRecombine(g, empty): want error, got nil")
	}
	if _, _, err := OnePointRecombine(g, g3, nil); err == nil {
		t.Fatal("OnePointRecombine(different lengths): want error, got nil")
	}
}

// --- TwoPointRecombine tests ---

func TestTwoPointRecombine_ValidGenes(t *testing.T) {
	cat := newIntCatalog(t)
	g1 := makeGene(t, "+.*.-.d0.d1.d0.d1", cat, nil)
	g2 := makeGene(t, "-.+.*.d1.d0.d1.d0", cat, nil)
	rng := rand.New(rand.NewSource(77))
	for i := 0; i < 20; i++ {
		c1, c2, err := TwoPointRecombine(g1, g2, rng)
		if err != nil {
			t.Fatalf("TwoPointRecombine: %v", err)
		}
		if len(c1.Symbols) != len(g1.Symbols) {
			t.Errorf("c1 length changed: %d vs %d", len(c1.Symbols), len(g1.Symbols))
		}
		origNames := countNames(g1, g2)
		childNames := countNames(c1, c2)
		for k, v := range origNames {
			if childNames[k] != v {
				t.Errorf("symbol %q count after two-point recombination: orig=%d child=%d", k, v, childNames[k])
			}
		}
	}
}

func TestTwoPointRecombine_Errors(t *testing.T) {
	cat := newIntCatalog(t)
	empty := Gene[int]{}
	g := makeGene(t, "+.d0.d1", cat, nil)
	g3 := makeGene(t, "+.*.-.d0.d1.d0.d1", cat, nil)

	if _, _, err := TwoPointRecombine(empty, g, nil); err == nil {
		t.Fatal("TwoPointRecombine(empty, g): want error, got nil")
	}
	if _, _, err := TwoPointRecombine(g, empty, nil); err == nil {
		t.Fatal("TwoPointRecombine(g, empty): want error, got nil")
	}
	if _, _, err := TwoPointRecombine(g, g3, nil); err == nil {
		t.Fatal("TwoPointRecombine(different lengths): want error, got nil")
	}
}

// --- GeneTranspose tests ---

func TestGeneTranspose_SingleGene(t *testing.T) {
	cat := newIntCatalog(t)
	g := makeGene(t, "+.d0.d1", cat, nil)
	link, _ := NewLinkFunc[int]("id", func(v []int) int { return v[0] })
	genome := Genome[int]{Genes: []Gene[int]{g}, Link: link}
	dst, err := GeneTranspose(genome, nil)
	if err != nil {
		t.Fatalf("GeneTranspose(single gene): %v", err)
	}
	// Single gene → returned unchanged.
	if dst.KarvaString() != genome.KarvaString() {
		t.Errorf("GeneTranspose(single) changed genome: %q vs %q", dst.KarvaString(), genome.KarvaString())
	}
}

func TestGeneTranspose_MultiGene(t *testing.T) {
	cat := newIntCatalog(t)
	g0 := makeGene(t, "+.d0.d1", cat, nil)
	g1 := makeGene(t, "*.d0.d1", cat, nil)
	g2 := makeGene(t, "-.d0.d1", cat, nil)
	link, _ := NewLinkFunc[int]("+", func(v []int) int {
		sum := 0
		for _, v := range v {
			sum += v
		}
		return sum
	})
	genome := Genome[int]{Genes: []Gene[int]{g0, g1, g2}, Link: link}

	rng := rand.New(rand.NewSource(13))
	for i := 0; i < 20; i++ {
		dst, err := GeneTranspose(genome, rng)
		if err != nil {
			t.Fatalf("GeneTranspose: %v", err)
		}
		if len(dst.Genes) != len(genome.Genes) {
			t.Errorf("GeneTranspose changed gene count: %d vs %d", len(dst.Genes), len(genome.Genes))
		}
		// The first gene of dst must have been one of the non-first genes of original.
		firstKarva := dst.Genes[0].KarvaString()
		if firstKarva == g0.KarvaString() {
			t.Errorf("GeneTranspose: first gene is still the original first gene %q", firstKarva)
		}
		// All original karva strings must still be present.
		origSet := map[string]int{
			g0.KarvaString(): 1,
			g1.KarvaString(): 1,
			g2.KarvaString(): 1,
		}
		for _, gene := range dst.Genes {
			origSet[gene.KarvaString()]--
		}
		for k, v := range origSet {
			if v != 0 {
				t.Errorf("gene %q count wrong after GeneTranspose: excess=%d", k, -v)
			}
		}
	}
}

func TestGeneTranspose_Deterministic(t *testing.T) {
	cat := newIntCatalog(t)
	g0 := makeGene(t, "+.d0.d1", cat, nil)
	g1 := makeGene(t, "*.d0.d1", cat, nil)
	link, _ := NewLinkFunc[int]("+", func(v []int) int { return v[0] + v[1] })
	genome := Genome[int]{Genes: []Gene[int]{g0, g1}, Link: link}

	r1, _ := GeneTranspose(genome, rand.New(rand.NewSource(42)))
	r2, _ := GeneTranspose(genome, rand.New(rand.NewSource(42)))
	if r1.KarvaString() != r2.KarvaString() {
		t.Errorf("GeneTranspose not deterministic: %q vs %q", r1.KarvaString(), r2.KarvaString())
	}
}

// --- integration: mutated genome still evaluates ---

func TestPointMutate_EvalAfterMutation(t *testing.T) {
	cat := newIntCatalog(t)
	g, err := NewRandomGene(cat, 5, 2, 0, rand.New(rand.NewSource(55)))
	if err != nil {
		t.Fatalf("NewRandomGene: %v", err)
	}
	rng := rand.New(rand.NewSource(55))
	for i := 0; i < 30; i++ {
		mut, err := PointMutate(g, cat, 5, 2, 0, rng)
		if err != nil {
			t.Fatalf("PointMutate: %v", err)
		}
		if _, err := mut.Eval([]int{3, 7}); err != nil {
			t.Errorf("mutated gene failed to evaluate: %v", err)
		}
	}
}

// --- helpers ---

// countNames returns a map of symbol name → occurrence count across given genes.
func countNames[T any](genes ...Gene[T]) map[string]int {
	m := make(map[string]int)
	for _, g := range genes {
		for _, sym := range g.Symbols {
			m[sym.Name]++
		}
	}
	return m
}

// TestKarvaRoundtrip checks that KarvaString→ParseSymbols roundtrips correctly.
func TestKarvaRoundtrip(t *testing.T) {
	cat := newIntCatalog(t)
	karva := "+.*.-.d0.d1.d0.d1"
	orig := makeGene(t, karva, cat, nil)
	roundtripped := strings.Split(orig.KarvaString(), ".")
	syms, err := ParseSymbols(roundtripped, cat)
	if err != nil {
		t.Fatalf("ParseSymbols on roundtripped karva: %v", err)
	}
	if len(syms) != len(orig.Symbols) {
		t.Fatalf("roundtrip symbol count: %d vs %d", len(syms), len(orig.Symbols))
	}
	for i, sym := range syms {
		if sym.Name != orig.Symbols[i].Name {
			t.Errorf("roundtrip symbol[%d]: got %q want %q", i, sym.Name, orig.Symbols[i].Name)
		}
		if sym.Kind != orig.Symbols[i].Kind {
			t.Errorf("roundtrip symbol[%d] kind: got %v want %v", i, sym.Kind, orig.Symbols[i].Kind)
		}
	}
}
