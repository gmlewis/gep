// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package recombination provides configurable crossover operators for the typed
// evolution engine.
//
// Operators are applied pairwise to adjacent genomes in the population:
// [0,1], [2,3], ... . An odd last genome is deep-copied unchanged.
// Input genomes are never modified in place.
package recombination

import (
	"math/rand"

	"github.com/gmlewis/gep/v2/core"
)

// Config controls which recombination operators are applied and at what rates.
// Zero values disable the corresponding operators.
type Config struct {
	// OnePointRate is the probability [0, 1] that each aligned gene pair
	// undergoes one-point crossover.
	OnePointRate float64

	// TwoPointRate is the probability [0, 1] that each aligned gene pair
	// undergoes two-point crossover.
	TwoPointRate float64
}

// Apply applies configured recombination operators and returns a new population.
//
// Operators are applied per aligned gene pair in this order:
//  1. One-point recombination
//  2. Two-point recombination
//
// rng may be nil; when nil the global math/rand source is used.
func Apply[T any](genomes []core.Genome[T], cfg Config, rng *rand.Rand) []core.Genome[T] {
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
	for i := 0; i < len(genomes); i += 2 {
		if i+1 >= len(genomes) {
			result[i] = genomes[i].Dup()
			continue
		}

		g1 := genomes[i].Dup()
		g2 := genomes[i+1].Dup()

		geneCount := len(g1.Genes)
		if len(g2.Genes) < geneCount {
			geneCount = len(g2.Genes)
		}
		for j := 0; j < geneCount; j++ {
			if cfg.OnePointRate > 0 && randFloat64() < cfg.OnePointRate {
				if c1, c2, err := core.OnePointRecombine(g1.Genes[j], g2.Genes[j], rng); err == nil {
					g1.Genes[j], g2.Genes[j] = c1, c2
				}
			}
			if cfg.TwoPointRate > 0 && randFloat64() < cfg.TwoPointRate {
				if c1, c2, err := core.TwoPointRecombine(g1.Genes[j], g2.Genes[j], rng); err == nil {
					g1.Genes[j], g2.Genes[j] = c1, c2
				}
			}
		}

		result[i] = g1
		result[i+1] = g2
	}

	return result
}
