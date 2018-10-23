// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package mathNodes

import (
	"testing"
)

func TestPlus(t *testing.T) {
	if g, w := Math["+"].Float64Function([]float64{2, 3, 4, 5}), 5.0; g != w {
		t.Errorf("Math[+]([2,3,4,5]) = %v, want %v", g, w)
	}
}

var result float64

func runBenchmark(b *testing.B, sym string) {
	x := []float64{2, 3, 4, 5}
	var v float64
	f := Math[sym]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(x)
	}
	result = v
}
