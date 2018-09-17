// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model

import (
	"testing"

	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
)

func TestMaxArity(t *testing.T) {
	funcs := []gene.FuncWeight{
		{"+", 1},
		{"-", 2},
		{"*", 3},
		{"/", 4},
	}
	if g, w := maxArity(funcs, functions.Float64), 2; g != w {
		t.Errorf("maxArity(%v, functions.Float64) = %v, want %v", funcs, g, w)
	}
	funcs = append(funcs, gene.FuncWeight{
		Symbol: "LT3A",
		Weight: 1,
	})
	if g, w := maxArity(funcs, functions.Float64), 3; g != w {
		t.Errorf("maxArity(%v, functions.Float64) = %v, want %v", funcs, g, w)
	}
}

func BenchmarkReplication(b *testing.B) {
	funcs := []gene.FuncWeight{
		{"+", 1},
		{"-", 1},
		{"*", 1},
	}
	e := New(funcs, functions.Float64, 30, 8, 4, 1, 0, "+", nil)
	for i := 0; i < b.N; i++ {
		e.replication()
	}
}

func BenchmarkMutation(b *testing.B) {
	funcs := []gene.FuncWeight{
		{"+", 1},
		{"-", 1},
		{"*", 1},
	}
	e := New(funcs, functions.Float64, 30, 8, 4, 1, 0, "+", nil)
	for i := 0; i < b.N; i++ {
		e.mutation()
	}
}
