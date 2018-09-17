// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"math"
	"reflect"
	"testing"

	"github.com/gmlewis/gep/v2/functions"
)

var nandTests = []struct {
	in   []bool
	want bool
}{
	{[]bool{false, false}, true},
	{[]bool{false, true}, true},
	{[]bool{true, false}, true},
	{[]bool{true, true}, false},
}

func validateNand(t *testing.T, g *Gene) {
	for i, n := range nandTests {
		got := g.EvalBool(n.in)
		if got != n.want {
			t.Errorf("%v: nand.EvalBool(%#v, BoolAllGates) => %v, want %v", i, n.in, got, n.want)
		}
	}
}

func TestNand(t *testing.T) {
	nand := New("Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0", functions.Bool)
	validateNand(t, nand)
	w := map[string]int{
		"And": 3,
		"Not": 2,
		"Or":  2,
		"d0":  2,
		"d1":  4,
	}
	if !reflect.DeepEqual(nand.SymbolMap, w) {
		t.Errorf("Gene %q SymbolMap=%v, want %v", nand, nand.SymbolMap, w)
	}
	nand = New("Or.And.Not.d0.Not.And.Or.d0.d0.d1.d1.d0.d1.d1.d1", functions.Bool)
	validateNand(t, nand)
	w = map[string]int{
		"And": 2,
		"Not": 2,
		"Or":  2,
		"d0":  3,
		"d1":  2,
	}
	if !reflect.DeepEqual(nand.SymbolMap, w) {
		t.Errorf("Gene %q SymbolMap=%v, want %v", nand, nand.SymbolMap, w)
	}
}

type intTest struct {
	in   []int
	want int
}

var intTests = []struct {
	gene     string
	tests    []intTest
	count    map[string]int
	argOrder [][]int
}{
	{
		gene: "+.d0.d1.+.+.+.+.d0.d1.d1.d1.d0.d1.d1.d0",
		tests: []intTest{
			{in: []int{1.0, 2.0}, want: 3.0},
		},
		count: map[string]int{
			"+":  1,
			"d0": 1,
			"d1": 1,
		},
		argOrder: [][]int{{1, 2}, nil, nil, {3, 4}, {5, 6}, {7, 8}, {9, 10}, nil, nil, nil, nil, nil, nil, nil, nil},
	},
	{
		gene: "-.+.+.-.-.*.d0.d0.d0.d0.d0.d0.d0",
		tests: []intTest{
			{in: []int{0}, want: 0},
			{in: []int{2}, want: -6},
			{in: []int{6}, want: -42},
			{in: []int{7}, want: -56},
			{in: []int{8}, want: -72},
			{in: []int{10}, want: -110},
			{in: []int{11}, want: -132},
			{in: []int{12}, want: -156},
			{in: []int{14}, want: -210},
			{in: []int{15}, want: -240},
			{in: []int{20}, want: -420},
		},
		count: map[string]int{
			"+":  2,
			"-":  3,
			"*":  1,
			"d0": 7,
		},
		argOrder: [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}, {11, 12}, nil, nil, nil, nil, nil, nil, nil},
	},
	{
		gene: "-.*.*.*.d0./.d0.d0.d0.d0.d0.d0.d0",
		tests: []intTest{
			{in: []int{20.0}, want: 7980.0},
		},
		count: map[string]int{
			"-":  1,
			"*":  3,
			"/":  1,
			"d0": 6,
		},
		argOrder: [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, nil, {9, 10}, nil, nil, nil, nil, nil, nil, nil},
	},
}

func validateInt(t *testing.T, g *Gene, in []int, want int) {
	got := g.EvalInt(in)
	if got != want {
		t.Errorf("%v: math.Eval(%#v) => %v, want %v", g, in, got, want)
	}
}

func TestInt(t *testing.T) {
	for _, test := range intTests {
		g := New(test.gene, functions.Float64)
		argOrder := g.getArgOrder()
		if !reflect.DeepEqual(argOrder, test.argOrder) {
			t.Errorf("Gene %q argOrder=%#v, want %#v", g, argOrder, test.argOrder)
		}
		for _, n := range test.tests {
			validateInt(t, g, n.in, n.want)
		}
		if !reflect.DeepEqual(g.SymbolMap, test.count) {
			t.Errorf("Gene %q SymbolMap=%v, want %v", g, g.SymbolMap, test.count)
		}
	}
}

type mathTest struct {
	in   []float64
	want float64
}

var mathTests = []struct {
	gene     string
	tests    []mathTest
	count    map[string]int
	argOrder [][]int
}{
	{
		gene: "+.d0.d1.+.+.+.+.d0.d1.d1.d1.d0.d1.d1.d0",
		tests: []mathTest{
			{in: []float64{1.0, 2.0}, want: 3.0},
		},
		count: map[string]int{
			"+":  1,
			"d0": 1,
			"d1": 1,
		},
		argOrder: [][]int{{1, 2}, nil, nil, {3, 4}, {5, 6}, {7, 8}, {9, 10}, nil, nil, nil, nil, nil, nil, nil, nil},
	},
	{
		gene: "-.+.+.-.-.*.d0.d0.d0.d0.d0.d0.d0",
		tests: []mathTest{
			{in: []float64{0}, want: 0},
			{in: []float64{2.81}, want: -10.7061},
			{in: []float64{6}, want: -42},
			{in: []float64{7.043}, want: -56.646849},
			{in: []float64{8}, want: -72},
			{in: []float64{10}, want: -110},
			{in: []float64{11.38}, want: -140.8844},
			{in: []float64{12}, want: -156},
			{in: []float64{14}, want: -210},
			{in: []float64{15}, want: -240},
			{in: []float64{20}, want: -420},
		},
		count: map[string]int{
			"+":  2,
			"-":  3,
			"*":  1,
			"d0": 7,
		},
		argOrder: [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}, {11, 12}, nil, nil, nil, nil, nil, nil, nil},
	},
	{
		gene: "-.*.*.*.d0./.d0.d0.d0.d0.d0.d0.d0",
		tests: []mathTest{
			{in: []float64{20.0}, want: 7980.0},
		},
		count: map[string]int{
			"-":  1,
			"*":  3,
			"/":  1,
			"d0": 6,
		},
		argOrder: [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, nil, {9, 10}, nil, nil, nil, nil, nil, nil, nil},
	},
}

func validateMath(t *testing.T, g *Gene, in []float64, want float64) {
	got := g.EvalMath(in)
	if math.Abs(got-want) > 1e-10 {
		t.Errorf("%v: math.Eval(%#v) => %v, want %v", g, in, got, want)
	}
}

func TestMath(t *testing.T) {
	for _, test := range mathTests {
		g := New(test.gene, functions.Float64)
		argOrder := g.getArgOrder()
		if !reflect.DeepEqual(argOrder, test.argOrder) {
			t.Errorf("Gene %q argOrder=%#v, want %#v", g, argOrder, test.argOrder)
		}
		for _, n := range test.tests {
			validateMath(t, g, n.in, n.want)
		}
		if !reflect.DeepEqual(g.SymbolMap, test.count) {
			t.Errorf("Gene %q SymbolMap=%v, want %v", g, g.SymbolMap, test.count)
		}
	}
}

func TestGetBoolArgOrder(t *testing.T) {
	nand := New("Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0", functions.Bool)
	got := nand.getArgOrder()
	want := [][]int{
		{1, 2}, {3, 4}, {5}, {6}, {7, 8}, {9, 10}, {11, 12}, nil, nil, nil, nil, nil, nil, nil, nil,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("nand.GetBoolArgOrder() got %#v, want %#v", got, want)
	}
}

func TestDup(t *testing.T) {
	nand := New("Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0", functions.Bool)
	validateNand(t, nand) // Force evaluation
	g1 := nand.Dup()
	if err := CheckEqual(g1, nand); err != nil {
		t.Errorf("TestDup after Dup failed: g1 != nand: %v\n", err)
	}
	validateNand(t, nand) // Force evaluation
	validateNand(t, g1)

	g1 = New(mathTests[0].gene, functions.Float64)
	test := mathTests[0].tests[0]
	validateMath(t, g1, test.in, test.want) // Force evaluation
	nand = g1.Dup()
	if err := CheckEqual(g1, nand); err != nil {
		t.Errorf("TestDup after Dup failed: g1 != nand: %v\n", err)
	}
	validateMath(t, g1, test.in, test.want) // Force evaluation
	validateMath(t, nand, test.in, test.want)
}

func TestMutate(t *testing.T) {
	headSize := 7
	maxArity := 2
	tailSize := headSize*(maxArity-1) + 1
	numTerminals := 5
	funcs := []FuncWeight{
		{"Not", 1},
		{"And", 5},
		{"Or", 5},
	}
	g1 := RandomNew(headSize, tailSize, numTerminals, 0, funcs, functions.Bool)
	gn := g1.Dup()
	g1.Mutate()
	if err := CheckEqual(gn, g1); err == nil {
		t.Errorf("TestMutate failed: g1 == mux\n")
	}
}

func BenchmarkMutate(b *testing.B) {
	headSize := 7
	maxArity := 2
	tailSize := headSize*(maxArity-1) + 1
	numTerminals := 5
	numConstants := 5
	funcs := []FuncWeight{
		{"+", 1},
		{"-", 5},
		{"*", 5},
	}
	g := RandomNew(headSize, tailSize, numTerminals, numConstants, funcs, functions.Float64)
	for i := 0; i < b.N; i++ {
		g.Mutate()
	}
}

var result *Gene

func BenchmarkDup(b *testing.B) {
	headSize := 7
	maxArity := 2
	tailSize := headSize*(maxArity-1) + 1
	numTerminals := 5
	numConstants := 5
	funcs := []FuncWeight{
		{"+", 1},
		{"-", 5},
		{"*", 5},
	}
	g := RandomNew(headSize, tailSize, numTerminals, numConstants, funcs, functions.Float64)
	var v *Gene
	for i := 0; i < b.N; i++ {
		v = g.Dup()
	}
	result = v
}
