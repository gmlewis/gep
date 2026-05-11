// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package artifacts emits deterministic serialized artifacts for circuit-domain
// candidates.
//
// The emitters in this package target the reusable [circuit.CircuitProgram]
// model and intentionally stay lightweight: they validate the program's graph
// structure, then render stable text/JSON output without invoking any external
// simulator or synthesis tool.
package artifacts
