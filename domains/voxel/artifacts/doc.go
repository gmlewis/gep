// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package artifacts emits deterministic serialized artifacts for voxel-domain
// candidates.
//
// The emitters in this package target the reusable [voxel.VoxelProgram] model
// and intentionally stay lightweight: they validate the program's design, then
// render stable output without invoking any external geometry kernel.
//
// Emitters:
//
//   - [JSON]: canonical indented JSON.
//   - [OBJ]: Wavefront OBJ mesh with one unit cube per occupied voxel cell.
//   - [Summary]: concise human-readable plain-text overview.
package artifacts
