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

func runBenchmark(b *testing.B, sym string) {
	var v bool
	f := BoolAllGates[sym]
	for i := 0; i < b.N; i++ {
		v = f.BoolFunction(false, true, true, false)
	}
	result = v
}

func BenchmarkNot(b *testing.B)   { runBenchmark(b, "Not") }
func BenchmarkAnd(b *testing.B)   { runBenchmark(b, "And") }
func BenchmarkOr(b *testing.B)    { runBenchmark(b, "Or") }
func BenchmarkNand(b *testing.B)  { runBenchmark(b, "Nand") }
func BenchmarkNor(b *testing.B)   { runBenchmark(b, "Nor") }
func BenchmarkXor(b *testing.B)   { runBenchmark(b, "Xor") }
func BenchmarkNxor(b *testing.B)  { runBenchmark(b, "Nxor") }
func BenchmarkAnd3(b *testing.B)  { runBenchmark(b, "And3") }
func BenchmarkOr3(b *testing.B)   { runBenchmark(b, "Or3") }
func BenchmarkNand3(b *testing.B) { runBenchmark(b, "Nand3") }
func BenchmarkNor3(b *testing.B)  { runBenchmark(b, "Nor3") }
func BenchmarkOdd3(b *testing.B)  { runBenchmark(b, "Odd3") }
func BenchmarkEven3(b *testing.B) { runBenchmark(b, "Even3") }
func BenchmarkAnd4(b *testing.B)  { runBenchmark(b, "And4") }
func BenchmarkOr4(b *testing.B)   { runBenchmark(b, "Or4") }
func BenchmarkNand4(b *testing.B) { runBenchmark(b, "Nand4") }
func BenchmarkNor4(b *testing.B)  { runBenchmark(b, "Nor4") }
func BenchmarkOdd4(b *testing.B)  { runBenchmark(b, "Odd4") }
func BenchmarkEven4(b *testing.B) { runBenchmark(b, "Even4") }
func BenchmarkId(b *testing.B)    { runBenchmark(b, "Id") }
func BenchmarkZero(b *testing.B)  { runBenchmark(b, "Zero") }
func BenchmarkOne(b *testing.B)   { runBenchmark(b, "One") }
func BenchmarkLT(b *testing.B)    { runBenchmark(b, "LT") }
func BenchmarkGT(b *testing.B)    { runBenchmark(b, "GT") }
func BenchmarkLOE(b *testing.B)   { runBenchmark(b, "LOE") }
func BenchmarkGOE(b *testing.B)   { runBenchmark(b, "GOE") }
func BenchmarkNotA(b *testing.B)  { runBenchmark(b, "NotA") }
func BenchmarkNotB(b *testing.B)  { runBenchmark(b, "NotB") }
func BenchmarkIdA(b *testing.B)   { runBenchmark(b, "IdA") }
func BenchmarkIdB(b *testing.B)   { runBenchmark(b, "IdB") }
func BenchmarkZero2(b *testing.B) { runBenchmark(b, "Zero2") }
func BenchmarkOne2(b *testing.B)  { runBenchmark(b, "One2") }
func BenchmarkLT3(b *testing.B)   { runBenchmark(b, "LT3") }
func BenchmarkGT3(b *testing.B)   { runBenchmark(b, "GT3") }
func BenchmarkLOE3(b *testing.B)  { runBenchmark(b, "LOE3") }
func BenchmarkGOE3(b *testing.B)  { runBenchmark(b, "GOE3") }
func BenchmarkMux(b *testing.B)   { runBenchmark(b, "Mux") }
func BenchmarkIf(b *testing.B)    { runBenchmark(b, "If") }
func BenchmarkMaj(b *testing.B)   { runBenchmark(b, "Maj") }
func BenchmarkMin(b *testing.B)   { runBenchmark(b, "Min") }
func Benchmark2Off(b *testing.B)  { runBenchmark(b, "2Off") }
func Benchmark2On(b *testing.B)   { runBenchmark(b, "2On") }
func BenchmarkLM3A1(b *testing.B) { runBenchmark(b, "LM3A1") }
func BenchmarkLM3A2(b *testing.B) { runBenchmark(b, "LM3A2") }
func BenchmarkLM3A3(b *testing.B) { runBenchmark(b, "LM3A3") }
func BenchmarkLM3A4(b *testing.B) { runBenchmark(b, "LM3A4") }
func BenchmarkLM3B1(b *testing.B) { runBenchmark(b, "LM3B1") }
func BenchmarkLM3B2(b *testing.B) { runBenchmark(b, "LM3B2") }
func BenchmarkLM3B3(b *testing.B) { runBenchmark(b, "LM3B3") }
func BenchmarkLM3B4(b *testing.B) { runBenchmark(b, "LM3B4") }
func BenchmarkLM3C1(b *testing.B) { runBenchmark(b, "LM3C1") }
func BenchmarkLM3C2(b *testing.B) { runBenchmark(b, "LM3C2") }
func BenchmarkLM3C3(b *testing.B) { runBenchmark(b, "LM3C3") }
func BenchmarkLM3C4(b *testing.B) { runBenchmark(b, "LM3C4") }
func BenchmarkLM3D1(b *testing.B) { runBenchmark(b, "LM3D1") }
func BenchmarkLM3D2(b *testing.B) { runBenchmark(b, "LM3D2") }
func BenchmarkLM3D3(b *testing.B) { runBenchmark(b, "LM3D3") }
func BenchmarkLM3D4(b *testing.B) { runBenchmark(b, "LM3D4") }
func BenchmarkLM3E1(b *testing.B) { runBenchmark(b, "LM3E1") }
func BenchmarkLM3E2(b *testing.B) { runBenchmark(b, "LM3E2") }
func BenchmarkLM3E3(b *testing.B) { runBenchmark(b, "LM3E3") }
func BenchmarkLM3F1(b *testing.B) { runBenchmark(b, "LM3F1") }
func BenchmarkLM3F2(b *testing.B) { runBenchmark(b, "LM3F2") }
func BenchmarkLM3F3(b *testing.B) { runBenchmark(b, "LM3F3") }
func BenchmarkLM3G1(b *testing.B) { runBenchmark(b, "LM3G1") }
func BenchmarkLM3G2(b *testing.B) { runBenchmark(b, "LM3G2") }
func BenchmarkLM3G3(b *testing.B) { runBenchmark(b, "LM3G3") }
func BenchmarkLM3G4(b *testing.B) { runBenchmark(b, "LM3G4") }
func BenchmarkLM3H1(b *testing.B) { runBenchmark(b, "LM3H1") }
func BenchmarkLM3H2(b *testing.B) { runBenchmark(b, "LM3H2") }
func BenchmarkLM3H3(b *testing.B) { runBenchmark(b, "LM3H3") }
func BenchmarkLM3H4(b *testing.B) { runBenchmark(b, "LM3H4") }
func BenchmarkLT3A(b *testing.B)  { runBenchmark(b, "LT3A") }
func BenchmarkGT3A(b *testing.B)  { runBenchmark(b, "GT3A") }
func BenchmarkLOE3A(b *testing.B) { runBenchmark(b, "LOE3A") }
func BenchmarkGOE3A(b *testing.B) { runBenchmark(b, "GOE3A") }
func BenchmarkET3A(b *testing.B)  { runBenchmark(b, "ET3A") }
func BenchmarkNET3A(b *testing.B) { runBenchmark(b, "NET3A") }
func BenchmarkLT3B(b *testing.B)  { runBenchmark(b, "LT3B") }
func BenchmarkGT3B(b *testing.B)  { runBenchmark(b, "GT3B") }
func BenchmarkLOE3B(b *testing.B) { runBenchmark(b, "LOE3B") }
func BenchmarkGOE3B(b *testing.B) { runBenchmark(b, "GOE3B") }
func BenchmarkET3B(b *testing.B)  { runBenchmark(b, "ET3B") }
func BenchmarkNET3B(b *testing.B) { runBenchmark(b, "NET3B") }
func BenchmarkLT3C(b *testing.B)  { runBenchmark(b, "LT3C") }
func BenchmarkGT3C(b *testing.B)  { runBenchmark(b, "GT3C") }
func BenchmarkLOE3C(b *testing.B) { runBenchmark(b, "LOE3C") }
func BenchmarkGOE3C(b *testing.B) { runBenchmark(b, "GOE3C") }
func BenchmarkET3C(b *testing.B)  { runBenchmark(b, "ET3C") }
func BenchmarkNET3C(b *testing.B) { runBenchmark(b, "NET3C") }
func BenchmarkT004(b *testing.B)  { runBenchmark(b, "T004") }
func BenchmarkT008(b *testing.B)  { runBenchmark(b, "T008") }
func BenchmarkT009(b *testing.B)  { runBenchmark(b, "T009") }
func BenchmarkT032(b *testing.B)  { runBenchmark(b, "T032") }
func BenchmarkT033(b *testing.B)  { runBenchmark(b, "T033") }
func BenchmarkT041(b *testing.B)  { runBenchmark(b, "T041") }
func BenchmarkT055(b *testing.B)  { runBenchmark(b, "T055") }
func BenchmarkT057(b *testing.B)  { runBenchmark(b, "T057") }
func BenchmarkT064(b *testing.B)  { runBenchmark(b, "T064") }
func BenchmarkT065(b *testing.B)  { runBenchmark(b, "T065") }
func BenchmarkT069(b *testing.B)  { runBenchmark(b, "T069") }
func BenchmarkT073(b *testing.B)  { runBenchmark(b, "T073") }
func BenchmarkT081(b *testing.B)  { runBenchmark(b, "T081") }
func BenchmarkT089(b *testing.B)  { runBenchmark(b, "T089") }
func BenchmarkT093(b *testing.B)  { runBenchmark(b, "T093") }
func BenchmarkT096(b *testing.B)  { runBenchmark(b, "T096") }
func BenchmarkT101(b *testing.B)  { runBenchmark(b, "T101") }
func BenchmarkT109(b *testing.B)  { runBenchmark(b, "T109") }
func BenchmarkT111(b *testing.B)  { runBenchmark(b, "T111") }
func BenchmarkT121(b *testing.B)  { runBenchmark(b, "T121") }
func BenchmarkT123(b *testing.B)  { runBenchmark(b, "T123") }
func BenchmarkT125(b *testing.B)  { runBenchmark(b, "T125") }
func BenchmarkT154(b *testing.B)  { runBenchmark(b, "T154") }
func BenchmarkT223(b *testing.B)  { runBenchmark(b, "T223") }
func BenchmarkT239(b *testing.B)  { runBenchmark(b, "T239") }
func BenchmarkT249(b *testing.B)  { runBenchmark(b, "T249") }
func BenchmarkT251(b *testing.B)  { runBenchmark(b, "T251") }
func BenchmarkT253(b *testing.B)  { runBenchmark(b, "T253") }
func BenchmarkLT4(b *testing.B)   { runBenchmark(b, "LT4") }
func BenchmarkGT4(b *testing.B)   { runBenchmark(b, "GT4") }
func BenchmarkLOE4(b *testing.B)  { runBenchmark(b, "LOE4") }
func BenchmarkGOE4(b *testing.B)  { runBenchmark(b, "GOE4") }
func BenchmarkTie(b *testing.B)   { runBenchmark(b, "Tie") }
func BenchmarkNtie(b *testing.B)  { runBenchmark(b, "Ntie") }
func Benchmark3Off(b *testing.B)  { runBenchmark(b, "3Off") }
func Benchmark3On(b *testing.B)   { runBenchmark(b, "3On") }
func BenchmarkLM4A1(b *testing.B) { runBenchmark(b, "LM4A1") }
func BenchmarkLM4A2(b *testing.B) { runBenchmark(b, "LM4A2") }
func BenchmarkLM4A3(b *testing.B) { runBenchmark(b, "LM4A3") }
func BenchmarkLM4A4(b *testing.B) { runBenchmark(b, "LM4A4") }
func BenchmarkLM4A5(b *testing.B) { runBenchmark(b, "LM4A5") }
func BenchmarkLM4A6(b *testing.B) { runBenchmark(b, "LM4A6") }
func BenchmarkLM4A7(b *testing.B) { runBenchmark(b, "LM4A7") }
func BenchmarkLM4A8(b *testing.B) { runBenchmark(b, "LM4A8") }
func BenchmarkLM4B1(b *testing.B) { runBenchmark(b, "LM4B1") }
func BenchmarkLM4B2(b *testing.B) { runBenchmark(b, "LM4B2") }
func BenchmarkLM4B3(b *testing.B) { runBenchmark(b, "LM4B3") }
func BenchmarkLM4B4(b *testing.B) { runBenchmark(b, "LM4B4") }
func BenchmarkLM4B5(b *testing.B) { runBenchmark(b, "LM4B5") }
func BenchmarkLM4B6(b *testing.B) { runBenchmark(b, "LM4B6") }
func BenchmarkLM4B7(b *testing.B) { runBenchmark(b, "LM4B7") }
func BenchmarkLM4B8(b *testing.B) { runBenchmark(b, "LM4B8") }
func BenchmarkLM4C1(b *testing.B) { runBenchmark(b, "LM4C1") }
func BenchmarkLM4C2(b *testing.B) { runBenchmark(b, "LM4C2") }
func BenchmarkLM4C3(b *testing.B) { runBenchmark(b, "LM4C3") }
func BenchmarkLM4C4(b *testing.B) { runBenchmark(b, "LM4C4") }
func BenchmarkLM4C5(b *testing.B) { runBenchmark(b, "LM4C5") }
func BenchmarkLM4C6(b *testing.B) { runBenchmark(b, "LM4C6") }
func BenchmarkLM4C7(b *testing.B) { runBenchmark(b, "LM4C7") }
func BenchmarkLM4C8(b *testing.B) { runBenchmark(b, "LM4C8") }
func BenchmarkLM4D1(b *testing.B) { runBenchmark(b, "LM4D1") }
func BenchmarkLM4D2(b *testing.B) { runBenchmark(b, "LM4D2") }
func BenchmarkLM4D3(b *testing.B) { runBenchmark(b, "LM4D3") }
func BenchmarkLM4D4(b *testing.B) { runBenchmark(b, "LM4D4") }
func BenchmarkLM4D5(b *testing.B) { runBenchmark(b, "LM4D5") }
func BenchmarkLM4D6(b *testing.B) { runBenchmark(b, "LM4D6") }
func BenchmarkLM4D7(b *testing.B) { runBenchmark(b, "LM4D7") }
func BenchmarkLM4D8(b *testing.B) { runBenchmark(b, "LM4D8") }
func BenchmarkLM4E1(b *testing.B) { runBenchmark(b, "LM4E1") }
func BenchmarkLM4E2(b *testing.B) { runBenchmark(b, "LM4E2") }
func BenchmarkLM4E3(b *testing.B) { runBenchmark(b, "LM4E3") }
func BenchmarkLM4E4(b *testing.B) { runBenchmark(b, "LM4E4") }
func BenchmarkLM4E5(b *testing.B) { runBenchmark(b, "LM4E5") }
func BenchmarkLM4E6(b *testing.B) { runBenchmark(b, "LM4E6") }
func BenchmarkLM4E7(b *testing.B) { runBenchmark(b, "LM4E7") }
func BenchmarkLM4E8(b *testing.B) { runBenchmark(b, "LM4E8") }
func BenchmarkLM4F1(b *testing.B) { runBenchmark(b, "LM4F1") }
func BenchmarkLM4F2(b *testing.B) { runBenchmark(b, "LM4F2") }
func BenchmarkLM4F3(b *testing.B) { runBenchmark(b, "LM4F3") }
func BenchmarkLM4F4(b *testing.B) { runBenchmark(b, "LM4F4") }
func BenchmarkLM4F5(b *testing.B) { runBenchmark(b, "LM4F5") }
func BenchmarkLM4F6(b *testing.B) { runBenchmark(b, "LM4F6") }
func BenchmarkLM4F7(b *testing.B) { runBenchmark(b, "LM4F7") }
func BenchmarkLM4F8(b *testing.B) { runBenchmark(b, "LM4F8") }
func BenchmarkLM4G1(b *testing.B) { runBenchmark(b, "LM4G1") }
func BenchmarkLM4G2(b *testing.B) { runBenchmark(b, "LM4G2") }
func BenchmarkLM4G3(b *testing.B) { runBenchmark(b, "LM4G3") }
func BenchmarkLM4G4(b *testing.B) { runBenchmark(b, "LM4G4") }
func BenchmarkLM4G5(b *testing.B) { runBenchmark(b, "LM4G5") }
func BenchmarkLM4G6(b *testing.B) { runBenchmark(b, "LM4G6") }
func BenchmarkLM4G7(b *testing.B) { runBenchmark(b, "LM4G7") }
func BenchmarkLM4G8(b *testing.B) { runBenchmark(b, "LM4G8") }
func BenchmarkLM4H1(b *testing.B) { runBenchmark(b, "LM4H1") }
func BenchmarkLM4H2(b *testing.B) { runBenchmark(b, "LM4H2") }
func BenchmarkLM4H3(b *testing.B) { runBenchmark(b, "LM4H3") }
func BenchmarkLM4H4(b *testing.B) { runBenchmark(b, "LM4H4") }
func BenchmarkLM4H5(b *testing.B) { runBenchmark(b, "LM4H5") }
func BenchmarkLM4H6(b *testing.B) { runBenchmark(b, "LM4H6") }
func BenchmarkLM4H7(b *testing.B) { runBenchmark(b, "LM4H7") }
func BenchmarkLM4H8(b *testing.B) { runBenchmark(b, "LM4H8") }
func BenchmarkLM4I1(b *testing.B) { runBenchmark(b, "LM4I1") }
func BenchmarkLM4I2(b *testing.B) { runBenchmark(b, "LM4I2") }
func BenchmarkLM4I3(b *testing.B) { runBenchmark(b, "LM4I3") }
func BenchmarkLM4I4(b *testing.B) { runBenchmark(b, "LM4I4") }
func BenchmarkLM4I5(b *testing.B) { runBenchmark(b, "LM4I5") }
func BenchmarkLM4I6(b *testing.B) { runBenchmark(b, "LM4I6") }
func BenchmarkLM4I7(b *testing.B) { runBenchmark(b, "LM4I7") }
func BenchmarkLM4I8(b *testing.B) { runBenchmark(b, "LM4I8") }
func BenchmarkLT4A(b *testing.B)  { runBenchmark(b, "LT4A") }
func BenchmarkGT4A(b *testing.B)  { runBenchmark(b, "GT4A") }
func BenchmarkLOE4A(b *testing.B) { runBenchmark(b, "LOE4A") }
func BenchmarkGOE4A(b *testing.B) { runBenchmark(b, "GOE4A") }
func BenchmarkET4A(b *testing.B)  { runBenchmark(b, "ET4A") }
func BenchmarkNET4A(b *testing.B) { runBenchmark(b, "NET4A") }
func BenchmarkLT4B(b *testing.B)  { runBenchmark(b, "LT4B") }
func BenchmarkGT4B(b *testing.B)  { runBenchmark(b, "GT4B") }
func BenchmarkLOE4B(b *testing.B) { runBenchmark(b, "LOE4B") }
func BenchmarkGOE4B(b *testing.B) { runBenchmark(b, "GOE4B") }
func BenchmarkET4B(b *testing.B)  { runBenchmark(b, "ET4B") }
func BenchmarkNET4B(b *testing.B) { runBenchmark(b, "NET4B") }
func BenchmarkLT4C(b *testing.B)  { runBenchmark(b, "LT4C") }
func BenchmarkGT4C(b *testing.B)  { runBenchmark(b, "GT4C") }
func BenchmarkLOE4C(b *testing.B) { runBenchmark(b, "LOE4C") }
func BenchmarkGOE4C(b *testing.B) { runBenchmark(b, "GOE4C") }
func BenchmarkET4C(b *testing.B)  { runBenchmark(b, "ET4C") }
func BenchmarkNET4C(b *testing.B) { runBenchmark(b, "NET4C") }
func BenchmarkLT4D(b *testing.B)  { runBenchmark(b, "LT4D") }
func BenchmarkGT4D(b *testing.B)  { runBenchmark(b, "GT4D") }
func BenchmarkLOE4D(b *testing.B) { runBenchmark(b, "LOE4D") }
func BenchmarkGOE4D(b *testing.B) { runBenchmark(b, "GOE4D") }
func BenchmarkET4D(b *testing.B)  { runBenchmark(b, "ET4D") }
func BenchmarkNET4D(b *testing.B) { runBenchmark(b, "NET4D") }
func BenchmarkLT4E(b *testing.B)  { runBenchmark(b, "LT4E") }
func BenchmarkGT4E(b *testing.B)  { runBenchmark(b, "GT4E") }
func BenchmarkLOE4E(b *testing.B) { runBenchmark(b, "LOE4E") }
func BenchmarkGOE4E(b *testing.B) { runBenchmark(b, "GOE4E") }
func BenchmarkET4E(b *testing.B)  { runBenchmark(b, "ET4E") }
func BenchmarkNET4E(b *testing.B) { runBenchmark(b, "NET4E") }
func BenchmarkQ0002(b *testing.B) { runBenchmark(b, "Q0002") }
func BenchmarkQ001C(b *testing.B) { runBenchmark(b, "Q001C") }
func BenchmarkQ0048(b *testing.B) { runBenchmark(b, "Q0048") }
func BenchmarkQ0800(b *testing.B) { runBenchmark(b, "Q0800") }
func BenchmarkQ3378(b *testing.B) { runBenchmark(b, "Q3378") }
func BenchmarkQ3475(b *testing.B) { runBenchmark(b, "Q3475") }
func BenchmarkQ3CB0(b *testing.B) { runBenchmark(b, "Q3CB0") }
func BenchmarkQ3DEF(b *testing.B) { runBenchmark(b, "Q3DEF") }
func BenchmarkQ3DFF(b *testing.B) { runBenchmark(b, "Q3DFF") }
func BenchmarkQ4200(b *testing.B) { runBenchmark(b, "Q4200") }
func BenchmarkQ4C11(b *testing.B) { runBenchmark(b, "Q4C11") }
func BenchmarkQ5100(b *testing.B) { runBenchmark(b, "Q5100") }
func BenchmarkQ5EEF(b *testing.B) { runBenchmark(b, "Q5EEF") }
func BenchmarkQ5EFF(b *testing.B) { runBenchmark(b, "Q5EFF") }
func BenchmarkQ6A6D(b *testing.B) { runBenchmark(b, "Q6A6D") }
func BenchmarkQ6F75(b *testing.B) { runBenchmark(b, "Q6F75") }
func BenchmarkQ74C4(b *testing.B) { runBenchmark(b, "Q74C4") }
func BenchmarkQ7DA3(b *testing.B) { runBenchmark(b, "Q7DA3") }
func BenchmarkQ8304(b *testing.B) { runBenchmark(b, "Q8304") }
func BenchmarkQ8430(b *testing.B) { runBenchmark(b, "Q8430") }
func BenchmarkQ8543(b *testing.B) { runBenchmark(b, "Q8543") }
func BenchmarkQ9D80(b *testing.B) { runBenchmark(b, "Q9D80") }
func BenchmarkQA092(b *testing.B) { runBenchmark(b, "QA092") }
func BenchmarkQB36A(b *testing.B) { runBenchmark(b, "QB36A") }
func BenchmarkQCBCF(b *testing.B) { runBenchmark(b, "QCBCF") }
func BenchmarkQEEB1(b *testing.B) { runBenchmark(b, "QEEB1") }
func BenchmarkQEFFF(b *testing.B) { runBenchmark(b, "QEFFF") }
func BenchmarkQFF7B(b *testing.B) { runBenchmark(b, "QFF7B") }
func BenchmarkQFFF6(b *testing.B) { runBenchmark(b, "QFFF6") }
func BenchmarkQFFFB(b *testing.B) { runBenchmark(b, "QFFFB") }
