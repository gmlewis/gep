// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package regression is a cross-domain discovery regression suite that proves
// the applied-design pipeline contract across all three domain pilots:
//
//   - experiments/circuit/half_adder  (boolean logic, SPICE/Verilog export)
//   - experiments/voxel/bracket       (voxel geometry, JSON/OBJ export)
//   - experiments/control/mass_spring_damper (plant simulation, policy export)
//
// Each pilot exercises the full flow:
//
//	evolve → decode → constrain → validate → promote → export → checkpoint
//
// The regression suite verifies that all three pilots complete that flow
// successfully with a fixed seed, using the shared design infrastructure
// (design.RunManifest, promotion.PromotionReport, checkpoint.Snapshot).
//
// Run the regression suite with:
//
//	go test ./experiments/regression/...
package regression
