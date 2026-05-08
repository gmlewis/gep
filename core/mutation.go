// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"math/rand"
)

// randIntn returns a random int in [0, n) using rng when non-nil, else the
// global math/rand source.
func randIntn(n int, rng *rand.Rand) int {
	if rng != nil {
		return rng.Intn(n)
	}
	return rand.Intn(n) //nolint:gosec
}

// buildTermChoices builds the terminal symbol slice (d0..d(numTerminals-1)
// followed by c0..c(numConstants-1)) used by mutation operators.
func buildTermChoices[T any](numTerminals, numConstants int) []Symbol[T] {
	choices := make([]Symbol[T], 0, numTerminals+numConstants)
	for i := 0; i < numTerminals; i++ {
		sym, _ := NewTerminalSymbol[T](i)
		choices = append(choices, sym)
	}
	for i := 0; i < numConstants; i++ {
		sym, _ := NewConstantSymbol[T](i)
		choices = append(choices, sym)
	}
	return choices
}

// PointMutate performs a single point mutation on a copy of g.
//
// A random position is chosen.  Positions in the head ([0, headSize)) may be
// replaced with any symbol from the catalog or with a terminal/constant.
// Positions in the tail ([headSize, len)) are replaced with a terminal or
// constant only, preserving the Karva syntactic validity guarantee.
//
// headSize must be >= 0.  numTerminals and numConstants must be >= 0 and their
// sum must be > 0.  cat must be non-nil.  rng may be nil (global source used).
func PointMutate[T any](g Gene[T], cat *Catalog[T], headSize, numTerminals, numConstants int, rng *rand.Rand) (Gene[T], error) {
	if cat == nil {
		return Gene[T]{}, fmt.Errorf("core.PointMutate: catalog cannot be nil")
	}
	if headSize < 0 {
		return Gene[T]{}, fmt.Errorf("core.PointMutate: headSize must be >= 0")
	}
	if numTerminals < 0 {
		return Gene[T]{}, fmt.Errorf("core.PointMutate: numTerminals must be >= 0")
	}
	if numConstants < 0 {
		return Gene[T]{}, fmt.Errorf("core.PointMutate: numConstants must be >= 0")
	}
	if numTerminals+numConstants == 0 {
		return Gene[T]{}, fmt.Errorf("core.PointMutate: numTerminals+numConstants must be > 0")
	}
	if len(g.Symbols) == 0 {
		return Gene[T]{}, fmt.Errorf("core.PointMutate: gene has no symbols")
	}

	termChoices := buildTermChoices[T](numTerminals, numConstants)
	funcNames := cat.Symbols()

	dst := g.Dup()
	pos := randIntn(len(dst.Symbols), rng)

	if pos < headSize {
		// Head: any function or terminal
		headChoiceCount := len(funcNames) + len(termChoices)
		if headChoiceCount == 0 {
			return dst, nil
		}
		choice := randIntn(headChoiceCount, rng)
		if choice < len(funcNames) {
			node, _ := cat.Lookup(funcNames[choice])
			sym, _ := NewFunctionSymbol[T](node)
			dst.Symbols[pos] = sym
		} else {
			dst.Symbols[pos] = termChoices[choice-len(funcNames)]
		}
	} else {
		// Tail: terminals only
		dst.Symbols[pos] = termChoices[randIntn(len(termChoices), rng)]
	}
	return dst, nil
}

// Inversion reverses a random contiguous subsequence within the head of a copy
// of g.  If headSize <= 1 or the gene has no symbols, the gene is returned
// unchanged (deep-copied).
//
// headSize must be >= 0.  rng may be nil.
func Inversion[T any](g Gene[T], headSize int, rng *rand.Rand) (Gene[T], error) {
	if headSize < 0 {
		return Gene[T]{}, fmt.Errorf("core.Inversion: headSize must be >= 0")
	}
	dst := g.Dup()
	if headSize <= 1 || len(dst.Symbols) == 0 {
		return dst, nil
	}
	// Clamp headSize to actual symbol count.
	hs := headSize
	if hs > len(dst.Symbols) {
		hs = len(dst.Symbols)
	}
	start := randIntn(hs, rng)
	end := randIntn(hs, rng)
	if start == end {
		return dst, nil
	}
	if start > end {
		start, end = end, start
	}
	for start < end {
		dst.Symbols[start], dst.Symbols[end] = dst.Symbols[end], dst.Symbols[start]
		start++
		end--
	}
	return dst, nil
}

// ISTransposition performs IS (Insertion Sequence) element transposition on a
// copy of g.
//
// In GEP, an IS element is an arbitrary subsequence of symbols taken from
// anywhere in the gene.  The transposition operator extracts a random IS element
// of length in [1, maxISLen] and inserts it at a random position in
// [1, headSize), shifting the existing head symbols to the right; any symbols
// pushed past the head boundary are discarded.  The tail is never modified.
//
// headSize must be >= 0.  maxISLen must be >= 1.  rng may be nil.
// If the gene is too short for IS transposition (headSize <= 1 or no symbols),
// the gene is returned unchanged (deep-copied).
func ISTransposition[T any](g Gene[T], headSize, maxISLen int, rng *rand.Rand) (Gene[T], error) {
	if headSize < 0 {
		return Gene[T]{}, fmt.Errorf("core.ISTransposition: headSize must be >= 0")
	}
	if maxISLen < 1 {
		return Gene[T]{}, fmt.Errorf("core.ISTransposition: maxISLen must be >= 1")
	}
	dst := g.Dup()
	if headSize <= 1 || len(dst.Symbols) == 0 {
		return dst, nil
	}

	// Pick IS (Insertion Sequence) element.
	start := randIntn(len(dst.Symbols), rng)
	length := 1 + randIntn(maxISLen, rng)
	if start+length > len(dst.Symbols) {
		length = len(dst.Symbols) - start
	}
	isElem := make([]Symbol[T], length)
	copy(isElem, dst.Symbols[start:start+length])

	// Pick insertion position in [1, headSize).
	insPos := 1 + randIntn(headSize-1, rng)

	// Build new head: old[0:insPos] + isElem + old[insPos:] truncated to headSize.
	hs := headSize
	if hs > len(dst.Symbols) {
		hs = len(dst.Symbols)
	}
	oldHead := make([]Symbol[T], hs)
	copy(oldHead, dst.Symbols[:hs])

	newHead := make([]Symbol[T], hs)
	copy(newHead[:insPos], oldHead[:insPos])
	oldSrc := insPos
	newDst := insPos
	for i := 0; i < length && newDst < hs; i++ {
		newHead[newDst] = isElem[i]
		newDst++
	}
	for newDst < hs && oldSrc < hs {
		newHead[newDst] = oldHead[oldSrc]
		oldSrc++
		newDst++
	}
	copy(dst.Symbols[:hs], newHead)
	// Tail positions (>= headSize) are unchanged.
	return dst, nil
}

// RISTransposition performs RIS (Root Insertion Sequence) element transposition
// on a copy of g.
//
// In GEP, an RIS element is a subsequence that must begin with a function
// symbol.  Unlike IS transposition, the RIS element is always inserted at
// position 0 — the root — of the gene head, which guarantees that the head
// continues to start with a function.  Existing head symbols are shifted to the
// right; symbols pushed past the head boundary are discarded.  The tail is
// never modified.
//
// The operator scans the head for function symbols and picks one at random as
// the starting point for the RIS element.  If no function symbol exists in the
// head, or headSize == 0, the gene is returned unchanged (deep-copied).
//
// headSize must be >= 0.  maxISLen must be >= 1.  rng may be nil.
func RISTransposition[T any](g Gene[T], headSize, maxISLen int, rng *rand.Rand) (Gene[T], error) {
	if headSize < 0 {
		return Gene[T]{}, fmt.Errorf("core.RISTransposition: headSize must be >= 0")
	}
	if maxISLen < 1 {
		return Gene[T]{}, fmt.Errorf("core.RISTransposition: maxISLen must be >= 1")
	}
	dst := g.Dup()
	if headSize == 0 || len(dst.Symbols) == 0 {
		return dst, nil
	}

	// Collect function-symbol positions within the head.
	hs := headSize
	if hs > len(dst.Symbols) {
		hs = len(dst.Symbols)
	}
	var funcPositions []int
	for i := 0; i < hs; i++ {
		if dst.Symbols[i].Kind == SymbolKindFunction {
			funcPositions = append(funcPositions, i)
		}
	}
	if len(funcPositions) == 0 {
		return dst, nil
	}

	// Pick a random starting function position.
	start := funcPositions[randIntn(len(funcPositions), rng)]
	length := 1 + randIntn(maxISLen, rng)
	if start+length > len(dst.Symbols) {
		length = len(dst.Symbols) - start
	}
	isElem := make([]Symbol[T], length)
	copy(isElem, dst.Symbols[start:start+length])

	// Build new head: isElem at [0..length-1] then old head, truncated to hs.
	oldHead := make([]Symbol[T], hs)
	copy(oldHead, dst.Symbols[:hs])

	newHead := make([]Symbol[T], hs)
	newDst := 0
	for i := 0; i < length && newDst < hs; i++ {
		newHead[newDst] = isElem[i]
		newDst++
	}
	oldSrc := 0
	for newDst < hs && oldSrc < hs {
		newHead[newDst] = oldHead[oldSrc]
		oldSrc++
		newDst++
	}
	copy(dst.Symbols[:hs], newHead)
	return dst, nil
}

// OnePointRecombine performs one-point recombination (crossover) on copies of
// g1 and g2.  A random crossover point in [0, len) is chosen.  Symbols from
// the crossover point onwards are swapped between the two children.
//
// g1 and g2 must have the same number of symbols and must be non-empty.
// rng may be nil.
func OnePointRecombine[T any](g1, g2 Gene[T], rng *rand.Rand) (Gene[T], Gene[T], error) {
	if len(g1.Symbols) == 0 || len(g2.Symbols) == 0 {
		return Gene[T]{}, Gene[T]{}, fmt.Errorf("core.OnePointRecombine: genes must be non-empty")
	}
	if len(g1.Symbols) != len(g2.Symbols) {
		return Gene[T]{}, Gene[T]{}, fmt.Errorf("core.OnePointRecombine: genes must have equal symbol length (%d vs %d)",
			len(g1.Symbols), len(g2.Symbols))
	}
	point := randIntn(len(g1.Symbols), rng)
	c1 := g1.Dup()
	c2 := g2.Dup()
	for i := point; i < len(c1.Symbols); i++ {
		c1.Symbols[i], c2.Symbols[i] = c2.Symbols[i], c1.Symbols[i]
	}
	return c1, c2, nil
}

// TwoPointRecombine performs two-point recombination (crossover) on copies of
// g1 and g2.  Two random crossover points p1 <= p2 in [0, len) are chosen.
// Symbols in the range [p1, p2] are swapped between the two children.
//
// g1 and g2 must have the same number of symbols and must be non-empty.
// rng may be nil.
func TwoPointRecombine[T any](g1, g2 Gene[T], rng *rand.Rand) (Gene[T], Gene[T], error) {
	if len(g1.Symbols) == 0 || len(g2.Symbols) == 0 {
		return Gene[T]{}, Gene[T]{}, fmt.Errorf("core.TwoPointRecombine: genes must be non-empty")
	}
	if len(g1.Symbols) != len(g2.Symbols) {
		return Gene[T]{}, Gene[T]{}, fmt.Errorf("core.TwoPointRecombine: genes must have equal symbol length (%d vs %d)",
			len(g1.Symbols), len(g2.Symbols))
	}
	p1 := randIntn(len(g1.Symbols), rng)
	p2 := randIntn(len(g1.Symbols), rng)
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	c1 := g1.Dup()
	c2 := g2.Dup()
	for i := p1; i <= p2; i++ {
		c1.Symbols[i], c2.Symbols[i] = c2.Symbols[i], c1.Symbols[i]
	}
	return c1, c2, nil
}

// GeneTranspose moves a randomly chosen non-first gene to position 0 of a copy
// of genome g, shifting the other genes one position to the right.
//
// If the genome has fewer than two genes it is returned unchanged (deep-copied).
// rng may be nil.
func GeneTranspose[T any](g Genome[T], rng *rand.Rand) (Genome[T], error) {
	if len(g.Genes) <= 1 {
		return g.Dup(), nil
	}
	// Pick a random gene index in [1, len(g.Genes)).
	idx := 1 + randIntn(len(g.Genes)-1, rng)
	dst := g.Dup()
	selected := dst.Genes[idx]
	// Shift genes [0, idx) one position to the right.
	copy(dst.Genes[1:idx+1], dst.Genes[0:idx])
	dst.Genes[0] = selected
	return dst, nil
}
