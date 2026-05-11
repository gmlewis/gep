// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package constraints

import (
	"reflect"
	"testing"
)

// --- deterministic fake constraints for testing ---

// passConstraint always passes without modifying the candidate.
type passConstraint struct {
	name string
}

func (c passConstraint) Name() string { return c.name }
func (c passConstraint) Check(v int) (int, ConstraintResult) {
	return v, ConstraintResult{Name: c.name, Decision: Pass, Message: "ok"}
}

// rejectConstraint always rejects.
type rejectConstraint struct {
	name string
}

func (c rejectConstraint) Name() string { return c.name }
func (c rejectConstraint) Check(v int) (int, ConstraintResult) {
	return v, ConstraintResult{Name: c.name, Decision: Reject, Message: "rejected"}
}

// repairConstraint adds delta to the candidate.
type repairConstraint struct {
	name  string
	delta int
}

func (c repairConstraint) Name() string { return c.name }
func (c repairConstraint) Check(v int) (int, ConstraintResult) {
	updated := v + c.delta
	return updated, ConstraintResult{Name: c.name, Decision: Repair, Message: "repaired"}
}

// penalizeConstraint adds a fixed penalty.
type penalizeConstraint struct {
	name    string
	penalty float64
}

func (c penalizeConstraint) Name() string { return c.name }
func (c penalizeConstraint) Check(v int) (int, ConstraintResult) {
	return v, ConstraintResult{Name: c.name, Decision: Penalize, Penalty: c.penalty, Message: "penalized"}
}

// --- tests ---

func TestValidateStopOnReject(t *testing.T) {
	cs := []Constraint[int]{
		passConstraint{name: "allow-all"},
		rejectConstraint{name: "hard-reject"},
		passConstraint{name: "should-not-run"},
	}

	report := Validate(42, cs)

	if !report.Rejected {
		t.Fatal("expected Rejected=true")
	}
	// Only the first two constraints should have run.
	if got, want := len(report.Results), 2; got != want {
		t.Fatalf("results count = %d, want %d", got, want)
	}
	if report.Results[1].Decision != Reject {
		t.Fatalf("second result decision = %v, want Reject", report.Results[1].Decision)
	}
	// Candidate must be unchanged.
	if report.Candidate != 42 {
		t.Fatalf("candidate = %d after reject, want 42", report.Candidate)
	}
}

func TestValidateRepairChaining(t *testing.T) {
	cs := []Constraint[int]{
		repairConstraint{name: "add-10", delta: 10},
		repairConstraint{name: "add-5", delta: 5},
		passConstraint{name: "final-pass"},
	}

	report := Validate(0, cs)

	if report.Rejected {
		t.Fatal("expected Rejected=false")
	}
	// 0 + 10 + 5 = 15
	if report.Candidate != 15 {
		t.Fatalf("candidate = %d after repair chain, want 15", report.Candidate)
	}
	if got, want := len(report.Results), 3; got != want {
		t.Fatalf("results count = %d, want %d", got, want)
	}
	if report.Results[0].Decision != Repair {
		t.Fatalf("result[0].Decision = %v, want Repair", report.Results[0].Decision)
	}
	if report.Results[1].Decision != Repair {
		t.Fatalf("result[1].Decision = %v, want Repair", report.Results[1].Decision)
	}
}

func TestValidatePenaltyAccumulation(t *testing.T) {
	cs := []Constraint[int]{
		penalizeConstraint{name: "cost-a", penalty: 2.5},
		penalizeConstraint{name: "cost-b", penalty: 1.0},
		passConstraint{name: "ok"},
		penalizeConstraint{name: "cost-c", penalty: 0.5},
	}

	report := Validate(7, cs)

	if report.Rejected {
		t.Fatal("expected Rejected=false")
	}
	// Candidate must be unchanged by penalize constraints.
	if report.Candidate != 7 {
		t.Fatalf("candidate = %d, want 7", report.Candidate)
	}
	const wantPenalty = 4.0
	if report.TotalPenalty != wantPenalty {
		t.Fatalf("TotalPenalty = %v, want %v", report.TotalPenalty, wantPenalty)
	}
	if got, want := len(report.Results), 4; got != want {
		t.Fatalf("results count = %d, want %d", got, want)
	}
}

func TestValidateRejectSkipsPenalties(t *testing.T) {
	cs := []Constraint[int]{
		penalizeConstraint{name: "cost-a", penalty: 3.0},
		rejectConstraint{name: "hard-reject"},
		penalizeConstraint{name: "cost-b", penalty: 99.0},
	}

	report := Validate(1, cs)

	if !report.Rejected {
		t.Fatal("expected Rejected=true")
	}
	// Only the first two constraints ran; penalty from cost-b must not appear.
	const wantPenalty = 3.0
	if report.TotalPenalty != wantPenalty {
		t.Fatalf("TotalPenalty = %v, want %v", report.TotalPenalty, wantPenalty)
	}
	if got, want := len(report.Results), 2; got != want {
		t.Fatalf("results count = %d, want %d", got, want)
	}
}

func TestValidateRepairThenReject(t *testing.T) {
	cs := []Constraint[int]{
		repairConstraint{name: "add-1", delta: 1},
		rejectConstraint{name: "still-bad"},
		repairConstraint{name: "should-not-run", delta: 100},
	}

	report := Validate(0, cs)

	if !report.Rejected {
		t.Fatal("expected Rejected=true")
	}
	// Repair ran before reject; candidate reflects the repair.
	if report.Candidate != 1 {
		t.Fatalf("candidate = %d, want 1", report.Candidate)
	}
	if got, want := len(report.Results), 2; got != want {
		t.Fatalf("results count = %d, want %d", got, want)
	}
}

func TestValidateEmptyConstraints(t *testing.T) {
	report := Validate(42, []Constraint[int]{})

	if report.Rejected {
		t.Fatal("expected Rejected=false for empty constraints")
	}
	if report.Candidate != 42 {
		t.Fatalf("candidate = %d, want 42", report.Candidate)
	}
	if len(report.Results) != 0 {
		t.Fatalf("results count = %d, want 0", len(report.Results))
	}
	if report.TotalPenalty != 0 {
		t.Fatalf("TotalPenalty = %v, want 0", report.TotalPenalty)
	}
}

func TestValidateIsDeterministicForFixedInput(t *testing.T) {
	cs := []Constraint[int]{
		repairConstraint{name: "add-3", delta: 3},
		penalizeConstraint{name: "soft-cost", penalty: 1.5},
		passConstraint{name: "final"},
	}

	got1 := Validate(10, cs)
	got2 := Validate(10, cs)

	if !reflect.DeepEqual(got1, got2) {
		t.Fatalf("determinism mismatch:\n got1=%#v\n got2=%#v", got1, got2)
	}
	if got1.Candidate != 13 {
		t.Fatalf("candidate = %d, want 13", got1.Candidate)
	}
	if got1.TotalPenalty != 1.5 {
		t.Fatalf("TotalPenalty = %v, want 1.5", got1.TotalPenalty)
	}
}

func TestDecisionString(t *testing.T) {
	cases := []struct {
		d    Decision
		want string
	}{
		{Pass, "pass"},
		{Repair, "repair"},
		{Penalize, "penalize"},
		{Reject, "reject"},
		{Decision(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.d.String(); got != tc.want {
			t.Errorf("Decision(%d).String() = %q, want %q", tc.d, got, tc.want)
		}
	}
}
