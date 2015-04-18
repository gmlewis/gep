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

var result bool

func BenchmarkNot(b *testing.B) {
	var v bool
	f := BoolAllGates["Not"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkAnd(b *testing.B) {
	var v bool
	f := BoolAllGates["And"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkOr(b *testing.B) {
	var v bool
	f := BoolAllGates["Or"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNand(b *testing.B) {
	var v bool
	f := BoolAllGates["Nand"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNor(b *testing.B) {
	var v bool
	f := BoolAllGates["Nor"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkXor(b *testing.B) {
	var v bool
	f := BoolAllGates["Xor"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNxor(b *testing.B) {
	var v bool
	f := BoolAllGates["Nxor"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkAnd3(b *testing.B) {
	var v bool
	f := BoolAllGates["And3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkOr3(b *testing.B) {
	var v bool
	f := BoolAllGates["Or3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNand3(b *testing.B) {
	var v bool
	f := BoolAllGates["Nand3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNor3(b *testing.B) {
	var v bool
	f := BoolAllGates["Nor3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkOdd3(b *testing.B) {
	var v bool
	f := BoolAllGates["Odd3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkEven3(b *testing.B) {
	var v bool
	f := BoolAllGates["Even3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkAnd4(b *testing.B) {
	var v bool
	f := BoolAllGates["And4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkOr4(b *testing.B) {
	var v bool
	f := BoolAllGates["Or4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNand4(b *testing.B) {
	var v bool
	f := BoolAllGates["Nand4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNor4(b *testing.B) {
	var v bool
	f := BoolAllGates["Nor4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkOdd4(b *testing.B) {
	var v bool
	f := BoolAllGates["Odd4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkEven4(b *testing.B) {
	var v bool
	f := BoolAllGates["Even4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkId(b *testing.B) {
	var v bool
	f := BoolAllGates["Id"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkZero(b *testing.B) {
	var v bool
	f := BoolAllGates["Zero"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkOne(b *testing.B) {
	var v bool
	f := BoolAllGates["One"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLT(b *testing.B) {
	var v bool
	f := BoolAllGates["LT"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGT(b *testing.B) {
	var v bool
	f := BoolAllGates["GT"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLOE(b *testing.B) {
	var v bool
	f := BoolAllGates["LOE"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGOE(b *testing.B) {
	var v bool
	f := BoolAllGates["GOE"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNotA(b *testing.B) {
	var v bool
	f := BoolAllGates["NotA"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNotB(b *testing.B) {
	var v bool
	f := BoolAllGates["NotB"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkIdA(b *testing.B) {
	var v bool
	f := BoolAllGates["IdA"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkIdB(b *testing.B) {
	var v bool
	f := BoolAllGates["IdB"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkZero2(b *testing.B) {
	var v bool
	f := BoolAllGates["Zero2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkOne2(b *testing.B) {
	var v bool
	f := BoolAllGates["One2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLT3(b *testing.B) {
	var v bool
	f := BoolAllGates["LT3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGT3(b *testing.B) {
	var v bool
	f := BoolAllGates["GT3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLOE3(b *testing.B) {
	var v bool
	f := BoolAllGates["LOE3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGOE3(b *testing.B) {
	var v bool
	f := BoolAllGates["GOE3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkMux(b *testing.B) {
	var v bool
	f := BoolAllGates["Mux"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkIf(b *testing.B) {
	var v bool
	f := BoolAllGates["If"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkMaj(b *testing.B) {
	var v bool
	f := BoolAllGates["Maj"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkMin(b *testing.B) {
	var v bool
	f := BoolAllGates["Min"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func Benchmark2Off(b *testing.B) {
	var v bool
	f := BoolAllGates["2Off"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func Benchmark2On(b *testing.B) {
	var v bool
	f := BoolAllGates["2On"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3A1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3A1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3A2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3A2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3A3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3A3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3A4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3A4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3B1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3B1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3B2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3B2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3B3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3B3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3B4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3B4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3C1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3C1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3C2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3C2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3C3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3C3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3C4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3C4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3D1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3D1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3D2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3D2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3D3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3D3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3D4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3D4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3E1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3E1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3E2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3E2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3E3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3E3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3F1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3F1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3F2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3F2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3F3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3F3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3G1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3G1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3G2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3G2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3G3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3G3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3G4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3G4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3H1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3H1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3H2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3H2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3H3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3H3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM3H4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM3H4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLT3A(b *testing.B) {
	var v bool
	f := BoolAllGates["LT3A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGT3A(b *testing.B) {
	var v bool
	f := BoolAllGates["GT3A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLOE3A(b *testing.B) {
	var v bool
	f := BoolAllGates["LOE3A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGOE3A(b *testing.B) {
	var v bool
	f := BoolAllGates["GOE3A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkET3A(b *testing.B) {
	var v bool
	f := BoolAllGates["ET3A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNET3A(b *testing.B) {
	var v bool
	f := BoolAllGates["NET3A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLT3B(b *testing.B) {
	var v bool
	f := BoolAllGates["LT3B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGT3B(b *testing.B) {
	var v bool
	f := BoolAllGates["GT3B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLOE3B(b *testing.B) {
	var v bool
	f := BoolAllGates["LOE3B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGOE3B(b *testing.B) {
	var v bool
	f := BoolAllGates["GOE3B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkET3B(b *testing.B) {
	var v bool
	f := BoolAllGates["ET3B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNET3B(b *testing.B) {
	var v bool
	f := BoolAllGates["NET3B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLT3C(b *testing.B) {
	var v bool
	f := BoolAllGates["LT3C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGT3C(b *testing.B) {
	var v bool
	f := BoolAllGates["GT3C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLOE3C(b *testing.B) {
	var v bool
	f := BoolAllGates["LOE3C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGOE3C(b *testing.B) {
	var v bool
	f := BoolAllGates["GOE3C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkET3C(b *testing.B) {
	var v bool
	f := BoolAllGates["ET3C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNET3C(b *testing.B) {
	var v bool
	f := BoolAllGates["NET3C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT004(b *testing.B) {
	var v bool
	f := BoolAllGates["T004"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT008(b *testing.B) {
	var v bool
	f := BoolAllGates["T008"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT009(b *testing.B) {
	var v bool
	f := BoolAllGates["T009"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT032(b *testing.B) {
	var v bool
	f := BoolAllGates["T032"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT033(b *testing.B) {
	var v bool
	f := BoolAllGates["T033"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT041(b *testing.B) {
	var v bool
	f := BoolAllGates["T041"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT055(b *testing.B) {
	var v bool
	f := BoolAllGates["T055"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT057(b *testing.B) {
	var v bool
	f := BoolAllGates["T057"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT064(b *testing.B) {
	var v bool
	f := BoolAllGates["T064"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT065(b *testing.B) {
	var v bool
	f := BoolAllGates["T065"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT069(b *testing.B) {
	var v bool
	f := BoolAllGates["T069"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT073(b *testing.B) {
	var v bool
	f := BoolAllGates["T073"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT081(b *testing.B) {
	var v bool
	f := BoolAllGates["T081"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT089(b *testing.B) {
	var v bool
	f := BoolAllGates["T089"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT093(b *testing.B) {
	var v bool
	f := BoolAllGates["T093"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT096(b *testing.B) {
	var v bool
	f := BoolAllGates["T096"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT101(b *testing.B) {
	var v bool
	f := BoolAllGates["T101"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT109(b *testing.B) {
	var v bool
	f := BoolAllGates["T109"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT111(b *testing.B) {
	var v bool
	f := BoolAllGates["T111"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT121(b *testing.B) {
	var v bool
	f := BoolAllGates["T121"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT123(b *testing.B) {
	var v bool
	f := BoolAllGates["T123"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT125(b *testing.B) {
	var v bool
	f := BoolAllGates["T125"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT154(b *testing.B) {
	var v bool
	f := BoolAllGates["T154"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT223(b *testing.B) {
	var v bool
	f := BoolAllGates["T223"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT239(b *testing.B) {
	var v bool
	f := BoolAllGates["T239"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT249(b *testing.B) {
	var v bool
	f := BoolAllGates["T249"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT251(b *testing.B) {
	var v bool
	f := BoolAllGates["T251"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkT253(b *testing.B) {
	var v bool
	f := BoolAllGates["T253"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLT4(b *testing.B) {
	var v bool
	f := BoolAllGates["LT4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGT4(b *testing.B) {
	var v bool
	f := BoolAllGates["GT4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLOE4(b *testing.B) {
	var v bool
	f := BoolAllGates["LOE4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGOE4(b *testing.B) {
	var v bool
	f := BoolAllGates["GOE4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkTie(b *testing.B) {
	var v bool
	f := BoolAllGates["Tie"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNtie(b *testing.B) {
	var v bool
	f := BoolAllGates["Ntie"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func Benchmark3Off(b *testing.B) {
	var v bool
	f := BoolAllGates["3Off"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func Benchmark3On(b *testing.B) {
	var v bool
	f := BoolAllGates["3On"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4A1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4A1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4A2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4A2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4A3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4A3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4A4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4A4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4A5(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4A5"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4A6(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4A6"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4A7(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4A7"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4A8(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4A8"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4B1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4B1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4B2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4B2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4B3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4B3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4B4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4B4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4B5(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4B5"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4B6(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4B6"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4B7(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4B7"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4B8(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4B8"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4C1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4C1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4C2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4C2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4C3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4C3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4C4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4C4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4C5(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4C5"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4C6(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4C6"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4C7(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4C7"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4C8(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4C8"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4D1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4D1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4D2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4D2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4D3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4D3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4D4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4D4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4D5(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4D5"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4D6(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4D6"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4D7(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4D7"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4D8(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4D8"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4E1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4E1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4E2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4E2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4E3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4E3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4E4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4E4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4E5(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4E5"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4E6(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4E6"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4E7(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4E7"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4E8(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4E8"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4F1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4F1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4F2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4F2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4F3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4F3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4F4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4F4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4F5(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4F5"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4F6(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4F6"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4F7(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4F7"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4F8(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4F8"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4G1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4G1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4G2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4G2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4G3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4G3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4G4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4G4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4G5(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4G5"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4G6(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4G6"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4G7(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4G7"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4G8(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4G8"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4H1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4H1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4H2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4H2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4H3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4H3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4H4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4H4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4H5(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4H5"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4H6(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4H6"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4H7(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4H7"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4H8(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4H8"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4I1(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4I1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4I2(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4I2"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4I3(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4I3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4I4(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4I4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4I5(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4I5"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4I6(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4I6"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4I7(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4I7"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLM4I8(b *testing.B) {
	var v bool
	f := BoolAllGates["LM4I8"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLT4A(b *testing.B) {
	var v bool
	f := BoolAllGates["LT4A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGT4A(b *testing.B) {
	var v bool
	f := BoolAllGates["GT4A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLOE4A(b *testing.B) {
	var v bool
	f := BoolAllGates["LOE4A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGOE4A(b *testing.B) {
	var v bool
	f := BoolAllGates["GOE4A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkET4A(b *testing.B) {
	var v bool
	f := BoolAllGates["ET4A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNET4A(b *testing.B) {
	var v bool
	f := BoolAllGates["NET4A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLT4B(b *testing.B) {
	var v bool
	f := BoolAllGates["LT4B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGT4B(b *testing.B) {
	var v bool
	f := BoolAllGates["GT4B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLOE4B(b *testing.B) {
	var v bool
	f := BoolAllGates["LOE4B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGOE4B(b *testing.B) {
	var v bool
	f := BoolAllGates["GOE4B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkET4B(b *testing.B) {
	var v bool
	f := BoolAllGates["ET4B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNET4B(b *testing.B) {
	var v bool
	f := BoolAllGates["NET4B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLT4C(b *testing.B) {
	var v bool
	f := BoolAllGates["LT4C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGT4C(b *testing.B) {
	var v bool
	f := BoolAllGates["GT4C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLOE4C(b *testing.B) {
	var v bool
	f := BoolAllGates["LOE4C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGOE4C(b *testing.B) {
	var v bool
	f := BoolAllGates["GOE4C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkET4C(b *testing.B) {
	var v bool
	f := BoolAllGates["ET4C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNET4C(b *testing.B) {
	var v bool
	f := BoolAllGates["NET4C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLT4D(b *testing.B) {
	var v bool
	f := BoolAllGates["LT4D"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGT4D(b *testing.B) {
	var v bool
	f := BoolAllGates["GT4D"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLOE4D(b *testing.B) {
	var v bool
	f := BoolAllGates["LOE4D"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGOE4D(b *testing.B) {
	var v bool
	f := BoolAllGates["GOE4D"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkET4D(b *testing.B) {
	var v bool
	f := BoolAllGates["ET4D"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNET4D(b *testing.B) {
	var v bool
	f := BoolAllGates["NET4D"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLT4E(b *testing.B) {
	var v bool
	f := BoolAllGates["LT4E"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGT4E(b *testing.B) {
	var v bool
	f := BoolAllGates["GT4E"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkLOE4E(b *testing.B) {
	var v bool
	f := BoolAllGates["LOE4E"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkGOE4E(b *testing.B) {
	var v bool
	f := BoolAllGates["GOE4E"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkET4E(b *testing.B) {
	var v bool
	f := BoolAllGates["ET4E"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNET4E(b *testing.B) {
	var v bool
	f := BoolAllGates["NET4E"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ0002(b *testing.B) {
	var v bool
	f := BoolAllGates["Q0002"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ001C(b *testing.B) {
	var v bool
	f := BoolAllGates["Q001C"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ0048(b *testing.B) {
	var v bool
	f := BoolAllGates["Q0048"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ0800(b *testing.B) {
	var v bool
	f := BoolAllGates["Q0800"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ3378(b *testing.B) {
	var v bool
	f := BoolAllGates["Q3378"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ3475(b *testing.B) {
	var v bool
	f := BoolAllGates["Q3475"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ3CB0(b *testing.B) {
	var v bool
	f := BoolAllGates["Q3CB0"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ3DEF(b *testing.B) {
	var v bool
	f := BoolAllGates["Q3DEF"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ3DFF(b *testing.B) {
	var v bool
	f := BoolAllGates["Q3DFF"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ4200(b *testing.B) {
	var v bool
	f := BoolAllGates["Q4200"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ4C11(b *testing.B) {
	var v bool
	f := BoolAllGates["Q4C11"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ5100(b *testing.B) {
	var v bool
	f := BoolAllGates["Q5100"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ5EEF(b *testing.B) {
	var v bool
	f := BoolAllGates["Q5EEF"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ5EFF(b *testing.B) {
	var v bool
	f := BoolAllGates["Q5EFF"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ6A6D(b *testing.B) {
	var v bool
	f := BoolAllGates["Q6A6D"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ6F75(b *testing.B) {
	var v bool
	f := BoolAllGates["Q6F75"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ74C4(b *testing.B) {
	var v bool
	f := BoolAllGates["Q74C4"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ7DA3(b *testing.B) {
	var v bool
	f := BoolAllGates["Q7DA3"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ8304(b *testing.B) {
	var v bool
	f := BoolAllGates["Q8304"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ8430(b *testing.B) {
	var v bool
	f := BoolAllGates["Q8430"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ8543(b *testing.B) {
	var v bool
	f := BoolAllGates["Q8543"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQ9D80(b *testing.B) {
	var v bool
	f := BoolAllGates["Q9D80"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQA092(b *testing.B) {
	var v bool
	f := BoolAllGates["QA092"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQB36A(b *testing.B) {
	var v bool
	f := BoolAllGates["QB36A"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQCBCF(b *testing.B) {
	var v bool
	f := BoolAllGates["QCBCF"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQEEB1(b *testing.B) {
	var v bool
	f := BoolAllGates["QEEB1"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQEFFF(b *testing.B) {
	var v bool
	f := BoolAllGates["QEFFF"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQFF7B(b *testing.B) {
	var v bool
	f := BoolAllGates["QFF7B"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQFFF6(b *testing.B) {
	var v bool
	f := BoolAllGates["QFFF6"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkQFFFB(b *testing.B) {
	var v bool
	f := BoolAllGates["QFFFB"]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}
