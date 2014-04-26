package model

import (
	"testing"

	mn "github.com/gmlewis/gep/functions/math_nodes"
	"github.com/gmlewis/gep/gene"
)

func TestMaxArity(t *testing.T) {
	funcs := []gene.FuncWeight{
		{"+", 1},
		{"-", 2},
		{"*", 3},
		{"/", 4},
	}
	if g, w := maxArity(funcs, mn.Math), 2; g != w {
		t.Errorf("maxArity(%v, mn.Math) = %v, want %v", funcs, g, w)
	}
	funcs = append(funcs, gene.FuncWeight{
		Symbol: "LT3A",
		Weight: 1,
	})
	if g, w := maxArity(funcs, mn.Math), 3; g != w {
		t.Errorf("maxArity(%v, mn.Math) = %v, want %v", funcs, g, w)
	}
}
