// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model

import (
	"testing"

	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/genome"
)

func TestGetBestHandlesAllNegativeScores(t *testing.T) {
	g := &Generation{
		Individuals: []*genome.Genome{
			{},
			{},
			{},
		},
	}
	g.ScoringFunc = func(gn *genome.Genome) float64 {
		switch gn {
		case g.Individuals[0]:
			return -10
		case g.Individuals[1]:
			return -3
		default:
			return -7
		}
	}

	best := g.getBest()
	if best != g.Individuals[1] {
		t.Fatalf("getBest() returned %p, want %p", best, g.Individuals[1])
	}
	if best.Score != -3 {
		t.Fatalf("best score = %v, want -3", best.Score)
	}
}

func TestMaxArity(t *testing.T) {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 2},
		{Symbol: "*", Weight: 3},
		{Symbol: "/", Weight: 4},
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
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
	}
	e := New(funcs, functions.Float64, 30, 8, 4, 1, 0, "+", nil, false)
	for i := 0; i < b.N; i++ {
		e.replication()
	}
}

func BenchmarkMutation(b *testing.B) {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
	}
	e := New(funcs, functions.Float64, 30, 8, 4, 1, 0, "+", nil, false)
	for i := 0; i < b.N; i++ {
		e.mutation()
	}
}

func TestReplication(t *testing.T) {
	g := &Generation{
		Individuals: []*genome.Genome{
			{Score: -1000},
			{Score: -500},
			{Score: -100},
			{Score: -50},
			{Score: -10},
			{Score: -5},
			{Score: -1},
			{Score: 1},
			{Score: 5},
		},
	}

	before := len(g.Individuals)
	g.replication()
	got := len(g.Individuals)
	if want := before; got != want {
		t.Errorf("replication = %v individuals, want %v", got, want)
	}
}
