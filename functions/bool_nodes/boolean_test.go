package boolNodes

import (
	"testing"
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

func TestNand(t *testing.T) {
	for _, test := range nandTests {
		if g, w := BoolAllGates["Nand"].BoolFunction(test.in[0], test.in[1], false, false), test.out; g != w {
			t.Errorf("BoolAllGates[Nand](%v,%v,false,false) = %v, want %v", test.in[0], test.in[1], g, w)
		}
	}
}
