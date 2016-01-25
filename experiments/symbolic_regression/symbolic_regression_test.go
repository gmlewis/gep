// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	mn "github.com/gmlewis/gep/functions/math_nodes"
	"github.com/gmlewis/gep/gene"
	"github.com/gmlewis/gep/model"
)

var result float64

func BenchmarkValidateFunc(b *testing.B) {
	funcs := []gene.FuncWeight{
		{"+", 1},
		{"-", 1},
		{"*", 1},
	}
	e := model.New(funcs, mn.Math, 30, 8, 4, 1, 0, "+", nil)
	var v float64
	for i := 0; i < b.N; i++ {
		v = validateFunc(e.Genomes[0])
	}
	result = v
}
