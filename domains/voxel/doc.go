// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package voxel defines the reusable core data model for voxel-domain
// applied-design workflows.
//
// It provides serializable candidate representation types and deterministic
// validation helpers for structural checks only. This package intentionally does
// not include geometry kernels, simulation, or pilot-specific evaluation logic.
//
// # Core types
//
//   - [VoxelProgram]: serializable candidate package (ID + design + spec).
//   - [VoxelDesign]: occupied cells within one bounded design volume.
//   - [DesignVolume]: axis-aligned voxel bounds plus forbidden/interface
//     regions.
//   - [Material]: reusable material metadata for a voxel candidate.
//   - [LoadCase]: reusable load-case metadata for later engineering pilots.
//   - [InterfaceRegion]: axis-aligned region reserved for interfaces or
//     forbidden occupancy.
//   - [VoxelSpec]: high-level domain metadata for a candidate/program.
//
// # Validation
//
// [DesignVolume.Validate] checks:
//
//   - empty or malformed design volumes,
//   - malformed/out-of-bounds forbidden and interface regions,
//   - overlap between forbidden and interface regions.
//
// [VoxelDesign.Validate] additionally checks:
//
//   - out-of-bounds occupied cells.
package voxel
