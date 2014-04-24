package math_nodes

import (
	"testing"
)

func TestPlus(t *testing.T) {
	if g, w := Math["+"].Float64Function(2, 3, 4, 5), 5.0; g != w {
		t.Errorf("Math[+](2,3,4,5) = %v, want %v", g, w)
	}
}
