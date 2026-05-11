// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package novelty provides reusable novelty/diversity infrastructure for
// applied-design workflows.
//
// Novelty search rewards candidates that behave differently from those already
// archived, rather than rewarding proximity to a fixed fitness target. This
// prevents premature convergence and encourages exploration.
//
// # Core types
//
//   - BehaviorVector: a slice of float64 values that encodes the observable
//     behavior of one candidate design.
//   - ArchiveEntry: pairs a BehaviorVector with an optional string label for
//     book-keeping (e.g. a candidate ID).
//   - DistanceFunc: a function that measures the dissimilarity between two
//     BehaviorVectors; the default is squared Euclidean distance.
//   - ArchiveConfig: controls how the archive grows and how k-nearest-neighbor
//     novelty is computed (K, maximum archive size).
//   - NoveltyResult: the outcome of one novelty-score query — the scalar
//     novelty score together with the individual neighbor distances used to
//     compute it.
//   - Archive: the live novelty store; supports Add and Score operations.
//
// # Novelty score calculation
//
// The novelty score of a query BehaviorVector is the mean distance to its K
// nearest neighbors in the archive. Smaller distances mean the behavior is
// common; larger distances mean it is novel.
//
// When the archive contains fewer than K entries the mean is computed over all
// available entries. When the archive is empty Score returns a NoveltyResult
// with a score of 0 and no neighbor distances.
//
// # Determinism
//
// Score is a pure function of the archive contents and the query vector.
// Inserting the same sequence of entries and then calling Score with the same
// query always returns the same NoveltyResult.
package novelty
