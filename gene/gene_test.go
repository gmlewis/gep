// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"io"
	"math"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/gmlewis/gep/v2/functions"
)

func mustNew(t *testing.T, karva string, funcType functions.FuncType) *Gene {
	t.Helper()
	g, err := New(karva, funcType)
	if err != nil {
		t.Fatalf("New(%q) error: %v", karva, err)
	}
	return g
}

func mustEvalBool(t *testing.T, g *Gene, in []bool) bool {
	t.Helper()
	v, err := g.EvalBool(in)
	if err != nil {
		t.Fatalf("EvalBool() error: %v", err)
	}
	return v
}

func mustEvalInt(t *testing.T, g *Gene, in []int) int {
	t.Helper()
	v, err := g.EvalInt(in)
	if err != nil {
		t.Fatalf("EvalInt() error: %v", err)
	}
	return v
}

func mustEvalMath(t *testing.T, g *Gene, in []float64) float64 {
	t.Helper()
	v, err := g.EvalMath(in)
	if err != nil {
		t.Fatalf("EvalMath() error: %v", err)
	}
	return v
}

func mustEvalVectorInt(t *testing.T, g *Gene, in []VectorInt) VectorInt {
	t.Helper()
	v, err := g.EvalVectorInt(in)
	if err != nil {
		t.Fatalf("EvalVectorInt() error: %v", err)
	}
	return v
}

func mustArgOrder(t *testing.T, g *Gene) [][]int {
	t.Helper()
	v, err := g.getArgOrder()
	if err != nil {
		t.Fatalf("getArgOrder() error: %v", err)
	}
	return v
}

func mustMutate(t *testing.T, g *Gene) {
	t.Helper()
	if err := g.Mutate(); err != nil {
		t.Fatalf("Mutate() error: %v", err)
	}
}

func mustDup(t *testing.T, g *Gene) *Gene {
	t.Helper()
	v, err := g.Dup()
	if err != nil {
		t.Fatalf("Dup() error: %v", err)
	}
	return v
}

func TestNew_InvalidSymbolIndexesDoNotExit(t *testing.T) {
	g, err := New("dX.cY", functions.Float64)
	if err == nil {
		t.Fatalf("New() error = nil, want non-nil")
	}
	if g == nil {
		t.Fatalf("New() = nil, want non-nil")
	}
	if got := len(g.Constants); got != 0 {
		t.Fatalf("len(g.Constants) = %v, want 0", got)
	}
}

func TestNew_InvalidSymbolIndexesReturnsError(t *testing.T) {
	g, err := New("dX.cY", functions.Float64)
	if g == nil {
		t.Fatalf("New() = nil, want non-nil")
	}
	if err == nil {
		t.Fatalf("New() error = nil, want non-nil")
	}
}

func TestNew_InvalidSymbolIndexesDoNotWriteStderr(t *testing.T) {
	orig := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error: %v", err)
	}
	os.Stderr = w
	defer func() { os.Stderr = orig }()

	_, _ = New("dX.cY", functions.Float64)

	if err := w.Close(); err != nil {
		t.Fatalf("stderr close error: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("stderr read error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("New wrote to stderr: %q", string(out))
	}
}

func TestSymbolCount_UnknownFuncTypeReturnsError(t *testing.T) {
	g := mustNew(t, "d0", functions.FuncType(999))
	if _, err := g.SymbolCount("d0"); err == nil {
		t.Fatalf("SymbolCount() error = nil, want non-nil")
	}
}

func TestAllSymbolsEqualWeights_UnknownFuncTypeReturnsNil(t *testing.T) {
	if got, err := AllSymbolsEqualWeights(functions.FuncType(999)); err == nil || got != nil {
		t.Fatalf("AllSymbolsEqualWeights(unknown) = %v, want nil", got)
	}
}

func TestAllSymbolsEqualWeights_UnknownFuncTypeReturnsError(t *testing.T) {
	if _, err := AllSymbolsEqualWeights(functions.FuncType(999)); err == nil {
		t.Fatal("AllSymbolsEqualWeights(unknown) error = nil, want non-nil")
	}
}

func TestEvalInt_UnknownSymbolReturnsError(t *testing.T) {
	g := mustNew(t, "UNKNOWN.d0", functions.Int)
	if _, err := g.EvalInt([]int{1}); err == nil {
		t.Fatal("EvalInt() error = nil, want non-nil")
	}
}

func TestMutate_TooFewChoicesReturnsError(t *testing.T) {
	g := &Gene{
		Symbols:      []string{"d0"},
		HeadSize:     1,
		choiceSlice:  []string{"d0"},
		numTerminals: 1,
	}
	if err := g.Mutate(); err == nil {
		t.Fatal("Mutate() error = nil, want non-nil")
	}
}

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
		got := mustEvalBool(t, g, n.in)
		if got != n.want {
			t.Errorf("%v: mustEvalBool(t, nand, %#v, BoolAllGates) => %v, want %v", i, n.in, got, n.want)
		}
	}
}

func TestNand(t *testing.T) {
	nand := mustNew(t, "Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0", functions.Bool)
	validateNand(t, nand)
	w := map[string]int{
		"And": 3,
		"Not": 2,
		"Or":  2,
		"d0":  2,
		"d1":  4,
	}
	if !reflect.DeepEqual(nand.SymbolMap, w) {
		t.Errorf("Gene %v SymbolMap=%v, want %v", nand, nand.SymbolMap, w)
	}
	nand = mustNew(t, "Or.And.Not.d0.Not.And.Or.d0.d0.d1.d1.d0.d1.d1.d1", functions.Bool)
	validateNand(t, nand)
	w = map[string]int{
		"And": 2,
		"Not": 2,
		"Or":  2,
		"d0":  3,
		"d1":  2,
	}
	if !reflect.DeepEqual(nand.SymbolMap, w) {
		t.Errorf("Gene %v SymbolMap=%v, want %v", nand, nand.SymbolMap, w)
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
			{in: []int{1, 2}, want: 3},
		},
		count: map[string]int{
			"+":  1,
			"d0": 1,
			"d1": 1,
		},
		argOrder: [][]int{{1, 2}, nil, nil, {3, 4}, {5, 6}, {7, 8}, {9, 10}, nil, nil, nil, nil, nil, nil, nil, nil},
	},
	{
		gene: "-.+.+.-.-.*.d0.d0.d0.d0.d0.d0.d0", // -(x^2 + x)
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
		gene: "-.*.*.*.d0./.d0.d0.d0.d0.d0.d0.d0", // x^3 - x
		tests: []intTest{
			{in: []int{20}, want: 7980},
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
	got := mustEvalInt(t, g, in)
	if got != want {
		t.Errorf("%v: math.Eval(%#v) => %v, want %v", g, in, got, want)
	}
}

func TestInt(t *testing.T) {
	for _, test := range intTests {
		g := mustNew(t, test.gene, functions.Float64)
		argOrder := mustArgOrder(t, g)
		if !reflect.DeepEqual(argOrder, test.argOrder) {
			t.Errorf("Gene %v argOrder=%#v, want %#v", g, argOrder, test.argOrder)
		}
		for _, n := range test.tests {
			validateInt(t, g, n.in, n.want)
		}
		if !reflect.DeepEqual(g.SymbolMap, test.count) {
			t.Errorf("Gene %v SymbolMap=%v, want %v", g, g.SymbolMap, test.count)
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
	got := mustEvalMath(t, g, in)
	if math.Abs(got-want) > 1e-10 {
		t.Errorf("%v: math.Eval(%#v) => %v, want %v", g, in, got, want)
	}
}

func TestMath(t *testing.T) {
	for _, test := range mathTests {
		g := mustNew(t, test.gene, functions.Float64)
		argOrder := mustArgOrder(t, g)
		if !reflect.DeepEqual(argOrder, test.argOrder) {
			t.Errorf("Gene %v argOrder=%#v, want %#v", g, argOrder, test.argOrder)
		}
		for _, n := range test.tests {
			validateMath(t, g, n.in, n.want)
		}
		if !reflect.DeepEqual(g.SymbolMap, test.count) {
			t.Errorf("Gene %v SymbolMap=%v, want %v", g, g.SymbolMap, test.count)
		}
	}
}

type vectorIntTest struct {
	in   []VectorInt
	want VectorInt
}

var vectorIntTests = []struct {
	gene     string
	tests    []vectorIntTest
	count    map[string]int
	argOrder [][]int
}{
	{
		gene: "+.d0.d1.+.+.+.+.d0.d1.d1.d1.d0.d1.d1.d0",
		tests: []vectorIntTest{
			{in: []VectorInt{{1, 2, 3, 4}, {5, 6, 7, 8}}, want: VectorInt{6, 8, 10, 12}},
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
		tests: []vectorIntTest{
			{in: []VectorInt{{0, 2, 6, 7, 8, 10, 11, 12, 14, 15, 20}}, want: VectorInt{0, -6, -42, -56, -72, -110, -132, -156, -210, -240, -420}},
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
		tests: []vectorIntTest{
			{in: []VectorInt{{0, 20}}, want: VectorInt{0, 7980}},
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

func validateVectorInt(t *testing.T, g *Gene, in []VectorInt, want VectorInt) {
	got := mustEvalVectorInt(t, g, in)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v: math.Eval(%#v) => %v, want %v", g, in, got, want)
	}
}

func TestVectorInt(t *testing.T) {
	for _, test := range vectorIntTests {
		g := mustNew(t, test.gene, functions.Float64)
		argOrder := mustArgOrder(t, g)
		if !reflect.DeepEqual(argOrder, test.argOrder) {
			t.Errorf("Gene %v argOrder=%#v, want %#v", g, argOrder, test.argOrder)
		}
		for _, n := range test.tests {
			validateVectorInt(t, g, n.in, n.want)
		}
		if !reflect.DeepEqual(g.SymbolMap, test.count) {
			t.Errorf("Gene %v SymbolMap=%v, want %v", g, g.SymbolMap, test.count)
		}
	}
}

func TestBuildMathTree_OutOfBoundsSymbolIndex(t *testing.T) {
	g := mustNew(t, "d0", functions.Float64)
	f := g.buildMathTree(len(g.Symbols), mustArgOrder(t, g))
	if got := f([]float64{42}); got != 0 {
		t.Fatalf("buildMathTree out-of-bounds fallback = %v, want 0", got)
	}
}

func TestBuildVectorIntTree_OutOfBoundsSymbolIndex(t *testing.T) {
	g := mustNew(t, "d0", functions.VectorInts)
	f := g.buildVectorIntTree(len(g.Symbols), mustArgOrder(t, g))
	if got := f([]VectorInt{{1, 2, 3}}); len(got) != 0 {
		t.Fatalf("buildVectorIntTree out-of-bounds fallback = %v, want empty vector", got)
	}
}

func TestGetBoolArgOrder(t *testing.T) {
	nand := mustNew(t, "Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0", functions.Bool)
	got := mustArgOrder(t, nand)
	want := [][]int{
		{1, 2}, {3, 4}, {5}, {6}, {7, 8}, {9, 10}, {11, 12}, nil, nil, nil, nil, nil, nil, nil, nil,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("nand.GetBoolArgOrder() got %#v, want %#v", got, want)
	}
}

func TestDup(t *testing.T) {
	nand := mustNew(t, "Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0", functions.Bool)
	validateNand(t, nand) // Force evaluation
	g1 := mustDup(t, nand)
	if err := CheckEqual(g1, nand); err != nil {
		t.Errorf("TestDup after Dup failed: g1 != nand: %v\n", err)
	}
	validateNand(t, nand) // Force evaluation
	validateNand(t, g1)

	g1 = mustNew(t, mathTests[0].gene, functions.Float64)
	test := mathTests[0].tests[0]
	validateMath(t, g1, test.in, test.want) // Force evaluation
	nand = mustDup(t, g1)
	if err := CheckEqual(g1, nand); err != nil {
		t.Errorf("TestDup after Dup failed: g1 != nand: %v\n", err)
	}
	validateMath(t, g1, test.in, test.want) // Force evaluation
	validateMath(t, nand, test.in, test.want)
}

func TestDup_DoesNotReuseSourceCachedFunctions(t *testing.T) {
	g := mustNew(t, "c0", functions.Float64)
	g.Constants = []float64{2}
	if got := mustEvalMath(t, g, nil); got != 2 {
		t.Fatalf("mustEvalMath(t, g, nil) = %v, want 2", got)
	}

	dup := mustDup(t, g)
	dup.Constants[0] = 7
	if got := mustEvalMath(t, dup, nil); got != 7 {
		t.Fatalf("mustEvalMath(t, dup, nil) = %v, want 7", got)
	}
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
	gn := mustDup(t, g1)
	mustMutate(t, g1)
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
		if err := g.Mutate(); err != nil {
			b.Fatalf("Mutate() error: %v", err)
		}
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
		var err error
		v, err = g.Dup()
		if err != nil {
			b.Fatalf("Dup() error: %v", err)
		}
	}
	result = v
}

// TestRandomNew_SeededIsDeterministic verifies that two calls to RandomNew with
// the same seeded *rand.Rand produce identical genes.
func TestRandomNew_SeededIsDeterministic(t *testing.T) {
	headSize := 6
	maxArity := 2
	tailSize := headSize*(maxArity-1) + 1
	numTerminals := 3
	numConstants := 2
	funcs := []FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
	}
	rng1 := rand.New(rand.NewSource(42))
	rng2 := rand.New(rand.NewSource(42))

	g1 := RandomNew(headSize, tailSize, numTerminals, numConstants, funcs, functions.Float64, rng1)
	g2 := RandomNew(headSize, tailSize, numTerminals, numConstants, funcs, functions.Float64, rng2)

	if len(g1.Symbols) != len(g2.Symbols) {
		t.Fatalf("symbol count mismatch: %v vs %v", len(g1.Symbols), len(g2.Symbols))
	}
	for i, sym := range g1.Symbols {
		if sym != g2.Symbols[i] {
			t.Fatalf("symbol[%v]: %q != %q", i, sym, g2.Symbols[i])
		}
	}
}

// TestRandomNew_DifferentSeedsProduceDifferentGenes verifies that different seeds
// produce different genes (with overwhelming probability given a non-trivial gene size).
func TestRandomNew_DifferentSeedsProduceDifferentGenes(t *testing.T) {
	headSize := 6
	maxArity := 2
	tailSize := headSize*(maxArity-1) + 1
	numTerminals := 3
	numConstants := 0
	funcs := []FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
	}
	rng1 := rand.New(rand.NewSource(1))
	rng2 := rand.New(rand.NewSource(2))

	g1 := RandomNew(headSize, tailSize, numTerminals, numConstants, funcs, functions.Float64, rng1)
	g2 := RandomNew(headSize, tailSize, numTerminals, numConstants, funcs, functions.Float64, rng2)

	different := false
	for i, sym := range g1.Symbols {
		if i < len(g2.Symbols) && sym != g2.Symbols[i] {
			different = true
			break
		}
	}
	if !different {
		t.Fatal("genes from different seeds are identical; expected differences")
	}
}

// TestRandomNew_GlobalRandUsedWhenNoRngProvided verifies that RandomNew works
// correctly when no rng is provided (falls back to global math/rand source).
func TestRandomNew_GlobalRandUsedWhenNoRngProvided(t *testing.T) {
	headSize := 4
	maxArity := 2
	tailSize := headSize*(maxArity-1) + 1
	numTerminals := 2
	numConstants := 0
	funcs := []FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
	}
	g := RandomNew(headSize, tailSize, numTerminals, numConstants, funcs, functions.Float64)
	if g == nil {
		t.Fatal("RandomNew() returned nil")
	}
	want := headSize + tailSize
	if got := len(g.Symbols); got != want {
		t.Fatalf("len(Symbols) = %v, want %v", got, want)
	}
}
