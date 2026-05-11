// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package promotion provides shared promotion and acceptance-criteria
// evaluation for applied-design workflows.
//
// It formalizes what it means for a candidate to be "good enough" to promote
// after train/validation evaluation. The key workflow is:
//
//  1. Collect [objectives.AggregateResult] values for each evaluation split.
//  2. Summarize each split with [SummarizeSplit].
//  3. Define per-split thresholds using [AcceptanceCriterion] values.
//  4. Evaluate all criteria against the summaries with [Evaluate], which
//     returns a [PromotionReport] that records each per-split decision and
//     the final promoted flag.
//
// # Core types
//
//   - [AcceptanceCriterion]: names a split, a minimum mean aggregate score,
//     and an optional hard-fail gate.
//   - [SplitEvalSummary]: a per-split aggregate summary — count, mean, min,
//     max aggregate score, and hard-fail count.
//   - [PromotionDecision]: the pass/fail outcome of one criterion together
//     with a human-readable reason string.
//   - [PromotionReport]: the complete promotion result — all summaries, all
//     per-criterion decisions, and the final [PromotionReport.Promoted] flag.
//
// # Helpers
//
//   - [SummarizeSplit] computes a [SplitEvalSummary] from a slice of
//     [objectives.AggregateResult] values.
//   - [Decide] evaluates one [AcceptanceCriterion] against a set of summaries
//     and returns a [PromotionDecision].
//   - [Evaluate] runs all criteria, collects per-split decisions, and sets the
//     [PromotionReport.Promoted] flag to true only when every criterion passes.
//
// # Determinism
//
// [SummarizeSplit], [Decide], and [Evaluate] are pure functions: identical
// inputs always produce identical outputs. Decisions in a [PromotionReport]
// are stored in the same order as the supplied criteria slice.
package promotion
