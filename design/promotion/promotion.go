// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package promotion

import (
	"fmt"
	"math"

	"github.com/gmlewis/gep/v2/design/objectives"
	"github.com/gmlewis/gep/v2/design/scenarios"
)

// AcceptanceCriterion defines the minimum requirements for a candidate to be
// considered for promotion after evaluation on one split.
type AcceptanceCriterion struct {
	// Split is the evaluation partition this criterion applies to.
	Split scenarios.ScenarioSplit `json:"split"`
	// MinAggregateScore is the minimum mean aggregate score required across all
	// evaluated scenarios in the split. The candidate passes when
	// SplitEvalSummary.MeanAggregateScore >= MinAggregateScore.
	MinAggregateScore float64 `json:"min_aggregate_score"`
	// RequireNoHardFail, when true, requires that no scenario in the split
	// produced a hard failure. The candidate fails when
	// SplitEvalSummary.HardFailCount > 0.
	RequireNoHardFail bool `json:"require_no_hard_fail,omitempty"`
}

// SplitEvalSummary is a per-split summary of evaluation results for one
// candidate. It is computed from a slice of [objectives.AggregateResult]
// values via [SummarizeSplit].
type SplitEvalSummary struct {
	// Split identifies the evaluation partition.
	Split scenarios.ScenarioSplit `json:"split"`
	// Count is the total number of scenarios evaluated in this split.
	Count int `json:"count"`
	// MeanAggregateScore is the mean aggregate score over all evaluated
	// scenarios in this split.
	MeanAggregateScore float64 `json:"mean_aggregate_score"`
	// MinAggregateScore is the lowest aggregate score seen in this split.
	MinAggregateScore float64 `json:"min_aggregate_score"`
	// MaxAggregateScore is the highest aggregate score seen in this split.
	MaxAggregateScore float64 `json:"max_aggregate_score"`
	// HardFailCount is the number of scenarios on which the candidate
	// hard-failed (i.e. AggregateResult.Breakdown.HardFailed was true).
	HardFailCount int `json:"hard_fail_count,omitempty"`
}

// PromotionDecision is the pass/fail outcome of evaluating one
// [AcceptanceCriterion].
type PromotionDecision struct {
	// Split is the evaluation partition this decision covers.
	Split scenarios.ScenarioSplit `json:"split"`
	// Passed is true when the candidate met all requirements of the criterion.
	Passed bool `json:"passed"`
	// Reason is a human-readable explanation of the decision.
	Reason string `json:"reason,omitempty"`
}

// PromotionReport is the complete promotion evaluation result for one
// candidate. It is produced by [Evaluate].
type PromotionReport struct {
	// CandidateID is the stable identifier of the candidate under evaluation.
	CandidateID string `json:"candidate_id"`
	// Summaries holds the per-split evaluation summaries that were supplied to
	// [Evaluate].
	Summaries []SplitEvalSummary `json:"summaries,omitempty"`
	// Decisions holds the per-criterion promotion decisions in the same order
	// as the criteria supplied to [Evaluate].
	Decisions []PromotionDecision `json:"decisions,omitempty"`
	// Promoted is true when every criterion in Decisions passed.
	Promoted bool `json:"promoted"`
}

// SummarizeSplit computes a [SplitEvalSummary] from a slice of
// [objectives.AggregateResult] values for the named split. When results is
// empty, Count is zero and all score fields are zero.
func SummarizeSplit(split scenarios.ScenarioSplit, results []objectives.AggregateResult) SplitEvalSummary {
	s := SplitEvalSummary{
		Split:             split,
		Count:             len(results),
		MinAggregateScore: math.Inf(+1),
		MaxAggregateScore: math.Inf(-1),
	}
	if len(results) == 0 {
		s.MinAggregateScore = 0
		s.MaxAggregateScore = 0
		return s
	}
	var total float64
	for _, r := range results {
		total += r.AggregateScore
		if r.AggregateScore < s.MinAggregateScore {
			s.MinAggregateScore = r.AggregateScore
		}
		if r.AggregateScore > s.MaxAggregateScore {
			s.MaxAggregateScore = r.AggregateScore
		}
		if r.Breakdown.HardFailed {
			s.HardFailCount++
		}
	}
	s.MeanAggregateScore = total / float64(len(results))
	return s
}

// Decide evaluates one [AcceptanceCriterion] against a set of summaries and
// returns a [PromotionDecision].
//
// The summary for criterion.Split is located by scanning summaries in order.
// If no matching summary is found, the decision fails with an explanatory
// reason. When a matching summary is found:
//
//   - the candidate fails when Count == 0 (no scenarios were evaluated), or
//   - the candidate fails when the mean aggregate score is below
//     criterion.MinAggregateScore, or
//   - the candidate fails when criterion.RequireNoHardFail is true and
//     HardFailCount > 0.
//
// The decision passes when none of the above conditions hold.
func Decide(criterion AcceptanceCriterion, summaries []SplitEvalSummary) PromotionDecision {
	for _, s := range summaries {
		if s.Split != criterion.Split {
			continue
		}
		if s.Count == 0 {
			return PromotionDecision{
				Split:  criterion.Split,
				Passed: false,
				Reason: fmt.Sprintf("split %q: no scenarios evaluated", criterion.Split),
			}
		}
		if s.MeanAggregateScore < criterion.MinAggregateScore {
			return PromotionDecision{
				Split:  criterion.Split,
				Passed: false,
				Reason: fmt.Sprintf("split %q: mean aggregate score %.6g < required %.6g",
					criterion.Split, s.MeanAggregateScore, criterion.MinAggregateScore),
			}
		}
		if criterion.RequireNoHardFail && s.HardFailCount > 0 {
			return PromotionDecision{
				Split:  criterion.Split,
				Passed: false,
				Reason: fmt.Sprintf("split %q: %d hard failure(s) found, none allowed",
					criterion.Split, s.HardFailCount),
			}
		}
		return PromotionDecision{
			Split:  criterion.Split,
			Passed: true,
			Reason: fmt.Sprintf("split %q: mean aggregate score %.6g >= required %.6g",
				criterion.Split, s.MeanAggregateScore, criterion.MinAggregateScore),
		}
	}
	return PromotionDecision{
		Split:  criterion.Split,
		Passed: false,
		Reason: fmt.Sprintf("split %q: no evaluation summary found", criterion.Split),
	}
}

// Evaluate runs all criteria against summaries and returns a [PromotionReport]
// for the named candidate. The report's Promoted field is true only when every
// criterion passes.
//
// Decisions are recorded in the same order as criteria. Summaries are stored
// in the report as supplied.
func Evaluate(candidateID string, criteria []AcceptanceCriterion, summaries []SplitEvalSummary) PromotionReport {
	report := PromotionReport{
		CandidateID: candidateID,
		Summaries:   summaries,
		Promoted:    true,
	}
	for _, c := range criteria {
		d := Decide(c, summaries)
		report.Decisions = append(report.Decisions, d)
		if !d.Passed {
			report.Promoted = false
		}
	}
	return report
}
