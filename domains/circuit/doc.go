// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package circuit defines the reusable core data model for circuit-domain
// applied-design workflows.
//
// It provides serializable candidate representation types and deterministic
// validation helpers for structural checks only. This package intentionally does
// not include simulator execution or pilot logic.
//
// # Core types
//
//   - [NodeID]: stable string identifier for one graph node.
//   - [Port]: node/port reference used by components and wiring.
//   - [Component]: one typed graph node with named input/output ports.
//   - [CircuitGraph]: ordered component collection for one candidate topology.
//   - [CircuitProgram]: serializable candidate package (ID + graph + spec +
//     constraints).
//   - [CircuitSpec]: high-level domain metadata for a candidate/program.
//   - [CircuitConstraint]: serializable declarative constraint record.
//
// # Validation
//
// [CircuitGraph.Validate] checks:
//
//   - duplicate component node IDs,
//   - missing component names/types,
//   - illegal port references (ports referencing unknown node IDs).
package circuit
