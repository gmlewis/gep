// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package constraints

// Decision is the outcome of a single constraint check.
type Decision int

const (
	// Pass means the constraint was satisfied and no action was needed.
	Pass Decision = iota
	// Repair means the constraint modified the candidate to satisfy the rule.
	Repair
	// Penalize means the constraint added a soft penalty to the candidate's score.
	Penalize
	// Reject means the candidate violated a hard constraint and must be discarded.
	Reject
)

// String returns a human-readable name for the decision.
func (d Decision) String() string {
	switch d {
	case Pass:
		return "pass"
	case Repair:
		return "repair"
	case Penalize:
		return "penalize"
	case Reject:
		return "reject"
	default:
		return "unknown"
	}
}

// ConstraintResult records the outcome of one constraint check.
type ConstraintResult struct {
	// Name is the constraint's identifying name.
	Name string
	// Decision is what the constraint decided.
	Decision Decision
	// Penalty is the soft penalty added by a Penalize decision (zero otherwise).
	Penalty float64
	// Message is an optional human-readable explanation of the decision.
	Message string
}

// Constraint evaluates one candidate and returns the (possibly updated)
// candidate and a result describing what it decided.
//
// For a Repair decision the returned candidate reflects the repaired value.
// For Pass, Penalize, or Reject decisions the returned candidate is unchanged.
type Constraint[T any] interface {
	// Name returns the constraint's stable identifier used in reports.
	Name() string
	// Check evaluates the candidate and returns the updated candidate together
	// with a ConstraintResult describing the decision.
	Check(candidate T) (T, ConstraintResult)
}

// ValidationReport records the aggregate result of running a set of constraints
// against one candidate.
type ValidationReport[T any] struct {
	// Candidate is the final (possibly repaired) candidate after all constraints
	// that ran have been applied.
	Candidate T
	// Results contains one entry per constraint that ran, in order.
	Results []ConstraintResult
	// TotalPenalty is the sum of penalties from all Penalize decisions.
	TotalPenalty float64
	// Rejected is true if any constraint issued a Reject decision.
	Rejected bool
}

// Validate runs each constraint in order against the candidate and returns a
// ValidationReport.
//
// Behavior:
//   - On Reject: the run stops immediately and Rejected is set to true.
//   - On Repair: the updated candidate is forwarded to the next constraint.
//   - On Penalize: the penalty is added to TotalPenalty.
//   - On Pass: the candidate is forwarded unchanged.
func Validate[T any](candidate T, cs []Constraint[T]) ValidationReport[T] {
	report := ValidationReport[T]{Candidate: candidate}
	for _, c := range cs {
		updated, result := c.Check(report.Candidate)
		report.Results = append(report.Results, result)
		switch result.Decision {
		case Repair:
			report.Candidate = updated
		case Penalize:
			report.TotalPenalty += result.Penalty
		case Reject:
			report.Rejected = true
			return report
		}
	}
	return report
}
