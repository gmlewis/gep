package gene

import (
	"math"
	"reflect"
	"testing"

	bn "github.com/gmlewis/gep/functions/bool_nodes"
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

func validateNand(t *testing.T, g Gene) {
	for i, n := range nandTests {
		r := g.EvalBool(n.in, bn.BoolAllGates)
		if r != n.out {
			t.Errorf("%v: nand.EvalBool(%#v, BoolAllGates) => %v, want %v", i, n.in, r, n.out)
		}
	}
}

func TestNand(t *testing.T) {
	nand := New("Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0")
	validateNand(t, nand)
	nand = New("Or.And.Not.d0.Not.And.Or.d0.d0.d1.d1.d0.d1.d1.d1")
	validateNand(t, nand)
}

type mathTest struct {
	in  []float64
	out float64
}

var mathTests = []struct {
	gene  string
	tests []mathTest
}{
	{
		gene: "+.d0.d1.+.+.+.+.d0.d1.d1.d1.d0.d1.d1.d0",
		tests: []mathTest{
			mathTest{in: []float64{1.0, 2.0}, out: 3.0},
		},
	},
	{
		gene: "-.+.+.-.-.*.d0.d0.d0.d0.d0.d0.d0",
		tests: []mathTest{
			mathTest{in: []float64{0}, out: 0},
			mathTest{in: []float64{2.81}, out: -10.7061},
			mathTest{in: []float64{6}, out: -42},
			mathTest{in: []float64{7.043}, out: -56.646849},
			mathTest{in: []float64{8}, out: -72},
			mathTest{in: []float64{10}, out: -110},
			mathTest{in: []float64{11.38}, out: -140.8844},
			mathTest{in: []float64{12}, out: -156},
			mathTest{in: []float64{14}, out: -210},
			mathTest{in: []float64{15}, out: -240},
			mathTest{in: []float64{20}, out: -420},
		},
	},
	{
		gene: "-.*.*.*.d0./.d0.d0.d0.d0.d0.d0.d0",
		tests: []mathTest{
			mathTest{in: []float64{20.0}, out: 7980.0},
		},
	},
}

func TestMath(t *testing.T) {
	for _, v := range mathTests {
		g := New(v.gene)
		for i, n := range v.tests {
			r := g.EvalMath(n.in)
			if math.Abs(r-n.out) > 1e-10 {
				t.Errorf("%v[%v]: math.Eval(%#v) => %v, want %v", v.gene, i, n.in, r, n.out)
			}
		}
	}
}

func TestGetBoolArgOrder(t *testing.T) {
	nand := New("Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0")
	got := nand.getBoolArgOrder(bn.BoolAllGates)
	want := [][]int{
		{1, 2}, {3, 4}, {5}, {6}, {7, 8}, {9, 10}, {11, 12}, nil, nil, nil, nil, nil, nil, nil, nil,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("nand.GetBoolArgOrder() got %#v, want %#v", got, want)
	}
}
