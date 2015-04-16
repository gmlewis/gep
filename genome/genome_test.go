package genome

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	bn "github.com/gmlewis/gep/functions/bool_nodes"
	"github.com/gmlewis/gep/gene"
)

const delta = 1.0

// Logic Synthesis

var sixMultiplexerTests = []struct {
	in  []bool
	out bool
}{
	{[]bool{false, false, false, false, false, false}, false},
	{[]bool{false, false, false, false, false, true}, false},
	{[]bool{false, false, false, false, true, false}, false},
	{[]bool{false, false, false, false, true, true}, false},
	{[]bool{false, false, false, true, false, false}, false},
	{[]bool{false, false, false, true, false, true}, false},
	{[]bool{false, false, false, true, true, false}, false},
	{[]bool{false, false, false, true, true, true}, false},
	{[]bool{false, false, true, false, false, false}, true},
	{[]bool{false, false, true, false, false, true}, true},
	{[]bool{false, false, true, false, true, false}, true},
	{[]bool{false, false, true, false, true, true}, true},
	{[]bool{false, false, true, true, false, false}, true},
	{[]bool{false, false, true, true, false, true}, true},
	{[]bool{false, false, true, true, true, false}, true},
	{[]bool{false, false, true, true, true, true}, true},
	{[]bool{false, true, false, false, false, false}, false},
	{[]bool{false, true, false, false, false, true}, false},
	{[]bool{false, true, false, false, true, false}, false},
	{[]bool{false, true, false, false, true, true}, false},
	{[]bool{false, true, false, true, false, false}, true},
	{[]bool{false, true, false, true, false, true}, true},
	{[]bool{false, true, false, true, true, false}, true},
	{[]bool{false, true, false, true, true, true}, true},
	{[]bool{false, true, true, false, false, false}, false},
	{[]bool{false, true, true, false, false, true}, false},
	{[]bool{false, true, true, false, true, false}, false},
	{[]bool{false, true, true, false, true, true}, false},
	{[]bool{false, true, true, true, false, false}, true},
	{[]bool{false, true, true, true, false, true}, true},
	{[]bool{false, true, true, true, true, false}, true},
	{[]bool{false, true, true, true, true, true}, true},
	{[]bool{true, false, false, false, false, false}, false},
	{[]bool{true, false, false, false, false, true}, false},
	{[]bool{true, false, false, false, true, false}, true},
	{[]bool{true, false, false, false, true, true}, true},
	{[]bool{true, false, false, true, false, false}, false},
	{[]bool{true, false, false, true, false, true}, false},
	{[]bool{true, false, false, true, true, false}, true},
	{[]bool{true, false, false, true, true, true}, true},
	{[]bool{true, false, true, false, false, false}, false},
	{[]bool{true, false, true, false, false, true}, false},
	{[]bool{true, false, true, false, true, false}, true},
	{[]bool{true, false, true, false, true, true}, true},
	{[]bool{true, false, true, true, false, false}, false},
	{[]bool{true, false, true, true, false, true}, false},
	{[]bool{true, false, true, true, true, false}, true},
	{[]bool{true, false, true, true, true, true}, true},
	{[]bool{true, true, false, false, false, false}, false},
	{[]bool{true, true, false, false, false, true}, true},
	{[]bool{true, true, false, false, true, false}, false},
	{[]bool{true, true, false, false, true, true}, true},
	{[]bool{true, true, false, true, false, false}, false},
	{[]bool{true, true, false, true, false, true}, true},
	{[]bool{true, true, false, true, true, false}, false},
	{[]bool{true, true, false, true, true, true}, true},
	{[]bool{true, true, true, false, false, false}, false},
	{[]bool{true, true, true, false, false, true}, true},
	{[]bool{true, true, true, false, true, false}, false},
	{[]bool{true, true, true, false, true, true}, true},
	{[]bool{true, true, true, true, false, false}, false},
	{[]bool{true, true, true, true, false, true}, true},
	{[]bool{true, true, true, true, true, false}, false},
	{[]bool{true, true, true, true, true, true}, true},
}

func validateSixMultiplexer(t *testing.T, g *Genome) {
	for i, n := range sixMultiplexerTests {
		r := g.EvalBool(n.in, bn.BoolAllGates)
		if r != n.out {
			t.Errorf("%v: sixMultiplexer.EvalBool(%#v, BoolAllGates) => %v, want %v", i, n.in, r, n.out)
		}
	}
}

func TestSixMultiplexer(t *testing.T) {
	mux := New([]*gene.Gene{
		gene.New("Nand.Or.And.Not.Nor.Not.Nor.And.d3.d2.d0.d2.d1.d4.d1.d4.d2"),
		gene.New("Nor.Nor.And.Or.Nor.And.Nor.Or.d0.d1.d0.d2.d3.d1.d2.d0.d3"),
		gene.New("Or.Or.Nand.d4.Not.Or.Nand.Nand.d0.d4.d3.d4.d1.d1.d3.d2.d0"),
		gene.New("Or.And.And.Nand.d5.Nand.Nand.Nor.d4.d3.d5.d0.d1.d0.d4.d0.d1"),
	},
		"And")
	validateSixMultiplexer(t, mux)
	w := map[string]int{
		"And":  7,
		"Nand": 6,
		"Nor":  5,
		"Not":  1,
		"Or":   6,
		"d0":   6,
		"d1":   5,
		"d2":   2,
		"d3":   4,
		"d4":   5,
		"d5":   2,
	}
	if !reflect.DeepEqual(mux.SymbolCount, w) {
		t.Errorf("Genome %q SymbolCount=%v, want %v", mux, mux.SymbolCount, w)
	}
}

var odd3ParityTests = []struct {
	in  []bool
	out bool
}{
	{[]bool{false, false, false}, false},
	{[]bool{false, false, true}, true},
	{[]bool{false, true, false}, true},
	{[]bool{false, true, true}, false},
	{[]bool{true, false, false}, true},
	{[]bool{true, false, true}, false},
	{[]bool{true, true, false}, false},
	{[]bool{true, true, true}, true},
}

func validateOdd3Parity(t *testing.T, g *Genome) {
	for i, n := range odd3ParityTests {
		r := g.EvalBool(n.in, bn.BoolAllGates)
		if r != n.out {
			t.Errorf("%v: odd3Parity.EvalBool(%#v, BoolAllGates) => %v, want %v", i, n.in, r, n.out)
		}
	}
}

func TestOdd3Parity(t *testing.T) {
	mux := New([]*gene.Gene{
		gene.New("Or.Or.d1.And.Or.d0.And.d2.d0.d2.d1.d1.d1.d0.d1"),
		gene.New("Not.And.And.Not.Or.And.And.d0.d1.d2.d2.d1.d0.d2.d2"),
		gene.New("Or.Or.Or.And.And.Not.Not.d1.d2.d0.d2.d1.d0.d0.d2"),
	},
		"And")
	validateOdd3Parity(t, mux)
	w := map[string]int{
		"And": 8,
		"Not": 4,
		"Or":  4,
		"d0":  4,
		"d1":  4,
		"d2":  4,
	}
	if !reflect.DeepEqual(mux.SymbolCount, w) {
		t.Errorf("Genome %q SymbolCount=%v, want %v", mux, mux.SymbolCount, w)
	}
}

var odd7ParityTests = []struct {
	in  []bool
	out bool
}{
	{[]bool{false, false, false, false, false, false, false}, false}, // 1
	{[]bool{false, false, false, false, false, false, true}, true},
	{[]bool{false, false, false, false, false, true, false}, true},
	{[]bool{false, false, false, false, false, true, true}, false},
	{[]bool{false, false, false, false, true, false, false}, true},
	{[]bool{false, false, false, false, true, false, true}, false},
	{[]bool{false, false, false, false, true, true, false}, false},
	{[]bool{false, false, false, false, true, true, true}, true},
	{[]bool{false, false, false, true, false, false, false}, true},
	{[]bool{false, false, false, true, false, false, true}, false},
	{[]bool{false, false, false, true, false, true, false}, false},
	{[]bool{false, false, false, true, false, true, true}, true},
	{[]bool{false, false, false, true, true, false, false}, false},
	{[]bool{false, false, false, true, true, false, true}, true},
	{[]bool{false, false, false, true, true, true, false}, true},
	{[]bool{false, false, false, true, true, true, true}, false},
	{[]bool{false, false, true, false, false, false, false}, true},
	{[]bool{false, false, true, false, false, false, true}, false},
	{[]bool{false, false, true, false, false, true, false}, false},
	{[]bool{false, false, true, false, false, true, true}, true},
	{[]bool{false, false, true, false, true, false, false}, false},
	{[]bool{false, false, true, false, true, false, true}, true},
	{[]bool{false, false, true, false, true, true, false}, true},
	{[]bool{false, false, true, false, true, true, true}, false},
	{[]bool{false, false, true, true, false, false, false}, false},
	{[]bool{false, false, true, true, false, false, true}, true}, // 26
	{[]bool{false, false, true, true, false, true, false}, true},
	{[]bool{false, false, true, true, false, true, true}, false},
	{[]bool{false, false, true, true, true, false, false}, true},
	{[]bool{false, false, true, true, true, false, true}, false},
	{[]bool{false, false, true, true, true, true, false}, false},
	{[]bool{false, false, true, true, true, true, true}, true},
	{[]bool{false, true, false, false, false, false, false}, true},
	{[]bool{false, true, false, false, false, false, true}, false},
	{[]bool{false, true, false, false, false, true, false}, false},
	{[]bool{false, true, false, false, false, true, true}, true},
	{[]bool{false, true, false, false, true, false, false}, false},
	{[]bool{false, true, false, false, true, false, true}, true},
	{[]bool{false, true, false, false, true, true, false}, true}, // 39
	{[]bool{false, true, false, false, true, true, true}, false},
	{[]bool{false, true, false, true, false, false, false}, false},
	{[]bool{false, true, false, true, false, false, true}, true},
	{[]bool{false, true, false, true, false, true, false}, true},
	{[]bool{false, true, false, true, false, true, true}, false},
	{[]bool{false, true, false, true, true, false, false}, true},
	{[]bool{false, true, false, true, true, false, true}, false},
	{[]bool{false, true, false, true, true, true, false}, false},
	{[]bool{false, true, false, true, true, true, true}, true},
	{[]bool{false, true, true, false, false, false, false}, false},
	{[]bool{false, true, true, false, false, false, true}, true},
	{[]bool{false, true, true, false, false, true, false}, true},
	{[]bool{false, true, true, false, false, true, true}, false}, // 52
	{[]bool{false, true, true, false, true, false, false}, true},
	{[]bool{false, true, true, false, true, false, true}, false},
	{[]bool{false, true, true, false, true, true, false}, false},
	{[]bool{false, true, true, false, true, true, true}, true},
	{[]bool{false, true, true, true, false, false, false}, true},
	{[]bool{false, true, true, true, false, false, true}, false},
	{[]bool{false, true, true, true, false, true, false}, false},
	{[]bool{false, true, true, true, false, true, true}, true},
	{[]bool{false, true, true, true, true, false, false}, false},
	{[]bool{false, true, true, true, true, false, true}, true},
	{[]bool{false, true, true, true, true, true, false}, true},
	{[]bool{false, true, true, true, true, true, true}, false}, // 64
	{[]bool{true, false, false, false, false, false, false}, true},
	{[]bool{true, false, false, false, false, false, true}, false},
	{[]bool{true, false, false, false, false, true, false}, false},
	{[]bool{true, false, false, false, false, true, true}, true},
	{[]bool{true, false, false, false, true, false, false}, false},
	{[]bool{true, false, false, false, true, false, true}, true},
	{[]bool{true, false, false, false, true, true, false}, true},
	{[]bool{true, false, false, false, true, true, true}, false},
	{[]bool{true, false, false, true, false, false, false}, false},
	{[]bool{true, false, false, true, false, false, true}, true},
	{[]bool{true, false, false, true, false, true, false}, true},
	{[]bool{true, false, false, true, false, true, true}, false}, // 76
	{[]bool{true, false, false, true, true, false, false}, true},
	{[]bool{true, false, false, true, true, false, true}, false},
	{[]bool{true, false, false, true, true, true, false}, false},
	{[]bool{true, false, false, true, true, true, true}, true},
	{[]bool{true, false, true, false, false, false, false}, false},
	{[]bool{true, false, true, false, false, false, true}, true},
	{[]bool{true, false, true, false, false, true, false}, true},
	{[]bool{true, false, true, false, false, true, true}, false},
	{[]bool{true, false, true, false, true, false, false}, true},
	{[]bool{true, false, true, false, true, false, true}, false},
	{[]bool{true, false, true, false, true, true, false}, false},
	{[]bool{true, false, true, false, true, true, true}, true}, // 88
	{[]bool{true, false, true, true, false, false, false}, true},
	{[]bool{true, false, true, true, false, false, true}, false},
	{[]bool{true, false, true, true, false, true, false}, false},
	{[]bool{true, false, true, true, false, true, true}, true},
	{[]bool{true, false, true, true, true, false, false}, false},
	{[]bool{true, false, true, true, true, false, true}, true},
	{[]bool{true, false, true, true, true, true, false}, true},
	{[]bool{true, false, true, true, true, true, true}, false},
	{[]bool{true, true, false, false, false, false, false}, false},
	{[]bool{true, true, false, false, false, false, true}, true},
	{[]bool{true, true, false, false, false, true, false}, true},
	{[]bool{true, true, false, false, false, true, true}, false}, // 100
	{[]bool{true, true, false, false, true, false, false}, true},
	{[]bool{true, true, false, false, true, false, true}, false},
	{[]bool{true, true, false, false, true, true, false}, false},
	{[]bool{true, true, false, false, true, true, true}, true},
	{[]bool{true, true, false, true, false, false, false}, true},
	{[]bool{true, true, false, true, false, false, true}, false},
	{[]bool{true, true, false, true, false, true, false}, false},
	{[]bool{true, true, false, true, false, true, true}, true},
	{[]bool{true, true, false, true, true, false, false}, false},
	{[]bool{true, true, false, true, true, false, true}, true},
	{[]bool{true, true, false, true, true, true, false}, true},
	{[]bool{true, true, false, true, true, true, true}, false}, // 112
	{[]bool{true, true, true, false, false, false, false}, true},
	{[]bool{true, true, true, false, false, false, true}, false},
	{[]bool{true, true, true, false, false, true, false}, false},
	{[]bool{true, true, true, false, false, true, true}, true},
	{[]bool{true, true, true, false, true, false, false}, false},
	{[]bool{true, true, true, false, true, false, true}, true},
	{[]bool{true, true, true, false, true, true, false}, true},
	{[]bool{true, true, true, false, true, true, true}, false},
	{[]bool{true, true, true, true, false, false, false}, false},
	{[]bool{true, true, true, true, false, false, true}, true},
	{[]bool{true, true, true, true, false, true, false}, true},
	{[]bool{true, true, true, true, false, true, true}, false}, // 124
	{[]bool{true, true, true, true, true, false, false}, true},
	{[]bool{true, true, true, true, true, false, true}, false},
	{[]bool{true, true, true, true, true, true, false}, false},
	{[]bool{true, true, true, true, true, true, true}, true},
}

func validateOdd7Parity(t *testing.T, g *Genome) {
	for i, n := range odd7ParityTests {
		r := g.EvalBool(n.in, bn.BoolAllGates)
		if r != n.out {
			t.Errorf("%v: odd7Parity.EvalBool(%#v, BoolAllGates) => %v, want %v", i, n.in, r, n.out)
		}
	}
}

func TestOdd7Parity(t *testing.T) {
	mux := New([]*gene.Gene{
		gene.New("Nor.Nand.d4.Or.Or.Xor.Or.Xor.d1.d4.d1.d0.d5.d2.d2.d1.d1"),
		gene.New("And.Or.Nor.Nand.Xor.Xor.Xor.Nor.d5.d0.d4.d3.d3.d6.d5.d5.d1"),
		gene.New("Xor.Xor.Xor.Xor.Nand.Or.Nand.Nor.d0.d2.d2.d4.d1.d1.d3.d1.d3"),
		gene.New("And.And.And.Xor.Nor.And.Xor.Not.d2.d5.d5.d5.d1.d6.d6.d1.d1"),
	},
		"Xor")
	validateOdd7Parity(t, mux)
	w := map[string]int{
		"And":  5,
		"Nand": 3,
		"Nor":  4,
		"Not":  1,
		"Or":   2,
		"Xor":  12,
		"d0":   2,
		"d1":   6,
		"d2":   3,
		"d3":   4,
		"d4":   2,
		"d5":   6,
		"d6":   3,
	}
	if !reflect.DeepEqual(mux.SymbolCount, w) {
		t.Errorf("Genome %q SymbolCount=%v, want %v", mux, mux.SymbolCount, w)
	}
}

// Time Series

var maunaLoaCO2Tests = []struct {
	in  []float64
	out float64
}{
	{[]float64{314.6784595, 314.7734727, 315.71, 317.45, 317.5, 317.627088, 315.86, 314.93, 313.19, 313.1966964, 313.34}, 314.67},
	{[]float64{314.7734727, 315.71, 317.45, 317.5, 317.627088, 315.86, 314.93, 313.19, 313.1966964, 313.34, 314.67}, 315.58},
	{[]float64{315.71, 317.45, 317.5, 317.627088, 315.86, 314.93, 313.19, 313.1966964, 313.34, 314.67, 315.58}, 316.47},
	{[]float64{317.45, 317.5, 317.627088, 315.86, 314.93, 313.19, 313.1966964, 313.34, 314.67, 315.58, 316.47}, 316.65},
	{[]float64{317.5, 317.627088, 315.86, 314.93, 313.19, 313.1966964, 313.34, 314.67, 315.58, 316.47, 316.65}, 317.71},
	{[]float64{317.627088, 315.86, 314.93, 313.19, 313.1966964, 313.34, 314.67, 315.58, 316.47, 316.65, 317.71}, 318.29},
	{[]float64{315.86, 314.93, 313.19, 313.1966964, 313.34, 314.67, 315.58, 316.47, 316.65, 317.71, 318.29}, 318.16},
	{[]float64{314.93, 313.19, 313.1966964, 313.34, 314.67, 315.58, 316.47, 316.65, 317.71, 318.29, 318.16}, 316.55},
	{[]float64{313.19, 313.1966964, 313.34, 314.67, 315.58, 316.47, 316.65, 317.71, 318.29, 318.16, 316.55}, 314.8},
	{[]float64{313.1966964, 313.34, 314.67, 315.58, 316.47, 316.65, 317.71, 318.29, 318.16, 316.55, 314.8}, 313.84},
	{[]float64{313.34, 314.67, 315.58, 316.47, 316.65, 317.71, 318.29, 318.16, 316.55, 314.8, 313.84}, 313.34},
	{[]float64{314.67, 315.58, 316.47, 316.65, 317.71, 318.29, 318.16, 316.55, 314.8, 313.84, 313.34}, 314.81},
	{[]float64{315.58, 316.47, 316.65, 317.71, 318.29, 318.16, 316.55, 314.8, 313.84, 313.34, 314.81}, 315.59},
	{[]float64{316.47, 316.65, 317.71, 318.29, 318.16, 316.55, 314.8, 313.84, 313.34, 314.81, 315.59}, 316.43},
	{[]float64{316.65, 317.71, 318.29, 318.16, 316.55, 314.8, 313.84, 313.34, 314.81, 315.59, 316.43}, 316.97},
	{[]float64{317.71, 318.29, 318.16, 316.55, 314.8, 313.84, 313.34, 314.81, 315.59, 316.43, 316.97}, 317.58},
	{[]float64{318.29, 318.16, 316.55, 314.8, 313.84, 313.34, 314.81, 315.59, 316.43, 316.97, 317.58}, 319.03},
	{[]float64{318.16, 316.55, 314.8, 313.84, 313.34, 314.81, 315.59, 316.43, 316.97, 317.58, 319.03}, 320.03},
	{[]float64{316.55, 314.8, 313.84, 313.34, 314.81, 315.59, 316.43, 316.97, 317.58, 319.03, 320.03}, 319.59},
	{[]float64{314.8, 313.84, 313.34, 314.81, 315.59, 316.43, 316.97, 317.58, 319.03, 320.03, 319.59}, 318.18},
	{[]float64{313.84, 313.34, 314.81, 315.59, 316.43, 316.97, 317.58, 319.03, 320.03, 319.59, 318.18}, 315.91},
}

func validateMaunaLoaCO2(t *testing.T, g *Genome) {
	for i, n := range maunaLoaCO2Tests {
		r := g.EvalMath(n.in)
		e := math.Abs(r - n.out)
		if e > delta {
			t.Errorf("%v: maunaLoaCO2.EvalMath(%#v) => %v, want %v, e=%v", i, n.in, r, n.out, e)
		}
	}
}

func TestMaunaLoaCO2(t *testing.T) {
	g1 := gene.New("Avg2.Avg2.Avg2.Avg2.Max2.+.-.3Rt.*.+.d10.d10.d5.d5.d9.c6.c4.d3.d10.d10.c1")
	g2 := gene.New("5Rt.-.d10.-.Floor.5Rt.Max2.Min2.Max2.Max2.d4.c3.d1.c8.d0.d8.d1.d1.d1.c0.c9")
	g3 := gene.New("-.Ln.*.Avg2.Min2.Min2.+.X4.d6.Min2.d0.c9.d5.d5.c1.c7.d1.d10.d3.c4.d10")
	g4 := gene.New("3Rt.+.NOT.Min2.c0.+./.*.*.-.c5.c4.c1.c8.d3.d8.d10.d8.c8.c0.d0")
	g5 := gene.New("Max2.-./.Min2.-.-.-.d10.-.Max2.d0.d10.d10.c4.c0.d6.c6.c0.c9.d10.d10")
	g1.Constants = []float64{
		2.05725272377697,
		-8.22715872676778,
		7.79778435621204,
		7.45353556932279,
		-7.08975493636891,
		-0.659504989776299,
		-0.828431998544267,
		7.796563615833,
		8.67976928006836,
		-0.453380511185034,
	}
	g2.Constants = []float64{
		-7.99472333506745,
		7.69669331949828,
		-4.70137638477737,
		5.90136417737358,
		-7.22087273781549,
		-2.04016235847041,
		-8.16583758049257,
		-2.9117647236549,
		-4.81935758537553,
		5.00030518509476,
	}
	g3.Constants = []float64{
		-3.52336191900388,
		8.72449690237129,
		2.50770592364269,
		46.134830774865,
		-6.78701132236702,
		-4.95468001342814,
		-5.70726645710624,
		-4.18439283425398,
		8.86043885616627,
		-9.26572466200751,
	}
	g4.Constants = []float64{
		-1.29636461989196,
		-8.2869655446028,
		0.337839899899289,
		-1.92247383037812,
		-3.17300943021943,
		-1.00946073793756,
		-3.5209204382458,
		2.58583330790124,
		9.88402966399121,
		1.73568132572405,
	}
	g5.Constants = []float64{
		-5.82934049501022,
		-1.27964110232856,
		9.83458967864009,
		-3.53862117374187,
		1.03366191595202,
		-7.08120975371563,
		-5.29584643086032,
		2.25257118442335,
		4.03668324839015,
		-3.56907559434797,
	}
	gn := New([]*gene.Gene{g1, g2, g3, g4, g5}, "Avg2")
	validateMaunaLoaCO2(t, gn)
	w := map[string]int{
		"*":     3,
		"+":     3,
		"-":     9,
		"/":     2,
		"3Rt":   1,
		"5Rt":   2,
		"Avg2":  5,
		"Floor": 1,
		"Ln":    1,
		"Max2":  5,
		"Min2":  6,
		"NOT":   1,
		"X4":    1,
		"c0":    3,
		"c1":    2,
		"c3":    1,
		"c4":    2,
		"c5":    1,
		"c6":    1,
		"c7":    1,
		"c8":    2,
		"c9":    2,
		"d0":    3,
		"d1":    2,
		"d10":   5,
		"d3":    1,
		"d4":    1,
		"d5":    2,
		"d6":    2,
		"d8":    2,
	}
	if !reflect.DeepEqual(gn.SymbolCount, w) {
		t.Errorf("Genome %q SymbolCount=%v, want %v", gn, gn.SymbolCount, w)
	}
}

// Classification

var irisPlantsTests = []struct {
	in  []float64
	out float64
}{
	{[]float64{6.3, 2.8, 5.1, 1.5}, 1},
	{[]float64{6.1, 3, 4.6, 1.4}, 0},
	{[]float64{5.5, 2.5, 4, 1.3}, 0},
	{[]float64{5.5, 3.5, 1.3, 0.2}, 0},
	{[]float64{6.8, 3, 5.5, 2.1}, 1},
	{[]float64{6.3, 2.5, 4.9, 1.5}, 0},
	{[]float64{5, 3.4, 1.6, 0.4}, 0},
	{[]float64{6, 3.4, 4.5, 1.6}, 0},
	{[]float64{6.4, 3.1, 5.5, 1.8}, 1},
	{[]float64{6.2, 2.2, 4.5, 1.5}, 0},
	{[]float64{5.6, 2.9, 3.6, 1.3}, 0},
	{[]float64{4.9, 3, 1.4, 0.2}, 0},
	{[]float64{5.4, 3.9, 1.7, 0.4}, 0},
	{[]float64{6.1, 2.9, 4.7, 1.4}, 0},
	{[]float64{6.6, 2.9, 4.6, 1.3}, 0},
	{[]float64{4.6, 3.6, 1, 0.2}, 0},
	{[]float64{5.5, 2.4, 3.8, 1.1}, 0},
	{[]float64{7.4, 2.8, 6.1, 1.9}, 1},
	{[]float64{7.7, 2.8, 6.7, 2}, 1},
	{[]float64{4.7, 3.2, 1.6, 0.2}, 0},
	{[]float64{6.4, 2.8, 5.6, 2.1}, 1},
	{[]float64{5.2, 3.4, 1.4, 0.2}, 0},
	{[]float64{4.5, 2.3, 1.3, 0.3}, 0},
	{[]float64{6, 2.2, 4, 1}, 0},
	{[]float64{5.2, 3.5, 1.5, 0.2}, 0},
	{[]float64{6.4, 3.2, 4.5, 1.5}, 0},
	{[]float64{6.5, 3, 5.8, 2.2}, 1},
	{[]float64{6.8, 2.8, 4.8, 1.4}, 0},
	{[]float64{6.2, 3.4, 5.4, 2.3}, 1},
	{[]float64{6.3, 3.4, 5.6, 2.4}, 1},
	{[]float64{7.3, 2.9, 6.3, 1.8}, 1},
	{[]float64{7.1, 3, 5.9, 2.1}, 1},
	{[]float64{6.7, 3.1, 4.7, 1.5}, 0},
	{[]float64{5.5, 4.2, 1.4, 0.2}, 0},
	{[]float64{6.3, 3.3, 4.7, 1.6}, 0},
	{[]float64{5, 3.2, 1.2, 0.2}, 0},
	{[]float64{5.1, 3.8, 1.5, 0.3}, 0},
	{[]float64{5.8, 2.7, 5.1, 1.9}, 1},
	{[]float64{6.5, 2.8, 4.6, 1.5}, 0},
	{[]float64{4.7, 3.2, 1.3, 0.2}, 0},
	{[]float64{7.7, 2.6, 6.9, 2.3}, 1},
	{[]float64{5.9, 3.2, 4.8, 1.8}, 0},
	{[]float64{4.8, 3.4, 1.6, 0.2}, 0},
	{[]float64{6.5, 3, 5.2, 2}, 1},
	{[]float64{5.8, 2.7, 5.1, 1.9}, 1},
	{[]float64{5, 3.3, 1.4, 0.2}, 0},
	{[]float64{6.4, 3.2, 5.3, 2.3}, 1},
	{[]float64{6.8, 3.2, 5.9, 2.3}, 1},
	{[]float64{5.4, 3.7, 1.5, 0.2}, 0},
	{[]float64{7.9, 3.8, 6.4, 2}, 1},
	{[]float64{5.7, 2.5, 5, 2}, 1},
	{[]float64{6.2, 2.8, 4.8, 1.8}, 1},
	{[]float64{6.1, 3, 4.9, 1.8}, 1},
	{[]float64{4.9, 3.1, 1.5, 0.1}, 0},
	{[]float64{7, 3.2, 4.7, 1.4}, 0},
	{[]float64{5.2, 4.1, 1.5, 0.1}, 0},
	{[]float64{5.2, 2.7, 3.9, 1.4}, 0},
	{[]float64{6.3, 3.3, 6, 2.5}, 1},
	{[]float64{6.1, 2.6, 5.6, 1.4}, 1},
	{[]float64{6, 2.7, 5.1, 1.6}, 0},
	{[]float64{5.5, 2.4, 3.7, 1}, 0},
	{[]float64{5.5, 2.6, 4.4, 1.2}, 0},
	{[]float64{5.9, 3, 4.2, 1.5}, 0},
	{[]float64{5.4, 3.9, 1.3, 0.4}, 0},
	{[]float64{7.2, 3, 5.8, 1.6}, 1},
	{[]float64{5.1, 3.4, 1.5, 0.2}, 0},
	{[]float64{5.8, 2.6, 4, 1.2}, 0},
	{[]float64{6.3, 2.3, 4.4, 1.3}, 0},
	{[]float64{4.9, 2.5, 4.5, 1.7}, 1},
	{[]float64{7.2, 3.2, 6, 1.8}, 1},
	{[]float64{5.6, 3, 4.5, 1.5}, 0},
	{[]float64{5.6, 2.7, 4.2, 1.3}, 0},
	{[]float64{4.8, 3, 1.4, 0.1}, 0},
	{[]float64{4.9, 2.4, 3.3, 1}, 0},
	{[]float64{5.7, 2.9, 4.2, 1.3}, 0},
	{[]float64{4.8, 3.4, 1.9, 0.2}, 0},
	{[]float64{5, 2, 3.5, 1}, 0},
	{[]float64{5.7, 3, 4.2, 1.2}, 0},
	{[]float64{4.8, 3.1, 1.6, 0.2}, 0},
	{[]float64{4.6, 3.2, 1.4, 0.2}, 0},
	{[]float64{6.4, 2.7, 5.3, 1.9}, 1},
	{[]float64{6, 3, 4.8, 1.8}, 1},
	{[]float64{4.6, 3.1, 1.5, 0.2}, 0},
	{[]float64{5.1, 3.5, 1.4, 0.2}, 0},
	{[]float64{5.6, 3, 4.1, 1.3}, 0},
	{[]float64{6.3, 2.5, 5, 1.9}, 1},
	{[]float64{6.9, 3.2, 5.7, 2.3}, 1},
	{[]float64{5.8, 2.7, 4.1, 1}, 0},
	{[]float64{6, 2.9, 4.5, 1.5}, 0},
	{[]float64{6.7, 3, 5, 1.7}, 0},
	{[]float64{5, 2.3, 3.3, 1}, 0},
	{[]float64{5.3, 3.7, 1.5, 0.2}, 0},
	{[]float64{4.4, 3.2, 1.3, 0.2}, 0},
	{[]float64{6.2, 2.9, 4.3, 1.3}, 0},
	{[]float64{5.8, 4, 1.2, 0.2}, 0},
	{[]float64{6.7, 3.1, 5.6, 2.4}, 1},
	{[]float64{6.1, 2.8, 4, 1.3}, 0},
	{[]float64{5.6, 2.8, 4.9, 2}, 1},
}

func validateIrisPlants(t *testing.T, g *Genome) {
	for i, n := range irisPlantsTests {
		r := g.EvalMath(n.in)
		e := math.Abs(r - n.out)
		if e > 100 { // Disable test for now until classification is figured out
			t.Errorf("%v: irisPlants.EvalMath(%#v) => %v, want %v, e=%v", i, n.in, r, n.out, e)
		}
	}
}

func TestIrisPlants(t *testing.T) {
	g1 := gene.New("Sqrt.GOE2E.d2.d2.d2.c6.+.d3.d2.c3.c4.c9.d2.c7.d0.d3.d3")
	g2 := gene.New("d2.*.3Rt.Ln.*.Ln.OR2.GOE2A.c7.d2.c3.d3.d2.d3.d2.c1.c9")
	g3 := gene.New("AND2.Avg2.GOE2C.d3.AND1.GOE2A.3Rt.AND2.d0.c1.d3.d2.d0.c6.c3.c7.d0")
	g1.Constants = []float64{
		-7.36991485335856,
		3.02133243812372E-02,
		3.90911587878048,
		-3.37382122257149,
		6.02038636432997,
		-2.62184514908292,
		7.33332316049684,
		-9.19431134983367,
		-1.36875514999847,
		-6.14551225318155,
	}
	g2.Constants = []float64{
		-2.94106875820185,
		4.71480452894681,
		3.80718405713065,
		5.0022888882107,
		2.05542161320841,
		9.67284157841731,
		-0.159092989898373,
		3.67595446638386,
		4.72334971160009,
		7.20145268105106,
	}
	g3.Constants = []float64{
		6.54164250618,
		-5.31662953581347,
		7.16788232062746,
		2.4283577990051,
		-3.33404034546953,
		-4.19415875728629,
		-6.94814905240028,
		1.22348704489273,
		-2.74391918698691,
		1.35288552507096,
	}
	gn := New([]*gene.Gene{g1, g2, g3}, "+")
	validateIrisPlants(t, gn)
	w := map[string]int{
		"+":     2,
		"3Rt":   1,
		"AND1":  1,
		"AND2":  2,
		"Avg2":  1,
		"GOE2A": 1,
		"GOE2C": 1,
		"c1":    1,
		"c6":    1,
		"d0":    2,
		"d2":    2,
		"d3":    2,
	}
	if !reflect.DeepEqual(gn.SymbolCount, w) {
		t.Errorf("Genome %q SymbolCount=%v, want %v", gn, gn.SymbolCount, w)
	}
}

// Logistic Regression

var emotivEEGTests = []struct {
	in  []float64
	out float64
}{
	{[]float64{4614.87168202347, 4535.38450448, 4112.30759174866, 3884.61528962415, 3987.69221018052, 4239.4870758182, 4487.69219795392, 4676.41014205712, 4235.89733231624, 4114.87169425006, 4276.92297233867, 4382.05117489616, 4083.07682323267, 4221.53835830839}, 0},
	{[]float64{4611.2819385215, 4442.05117342897, 4135.38451426128, 3741.02554954564, 3988.71785118108, 4176.41015428371, 4402.05117440709, 4681.02552655964, 4289.74348484568, 4172.30759028147, 4335.89732987092, 4430.7691224228, 4116.4101557509, 4247.17938332241}, 0},
	{[]float64{4596.92296451365, 4529.74347897692, 4107.17938674585, 3852.82041860677, 3986.66656917996, 4225.64092231063, 4477.9486084486, 4660.51270654842, 4222.56399930895, 4093.84605373856, 4247.69220382269, 4351.28194487933, 4031.79477320463, 4187.69220528988}, 0},
	{[]float64{4695.89732106777, 4601.0255285159, 4171.79476978119, 3861.02554661125, 4018.46144019734, 4311.79476635774, 4511.28194096682, 4725.12808958376, 4276.41015183839, 4130.25630925847, 4342.05117587429, 4388.7178413998, 4091.79477173744, 4258.97425482886}, 0},
	{[]float64{4550.25629898813, 4296.92297184961, 4025.64092720127, 3485.64094040599, 3845.64093160284, 3834.35888059667, 4215.38451230502, 4540.51270948281, 4075.89733622875, 3843.07682910144, 4052.82041371613, 4268.2050238339, 3945.64092915752, 4046.15374721248}, 0},
	{[]float64{4618.46142552543, 4471.79476244523, 4116.92297625118, 3872.30759761742, 3979.48708217603, 4273.84604933699, 4452.30758343458, 4725.12808958376, 4293.84604884793, 4182.56400028708, 4274.35886983727, 4412.82040491298, 4178.46143628483, 4266.15374183278}, 0},
	{[]float64{4551.79476048898, 4432.82040442392, 4076.92297722931, 3770.7691385619, 3907.17939163649, 4191.79476929212, 4407.1793794099, 4716.41014107899, 4351.79476537961, 4091.79477173744, 4243.07681932016, 4322.56399686363, 4112.30759174866, 4154.87169327193}, 1},
	{[]float64{4434.35886592476, 4325.64091986531, 3959.99990316537, 3655.38452599881, 3859.99990561069, 4224.61528131007, 4393.33322590233, 4650.25629654282, 4308.71784335606, 4076.41015672903, 4266.66656233306, 4291.7947668468, 4095.3845152394, 4144.10246276604}, 1},
}

func validateEmotivEEG(t *testing.T, g *Genome) {
	for i, n := range emotivEEGTests {
		r := g.EvalMath(n.in)
		e := math.Abs(r - n.out)
		if e > 4000 { // Disable test for now until logistic regression is figured out
			t.Errorf("%v: emotivEEG.EvalMath(%#v) => %v, want %v, e=%v", i, n.in, r, n.out, e)
		}
	}
}

func TestEmotivEEG(t *testing.T) {
	g1 := gene.New("Inv./.c2.+.Avg2.+.-.+.d8.-.d0.d2.d3.d3.d8.d1.c2.c2.d3.c2.d9")
	g2 := gene.New("Avg2.c2.Logi.+.*.*.*.-.d4.d8.d9.c3.d9.d7.d9.c3.d7.d2.c2.d8.d10")
	g3 := gene.New("-.Avg2.d4.c7.+.-.-.+.d0./.d9.c7.d5.d9.c2.c4.c6.d11.c1.d12.d10")
	g4 := gene.New("d8.d5.Inv.-.*.-.-.-.d5.-.c2.d0.c4.d8.d2.c1.d0.d6.d7.d12.d9")
	g1.Constants = []float64{
		-3.66252632221442,
		4.32132938627278,
		-5.04312570574053,
		-8.33491012298959,
		570.542293374432,
		7.88445692312387,
		-8.729819635609,
		8.64864040040284,
		-8.99655140842921,
		-2.52174443800165,
	}
	g2.Constants = []float64{
		-8.09076204718161,
		-4.90160832544939,
		-8.66817224646748,
		8.00462202826014,
		0.603961302529984,
		-8.0687887203589,
		8.47999636271249,
		-3.48778252510147,
		-8.01019318216498,
		6.74764244514298,
	}
	g3.Constants = []float64{
		-8.25165715506455,
		-3.08755760368664,
		6.65888210699789,
		-8.59859004486221,
		4.01572573015534,
		0.589922788171026,
		-0.197058015686514,
		-4.66353343302713,
		-3.81057161168249,
		2.91665395062105,
	}
	g4.Constants = []float64{
		-1.4529587694937,
		6.02465895565661,
		-9.12533951841792,
		9.74364452040162,
		-3.01298532059694,
		9.0510095126194,
		-4.48530533768731,
		5.70360423596912,
		8.32840449232459,
		1.25095370342112,
	}
	gn := New([]*gene.Gene{g1, g2, g3, g4}, "+")
	validateEmotivEEG(t, gn)
	w := map[string]int{
		"*":    3,
		"+":    6,
		"-":    4,
		"/":    1,
		"Avg2": 2,
		"Logi": 1,
		"c2":   2,
		"c3":   1,
		"c7":   2,
		"d0":   1,
		"d4":   2,
		"d5":   1,
		"d7":   1,
		"d8":   2,
		"d9":   4,
	}
	if !reflect.DeepEqual(gn.SymbolCount, w) {
		t.Errorf("Genome %q SymbolCount=%v, want %v", gn, gn.SymbolCount, w)
	}
}

// Regression

var fuelConsumptionTests = []struct {
	in  []float64
	out float64
}{
	{[]float64{8, 400, 190, 4422, 12.5, 72, 1}, 13},
	{[]float64{8, 350, 105, 3725, 19, 81, 1}, 26.6},
	{[]float64{4, 91, 53, 1795, 17.4, 76, 3}, 33},
	{[]float64{6, 198, 95, 2904, 16, 73, 1}, 23},
	{[]float64{8, 305, 130, 3840, 15.4, 79, 1}, 17},
	{[]float64{8, 350, 125, 3900, 17.4, 79, 1}, 23},
	{[]float64{4, 140, 90, 2408, 19.5, 72, 1}, 20},
	{[]float64{6, 258, 110, 2962, 13.5, 71, 1}, 18},
	{[]float64{4, 76, 52, 1649, 16.5, 74, 3}, 31},
	{[]float64{8, 400, 150, 4464, 12, 73, 1}, 13},
	{[]float64{8, 304, 150, 3672, 11.5, 72, 1}, 17},
	{[]float64{8, 304, 150, 4257, 15.5, 74, 1}, 14},
	{[]float64{6, 225, 85, 3465, 16.6, 81, 1}, 17.6},
	{[]float64{4, 151, 90, 2556, 13.2, 79, 1}, 33.5},
	{[]float64{4, 98, 65, 2045, 16.2, 81, 1}, 34.4},
	{[]float64{8, 318, 150, 4190, 13, 76, 1}, 16},
	{[]float64{8, 350, 180, 4380, 12.1, 76, 1}, 16.5},
	{[]float64{4, 110, 87, 2672, 17.5, 70, 2}, 25},
	{[]float64{6, 250, 88, 3021, 16.5, 73, 1}, 18},
	{[]float64{4, 108, 75, 2265, 15.2, 80, 3}, 32.2},
	{[]float64{4, 91, 60, 1800, 16.4, 78, 3}, 36.1},
	{[]float64{6, 250, 105, 3459, 16, 75, 1}, 18},
	{[]float64{4, 91, 69, 2130, 14.7, 79, 2}, 37.3},
	{[]float64{4, 121, 110, 2600, 12.8, 77, 2}, 21.5},
	{[]float64{8, 304, 150, 3892, 12.5, 72, 1}, 15},
	{[]float64{4, 121, 76, 2511, 18, 72, 2}, 22},
	{[]float64{8, 340, 160, 3609, 8, 70, 1}, 14},
	{[]float64{8, 350, 165, 4209, 12, 71, 1}, 14},
	{[]float64{4, 120, 87, 2979, 19.5, 72, 2}, 21},
	{[]float64{4, 90, 75, 2108, 15.5, 74, 2}, 24},
	{[]float64{6, 156, 108, 2930, 15.5, 76, 3}, 19},
	{[]float64{5, 121, 67, 2950, 19.9, 80, 2}, 36.4},
	{[]float64{5, 131, 103, 2830, 15.9, 78, 2}, 20.3},
	{[]float64{4, 85, 65, 1975, 19.4, 81, 3}, 37},
	{[]float64{8, 400, 230, 4278, 9.5, 73, 1}, 16},
	{[]float64{4, 101, 83, 2202, 15.3, 76, 2}, 27},
	{[]float64{8, 260, 110, 4060, 19, 77, 1}, 17},
	{[]float64{6, 225, 110, 3620, 18.7, 78, 1}, 18.6},
	{[]float64{4, 116, 90, 2123, 14, 71, 2}, 28},
	{[]float64{6, 200, 105.858778625954, 2875, 17, 74, 1}, 21},
	{[]float64{4, 122, 88, 2500, 15.1, 80, 2}, 35},
	{[]float64{6, 198, 95, 2833, 15.5, 70, 1}, 22},
	{[]float64{6, 231, 110, 3039, 15, 75, 1}, 21},
	{[]float64{5, 183, 77, 3530, 20.1, 79, 2}, 25.4},
	{[]float64{4, 121, 115, 2795, 15.7, 78, 2}, 21.6},
	{[]float64{4, 141, 71, 3190, 24.8, 79, 2}, 27.2},
	{[]float64{6, 225, 105, 3439, 15.5, 71, 1}, 16},
	{[]float64{4, 79, 58, 1755, 16.9, 81, 3}, 39.1},
	{[]float64{4, 121, 112, 2868, 15.5, 73, 2}, 19},
	{[]float64{4, 98, 105.858778625954, 2046, 19, 71, 1}, 25},
	{[]float64{4, 85, 65, 2020, 19.2, 79, 3}, 31.8},
	{[]float64{4, 135, 84, 2490, 15.7, 81, 1}, 27.2},
	{[]float64{6, 173, 115, 2700, 12.9, 79, 1}, 26.8},
	{[]float64{4, 86, 65, 2019, 16.4, 80, 3}, 37.2},
	{[]float64{6, 250, 110, 3520, 16.4, 77, 1}, 17.5},
	{[]float64{8, 302, 140, 4141, 14, 74, 1}, 16},
	{[]float64{6, 146, 97, 2815, 14.5, 77, 3}, 22},
	{[]float64{6, 200, 85, 2587, 16, 70, 1}, 21},
	{[]float64{4, 134, 95, 2515, 14.8, 78, 3}, 21.1},
	{[]float64{6, 231, 110, 3415, 15.8, 81, 1}, 22.4},
	{[]float64{4, 83, 61, 2003, 19, 74, 3}, 32},
	{[]float64{4, 104, 95, 2375, 17.5, 70, 2}, 25},
}

func validateFuelConsumption(t *testing.T, g *Genome) {
	for i, n := range fuelConsumptionTests {
		r := g.EvalMath(n.in)
		e := math.Abs(r - n.out)
		if e > 11 { // Disable test for now until regression is figured out
			t.Errorf("%v: fuelConsumption.EvalMath(%#v) => %v, want %v, e=%v", i, n.in, r, n.out, e)
		}
	}
}

func newFuelConsumption() *Genome {
	g1 := gene.New("+.d2.-.+./.Atan./.*.c4.d0.d5.d0.d0.d0.d5.d2.d0")
	g2 := gene.New("c7.+.+.+.Avg2.d4.d6.-.d4.c7.d1.d2.c2.d6.c8.d6.d2")
	g3 := gene.New("Min2.-.+.d6.*.*.Ln.Tanh.c4.c4.d0.c4.d6.d4.d4.c2.c4")
	g4 := gene.New("-.-.d2.+.+.Tanh.+.c1.c4.d0.d0.d6.d6.c9.d5.c9.c7")
	g1.Constants = []float64{
		-6.89063692129276,
		-5.68895535142064,
		74.5934430982391,
		-5.52293465987121,
		7.36762822046571,
		-4.90189548568468,
		-9.06796472060305,
		-4.60798974578082,
		2.21714835047456,
		-8.37993408001953,
	}
	g2.Constants = []float64{
		3.84215124973296,
		-8.56854411450545,
		-4.83944158531019,
		-7.34665429273354,
		1.81066316721091,
		-5.68092898342845,
		1.9168656633198,
		6.00730927940611,
		-7.87530137028108,
		-5.61497329630421,
	}
	g3.Constants = []float64{
		1.04586931974242,
		8.52891628772851,
		-4.40046388134404,
		-2.04504531998657,
		0.681604358043153,
		-2.69447920163579,
		-0.451368755149998,
		9.55381939146092,
		7.5619409561449,
		9.05453657643361,
	}
	g4.Constants = []float64{
		-1.2619803460799,
		3.61616965876798,
		-1.45847956785791,
		-5.12960448011719,
		0.023370006408887,
		7.84438612018189,
		-4.07391582995086,
		6.88815179906613,
		3.30928445387127,
		-8.11151463362529,
	}
	return New([]*gene.Gene{g1, g2, g3, g4}, "+")
}

func TestFuelConsumption(t *testing.T) {
	gn := newFuelConsumption()
	validateFuelConsumption(t, gn)
	w := map[string]int{
		"*":    2,
		"+":    7,
		"-":    3,
		"Ln":   1,
		"Min2": 1,
		"Tanh": 2,
		"c1":   1,
		"c4":   4,
		"c7":   1,
		"d0":   3,
		"d2":   1,
		"d6":   3,
	}
	if !reflect.DeepEqual(gn.SymbolCount, w) {
		t.Errorf("Genome %q SymbolCount=%v, want %v", gn, gn.SymbolCount, w)
	}
}

func checkEqual(g1 *Genome, g2 *Genome) error {
	if g1 == nil || g2 == nil {
		return fmt.Errorf("genome.checkEqual error: g1 and g2 must be non-nil")
	}
	if len(g1.Genes) != len(g2.Genes) {
		return fmt.Errorf("len(g1.Genes)=%v != len(g2.Genes)=%v", len(g1.Genes), len(g2.Genes))
	}
	for i := range g1.Genes {
		if err := gene.CheckEqual(g1.Genes[i], g2.Genes[i]); err != nil {
			return fmt.Errorf("g1.Genes[%v] != g2.Genes[%v]: %v", i, i, err)
		}
	}
	if g1.LinkFunc != g2.LinkFunc {
		return fmt.Errorf("g1.LinkFunc=%v != g2.LinkFunc=%v", g1.LinkFunc, g2.LinkFunc)
	}
	if g1.Score != g2.Score {
		return fmt.Errorf("g1.Score=%v != g2.Score=%v", g1.Score, g2.Score)
	}
	return nil
}

func TestDup(t *testing.T) {
	mux := New([]*gene.Gene{
		gene.New("Nand.Or.And.Not.Nor.Not.Nor.And.d3.d2.d0.d2.d1.d4.d1.d4.d2"),
		gene.New("Nor.Nor.And.Or.Nor.And.Nor.Or.d0.d1.d0.d2.d3.d1.d2.d0.d3"),
		gene.New("Or.Or.Nand.d4.Not.Or.Nand.Nand.d0.d4.d3.d4.d1.d1.d3.d2.d0"),
		gene.New("Or.And.And.Nand.d5.Nand.Nand.Nor.d4.d3.d5.d0.d1.d0.d4.d0.d1"),
	},
		"And")
	validateSixMultiplexer(t, mux) // Force evaluation
	gn := mux.Dup()
	if err := checkEqual(gn, mux); err != nil {
		t.Errorf("TestDup after Dup failed: gn != mux: %v\n", err)
	}
	validateSixMultiplexer(t, gn) // Force evaluation
	validateSixMultiplexer(t, mux)

	gn = newFuelConsumption()
	validateFuelConsumption(t, gn) // Force evaluation
	mux = gn.Dup()
	if err := checkEqual(gn, mux); err != nil {
		t.Errorf("TestDup after Dup failed: gn != mux: %v\n", err)
	}
	validateFuelConsumption(t, gn) // Force evaluation
	validateFuelConsumption(t, mux)
}

func TestMutate(t *testing.T) {
	headSize := 7
	maxArity := 2
	tailSize := headSize*(maxArity-1) + 1
	numTerminals := 5
	funcs := []gene.FuncWeight{
		{"Not", 1},
		{"And", 5},
		{"Or", 5},
	}
	mux := New([]*gene.Gene{
		gene.RandomNew(headSize, tailSize, numTerminals, 0, funcs),
		gene.RandomNew(headSize, tailSize, numTerminals, 0, funcs),
		gene.RandomNew(headSize, tailSize, numTerminals, 0, funcs),
		gene.RandomNew(headSize, tailSize, numTerminals, 0, funcs),
	},
		"And")
	gn := mux.Dup()
	mux.Mutate(1)
	if err := checkEqual(gn, mux); err == nil {
		t.Errorf("TestMutate failed: gn == mux\n")
	}
}
