// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package eval provides concurrency-safe batched evaluation primitives for
// applied-design workflows.
//
// It defines typed request/result contracts plus a bounded worker dispatcher
// that future domain evaluators can reuse.
package eval
