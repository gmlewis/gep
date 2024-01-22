// -*- compile-command: "go test ./..."; -*-
// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/model"
)

var result float64

func BenchmarkValidateFunc(b *testing.B) {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
	}
	e := model.New(funcs, functions.Float64, 30, 8, 4, 1, 0, "+", nil)
	var v float64
	for i := 0; i < b.N; i++ {
		v = validateFunc(e.Genomes[0])
	}
	result = v
}
