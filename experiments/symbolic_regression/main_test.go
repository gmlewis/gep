// -*- compile-command: "go test ./..."; -*-
// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"testing"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/evolution"
	"github.com/gmlewis/gep/v2/functions"
	mathNodes "github.com/gmlewis/gep/v2/functions/math_nodes"
)

var result float64

func TestValidateFunc_TypedGenome(t *testing.T) {
	link, err := core.NewLinkFunc[float64]("id", func(v []float64) float64 { return v[0] })
	if err != nil {
		t.Fatalf("core.NewLinkFunc error: %v", err)
	}
	sym, err := core.NewTerminalSymbol[float64](0)
	if err != nil {
		t.Fatalf("core.NewTerminalSymbol error: %v", err)
	}
	g := core.Genome[float64]{
		Genes: []core.Gene[float64]{
			{Symbols: []core.Symbol[float64]{sym}},
		},
		Link: link,
	}
	score := validateFunc(g)
	if math.IsNaN(score) || math.IsInf(score, 0) {
		t.Fatalf("validateFunc returned non-finite score: %v", score)
	}
}

func BenchmarkValidateFunc(b *testing.B) {
	funcs := []string{
		"+",
		"-",
		"*",
	}
	fm := make(functions.FuncMap, len(funcs))
	for _, sym := range funcs {
		fn, ok := mathNodes.Math[sym]
		if !ok {
			b.Fatalf("unsupported math function %q", sym)
		}
		fm[sym] = fn
	}
	cat, err := mathNodes.CatalogFrom(fm)
	if err != nil {
		b.Fatalf("CatalogFrom error: %v", err)
	}
	linkNode, ok := mathNodes.Math["+"]
	if !ok {
		b.Fatal(`link function "+" not found`)
	}
	link, err := core.NewLinkFunc[float64]("+", linkNode.Float64Function)
	if err != nil {
		b.Fatalf("core.NewLinkFunc error: %v", err)
	}
	e, err := evolution.New(cat, 30, 8, 4, 1, 0, link, nil)
	if err != nil {
		b.Fatalf("evolution.New(...) error: %v", err)
	}
	var v float64
	for i := 0; i < b.N; i++ {
		v = validateFunc(e.Individuals[0].Genome)
	}
	result = v
}
