// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package objectives provides shared multi-objective scoring for applied-design
// workflows.
//
// It defines reusable types for combining hard constraints and soft objectives
// into a single score breakdown that domain pilots can share:
//
//   - ObjectiveDef: names a scoring objective, assigns a weight, and marks
//     whether it is soft (weighted contribution) or hard (failure gate).
//   - WeightedContribution: records the raw score, weight, and weighted
//     contribution of one objective.
//   - ScoreBreakdown: collects all per-objective contributions together with
//     a hard-failed flag and an accumulated constraint penalty.
//   - AggregateResult: combines the breakdown with a single aggregate score
//     suitable for ranking candidates.
//
// # Hard-fail gating
//
// A candidate is hard-failed and receives an AggregateScore of zero when
// either of the following is true:
//
//   - the rejected flag passed to Score is true (e.g. from a
//     constraints.ValidationReport.Rejected field), or
//   - any objective defined with Kind Hard has a raw score <= 0.
//
// # Deterministic ordering
//
// Score is a pure function: the same inputs always produce identical
// AggregateResult values. Contributions in the breakdown are recorded in
// definition order. The Less helper provides a deterministic tie-breaking
// comparator suitable for sorting slices of AggregateResult.
package objectives
