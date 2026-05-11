// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package promotion

import (
	"encoding/json"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/design/objectives"
	"github.com/gmlewis/gep/v2/design/scenarios"
)

// --- helpers ---

// makeResult constructs an AggregateResult with the given score and hard-fail flag.
func makeResult(score float64, hardFailed bool) objectives.AggregateResult {
	return objectives.AggregateResult{
		AggregateScore: score,
		Breakdown: objectives.ScoreBreakdown{
			HardFailed: hardFailed,
		},
	}
}

// makeSummary is a shorthand for building a SplitEvalSummary inline.
func makeSummary(split scenarios.ScenarioSplit, count int, mean, min, max float64, hardFails int) SplitEvalSummary {
	return SplitEvalSummary{
		Split:              split,
		Count:              count,
		MeanAggregateScore: mean,
		MinAggregateScore:  min,
		MaxAggregateScore:  max,
		HardFailCount:      hardFails,
	}
}

// --- SummarizeSplit tests ---

func TestSummarizeSplitEmpty(t *testing.T) {
	s := SummarizeSplit(scenarios.Train, nil)
	if s.Split != scenarios.Train {
		t.Errorf("Split = %q, want %q", s.Split, scenarios.Train)
	}
	if s.Count != 0 {
		t.Errorf("Count = %d, want 0", s.Count)
	}
	if s.MeanAggregateScore != 0 {
		t.Errorf("MeanAggregateScore = %v, want 0", s.MeanAggregateScore)
	}
	if s.MinAggregateScore != 0 {
		t.Errorf("MinAggregateScore = %v, want 0", s.MinAggregateScore)
	}
	if s.MaxAggregateScore != 0 {
		t.Errorf("MaxAggregateScore = %v, want 0", s.MaxAggregateScore)
	}
	if s.HardFailCount != 0 {
		t.Errorf("HardFailCount = %d, want 0", s.HardFailCount)
	}
}

func TestSummarizeSplitSingleResult(t *testing.T) {
	results := []objectives.AggregateResult{makeResult(5.0, false)}
	s := SummarizeSplit(scenarios.Validation, results)

	if s.Split != scenarios.Validation {
		t.Errorf("Split = %q, want %q", s.Split, scenarios.Validation)
	}
	if s.Count != 1 {
		t.Errorf("Count = %d, want 1", s.Count)
	}
	if s.MeanAggregateScore != 5.0 {
		t.Errorf("MeanAggregateScore = %v, want 5.0", s.MeanAggregateScore)
	}
	if s.MinAggregateScore != 5.0 {
		t.Errorf("MinAggregateScore = %v, want 5.0", s.MinAggregateScore)
	}
	if s.MaxAggregateScore != 5.0 {
		t.Errorf("MaxAggregateScore = %v, want 5.0", s.MaxAggregateScore)
	}
	if s.HardFailCount != 0 {
		t.Errorf("HardFailCount = %d, want 0", s.HardFailCount)
	}
}

func TestSummarizeSplitMultipleResults(t *testing.T) {
	results := []objectives.AggregateResult{
		makeResult(2.0, false),
		makeResult(4.0, false),
		makeResult(6.0, false),
	}
	s := SummarizeSplit(scenarios.Train, results)

	if s.Count != 3 {
		t.Errorf("Count = %d, want 3", s.Count)
	}
	// mean = (2+4+6)/3 = 4.0
	if s.MeanAggregateScore != 4.0 {
		t.Errorf("MeanAggregateScore = %v, want 4.0", s.MeanAggregateScore)
	}
	if s.MinAggregateScore != 2.0 {
		t.Errorf("MinAggregateScore = %v, want 2.0", s.MinAggregateScore)
	}
	if s.MaxAggregateScore != 6.0 {
		t.Errorf("MaxAggregateScore = %v, want 6.0", s.MaxAggregateScore)
	}
	if s.HardFailCount != 0 {
		t.Errorf("HardFailCount = %d, want 0", s.HardFailCount)
	}
}

func TestSummarizeSplitCountsHardFails(t *testing.T) {
	results := []objectives.AggregateResult{
		makeResult(0, true),
		makeResult(5.0, false),
		makeResult(0, true),
	}
	s := SummarizeSplit(scenarios.Test, results)

	if s.HardFailCount != 2 {
		t.Errorf("HardFailCount = %d, want 2", s.HardFailCount)
	}
	if s.Count != 3 {
		t.Errorf("Count = %d, want 3", s.Count)
	}
}

func TestSummarizeSplitIsDeterministic(t *testing.T) {
	results := []objectives.AggregateResult{
		makeResult(1.0, false),
		makeResult(3.0, true),
		makeResult(2.0, false),
	}
	s1 := SummarizeSplit(scenarios.Train, results)
	s2 := SummarizeSplit(scenarios.Train, results)
	if !reflect.DeepEqual(s1, s2) {
		t.Fatal("SummarizeSplit is not deterministic")
	}
}

func TestSummarizeSplitNoFiniteInfInResult(t *testing.T) {
	// Ensure that the initial +Inf/-Inf sentinels are not leaked when results
	// is non-empty.
	results := []objectives.AggregateResult{makeResult(3.0, false)}
	s := SummarizeSplit(scenarios.Train, results)
	if math.IsInf(s.MinAggregateScore, 0) {
		t.Errorf("MinAggregateScore = %v, should not be Inf", s.MinAggregateScore)
	}
	if math.IsInf(s.MaxAggregateScore, 0) {
		t.Errorf("MaxAggregateScore = %v, should not be Inf", s.MaxAggregateScore)
	}
}

// --- Decide tests ---

func TestDecidePassAboveThreshold(t *testing.T) {
	criterion := AcceptanceCriterion{
		Split:             scenarios.Validation,
		MinAggregateScore: 0.8,
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Validation, 3, 0.9, 0.8, 1.0, 0),
	}
	d := Decide(criterion, summaries)
	if !d.Passed {
		t.Errorf("expected Passed=true, got Reason=%q", d.Reason)
	}
	if d.Split != scenarios.Validation {
		t.Errorf("Split = %q, want %q", d.Split, scenarios.Validation)
	}
}

func TestDecidePassAtExactThreshold(t *testing.T) {
	// mean == threshold should pass (>= semantics).
	criterion := AcceptanceCriterion{
		Split:             scenarios.Validation,
		MinAggregateScore: 0.75,
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Validation, 2, 0.75, 0.75, 0.75, 0),
	}
	d := Decide(criterion, summaries)
	if !d.Passed {
		t.Errorf("expected Passed=true at exact threshold, got Reason=%q", d.Reason)
	}
}

func TestDecideFailBelowThreshold(t *testing.T) {
	criterion := AcceptanceCriterion{
		Split:             scenarios.Train,
		MinAggregateScore: 1.0,
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Train, 5, 0.7, 0.5, 0.9, 0),
	}
	d := Decide(criterion, summaries)
	if d.Passed {
		t.Errorf("expected Passed=false when mean score below threshold")
	}
	if !strings.Contains(d.Reason, "train") {
		t.Errorf("Reason %q should mention split name", d.Reason)
	}
}

func TestDecideFailNoSummaryForSplit(t *testing.T) {
	criterion := AcceptanceCriterion{
		Split:             scenarios.Test,
		MinAggregateScore: 0.5,
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Train, 3, 0.9, 0.8, 1.0, 0),
	}
	d := Decide(criterion, summaries)
	if d.Passed {
		t.Errorf("expected Passed=false when no matching summary")
	}
	if !strings.Contains(d.Reason, "test") {
		t.Errorf("Reason %q should mention split name", d.Reason)
	}
}

func TestDecideFailEmptySummary(t *testing.T) {
	criterion := AcceptanceCriterion{
		Split:             scenarios.Validation,
		MinAggregateScore: 0.0,
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Validation, 0, 0, 0, 0, 0),
	}
	d := Decide(criterion, summaries)
	if d.Passed {
		t.Errorf("expected Passed=false when Count=0")
	}
}

func TestDecideFailHardFailRequired(t *testing.T) {
	criterion := AcceptanceCriterion{
		Split:             scenarios.Validation,
		MinAggregateScore: 0.5,
		RequireNoHardFail: true,
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Validation, 3, 0.9, 0.7, 1.0, 1),
	}
	d := Decide(criterion, summaries)
	if d.Passed {
		t.Errorf("expected Passed=false when RequireNoHardFail=true and HardFailCount>0")
	}
	if !strings.Contains(d.Reason, "hard failure") {
		t.Errorf("Reason %q should mention hard failure", d.Reason)
	}
}

func TestDecidePassHardFailNotRequired(t *testing.T) {
	// RequireNoHardFail defaults to false; hard fails should not block.
	criterion := AcceptanceCriterion{
		Split:             scenarios.Validation,
		MinAggregateScore: 0.5,
		RequireNoHardFail: false,
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Validation, 3, 0.9, 0.7, 1.0, 2),
	}
	d := Decide(criterion, summaries)
	if !d.Passed {
		t.Errorf("expected Passed=true when RequireNoHardFail=false, got Reason=%q", d.Reason)
	}
}

func TestDecideIsDeterministic(t *testing.T) {
	criterion := AcceptanceCriterion{
		Split:             scenarios.Train,
		MinAggregateScore: 0.5,
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Train, 4, 0.8, 0.6, 1.0, 0),
	}
	d1 := Decide(criterion, summaries)
	d2 := Decide(criterion, summaries)
	if !reflect.DeepEqual(d1, d2) {
		t.Fatal("Decide is not deterministic")
	}
}

// --- Evaluate tests ---

func TestEvaluateNoCriteriaPromoted(t *testing.T) {
	// Vacuously true: no criteria → promoted.
	report := Evaluate("cand-1", nil, nil)
	if !report.Promoted {
		t.Fatal("expected Promoted=true with no criteria")
	}
	if report.CandidateID != "cand-1" {
		t.Errorf("CandidateID = %q, want %q", report.CandidateID, "cand-1")
	}
	if len(report.Decisions) != 0 {
		t.Errorf("Decisions count = %d, want 0", len(report.Decisions))
	}
}

func TestEvaluateAllCriteriaPass(t *testing.T) {
	criteria := []AcceptanceCriterion{
		{Split: scenarios.Train, MinAggregateScore: 0.5},
		{Split: scenarios.Validation, MinAggregateScore: 0.6},
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Train, 5, 0.8, 0.7, 0.9, 0),
		makeSummary(scenarios.Validation, 3, 0.75, 0.6, 0.9, 0),
	}

	report := Evaluate("cand-2", criteria, summaries)

	if !report.Promoted {
		t.Errorf("expected Promoted=true when all criteria pass")
	}
	if len(report.Decisions) != 2 {
		t.Fatalf("Decisions count = %d, want 2", len(report.Decisions))
	}
	for _, d := range report.Decisions {
		if !d.Passed {
			t.Errorf("criterion for split %q unexpectedly failed: %s", d.Split, d.Reason)
		}
	}
}

func TestEvaluateOneCriterionFails(t *testing.T) {
	criteria := []AcceptanceCriterion{
		{Split: scenarios.Train, MinAggregateScore: 0.5},
		{Split: scenarios.Validation, MinAggregateScore: 0.9}, // too high
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Train, 5, 0.8, 0.7, 0.9, 0),
		makeSummary(scenarios.Validation, 3, 0.7, 0.6, 0.8, 0), // mean 0.7 < 0.9
	}

	report := Evaluate("cand-3", criteria, summaries)

	if report.Promoted {
		t.Error("expected Promoted=false when one criterion fails")
	}
	if !report.Decisions[0].Passed {
		t.Errorf("train decision should have passed")
	}
	if report.Decisions[1].Passed {
		t.Errorf("validation decision should have failed")
	}
}

func TestEvaluateDecisionsInCriteriaOrder(t *testing.T) {
	// Ensure decisions are stored in the same order as criteria.
	criteria := []AcceptanceCriterion{
		{Split: scenarios.Test, MinAggregateScore: 0.5},
		{Split: scenarios.Train, MinAggregateScore: 0.5},
		{Split: scenarios.Validation, MinAggregateScore: 0.5},
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Train, 1, 0.9, 0.9, 0.9, 0),
		makeSummary(scenarios.Validation, 1, 0.9, 0.9, 0.9, 0),
		makeSummary(scenarios.Test, 1, 0.9, 0.9, 0.9, 0),
	}

	report := Evaluate("cand-4", criteria, summaries)

	wantOrder := []scenarios.ScenarioSplit{scenarios.Test, scenarios.Train, scenarios.Validation}
	if len(report.Decisions) != 3 {
		t.Fatalf("Decisions count = %d, want 3", len(report.Decisions))
	}
	for i, d := range report.Decisions {
		if d.Split != wantOrder[i] {
			t.Errorf("Decisions[%d].Split = %q, want %q", i, d.Split, wantOrder[i])
		}
	}
}

func TestEvaluateSummariesPreserved(t *testing.T) {
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Train, 2, 1.0, 0.9, 1.0, 0),
	}
	report := Evaluate("cand-5", nil, summaries)
	if !reflect.DeepEqual(report.Summaries, summaries) {
		t.Errorf("Summaries not preserved: got %+v, want %+v", report.Summaries, summaries)
	}
}

func TestEvaluateIsDeterministic(t *testing.T) {
	criteria := []AcceptanceCriterion{
		{Split: scenarios.Validation, MinAggregateScore: 0.7, RequireNoHardFail: true},
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Validation, 4, 0.85, 0.7, 1.0, 0),
	}

	r1 := Evaluate("cand-6", criteria, summaries)
	r2 := Evaluate("cand-6", criteria, summaries)

	if !reflect.DeepEqual(r1, r2) {
		t.Fatal("Evaluate is not deterministic")
	}
}

func TestEvaluateHardFailBlocksPromotion(t *testing.T) {
	criteria := []AcceptanceCriterion{
		{
			Split:             scenarios.Validation,
			MinAggregateScore: 0.5,
			RequireNoHardFail: true,
		},
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Validation, 5, 0.9, 0.7, 1.0, 3), // 3 hard fails
	}

	report := Evaluate("cand-7", criteria, summaries)

	if report.Promoted {
		t.Error("expected Promoted=false when hard fails present and RequireNoHardFail=true")
	}
}

func TestEvaluateEdgeThresholdAtBoundary(t *testing.T) {
	const threshold = 0.75
	criteria := []AcceptanceCriterion{
		{Split: scenarios.Test, MinAggregateScore: threshold},
	}

	// Exactly at threshold: should pass.
	atThreshold := []SplitEvalSummary{
		makeSummary(scenarios.Test, 2, threshold, threshold, threshold, 0),
	}
	above := Evaluate("edge-at", criteria, atThreshold)
	if !above.Promoted {
		t.Errorf("expected Promoted=true at exact threshold %.4f", threshold)
	}

	// Just below threshold: should fail.
	const epsilon = 1e-9
	belowThreshold := []SplitEvalSummary{
		makeSummary(scenarios.Test, 2, threshold-epsilon, threshold-epsilon, threshold-epsilon, 0),
	}
	below := Evaluate("edge-below", criteria, belowThreshold)
	if below.Promoted {
		t.Errorf("expected Promoted=false just below threshold %.4f", threshold)
	}
}

// --- JSON round-trip tests ---

func TestAcceptanceCriterionJSONRoundTrip(t *testing.T) {
	orig := AcceptanceCriterion{
		Split:             scenarios.Validation,
		MinAggregateScore: 0.85,
		RequireNoHardFail: true,
	}
	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var got AcceptanceCriterion
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !reflect.DeepEqual(orig, got) {
		t.Errorf("round-trip mismatch:\n orig=%+v\n  got=%+v", orig, got)
	}
}

func TestPromotionReportJSONRoundTrip(t *testing.T) {
	criteria := []AcceptanceCriterion{
		{Split: scenarios.Train, MinAggregateScore: 0.6},
		{Split: scenarios.Validation, MinAggregateScore: 0.7, RequireNoHardFail: true},
	}
	summaries := []SplitEvalSummary{
		makeSummary(scenarios.Train, 4, 0.8, 0.6, 0.9, 0),
		makeSummary(scenarios.Validation, 2, 0.75, 0.7, 0.8, 0),
	}
	orig := Evaluate("cand-rt", criteria, summaries)

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var got PromotionReport
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !reflect.DeepEqual(orig, got) {
		t.Errorf("round-trip mismatch:\n orig=%+v\n  got=%+v", orig, got)
	}
}

func TestSplitEvalSummaryJSONRoundTrip(t *testing.T) {
	orig := makeSummary(scenarios.Test, 6, 0.88, 0.7, 1.0, 1)
	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var got SplitEvalSummary
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !reflect.DeepEqual(orig, got) {
		t.Errorf("round-trip mismatch:\n orig=%+v\n  got=%+v", orig, got)
	}
}
