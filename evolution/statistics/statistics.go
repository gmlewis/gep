// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package statistics provides per-generation statistics collection for the
// typed evolution engine.
//
// Typical usage:
//
//	col := &statistics.Collector{}
//	gen.Statistics = col
//	result := gen.Evolve(500)
//	for _, s := range col.History {
//	    fmt.Printf("gen %d: best=%.4f mean=%.4f diversity=%.2f\n",
//	        s.Generation, s.BestScore, s.MeanScore, s.Diversity)
//	}
package statistics

// Stats holds aggregated statistics for a single generation.
type Stats struct {
	// Generation is the zero-based index of this generation within the Evolve
	// loop.
	Generation int

	// BestScore is the best score observed in the population for this
	// generation. "Best" respects the Minimize flag: when Minimize is false
	// (default) it is the maximum; when Minimize is true it is the minimum.
	BestScore float64

	// WorstScore is the worst score observed in the population. "Worst" is the
	// complement of "best": when Minimize is false it is the minimum; when
	// Minimize is true it is the maximum.
	WorstScore float64

	// MeanScore is the arithmetic mean of all individual scores.
	MeanScore float64

	// Diversity is the fraction of individuals with a unique Karva-string
	// representation, in the range [0, 1]. A value of 1.0 means every
	// individual has a distinct representation; 0.0 means all individuals are
	// identical (or no Karva strings were provided).
	Diversity float64
}

// Record computes per-generation statistics from the given scores and optional
// Karva-string representations.
//
// generation is the zero-based generation index.
// scores must be non-empty; all scores contribute to BestScore, WorstScore,
// and MeanScore.
// karvaStrings may be nil or empty, in which case Diversity is 0.
// minimize controls whether "best" means lowest (true) or highest (false).
//
// Record does not modify any internal state; it returns the computed Stats
// directly. Use Collector.Append to accumulate results over time.
func Record(generation int, scores []float64, karvaStrings []string, minimize bool) Stats {
	if len(scores) == 0 {
		return Stats{Generation: generation}
	}

	best := scores[0]
	worst := scores[0]
	sum := scores[0]
	for _, s := range scores[1:] {
		sum += s
		if minimize {
			if s < best {
				best = s
			}
			if s > worst {
				worst = s
			}
		} else {
			if s > best {
				best = s
			}
			if s < worst {
				worst = s
			}
		}
	}
	mean := sum / float64(len(scores))

	diversity := 0.0
	if len(karvaStrings) > 0 {
		seen := make(map[string]struct{}, len(karvaStrings))
		for _, k := range karvaStrings {
			seen[k] = struct{}{}
		}
		diversity = float64(len(seen)) / float64(len(karvaStrings))
	}

	return Stats{
		Generation: generation,
		BestScore:  best,
		WorstScore: worst,
		MeanScore:  mean,
		Diversity:  diversity,
	}
}

// Collector accumulates per-generation statistics produced by Record.
type Collector struct {
	// History holds one Stats entry per recorded generation, in the order they
	// were appended.
	History []Stats
}

// Append adds s to the collector's History.
func (c *Collector) Append(s Stats) {
	c.History = append(c.History, s)
}

// Last returns the most recently appended Stats. If History is empty, Last
// returns the zero Stats value.
func (c *Collector) Last() Stats {
	if len(c.History) == 0 {
		return Stats{}
	}
	return c.History[len(c.History)-1]
}
