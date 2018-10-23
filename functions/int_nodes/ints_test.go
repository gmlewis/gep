// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package intNodes

import (
	"testing"
)

func TestPlus(t *testing.T) {
	if g, w := Int["+"].IntFunction([]int{2, 3, 4, 5}), 5; g != w {
		t.Errorf("Int[+]([2,3,4,5]) = %v, want %v", g, w)
	}
}

var result int

func runBenchmark(b *testing.B, sym string) {
	x := []int{2, 3, 4, 5}
	var v int
	f := Int[sym]
	for i := 0; i < b.N; i++ {
		v = f.IntFunction(x)
	}
	result = v
}
