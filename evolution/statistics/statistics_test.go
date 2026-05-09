// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package statistics

import (
	"math"
	"testing"
)

// --- Record tests ---

func TestRecord_EmptyScores(t *testing.T) {
	s := Record(0, nil, nil, false)
	if s.Generation != 0 {
		t.Errorf("Generation=%d, want 0", s.Generation)
	}
	if s.BestScore != 0 || s.WorstScore != 0 || s.MeanScore != 0 || s.Diversity != 0 {
		t.Errorf("non-zero stats for empty scores: %+v", s)
	}
}

func TestRecord_SingleScore_Maximize(t *testing.T) {
	s := Record(3, []float64{7.5}, nil, false)
	if s.Generation != 3 {
		t.Errorf("Generation=%d, want 3", s.Generation)
	}
	if s.BestScore != 7.5 {
		t.Errorf("BestScore=%v, want 7.5", s.BestScore)
	}
	if s.WorstScore != 7.5 {
		t.Errorf("WorstScore=%v, want 7.5", s.WorstScore)
	}
	if s.MeanScore != 7.5 {
		t.Errorf("MeanScore=%v, want 7.5", s.MeanScore)
	}
	if s.Diversity != 0 {
		t.Errorf("Diversity=%v, want 0 (no karva strings)", s.Diversity)
	}
}

func TestRecord_Maximize_BestWorstMean(t *testing.T) {
	scores := []float64{1.0, 5.0, 3.0, 2.0, 4.0}
	s := Record(1, scores, nil, false)
	if s.BestScore != 5.0 {
		t.Errorf("BestScore=%v, want 5.0", s.BestScore)
	}
	if s.WorstScore != 1.0 {
		t.Errorf("WorstScore=%v, want 1.0", s.WorstScore)
	}
	wantMean := (1.0 + 5.0 + 3.0 + 2.0 + 4.0) / 5.0
	if math.Abs(s.MeanScore-wantMean) > 1e-9 {
		t.Errorf("MeanScore=%v, want %v", s.MeanScore, wantMean)
	}
}

func TestRecord_Minimize_BestWorstMean(t *testing.T) {
	scores := []float64{1.0, 5.0, 3.0, 2.0, 4.0}
	s := Record(2, scores, nil, true)
	if s.BestScore != 1.0 {
		t.Errorf("BestScore=%v, want 1.0 (minimize)", s.BestScore)
	}
	if s.WorstScore != 5.0 {
		t.Errorf("WorstScore=%v, want 5.0 (minimize)", s.WorstScore)
	}
	wantMean := (1.0 + 5.0 + 3.0 + 2.0 + 4.0) / 5.0
	if math.Abs(s.MeanScore-wantMean) > 1e-9 {
		t.Errorf("MeanScore=%v, want %v", s.MeanScore, wantMean)
	}
}

func TestRecord_NegativeScores(t *testing.T) {
	scores := []float64{-5.0, -2.0, -8.0, -1.0}
	s := Record(0, scores, nil, false)
	if s.BestScore != -1.0 {
		t.Errorf("BestScore=%v, want -1.0", s.BestScore)
	}
	if s.WorstScore != -8.0 {
		t.Errorf("WorstScore=%v, want -8.0", s.WorstScore)
	}
}

func TestRecord_AllEqualScores(t *testing.T) {
	scores := []float64{3.0, 3.0, 3.0}
	s := Record(0, scores, nil, false)
	if s.BestScore != 3.0 {
		t.Errorf("BestScore=%v, want 3.0", s.BestScore)
	}
	if s.WorstScore != 3.0 {
		t.Errorf("WorstScore=%v, want 3.0", s.WorstScore)
	}
	if s.MeanScore != 3.0 {
		t.Errorf("MeanScore=%v, want 3.0", s.MeanScore)
	}
}

func TestRecord_GenerationIndex(t *testing.T) {
	for _, gen := range []int{0, 1, 42, 999} {
		s := Record(gen, []float64{1.0}, nil, false)
		if s.Generation != gen {
			t.Errorf("gen=%d: Stats.Generation=%d, want %d", gen, s.Generation, gen)
		}
	}
}

// --- Diversity tests ---

func TestRecord_Diversity_AllUnique(t *testing.T) {
	karvas := []string{"a|b", "c|d", "e|f", "g|h"}
	s := Record(0, []float64{1, 2, 3, 4}, karvas, false)
	if s.Diversity != 1.0 {
		t.Errorf("Diversity=%v, want 1.0 (all unique)", s.Diversity)
	}
}

func TestRecord_Diversity_AllIdentical(t *testing.T) {
	karvas := []string{"a|b", "a|b", "a|b"}
	s := Record(0, []float64{1, 2, 3}, karvas, false)
	wantDiv := 1.0 / 3.0
	if math.Abs(s.Diversity-wantDiv) > 1e-9 {
		t.Errorf("Diversity=%v, want %v (all identical)", s.Diversity, wantDiv)
	}
}

func TestRecord_Diversity_Mixed(t *testing.T) {
	// 4 karva strings, 3 unique → diversity = 3/4
	karvas := []string{"a", "b", "a", "c"}
	s := Record(0, []float64{1, 2, 3, 4}, karvas, false)
	wantDiv := 3.0 / 4.0
	if math.Abs(s.Diversity-wantDiv) > 1e-9 {
		t.Errorf("Diversity=%v, want %v", s.Diversity, wantDiv)
	}
}

func TestRecord_Diversity_NilKarvaStrings(t *testing.T) {
	s := Record(0, []float64{1, 2, 3}, nil, false)
	if s.Diversity != 0 {
		t.Errorf("Diversity=%v, want 0 (nil karva strings)", s.Diversity)
	}
}

func TestRecord_Diversity_EmptyKarvaStrings(t *testing.T) {
	s := Record(0, []float64{1, 2, 3}, []string{}, false)
	if s.Diversity != 0 {
		t.Errorf("Diversity=%v, want 0 (empty karva strings)", s.Diversity)
	}
}

// --- Collector tests ---

func TestCollector_AppendAndHistory(t *testing.T) {
	col := &Collector{}
	if len(col.History) != 0 {
		t.Fatalf("new Collector has non-empty history: %d entries", len(col.History))
	}

	s0 := Record(0, []float64{1, 2, 3}, nil, false)
	s1 := Record(1, []float64{4, 5, 6}, nil, false)
	col.Append(s0)
	col.Append(s1)

	if len(col.History) != 2 {
		t.Fatalf("History len=%d, want 2", len(col.History))
	}
	if col.History[0].Generation != 0 {
		t.Errorf("History[0].Generation=%d, want 0", col.History[0].Generation)
	}
	if col.History[1].Generation != 1 {
		t.Errorf("History[1].Generation=%d, want 1", col.History[1].Generation)
	}
}

func TestCollector_Last_EmptyHistory(t *testing.T) {
	col := &Collector{}
	last := col.Last()
	if last != (Stats{}) {
		t.Errorf("Last() on empty Collector=%+v, want zero Stats", last)
	}
}

func TestCollector_Last_ReturnsNewest(t *testing.T) {
	col := &Collector{}
	col.Append(Record(0, []float64{1.0}, nil, false))
	col.Append(Record(1, []float64{2.0}, nil, false))
	col.Append(Record(2, []float64{9.0}, nil, false))

	last := col.Last()
	if last.Generation != 2 {
		t.Errorf("Last().Generation=%d, want 2", last.Generation)
	}
	if last.BestScore != 9.0 {
		t.Errorf("Last().BestScore=%v, want 9.0", last.BestScore)
	}
}

func TestCollector_OrderPreserved(t *testing.T) {
	col := &Collector{}
	for i := 0; i < 10; i++ {
		col.Append(Record(i, []float64{float64(i)}, nil, false))
	}
	for i, s := range col.History {
		if s.Generation != i {
			t.Errorf("History[%d].Generation=%d, want %d", i, s.Generation, i)
		}
	}
}
