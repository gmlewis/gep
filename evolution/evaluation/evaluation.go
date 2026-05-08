// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package evaluation provides configurable genome scoring helpers for the
// typed evolution engine.
package evaluation

import (
	"sync"

	"github.com/gmlewis/gep/v2/core"
)

// ScoringFunc computes a fitness score for a genome.
type ScoringFunc[T any] func(core.Genome[T]) float64

// Config controls the evaluation worker pool.
//
// Workers <= 0 means "auto", which uses one worker per genome.
type Config struct {
	Workers int
}

// ScoreAll scores all genomes and returns one score per input genome.
//
// Results are index-aligned with the input slice. The function is safe for
// concurrent execution because each worker writes to a distinct result index.
// If sf is nil, all returned scores are zero values.
func ScoreAll[T any](genomes []core.Genome[T], sf ScoringFunc[T], cfg Config) []float64 {
	scores := make([]float64, len(genomes))
	if len(genomes) == 0 || sf == nil {
		return scores
	}

	workers := cfg.Workers
	if workers <= 0 || workers > len(genomes) {
		workers = len(genomes)
	}

	jobs := make(chan int, len(genomes))
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				scores[idx] = sf(genomes[idx])
			}
		}()
	}

	for i := range genomes {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return scores
}
