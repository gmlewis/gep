// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package transposition provides configurable genetic transposition operators
// for the typed evolution engine.
//
// Each operator is applied probabilistically according to the rates in Config.
// Gene-level operators (IS transposition and RIS transposition) are applied
// independently to every gene in every genome. Gene transposition is applied to
// the genome as a whole. All operators produce deep copies; the input genomes
// are never modified in place.
package transposition

import (
	"math/rand"

	"github.com/gmlewis/gep/v2/core"
)

// Config controls which transposition operators are applied and at what rates.
// Zero values disable the corresponding operators.
type Config struct {
	// HeadSize is the length of each gene's head region. This is required by
	// operators that need the head/tail boundary (IS transposition and RIS
	// transposition).
	HeadSize int

	// ISTranspositionRate is the probability [0, 1] that each gene undergoes
	// IS-element transposition. Zero disables the operator.
	ISTranspositionRate float64

	// MaxISLen is the maximum IS element length used by IS transposition.
	// Defaults to 1 when zero.
	MaxISLen int

	// RISTranspositionRate is the probability [0, 1] that each gene undergoes
	// RIS-element transposition. Zero disables the operator.
	RISTranspositionRate float64

	// MaxRISLen is the maximum RIS element length used by RIS transposition.
	// Defaults to 1 when zero.
	MaxRISLen int

	// GeneTranspositionRate is the probability [0, 1] that each genome
	// undergoes gene transposition. Zero disables the operator.
	GeneTranspositionRate float64
}

// Apply applies the configured transposition operators to each genome and
// returns a new population of transposed genomes. The returned slice is
// index-aligned with the input slice.
//
// Operators are applied in the following order per genome:
//  1. IS transposition (per gene)
//  2. RIS transposition (per gene)
//  3. Gene transposition (per genome)
//
// rng may be nil; when nil the global math/rand source is used.
func Apply[T any](genomes []core.Genome[T], cfg Config, rng *rand.Rand) []core.Genome[T] {
	if len(genomes) == 0 {
		return nil
	}

	maxISLen := cfg.MaxISLen
	if maxISLen < 1 {
		maxISLen = 1
	}
	maxRISLen := cfg.MaxRISLen
	if maxRISLen < 1 {
		maxRISLen = 1
	}

	randFloat64 := func() float64 {
		if rng != nil {
			return rng.Float64()
		}
		return rand.Float64() //nolint:gosec
	}

	result := make([]core.Genome[T], len(genomes))
	for i, g := range genomes {
		current := g.Dup()

		for j := range current.Genes {
			gene := current.Genes[j]

			if cfg.ISTranspositionRate > 0 && randFloat64() < cfg.ISTranspositionRate {
				if t, err := core.ISTransposition(gene, cfg.HeadSize, maxISLen, rng); err == nil {
					gene = t
				}
			}

			if cfg.RISTranspositionRate > 0 && randFloat64() < cfg.RISTranspositionRate {
				if t, err := core.RISTransposition(gene, cfg.HeadSize, maxRISLen, rng); err == nil {
					gene = t
				}
			}

			current.Genes[j] = gene
		}

		if cfg.GeneTranspositionRate > 0 && randFloat64() < cfg.GeneTranspositionRate {
			if t, err := core.GeneTranspose(current, rng); err == nil {
				current = t
			}
		}

		result[i] = current
	}
	return result
}
