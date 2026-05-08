// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package selection provides configurable population-selection helpers for the
// typed evolution engine.
package selection

import (
	"math"
	"math/rand"

	"github.com/gmlewis/gep/v2/core"
)

// Candidate is a scored genome used by selection operators.
type Candidate[T any] struct {
	Genome core.Genome[T]
	Score  float64
}

// Config controls selection behavior.
type Config struct {
	// MinimizeScore controls score orientation. When false (default), higher
	// scores are better. When true, lower scores are better.
	MinimizeScore bool
}

// Roulette selects a new population using roulette-wheel
// (fitness-proportionate) selection.
//
// The result has one selected candidate per input candidate and is index
// independent (draw order does not preserve source indices).
// Returned candidates are deep copies of selected genomes.
func Roulette[T any](population []Candidate[T], cfg Config, rng *rand.Rand) []Candidate[T] {
	if len(population) == 0 {
		return nil
	}

	effectiveScore := func(score float64) float64 {
		if cfg.MinimizeScore {
			return -score
		}
		return score
	}

	// Map effective scores into [0.1, 1.0] so all weights are strictly positive.
	minEff := math.Inf(1)
	maxEff := math.Inf(-1)
	for _, ind := range population {
		e := effectiveScore(ind.Score)
		if e < minEff {
			minEff = e
		}
		if e > maxEff {
			maxEff = e
		}
	}
	scale := maxEff - minEff
	if scale <= 0 {
		scale = 1.0
	}
	scaledScore := func(score float64) float64 {
		return 0.1 + (effectiveScore(score)-minEff)/scale
	}

	randIntn := func(n int) int {
		if rng != nil {
			return rng.Intn(n)
		}
		return rand.Intn(n) //nolint:gosec
	}
	randFloat64 := func() float64 {
		if rng != nil {
			return rng.Float64()
		}
		return rand.Float64() //nolint:gosec
	}

	result := make([]Candidate[T], 0, len(population))
	idx := randIntn(len(population))
	beta := 0.0
	for range population {
		beta += randFloat64() * 2.0
		s := scaledScore(population[idx].Score)
		for beta > s {
			beta -= s
			idx = (idx + 1) % len(population)
			s = scaledScore(population[idx].Score)
		}
		result = append(result, Candidate[T]{
			Genome: population[idx].Genome.Dup(),
			Score:  population[idx].Score,
		})
	}
	return result
}
