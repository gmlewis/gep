// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package constraints provides typed constraint interfaces and validation
// reporting for applied-design workflows.
//
// It defines the three fundamental constraint behaviors — reject, repair, and
// penalize — as well as a ValidationReport aggregate that records which
// constraints ran, what decisions they made, and what repairs or penalties were
// applied.
//
// Constraints are evaluated in order. The first rejection terminates the run
// immediately. Repairs are chained so each subsequent constraint sees the
// updated candidate. Penalties accumulate across all constraints that ran.
package constraints
