package genome

import (
	"bytes"
	"testing"

	"github.com/gmlewis/gep/gene"
	"github.com/gmlewis/gep/grammars"
)

func TestWriteNand(t *testing.T) {
	want := `package gepModel

func gepModel(d []bool) bool {
	y := false

	y = (((!(d[0] && d[1])) && (d[0] || d[1])) || (!(d[1] && d[1])))

	return y
}
`

	g1 := gene.New("Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0")
	gn := New([]gene.Gene{g1}, "Or")
	grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		t.Fatalf("unable to LoadGoBooleanAllGatesGrammar(): %v", err)
	}

	b := new(bytes.Buffer)
	gn.Write(b, grammar)
	if b.String() != want {
		t.Errorf("gen.Write() got %v, want %v", b.String(), want)
	}
}

func TestWriteMath(t *testing.T) {
	want := `package gepModel

import (
	"math"
)

func gepModel(d []float64) float64 {
	y := 0.0

	y = (d[0] * d[0])
	y += d[0]
	y += (d[0] * (d[0] * d[0]))
	y += (((d[0] * d[0]) * d[0]) * d[0])

	return y
}
`

	g1 := gene.New("*.d0.d0.*.d0.*.*.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0")
	g2 := gene.New("d0.*.d0.*.*.d0.d0.*.d0.d0.d0.d0.d0.d0.d0.d0.d0")
	g3 := gene.New("*.d0.*.d0.d0.*.d0.*.d0.d0.d0.d0.d0.d0.d0.d0.d0")
	g4 := gene.New("*.*.d0.*.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0")
	gn := New([]gene.Gene{g1, g2, g3, g4}, "+")
	grammar, err := grammars.LoadGoMathGrammar()
	if err != nil {
		t.Fatalf("unable to LoadGoMathGrammar(): %v", err)
	}

	b := new(bytes.Buffer)
	gn.Write(b, grammar)
	if b.String() != want {
		t.Errorf("gen.Write() got %v, want %v", b.String(), want)
	}
}

func TestWrite6Multiplier(t *testing.T) {
	want := `package gepModel

func gepModel(d []bool) bool {
	y := false

	y = (d[3] || ((gepNor(d[1], d[3]) || (d[2] && d[0])) || gepNand((d[1] || d[0]), d[2])))
	y = y && gepNand((gepNor((d[5] || d[5]), gepNand(d[1], d[0])) || gepNor((d[2] || d[0]), d[3])), d[0])
	y = y && gepNor((((d[1] && d[2]) && d[4]) && gepNor(d[0], d[4])), gepNor((d[1] || d[4]), gepNor(d[4], d[0])))
	y = y && (gepNor(gepNor(gepNand(d[2], d[4]), d[1]), (!(d[2]))) || gepNand(gepNand(d[1], d[3]), gepNor(d[0], d[2])))

	return y
}

func gepNand(x, y bool) bool {
	return (!(x && y))
}

func gepNor(x, y bool) bool {
	return (!(x || y))
}
`

	g1 := gene.New("Or.d3.Or.Or.Nand.Nor.And.Or.d2.d1.d3.d2.d0.d1.d0.d0.d0")
	g2 := gene.New("Nand.Or.d0.Nor.Nor.Or.Nand.Or.d3.d5.d5.d1.d0.d2.d0.d1.d3")
	g3 := gene.New("Nor.And.Nor.And.Nor.Or.Nor.And.d4.d0.d4.d1.d4.d4.d0.d1.d2")
	g4 := gene.New("Or.Nor.Nand.Nor.Not.Nand.Nor.Nand.d1.d2.d1.d3.d0.d2.d2.d4.d4")
	gn := New([]gene.Gene{g1, g2, g3, g4}, "And")
	grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		t.Fatalf("unable to LoadGoBooleanAllGatesGrammar(): %v", err)
	}

	b := new(bytes.Buffer)
	gn.Write(b, grammar)
	if b.String() != want {
		t.Errorf("gen.Write() got %v, want %v", b.String(), want)
	}
}
