// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package vectorIntNodes

import (
	"reflect"
	"testing"
)

func TestPlus(t *testing.T) {
	v1 := VectorInt{2, 3, 4, 5}
	v2 := VectorInt{6, 7, 8, 9}
	v3 := VectorInt{10, 11, 12, 13}
	v4 := VectorInt{14, 15, 16, 17}
	x := []VectorInt{v1, v2, v3, v4}
	want := VectorInt{8, 10, 12, 14}
	if g := VectorIntFuncs["+"].VectorIntFunction(x); !reflect.DeepEqual(g, want) {
		t.Errorf("VectorInt[+](%v,%v,%v,%v) = %v, want %v", v1, v2, v3, v4, g, want)
	}
}

func TestMinus(t *testing.T) {
	v1 := VectorInt{2, 3, 4, 5}
	v2 := VectorInt{6, 7, 8, 9}
	v3 := VectorInt{10, 11, 12, 13}
	v4 := VectorInt{14, 15, 16, 17}
	x := []VectorInt{v1, v2, v3, v4}
	want := VectorInt{-4, -4, -4, -4}
	if g := VectorIntFuncs["-"].VectorIntFunction(x); !reflect.DeepEqual(g, want) {
		t.Errorf("VectorInt[-](%v,%v,%v,%v) = %v, want %v", v1, v2, v3, v4, g, want)
	}
}

var result VectorInt

func runBenchmark(b *testing.B, sym string) {
	v1 := VectorInt{2, 3, 4, 5}
	v2 := VectorInt{6, 7, 8, 9}
	v3 := VectorInt{10, 11, 12, 13}
	v4 := VectorInt{14, 15, 16, 17}
	x := []VectorInt{v1, v2, v3, v4}
	var v VectorInt
	f := VectorIntFuncs[sym]
	for i := 0; i < b.N; i++ {
		v = f.VectorIntFunction(x)
	}
	result = v
}
