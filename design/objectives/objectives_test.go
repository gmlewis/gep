// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package objectives

import (
	"reflect"
	"sort"
	"testing"
)

// --- helpers ---

func softDef(name string, weight float64) ObjectiveDef {
	return ObjectiveDef{Name: name, Weight: weight, Kind: Soft}
}

func hardDef(name string, weight float64) ObjectiveDef {
	return ObjectiveDef{Name: name, Weight: weight, Kind: Hard}
}

// --- tests ---

func TestScoreWeightedAggregation(t *testing.T) {
	defs := []ObjectiveDef{
		softDef("accuracy", 2.0),
		softDef("speed", 1.0),
	}
	raw := map[string]float64{"accuracy": 0.9, "speed": 0.5}

	got := Score(defs, raw, false, 0)

	// 0.9*2.0 + 0.5*1.0 = 1.8 + 0.5 = 2.3
	const wantScore = 2.3
	if got.AggregateScore != wantScore {
		t.Fatalf("AggregateScore = %v, want %v", got.AggregateScore, wantScore)
	}
	if got.Breakdown.HardFailed {
		t.Fatal("expected HardFailed=false")
	}
	if len(got.Breakdown.Contributions) != 2 {
		t.Fatalf("contributions count = %d, want 2", len(got.Breakdown.Contributions))
	}
	if got.Breakdown.Contributions[0].Name != "accuracy" {
		t.Errorf("contributions[0].Name = %q, want \"accuracy\"", got.Breakdown.Contributions[0].Name)
	}
	if got.Breakdown.Contributions[0].WeightedScore != 1.8 {
		t.Errorf("contributions[0].WeightedScore = %v, want 1.8", got.Breakdown.Contributions[0].WeightedScore)
	}
}

func TestScoreHardFailFromRejected(t *testing.T) {
	defs := []ObjectiveDef{
		softDef("accuracy", 2.0),
		softDef("speed", 1.0),
	}
	raw := map[string]float64{"accuracy": 0.9, "speed": 0.5}

	got := Score(defs, raw, true, 0)

	if !got.Breakdown.HardFailed {
		t.Fatal("expected HardFailed=true when rejected=true")
	}
	if got.AggregateScore != 0 {
		t.Fatalf("AggregateScore = %v, want 0 when hard-failed", got.AggregateScore)
	}
}

func TestScoreHardFailFromHardObjectiveZero(t *testing.T) {
	defs := []ObjectiveDef{
		hardDef("feasibility", 1.0),
		softDef("accuracy", 2.0),
	}
	raw := map[string]float64{"feasibility": 0, "accuracy": 0.9}

	got := Score(defs, raw, false, 0)

	if !got.Breakdown.HardFailed {
		t.Fatal("expected HardFailed=true when Hard objective score=0")
	}
	if got.AggregateScore != 0 {
		t.Fatalf("AggregateScore = %v, want 0 when hard-failed", got.AggregateScore)
	}
}

func TestScoreHardFailFromHardObjectiveNegative(t *testing.T) {
	defs := []ObjectiveDef{
		hardDef("feasibility", 1.0),
	}
	raw := map[string]float64{"feasibility": -5.0}

	got := Score(defs, raw, false, 0)

	if !got.Breakdown.HardFailed {
		t.Fatal("expected HardFailed=true when Hard objective score<0")
	}
	if got.AggregateScore != 0 {
		t.Fatalf("AggregateScore = %v, want 0", got.AggregateScore)
	}
}

func TestScoreHardFailDominatesSoftObjectives(t *testing.T) {
	// Even very high soft-objective scores must be dominated by a hard failure.
	defs := []ObjectiveDef{
		hardDef("feasibility", 1.0),
		softDef("accuracy", 100.0),
		softDef("speed", 100.0),
	}
	raw := map[string]float64{
		"feasibility": 0,    // hard fail
		"accuracy":    1000, // would have been huge
		"speed":       1000,
	}

	got := Score(defs, raw, false, 0)

	if !got.Breakdown.HardFailed {
		t.Fatal("expected HardFailed=true")
	}
	if got.AggregateScore != 0 {
		t.Fatalf("AggregateScore = %v, want 0 (hard failure dominates)", got.AggregateScore)
	}
}

func TestScorePenaltyReducesAggregate(t *testing.T) {
	defs := []ObjectiveDef{
		softDef("quality", 1.0),
	}
	raw := map[string]float64{"quality": 10.0}
	const penalty = 3.0

	got := Score(defs, raw, false, penalty)

	// 10*1.0 - 3.0 = 7.0
	const wantScore = 7.0
	if got.AggregateScore != wantScore {
		t.Fatalf("AggregateScore = %v, want %v", got.AggregateScore, wantScore)
	}
	if got.Breakdown.TotalPenalty != penalty {
		t.Fatalf("TotalPenalty = %v, want %v", got.Breakdown.TotalPenalty, penalty)
	}
}

func TestScoreDeterministicForFixedInput(t *testing.T) {
	defs := []ObjectiveDef{
		softDef("a", 1.5),
		softDef("b", 0.5),
	}
	raw := map[string]float64{"a": 3.0, "b": 2.0}

	got1 := Score(defs, raw, false, 0.25)
	got2 := Score(defs, raw, false, 0.25)

	if !reflect.DeepEqual(got1, got2) {
		t.Fatalf("determinism mismatch:\n got1=%#v\n got2=%#v", got1, got2)
	}
}

func TestScoreContributionsInDefinitionOrder(t *testing.T) {
	defs := []ObjectiveDef{
		softDef("z", 1.0),
		softDef("a", 1.0),
		softDef("m", 1.0),
	}
	raw := map[string]float64{"z": 1, "a": 2, "m": 3}

	got := Score(defs, raw, false, 0)

	names := make([]string, len(got.Breakdown.Contributions))
	for i, c := range got.Breakdown.Contributions {
		names[i] = c.Name
	}
	want := []string{"z", "a", "m"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("contribution order = %v, want %v", names, want)
	}
}

func TestScoreEmptyDefs(t *testing.T) {
	got := Score(nil, nil, false, 0)

	if got.Breakdown.HardFailed {
		t.Fatal("expected HardFailed=false for empty defs")
	}
	if got.AggregateScore != 0 {
		t.Fatalf("AggregateScore = %v, want 0", got.AggregateScore)
	}
	if len(got.Breakdown.Contributions) != 0 {
		t.Fatalf("contributions count = %d, want 0", len(got.Breakdown.Contributions))
	}
}

func TestScoreMissingRawScoreDefaultsToZero(t *testing.T) {
	defs := []ObjectiveDef{softDef("known", 2.0), softDef("missing", 3.0)}
	raw := map[string]float64{"known": 5.0}

	got := Score(defs, raw, false, 0)

	// 5*2 + 0*3 = 10
	if got.AggregateScore != 10 {
		t.Fatalf("AggregateScore = %v, want 10", got.AggregateScore)
	}
	if got.Breakdown.Contributions[1].RawScore != 0 {
		t.Fatalf("missing contribution RawScore = %v, want 0", got.Breakdown.Contributions[1].RawScore)
	}
}

func TestLessBetterScoreFirst(t *testing.T) {
	defs := []ObjectiveDef{softDef("q", 1.0)}
	a := Score(defs, map[string]float64{"q": 5}, false, 0) // score 5
	b := Score(defs, map[string]float64{"q": 3}, false, 0) // score 3

	if !Less(a, b) {
		t.Fatal("Less(a, b) should be true when a has higher score")
	}
	if Less(b, a) {
		t.Fatal("Less(b, a) should be false when b has lower score")
	}
}

func TestLessTieBreakingByContribution(t *testing.T) {
	// Same aggregate score through different distributions.
	defs := []ObjectiveDef{
		softDef("x", 1.0),
		softDef("y", 1.0),
	}
	// a: x=3, y=1 -> total=4; b: x=1, y=3 -> total=4
	a := Score(defs, map[string]float64{"x": 3, "y": 1}, false, 0)
	b := Score(defs, map[string]float64{"x": 1, "y": 3}, false, 0)

	// a has higher WeightedScore on first contribution, so a < b (ranks better).
	if !Less(a, b) {
		t.Fatal("Less(a, b) should be true: a has higher first-contribution score")
	}
	if Less(b, a) {
		t.Fatal("Less(b, a) should be false")
	}
}

func TestLessEqualCandidatesReturnFalse(t *testing.T) {
	defs := []ObjectiveDef{softDef("q", 1.0)}
	a := Score(defs, map[string]float64{"q": 5}, false, 0)
	b := Score(defs, map[string]float64{"q": 5}, false, 0)

	if Less(a, b) || Less(b, a) {
		t.Fatal("equal candidates should return false from Less in both directions")
	}
}

func TestLessUsableForSort(t *testing.T) {
	defs := []ObjectiveDef{softDef("q", 1.0)}
	candidates := []AggregateResult{
		Score(defs, map[string]float64{"q": 2}, false, 0),
		Score(defs, map[string]float64{"q": 5}, false, 0),
		Score(defs, map[string]float64{"q": 1}, false, 0),
		Score(defs, map[string]float64{"q": 4}, false, 0),
	}

	sort.Slice(candidates, func(i, j int) bool {
		return Less(candidates[i], candidates[j])
	})

	scores := make([]float64, len(candidates))
	for i, c := range candidates {
		scores[i] = c.AggregateScore
	}
	want := []float64{5, 4, 2, 1}
	if !reflect.DeepEqual(scores, want) {
		t.Fatalf("sorted scores = %v, want %v", scores, want)
	}
}

func TestObjectiveKindString(t *testing.T) {
	cases := []struct {
		k    ObjectiveKind
		want string
	}{
		{Soft, "soft"},
		{Hard, "hard"},
		{ObjectiveKind(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.k.String(); got != tc.want {
			t.Errorf("ObjectiveKind(%d).String() = %q, want %q", tc.k, got, tc.want)
		}
	}
}
