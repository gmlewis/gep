// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package boolNodes

import (
	"testing"
)

var nandTests = []struct {
	in  []bool
	out bool
}{
	{[]bool{false, false}, true},
	{[]bool{false, true}, true},
	{[]bool{true, false}, true},
	{[]bool{true, true}, false},
}

func TestNand(t *testing.T) {
	for _, test := range nandTests {
		if g, w := BoolAllGates["Nand"].BoolFunction(test.in), test.out; g != w {
			t.Errorf("BoolAllGates[Nand]([%v,%v]) = %v, want %v", test.in[0], test.in[1], g, w)
		}
	}
}

var result bool

func runBenchmark(b *testing.B, sym string) {
	x := []bool{false, true, true, false}
	var v bool
	f := BoolAllGates[sym]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(x)
	}
	result = v
}
