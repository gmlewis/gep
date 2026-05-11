// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package checkpoint provides checkpoint and replay support for applied-design
// runs.
//
// A checkpoint captures enough state to pause, inspect, and replay a run
// without re-inventing ad hoc storage in each domain. Checkpoints are written
// as versioned, indented JSON files so they are both machine-readable and
// human-inspectable.
//
// # Snapshot types
//
// A [Snapshot] combines all durable run state into one value:
//
//   - [design.RunManifest]: run-level configuration, seeds, artifact refs, and
//     scenario-split summaries.
//   - [EliteRecord]: metadata for each elite candidate retained at the time of
//     the checkpoint, including its aggregate score and behavior vector.
//   - [AggregateSnapshot]: aggregate scores and per-objective breakdowns.
//   - [NoveltySnapshot]: the serializable state of the novelty archive so that
//     diversity metrics can be restored alongside other run metadata.
//   - ArtifactRefs: the subset of artifact references associated with the
//     checkpoint.
//
// # Versioning
//
// Every saved [Snapshot] includes a SchemaVersion field. The current version
// is 1. Load helpers validate this field and return an error for unknown
// versions. This ensures that stale on-disk checkpoints are detected early
// rather than silently misread.
//
// # Save and load
//
// Use [Save] / [SaveFile] to write a snapshot and [Load] / [LoadFile] to
// restore it. The helpers round-trip deterministically: a snapshot written
// with Save and read back with Load is identical to the original value when
// compared with reflect.DeepEqual.
//
// # Replay helpers
//
// [ReplayManifest] extracts the [design.RunManifest] from a saved snapshot so
// that a later run can reuse it even if the original evaluator is unavailable.
// [ReplayElites] returns the elite records, and [ReplayNoveltyEntries] returns
// the archived behavior vectors so that a resumed run can seed its novelty
// archive without re-evaluating old candidates.
package checkpoint
