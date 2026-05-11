// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package scenarios provides shared train/validation/test scenario-set
// handling for applied-design workflows.
//
// A scenario is a single reproducible evaluation input identified by a stable
// [ScenarioID]. Scenarios are grouped into named [ScenarioSet] values and each
// scenario is assigned to exactly one evaluation partition using the
// [ScenarioSplit] type. A [ScenarioRegistry] aggregates one or more sets and
// validates that no scenario ID is assigned to more than one split.
//
// # Core types
//
//   - [ScenarioID]: a stable string identifier for one scenario.
//   - [ScenarioSplit]: the name of an evaluation partition. Built-in constants
//     [Train], [Validation], and [Test] cover the standard three-way split.
//   - [Scenario]: one scenario record — an ID, a split assignment, optional
//     tags, and arbitrary JSON-serializable parameters.
//   - [ScenarioSet]: a named, ordered collection of [Scenario] records loaded
//     from one source (e.g. a testdata fixture file).
//   - [ScenarioRegistry]: aggregates one or more [ScenarioSet] values and
//     exposes [ScenarioRegistry.Validate], [ScenarioRegistry.ByID], and
//     [ScenarioRegistry.BySplit] helpers.
//
// # Validation
//
// A registry is valid when each [ScenarioID] appears in at most one split
// across all sets. [ScenarioRegistry.Validate] returns a descriptive error for
// the first conflict it finds, whether the conflict is across different splits
// or a repeated appearance within the same split.
//
// # Fixture loading
//
// [LoadScenarioSet] and [LoadScenarioSetFile] decode a [ScenarioSet] from JSON
// (reader or file path). Fixture JSON files may be placed in a testdata/
// directory alongside the domain package that uses them.
//
// # Determinism
//
// [ScenarioRegistry.ByID] and [ScenarioRegistry.BySplit] iterate in
// declaration order, so identical registry contents always yield identical
// results.
package scenarios
