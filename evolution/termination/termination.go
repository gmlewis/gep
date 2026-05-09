// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package termination provides configurable stopping criteria for the typed
// evolution engine.
//
// Each Criterion is evaluated after every generation. When any criterion
// returns true from ShouldStop, Evolve halts and returns the current best
// individual.
//
// Built-in criteria:
//
//   - ScoreThreshold – stop when the best score crosses a threshold.
//   - NoImprovement – stop when the best score has not improved for N
//     consecutive generations.
//   - Any – composite criterion that stops when at least one sub-criterion fires.
//   - All – composite criterion that stops only when every sub-criterion fires.
//
// Typical usage:
//
//	gen.TerminationCriteria = []termination.Criterion{
//	    termination.ScoreThreshold(0.99, false),
//	    termination.NoImprovement(50, false),
//	}
//	result := gen.Evolve(1000)
package termination

// Criterion is a predicate evaluated after each generation to decide whether
// evolution should stop.
//
// generation is the current zero-based generation index.
// bestScore is the best raw fitness score observed in the current population
// (before any minimize/maximize adjustment).
//
// Implementations may be stateful (e.g. NoImprovement tracks previous bests).
// Each Criterion is owned by one Generation; do not share a single Criterion
// instance across concurrent Generation.Evolve calls.
type Criterion interface {
	// ShouldStop returns true when the stopping condition has been met.
	ShouldStop(generation int, bestScore float64) bool
}

// ─── ScoreThreshold ────────────────────────────────────────────────────────

type scoreThresholdCriterion struct {
	threshold float64
	minimize  bool
}

// ScoreThreshold returns a Criterion that stops when the best score crosses
// the given threshold.
//
// When minimize is false (the default for maximization), evolution stops when
// bestScore >= threshold.
// When minimize is true, evolution stops when bestScore <= threshold.
func ScoreThreshold(threshold float64, minimize bool) Criterion {
	return &scoreThresholdCriterion{threshold: threshold, minimize: minimize}
}

func (c *scoreThresholdCriterion) ShouldStop(_ int, bestScore float64) bool {
	if c.minimize {
		return bestScore <= c.threshold
	}
	return bestScore >= c.threshold
}

// ─── NoImprovement ─────────────────────────────────────────────────────────

type noImprovementCriterion struct {
	patience    int
	minimize    bool
	bestSeen    float64
	stagnantFor int
	initialized bool
}

// NoImprovement returns a Criterion that stops when the best score has not
// strictly improved for at least patience consecutive generations.
//
// When minimize is false (maximization), an improvement means the new best
// score is strictly greater than the previously recorded best. When minimize
// is true, an improvement means the new score is strictly less.
//
// patience must be >= 1; if a value < 1 is supplied it is treated as 1.
func NoImprovement(patience int, minimize bool) Criterion {
	if patience < 1 {
		patience = 1
	}
	return &noImprovementCriterion{patience: patience, minimize: minimize}
}

func (c *noImprovementCriterion) ShouldStop(_ int, bestScore float64) bool {
	if !c.initialized {
		c.bestSeen = bestScore
		c.initialized = true
		c.stagnantFor = 0
		return false
	}

	improved := false
	if c.minimize {
		improved = bestScore < c.bestSeen
	} else {
		improved = bestScore > c.bestSeen
	}

	if improved {
		c.bestSeen = bestScore
		c.stagnantFor = 0
		return false
	}

	c.stagnantFor++
	return c.stagnantFor >= c.patience
}

// ─── Any ───────────────────────────────────────────────────────────────────

type anyCriterion struct {
	criteria []Criterion
}

// Any returns a composite Criterion that stops when at least one of the given
// sub-criteria returns true (logical OR). If criteria is empty, Any never
// signals a stop.
func Any(criteria ...Criterion) Criterion {
	return &anyCriterion{criteria: criteria}
}

func (c *anyCriterion) ShouldStop(generation int, bestScore float64) bool {
	for _, cr := range c.criteria {
		if cr.ShouldStop(generation, bestScore) {
			return true
		}
	}
	return false
}

// ─── All ───────────────────────────────────────────────────────────────────

type allCriterion struct {
	criteria []Criterion
}

// All returns a composite Criterion that stops only when every one of the
// given sub-criteria returns true (logical AND). If criteria is empty, All
// always signals a stop.
func All(criteria ...Criterion) Criterion {
	return &allCriterion{criteria: criteria}
}

func (c *allCriterion) ShouldStop(generation int, bestScore float64) bool {
	for _, cr := range c.criteria {
		if !cr.ShouldStop(generation, bestScore) {
			return false
		}
	}
	return true
}
