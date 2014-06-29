package gene

import (
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

var mathTests = []struct {
	in  []float64
	out float64
}{
	{[]float64{1.0, 2.0}, 3.0},
}

func TestMath(t *testing.T) {
	math := New("+.d0.d1.+.+.+.+.d0.d1.d1.d1.d0.d1.d1.d0")
	for i, n := range mathTests {
		r := math.EvalMath(n.in)
		if r != n.out {
			t.Errorf("%v: math.EvalFloat64(%#v, MathNodes) => %v, want %v", i, n.in, r, n.out)
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
