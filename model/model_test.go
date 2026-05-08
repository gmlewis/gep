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

func TestMaxArity_UnknownFuncTypeReturnsZero(t *testing.T) {
	if got := maxArity(nil, functions.FuncType(999)); got != 0 {
		t.Fatalf("maxArity(nil, unknown) = %v, want 0", got)
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

func TestGetBestHandlesAllNegativeScores(t *testing.T) {
	g1 := &genome.Genome{}
	g2 := &genome.Genome{}
	g3 := &genome.Genome{}
	scores := map[*genome.Genome]float64{
		g1: -100,
		g2: -1,
		g3: -50,
	}
	g := &Generation{
		Individuals: []*genome.Genome{g1, g2, g3},
		ScoringFunc: func(gn *genome.Genome) float64 {
			return scores[gn]
		},
	}

	got := g.getBest()
	if got != g2 {
		t.Fatalf("getBest() = %p (score=%v), want %p (score=%v)", got, scores[got], g2, scores[g2])
	}
	if scores[got] != -1 {
		t.Fatalf("getBest() score = %v, want -1", scores[got])
	}
}

func TestSingleCrossover_MismatchedGenesDoesNotMutate(t *testing.T) {
	gene1 := &gene.Gene{Symbols: []string{"d0"}, HeadSize: 1}
	gene2 := &gene.Gene{Symbols: []string{"d0", "d1"}, HeadSize: 2}
	g := &Generation{
		Individuals: []*genome.Genome{
			{Genes: []*gene.Gene{gene1}},
			{Genes: []*gene.Gene{gene2}},
		},
	}
	g.singleCrossover(0, 1)
	if got, want := g.Individuals[0].Genes[0].Symbols, []string{"d0"}; len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("genome[0] symbols changed: got %v, want %v", got, want)
	}
	if got, want := g.Individuals[1].Genes[0].Symbols, []string{"d0", "d1"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("genome[1] symbols changed: got %v, want %v", got, want)
	}
}
