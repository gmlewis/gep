// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package objectives

// ObjectiveKind indicates whether an objective contributes to soft weighted
// aggregation or acts as a hard-failure gate.
type ObjectiveKind int

const (
	// Soft means the objective contributes its weighted score to the aggregate.
	Soft ObjectiveKind = iota
	// Hard means the objective gates the entire result: a zero or negative raw
	// score forces HardFailed to true and AggregateScore to zero.
	Hard
)

// String returns a human-readable name for the kind.
func (k ObjectiveKind) String() string {
	switch k {
	case Soft:
		return "soft"
	case Hard:
		return "hard"
	default:
		return "unknown"
	}
}

// ObjectiveDef defines one named scoring objective.
type ObjectiveDef struct {
	// Name is the stable identifier for this objective.
	Name string `json:"name"`
	// Weight scales the raw score contribution during soft aggregation.
	Weight float64 `json:"weight"`
	// Kind determines whether this is a soft contribution or a hard gate.
	Kind ObjectiveKind `json:"kind"`
}

// WeightedContribution records the contribution of one objective to the
// aggregate score.
type WeightedContribution struct {
	// Name matches the ObjectiveDef.Name.
	Name string `json:"name"`
	// RawScore is the unscaled score for this objective.
	RawScore float64 `json:"raw_score"`
	// Weight is the multiplier used.
	Weight float64 `json:"weight"`
	// WeightedScore is RawScore * Weight.
	WeightedScore float64 `json:"weighted_score"`
}

// ScoreBreakdown holds the per-objective contributions and summary flags for
// one candidate.
type ScoreBreakdown struct {
	// Contributions contains one entry per ObjectiveDef, in definition order.
	Contributions []WeightedContribution `json:"contributions,omitempty"`
	// HardFailed is true when any hard gate fired or the candidate was rejected
	// by a validation constraint.
	HardFailed bool `json:"hard_failed,omitempty"`
	// TotalPenalty carries the accumulated soft penalty from constraint checks.
	TotalPenalty float64 `json:"total_penalty,omitempty"`
}

// AggregateResult is the final scoring outcome for one candidate.
type AggregateResult struct {
	// Breakdown holds the per-objective contributions and hard-fail flag.
	Breakdown ScoreBreakdown `json:"breakdown"`
	// AggregateScore is the final scalar score. It is zero when HardFailed is true.
	AggregateScore float64 `json:"aggregate_score"`
}

// Score computes the AggregateResult for a candidate.
//
// defs lists the objectives in evaluation order; the same order is preserved
// in ScoreBreakdown.Contributions so the breakdown is always deterministic.
//
// rawScores maps each objective name to its unscaled score. Missing names
// default to zero.
//
// rejected must be set to true when the candidate failed a hard validation
// constraint (e.g. constraints.ValidationReport.Rejected). totalPenalty
// carries the accumulated soft penalty from those checks.
//
// Hard-fail gating rules (AggregateScore forced to zero):
//   - rejected is true, or
//   - any objective with Kind Hard has a raw score <= 0.
func Score(defs []ObjectiveDef, rawScores map[string]float64, rejected bool, totalPenalty float64) AggregateResult {
	breakdown := ScoreBreakdown{
		TotalPenalty: totalPenalty,
	}
	hardFailed := rejected

	for _, def := range defs {
		raw := rawScores[def.Name]
		weighted := raw * def.Weight
		if def.Kind == Hard && raw <= 0 {
			hardFailed = true
		}
		breakdown.Contributions = append(breakdown.Contributions, WeightedContribution{
			Name:          def.Name,
			RawScore:      raw,
			Weight:        def.Weight,
			WeightedScore: weighted,
		})
	}

	breakdown.HardFailed = hardFailed
	if hardFailed {
		return AggregateResult{Breakdown: breakdown}
	}

	var aggregate float64
	for _, c := range breakdown.Contributions {
		aggregate += c.WeightedScore
	}
	aggregate -= totalPenalty

	return AggregateResult{
		Breakdown:      breakdown,
		AggregateScore: aggregate,
	}
}

// Less reports whether a should be ordered before b when ranking candidates
// from best to worst.
//
// Higher AggregateScore ranks better. When scores are equal the tie is broken
// by comparing WeightedScore values of contributions in definition order,
// using the first differing value. If all contributions are equal the result
// is false (a and b are equivalent for ranking purposes).
func Less(a, b AggregateResult) bool {
	if a.AggregateScore != b.AggregateScore {
		return a.AggregateScore > b.AggregateScore
	}
	n := len(a.Breakdown.Contributions)
	if len(b.Breakdown.Contributions) < n {
		n = len(b.Breakdown.Contributions)
	}
	for i := range n {
		ca := a.Breakdown.Contributions[i]
		cb := b.Breakdown.Contributions[i]
		if ca.WeightedScore != cb.WeightedScore {
			return ca.WeightedScore > cb.WeightedScore
		}
	}
	return false
}
