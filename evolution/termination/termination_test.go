// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package termination

import (
	"testing"
)

// ─── ScoreThreshold tests ──────────────────────────────────────────────────

func TestScoreThreshold_Maximize_BelowThreshold(t *testing.T) {
	c := ScoreThreshold(100.0, false)
	if c.ShouldStop(0, 99.99) {
		t.Fatal("ShouldStop(99.99) with threshold=100 (maximize): got true, want false")
	}
}

func TestScoreThreshold_Maximize_AtThreshold(t *testing.T) {
	c := ScoreThreshold(100.0, false)
	if !c.ShouldStop(0, 100.0) {
		t.Fatal("ShouldStop(100.0) with threshold=100 (maximize): got false, want true")
	}
}

func TestScoreThreshold_Maximize_AboveThreshold(t *testing.T) {
	c := ScoreThreshold(100.0, false)
	if !c.ShouldStop(0, 200.0) {
		t.Fatal("ShouldStop(200.0) with threshold=100 (maximize): got false, want true")
	}
}

func TestScoreThreshold_Minimize_AboveThreshold(t *testing.T) {
	c := ScoreThreshold(0.01, true)
	if c.ShouldStop(0, 0.5) {
		t.Fatal("ShouldStop(0.5) with threshold=0.01 (minimize): got true, want false")
	}
}

func TestScoreThreshold_Minimize_AtThreshold(t *testing.T) {
	c := ScoreThreshold(0.01, true)
	if !c.ShouldStop(0, 0.01) {
		t.Fatal("ShouldStop(0.01) with threshold=0.01 (minimize): got false, want true")
	}
}

func TestScoreThreshold_Minimize_BelowThreshold(t *testing.T) {
	c := ScoreThreshold(0.01, true)
	if !c.ShouldStop(0, 0.005) {
		t.Fatal("ShouldStop(0.005) with threshold=0.01 (minimize): got false, want true")
	}
}

func TestScoreThreshold_GenerationArgIgnored(t *testing.T) {
	c := ScoreThreshold(50.0, false)
	// generation index must have no effect on the result
	for _, gen := range []int{0, 1, 100, 9999} {
		if c.ShouldStop(gen, 60.0) != true {
			t.Errorf("ShouldStop(gen=%d, 60) with threshold=50: got false, want true", gen)
		}
	}
}

// ─── NoImprovement tests ───────────────────────────────────────────────────

func TestNoImprovement_FirstCallNeverStops(t *testing.T) {
	c := NoImprovement(1, false)
	if c.ShouldStop(0, 10.0) {
		t.Fatal("NoImprovement: first call must not stop")
	}
}

func TestNoImprovement_StopsAfterPatience_Maximize(t *testing.T) {
	const patience = 3
	c := NoImprovement(patience, false)
	score := 10.0
	c.ShouldStop(0, score) // initialize
	// patience consecutive non-improvements should trigger stop
	for i := 1; i <= patience; i++ {
		got := c.ShouldStop(i, score) // same score – no improvement
		wantStop := i == patience
		if got != wantStop {
			t.Errorf("ShouldStop at gen %d (stagnant=%d): got %v, want %v", i, i, got, wantStop)
		}
	}
}

func TestNoImprovement_StopsAfterPatience_Minimize(t *testing.T) {
	const patience = 2
	c := NoImprovement(patience, true)
	score := 5.0
	c.ShouldStop(0, score) // initialize
	for i := 1; i <= patience; i++ {
		got := c.ShouldStop(i, score) // same score – no improvement
		wantStop := i == patience
		if got != wantStop {
			t.Errorf("ShouldStop at gen %d: got %v, want %v", i, got, wantStop)
		}
	}
}

func TestNoImprovement_ImprovementResetsCounter_Maximize(t *testing.T) {
	const patience = 3
	c := NoImprovement(patience, false)
	c.ShouldStop(0, 10.0) // initialize
	c.ShouldStop(1, 10.0) // stagnant 1
	c.ShouldStop(2, 10.0) // stagnant 2
	// improvement: counter must reset
	if c.ShouldStop(3, 20.0) {
		t.Fatal("improvement should reset stagnation counter")
	}
	// Two more stagnant gens (need patience=3 to stop again)
	c.ShouldStop(4, 20.0)      // stagnant 1
	c.ShouldStop(5, 20.0)      // stagnant 2
	if c.ShouldStop(6, 20.0) { // stagnant 3 – should stop now
		// ok
	} else {
		t.Fatal("expected stop after patience=3 stagnant gens post-improvement")
	}
}

func TestNoImprovement_ImprovementResetsCounter_Minimize(t *testing.T) {
	const patience = 2
	c := NoImprovement(patience, true)
	c.ShouldStop(0, 10.0) // initialize
	c.ShouldStop(1, 10.0) // stagnant 1
	// improvement (lower is better)
	if c.ShouldStop(2, 5.0) {
		t.Fatal("improvement should reset stagnation counter")
	}
	c.ShouldStop(3, 5.0) // stagnant 1
	if !c.ShouldStop(4, 5.0) {
		t.Fatal("expected stop after patience=2 stagnant gens post-improvement")
	}
}

func TestNoImprovement_PatienceLessThanOneClamped(t *testing.T) {
	c := NoImprovement(0, false) // should be clamped to 1
	c.ShouldStop(0, 1.0)         // initialize
	if !c.ShouldStop(1, 1.0) {   // patience=1: one stagnant gen => stop
		t.Fatal("patience clamped to 1: one stagnant gen should stop")
	}
}

func TestNoImprovement_NegativePatienceClamped(t *testing.T) {
	c := NoImprovement(-5, false) // should be clamped to 1
	c.ShouldStop(0, 1.0)
	if !c.ShouldStop(1, 1.0) {
		t.Fatal("negative patience clamped to 1: one stagnant gen should stop")
	}
}

// ─── Any tests ─────────────────────────────────────────────────────────────

func TestAny_EmptyNeverStops(t *testing.T) {
	c := Any()
	for i := 0; i < 10; i++ {
		if c.ShouldStop(i, float64(i)) {
			t.Fatalf("Any() (empty) should never stop, but stopped at gen %d", i)
		}
	}
}

func TestAny_OneMatchStops(t *testing.T) {
	never := ScoreThreshold(9999.0, false) // won't fire
	fires := ScoreThreshold(10.0, false)   // fires at score >= 10
	c := Any(never, fires)
	if !c.ShouldStop(0, 10.0) {
		t.Fatal("Any(never, fires): expected stop when fires criterion is met")
	}
}

func TestAny_NoneMatchContinues(t *testing.T) {
	c := Any(ScoreThreshold(50.0, false), ScoreThreshold(100.0, false))
	if c.ShouldStop(0, 20.0) {
		t.Fatal("Any: no criterion fires at score=20, should not stop")
	}
}

func TestAny_AllMatchStops(t *testing.T) {
	c := Any(ScoreThreshold(5.0, false), ScoreThreshold(10.0, false))
	if !c.ShouldStop(0, 15.0) {
		t.Fatal("Any: all criteria fire at score=15, should stop")
	}
}

// ─── All tests ─────────────────────────────────────────────────────────────

func TestAll_EmptyAlwaysStops(t *testing.T) {
	c := All()
	if !c.ShouldStop(0, 0.0) {
		t.Fatal("All() (empty): expected true (vacuously all conditions met)")
	}
}

func TestAll_OneNotMatchContinues(t *testing.T) {
	fires := ScoreThreshold(10.0, false)   // fires
	never := ScoreThreshold(9999.0, false) // won't fire
	c := All(fires, never)
	if c.ShouldStop(0, 20.0) {
		t.Fatal("All(fires, never): should not stop because 'never' didn't fire")
	}
}

func TestAll_AllMatchStops(t *testing.T) {
	c := All(ScoreThreshold(5.0, false), ScoreThreshold(10.0, false))
	if !c.ShouldStop(0, 20.0) {
		t.Fatal("All: both criteria fire at score=20, should stop")
	}
}

func TestAll_NoneMatchContinues(t *testing.T) {
	c := All(ScoreThreshold(50.0, false), ScoreThreshold(100.0, false))
	if c.ShouldStop(0, 20.0) {
		t.Fatal("All: no criterion fires at score=20, should not stop")
	}
}

// ─── Composition tests ─────────────────────────────────────────────────────

func TestComposition_AnyOfAll(t *testing.T) {
	// Stop if (score >= 100 AND score <= 200) OR score >= 500
	// With score = 150: first All fires (100<=150<=200), Any should stop.
	c := Any(
		All(ScoreThreshold(100.0, false), ScoreThreshold(200.0, true)),
		ScoreThreshold(500.0, false),
	)
	if !c.ShouldStop(0, 150.0) {
		t.Fatal("AnyOfAll: expected stop at score=150 (100<=150<=200)")
	}
	if c.ShouldStop(0, 300.0) {
		t.Fatal("AnyOfAll: 300 fails the All group and doesn't meet 500 threshold, should not stop")
	}
	if !c.ShouldStop(0, 600.0) {
		t.Fatal("AnyOfAll: expected stop at score=600 (>=500)")
	}
}
