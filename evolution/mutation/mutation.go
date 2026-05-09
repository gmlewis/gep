// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package mutation provides configurable genetic mutation operators for the
// typed evolution engine.
//
// Each operator is applied probabilistically according to the rates in Config.
// Gene-level operators (point mutation and inversion) are applied
// independently to every gene in every genome. All operators produce deep
// copies; the input genomes are never modified in place.
package mutation

import (
	"math/rand"

	"github.com/gmlewis/gep/v2/core"
)

// Config controls which genetic operators are applied and at what rates.
// Zero values disable the corresponding operators.
type Config struct {
	// HeadSize is the length of each gene's head region. This is required by
	// operators that need the head/tail boundary (inversion and point
	// mutation).
	HeadSize int

	// NumTerminals is the number of terminal symbols available. Required by
	// PointMutate (NumTerminals + NumConstants must be > 0).
	NumTerminals int

	// NumConstants is the number of constant symbols available. Required by
	// PointMutate (NumTerminals + NumConstants must be > 0).
	NumConstants int

	// PointMutationRate is the probability [0, 1] that each gene undergoes
	// one point mutation. Zero disables the operator.
	PointMutationRate float64

	// InversionRate is the probability [0, 1] that each gene undergoes head
	// inversion. Zero disables the operator.
	InversionRate float64
}

// Apply applies the configured mutation operators to each genome and returns a
// new population of mutated genomes.  The returned slice is index-aligned with
// the input slice.
//
// Operators are applied in the following order per genome:
//  1. Point mutation  (per gene)
//  2. Inversion       (per gene)
//
// cat must be non-nil when cfg.PointMutationRate > 0.
// rng may be nil; when nil the global math/rand source is used.
func Apply[T any](genomes []core.Genome[T], cat *core.Catalog[T], cfg Config, rng *rand.Rand) []core.Genome[T] {
	if len(genomes) == 0 {
		return nil
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

		// Gene-level operators: applied independently to each gene.
		for j := range current.Genes {
			gene := current.Genes[j]

			// 1. Point mutation
			if cfg.PointMutationRate > 0 && cat != nil && randFloat64() < cfg.PointMutationRate {
				if mut, err := core.PointMutate(gene, cat, cfg.HeadSize, cfg.NumTerminals, cfg.NumConstants, rng); err == nil {
					gene = mut
				}
			}

			// 2. Inversion
			if cfg.InversionRate > 0 && randFloat64() < cfg.InversionRate {
				if inv, err := core.Inversion(gene, cfg.HeadSize, rng); err == nil {
					gene = inv
				}
			}

			current.Genes[j] = gene
		}

		result[i] = current
	}
	return result
}
