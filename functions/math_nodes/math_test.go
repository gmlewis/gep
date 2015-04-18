package mathNodes

import (
	"testing"
)

func TestPlus(t *testing.T) {
	if g, w := Math["+"].Float64Function(2, 3, 4, 5), 5.0; g != w {
		t.Errorf("Math[+](2,3,4,5) = %v, want %v", g, w)
	}
}

var result float64

func BenchmarkPlus(b *testing.B) {
	var v float64
	f := Math["+"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkMinus(b *testing.B) {
	var v float64
	f := Math["-"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkMultiply(b *testing.B) {
	var v float64
	f := Math["*"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkDivide(b *testing.B) {
	var v float64
	f := Math["/"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkMod(b *testing.B) {
	var v float64
	f := Math["Mod"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkPow(b *testing.B) {
	var v float64
	f := Math["Pow"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkSqrt(b *testing.B) {
	var v float64
	f := Math["Sqrt"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkExp(b *testing.B) {
	var v float64
	f := Math["Exp"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkPow10(b *testing.B) {
	var v float64
	f := Math["Pow10"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLn(b *testing.B) {
	var v float64
	f := Math["Ln"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLog(b *testing.B) {
	var v float64
	f := Math["Log"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLog2(b *testing.B) {
	var v float64
	f := Math["Log2"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkFloor(b *testing.B) {
	var v float64
	f := Math["Floor"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkCeil(b *testing.B) {
	var v float64
	f := Math["Ceil"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAbs(b *testing.B) {
	var v float64
	f := Math["Abs"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkInv(b *testing.B) {
	var v float64
	f := Math["Inv"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNeg(b *testing.B) {
	var v float64
	f := Math["Neg"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNop(b *testing.B) {
	var v float64
	f := Math["Nop"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkX2(b *testing.B) {
	var v float64
	f := Math["X2"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkX3(b *testing.B) {
	var v float64
	f := Math["X3"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkX4(b *testing.B) {
	var v float64
	f := Math["X4"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkX5(b *testing.B) {
	var v float64
	f := Math["X5"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func Benchmark3Rt(b *testing.B) {
	var v float64
	f := Math["3Rt"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func Benchmark4Rt(b *testing.B) {
	var v float64
	f := Math["4Rt"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func Benchmark5Rt(b *testing.B) {
	var v float64
	f := Math["5Rt"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAdd3(b *testing.B) {
	var v float64
	f := Math["Add3"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkSub3(b *testing.B) {
	var v float64
	f := Math["Sub3"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkMul3(b *testing.B) {
	var v float64
	f := Math["Mul3"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkDiv3(b *testing.B) {
	var v float64
	f := Math["Div3"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAdd4(b *testing.B) {
	var v float64
	f := Math["Add4"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkSub4(b *testing.B) {
	var v float64
	f := Math["Sub4"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkMul4(b *testing.B) {
	var v float64
	f := Math["Mul4"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkDiv4(b *testing.B) {
	var v float64
	f := Math["Div4"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkMin2(b *testing.B) {
	var v float64
	f := Math["Min2"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkMin3(b *testing.B) {
	var v float64
	f := Math["Min3"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkMin4(b *testing.B) {
	var v float64
	f := Math["Min4"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkMax2(b *testing.B) {
	var v float64
	f := Math["Max2"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkMax3(b *testing.B) {
	var v float64
	f := Math["Max3"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkMax4(b *testing.B) {
	var v float64
	f := Math["Max4"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAvg2(b *testing.B) {
	var v float64
	f := Math["Avg2"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAvg3(b *testing.B) {
	var v float64
	f := Math["Avg3"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAvg4(b *testing.B) {
	var v float64
	f := Math["Avg4"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLogi(b *testing.B) {
	var v float64
	f := Math["Logi"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLogi2(b *testing.B) {
	var v float64
	f := Math["Logi2"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLogi3(b *testing.B) {
	var v float64
	f := Math["Logi3"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLogi4(b *testing.B) {
	var v float64
	f := Math["Logi4"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGau(b *testing.B) {
	var v float64
	f := Math["Gau"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGau2(b *testing.B) {
	var v float64
	f := Math["Gau2"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGau3(b *testing.B) {
	var v float64
	f := Math["Gau3"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGau4(b *testing.B) {
	var v float64
	f := Math["Gau4"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkZero(b *testing.B) {
	var v float64
	f := Math["Zero"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkOne(b *testing.B) {
	var v float64
	f := Math["One"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkZero2(b *testing.B) {
	var v float64
	f := Math["Zero2"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkOne2(b *testing.B) {
	var v float64
	f := Math["One2"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkPi(b *testing.B) {
	var v float64
	f := Math["Pi"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkE(b *testing.B) {
	var v float64
	f := Math["E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkSin(b *testing.B) {
	var v float64
	f := Math["Sin"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkCos(b *testing.B) {
	var v float64
	f := Math["Cos"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkTan(b *testing.B) {
	var v float64
	f := Math["Tan"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkCsc(b *testing.B) {
	var v float64
	f := Math["Csc"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkSec(b *testing.B) {
	var v float64
	f := Math["Sec"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkCot(b *testing.B) {
	var v float64
	f := Math["Cot"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAsin(b *testing.B) {
	var v float64
	f := Math["Asin"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAcos(b *testing.B) {
	var v float64
	f := Math["Acos"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAtan(b *testing.B) {
	var v float64
	f := Math["Atan"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAcsc(b *testing.B) {
	var v float64
	f := Math["Acsc"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAsec(b *testing.B) {
	var v float64
	f := Math["Asec"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAcot(b *testing.B) {
	var v float64
	f := Math["Acot"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkSinh(b *testing.B) {
	var v float64
	f := Math["Sinh"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkCosh(b *testing.B) {
	var v float64
	f := Math["Cosh"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkTanh(b *testing.B) {
	var v float64
	f := Math["Tanh"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkCsch(b *testing.B) {
	var v float64
	f := Math["Csch"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkSech(b *testing.B) {
	var v float64
	f := Math["Sech"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkCoth(b *testing.B) {
	var v float64
	f := Math["Coth"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAsinh(b *testing.B) {
	var v float64
	f := Math["Asinh"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAcosh(b *testing.B) {
	var v float64
	f := Math["Acosh"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAtanh(b *testing.B) {
	var v float64
	f := Math["Atanh"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAcsch(b *testing.B) {
	var v float64
	f := Math["Acsch"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAsech(b *testing.B) {
	var v float64
	f := Math["Asech"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAcoth(b *testing.B) {
	var v float64
	f := Math["Acoth"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNOT(b *testing.B) {
	var v float64
	f := Math["NOT"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkOR1(b *testing.B) {
	var v float64
	f := Math["OR1"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkOR2(b *testing.B) {
	var v float64
	f := Math["OR2"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkOR3(b *testing.B) {
	var v float64
	f := Math["OR3"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkOR4(b *testing.B) {
	var v float64
	f := Math["OR4"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkOR5(b *testing.B) {
	var v float64
	f := Math["OR5"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkOR6(b *testing.B) {
	var v float64
	f := Math["OR6"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAND1(b *testing.B) {
	var v float64
	f := Math["AND1"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAND2(b *testing.B) {
	var v float64
	f := Math["AND2"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAND3(b *testing.B) {
	var v float64
	f := Math["AND3"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAND4(b *testing.B) {
	var v float64
	f := Math["AND4"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAND5(b *testing.B) {
	var v float64
	f := Math["AND5"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkAND6(b *testing.B) {
	var v float64
	f := Math["AND6"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT2A(b *testing.B) {
	var v float64
	f := Math["LT2A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT2A(b *testing.B) {
	var v float64
	f := Math["GT2A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE2A(b *testing.B) {
	var v float64
	f := Math["LOE2A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE2A(b *testing.B) {
	var v float64
	f := Math["GOE2A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET2A(b *testing.B) {
	var v float64
	f := Math["ET2A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET2A(b *testing.B) {
	var v float64
	f := Math["NET2A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT2B(b *testing.B) {
	var v float64
	f := Math["LT2B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT2B(b *testing.B) {
	var v float64
	f := Math["GT2B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE2B(b *testing.B) {
	var v float64
	f := Math["LOE2B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE2B(b *testing.B) {
	var v float64
	f := Math["GOE2B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET2B(b *testing.B) {
	var v float64
	f := Math["ET2B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET2B(b *testing.B) {
	var v float64
	f := Math["NET2B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT2C(b *testing.B) {
	var v float64
	f := Math["LT2C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT2C(b *testing.B) {
	var v float64
	f := Math["GT2C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE2C(b *testing.B) {
	var v float64
	f := Math["LOE2C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE2C(b *testing.B) {
	var v float64
	f := Math["GOE2C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET2C(b *testing.B) {
	var v float64
	f := Math["ET2C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET2C(b *testing.B) {
	var v float64
	f := Math["NET2C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT2D(b *testing.B) {
	var v float64
	f := Math["LT2D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT2D(b *testing.B) {
	var v float64
	f := Math["GT2D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE2D(b *testing.B) {
	var v float64
	f := Math["LOE2D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE2D(b *testing.B) {
	var v float64
	f := Math["GOE2D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET2D(b *testing.B) {
	var v float64
	f := Math["ET2D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET2D(b *testing.B) {
	var v float64
	f := Math["NET2D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT2E(b *testing.B) {
	var v float64
	f := Math["LT2E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT2E(b *testing.B) {
	var v float64
	f := Math["GT2E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE2E(b *testing.B) {
	var v float64
	f := Math["LOE2E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE2E(b *testing.B) {
	var v float64
	f := Math["GOE2E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET2E(b *testing.B) {
	var v float64
	f := Math["ET2E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET2E(b *testing.B) {
	var v float64
	f := Math["NET2E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT2F(b *testing.B) {
	var v float64
	f := Math["LT2F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT2F(b *testing.B) {
	var v float64
	f := Math["GT2F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE2F(b *testing.B) {
	var v float64
	f := Math["LOE2F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE2F(b *testing.B) {
	var v float64
	f := Math["GOE2F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET2F(b *testing.B) {
	var v float64
	f := Math["ET2F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET2F(b *testing.B) {
	var v float64
	f := Math["NET2F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT2G(b *testing.B) {
	var v float64
	f := Math["LT2G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT2G(b *testing.B) {
	var v float64
	f := Math["GT2G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE2G(b *testing.B) {
	var v float64
	f := Math["LOE2G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE2G(b *testing.B) {
	var v float64
	f := Math["GOE2G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET2G(b *testing.B) {
	var v float64
	f := Math["ET2G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET2G(b *testing.B) {
	var v float64
	f := Math["NET2G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT3A(b *testing.B) {
	var v float64
	f := Math["LT3A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT3A(b *testing.B) {
	var v float64
	f := Math["GT3A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE3A(b *testing.B) {
	var v float64
	f := Math["LOE3A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE3A(b *testing.B) {
	var v float64
	f := Math["GOE3A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET3A(b *testing.B) {
	var v float64
	f := Math["ET3A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET3A(b *testing.B) {
	var v float64
	f := Math["NET3A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT3B(b *testing.B) {
	var v float64
	f := Math["LT3B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT3B(b *testing.B) {
	var v float64
	f := Math["GT3B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE3B(b *testing.B) {
	var v float64
	f := Math["LOE3B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE3B(b *testing.B) {
	var v float64
	f := Math["GOE3B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET3B(b *testing.B) {
	var v float64
	f := Math["ET3B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET3B(b *testing.B) {
	var v float64
	f := Math["NET3B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT3C(b *testing.B) {
	var v float64
	f := Math["LT3C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT3C(b *testing.B) {
	var v float64
	f := Math["GT3C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE3C(b *testing.B) {
	var v float64
	f := Math["LOE3C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE3C(b *testing.B) {
	var v float64
	f := Math["GOE3C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET3C(b *testing.B) {
	var v float64
	f := Math["ET3C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET3C(b *testing.B) {
	var v float64
	f := Math["NET3C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT3D(b *testing.B) {
	var v float64
	f := Math["LT3D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT3D(b *testing.B) {
	var v float64
	f := Math["GT3D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE3D(b *testing.B) {
	var v float64
	f := Math["LOE3D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE3D(b *testing.B) {
	var v float64
	f := Math["GOE3D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET3D(b *testing.B) {
	var v float64
	f := Math["ET3D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET3D(b *testing.B) {
	var v float64
	f := Math["NET3D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT3E(b *testing.B) {
	var v float64
	f := Math["LT3E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT3E(b *testing.B) {
	var v float64
	f := Math["GT3E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE3E(b *testing.B) {
	var v float64
	f := Math["LOE3E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE3E(b *testing.B) {
	var v float64
	f := Math["GOE3E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET3E(b *testing.B) {
	var v float64
	f := Math["ET3E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET3E(b *testing.B) {
	var v float64
	f := Math["NET3E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT3F(b *testing.B) {
	var v float64
	f := Math["LT3F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT3F(b *testing.B) {
	var v float64
	f := Math["GT3F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE3F(b *testing.B) {
	var v float64
	f := Math["LOE3F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE3F(b *testing.B) {
	var v float64
	f := Math["GOE3F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET3F(b *testing.B) {
	var v float64
	f := Math["ET3F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET3F(b *testing.B) {
	var v float64
	f := Math["NET3F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT3G(b *testing.B) {
	var v float64
	f := Math["LT3G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT3G(b *testing.B) {
	var v float64
	f := Math["GT3G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE3G(b *testing.B) {
	var v float64
	f := Math["LOE3G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE3G(b *testing.B) {
	var v float64
	f := Math["GOE3G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET3G(b *testing.B) {
	var v float64
	f := Math["ET3G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET3G(b *testing.B) {
	var v float64
	f := Math["NET3G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT3H(b *testing.B) {
	var v float64
	f := Math["LT3H"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT3H(b *testing.B) {
	var v float64
	f := Math["GT3H"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE3H(b *testing.B) {
	var v float64
	f := Math["LOE3H"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE3H(b *testing.B) {
	var v float64
	f := Math["GOE3H"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET3H(b *testing.B) {
	var v float64
	f := Math["ET3H"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET3H(b *testing.B) {
	var v float64
	f := Math["NET3H"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT3I(b *testing.B) {
	var v float64
	f := Math["LT3I"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT3I(b *testing.B) {
	var v float64
	f := Math["GT3I"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE3I(b *testing.B) {
	var v float64
	f := Math["LOE3I"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE3I(b *testing.B) {
	var v float64
	f := Math["GOE3I"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET3I(b *testing.B) {
	var v float64
	f := Math["ET3I"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET3I(b *testing.B) {
	var v float64
	f := Math["NET3I"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT3J(b *testing.B) {
	var v float64
	f := Math["LT3J"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT3J(b *testing.B) {
	var v float64
	f := Math["GT3J"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE3J(b *testing.B) {
	var v float64
	f := Math["LOE3J"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE3J(b *testing.B) {
	var v float64
	f := Math["GOE3J"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET3J(b *testing.B) {
	var v float64
	f := Math["ET3J"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET3J(b *testing.B) {
	var v float64
	f := Math["NET3J"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT3K(b *testing.B) {
	var v float64
	f := Math["LT3K"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT3K(b *testing.B) {
	var v float64
	f := Math["GT3K"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE3K(b *testing.B) {
	var v float64
	f := Math["LOE3K"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE3K(b *testing.B) {
	var v float64
	f := Math["GOE3K"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET3K(b *testing.B) {
	var v float64
	f := Math["ET3K"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET3K(b *testing.B) {
	var v float64
	f := Math["NET3K"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT3L(b *testing.B) {
	var v float64
	f := Math["LT3L"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT3L(b *testing.B) {
	var v float64
	f := Math["GT3L"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE3L(b *testing.B) {
	var v float64
	f := Math["LOE3L"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE3L(b *testing.B) {
	var v float64
	f := Math["GOE3L"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET3L(b *testing.B) {
	var v float64
	f := Math["ET3L"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET3L(b *testing.B) {
	var v float64
	f := Math["NET3L"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT4A(b *testing.B) {
	var v float64
	f := Math["LT4A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT4A(b *testing.B) {
	var v float64
	f := Math["GT4A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE4A(b *testing.B) {
	var v float64
	f := Math["LOE4A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE4A(b *testing.B) {
	var v float64
	f := Math["GOE4A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET4A(b *testing.B) {
	var v float64
	f := Math["ET4A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET4A(b *testing.B) {
	var v float64
	f := Math["NET4A"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT4B(b *testing.B) {
	var v float64
	f := Math["LT4B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT4B(b *testing.B) {
	var v float64
	f := Math["GT4B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE4B(b *testing.B) {
	var v float64
	f := Math["LOE4B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE4B(b *testing.B) {
	var v float64
	f := Math["GOE4B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET4B(b *testing.B) {
	var v float64
	f := Math["ET4B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET4B(b *testing.B) {
	var v float64
	f := Math["NET4B"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT4C(b *testing.B) {
	var v float64
	f := Math["LT4C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT4C(b *testing.B) {
	var v float64
	f := Math["GT4C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE4C(b *testing.B) {
	var v float64
	f := Math["LOE4C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE4C(b *testing.B) {
	var v float64
	f := Math["GOE4C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET4C(b *testing.B) {
	var v float64
	f := Math["ET4C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET4C(b *testing.B) {
	var v float64
	f := Math["NET4C"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT4D(b *testing.B) {
	var v float64
	f := Math["LT4D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT4D(b *testing.B) {
	var v float64
	f := Math["GT4D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE4D(b *testing.B) {
	var v float64
	f := Math["LOE4D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE4D(b *testing.B) {
	var v float64
	f := Math["GOE4D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET4D(b *testing.B) {
	var v float64
	f := Math["ET4D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET4D(b *testing.B) {
	var v float64
	f := Math["NET4D"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT4E(b *testing.B) {
	var v float64
	f := Math["LT4E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT4E(b *testing.B) {
	var v float64
	f := Math["GT4E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE4E(b *testing.B) {
	var v float64
	f := Math["LOE4E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE4E(b *testing.B) {
	var v float64
	f := Math["GOE4E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET4E(b *testing.B) {
	var v float64
	f := Math["ET4E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET4E(b *testing.B) {
	var v float64
	f := Math["NET4E"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT4F(b *testing.B) {
	var v float64
	f := Math["LT4F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT4F(b *testing.B) {
	var v float64
	f := Math["GT4F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE4F(b *testing.B) {
	var v float64
	f := Math["LOE4F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE4F(b *testing.B) {
	var v float64
	f := Math["GOE4F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET4F(b *testing.B) {
	var v float64
	f := Math["ET4F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET4F(b *testing.B) {
	var v float64
	f := Math["NET4F"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT4G(b *testing.B) {
	var v float64
	f := Math["LT4G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT4G(b *testing.B) {
	var v float64
	f := Math["GT4G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE4G(b *testing.B) {
	var v float64
	f := Math["LOE4G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE4G(b *testing.B) {
	var v float64
	f := Math["GOE4G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET4G(b *testing.B) {
	var v float64
	f := Math["ET4G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET4G(b *testing.B) {
	var v float64
	f := Math["NET4G"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT4H(b *testing.B) {
	var v float64
	f := Math["LT4H"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT4H(b *testing.B) {
	var v float64
	f := Math["GT4H"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE4H(b *testing.B) {
	var v float64
	f := Math["LOE4H"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE4H(b *testing.B) {
	var v float64
	f := Math["GOE4H"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET4H(b *testing.B) {
	var v float64
	f := Math["ET4H"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET4H(b *testing.B) {
	var v float64
	f := Math["NET4H"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT4I(b *testing.B) {
	var v float64
	f := Math["LT4I"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT4I(b *testing.B) {
	var v float64
	f := Math["GT4I"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE4I(b *testing.B) {
	var v float64
	f := Math["LOE4I"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE4I(b *testing.B) {
	var v float64
	f := Math["GOE4I"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET4I(b *testing.B) {
	var v float64
	f := Math["ET4I"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET4I(b *testing.B) {
	var v float64
	f := Math["NET4I"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT4J(b *testing.B) {
	var v float64
	f := Math["LT4J"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT4J(b *testing.B) {
	var v float64
	f := Math["GT4J"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE4J(b *testing.B) {
	var v float64
	f := Math["LOE4J"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE4J(b *testing.B) {
	var v float64
	f := Math["GOE4J"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET4J(b *testing.B) {
	var v float64
	f := Math["ET4J"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET4J(b *testing.B) {
	var v float64
	f := Math["NET4J"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT4K(b *testing.B) {
	var v float64
	f := Math["LT4K"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT4K(b *testing.B) {
	var v float64
	f := Math["GT4K"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE4K(b *testing.B) {
	var v float64
	f := Math["LOE4K"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE4K(b *testing.B) {
	var v float64
	f := Math["GOE4K"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET4K(b *testing.B) {
	var v float64
	f := Math["ET4K"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET4K(b *testing.B) {
	var v float64
	f := Math["NET4K"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLT4L(b *testing.B) {
	var v float64
	f := Math["LT4L"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGT4L(b *testing.B) {
	var v float64
	f := Math["GT4L"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkLOE4L(b *testing.B) {
	var v float64
	f := Math["LOE4L"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkGOE4L(b *testing.B) {
	var v float64
	f := Math["GOE4L"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkET4L(b *testing.B) {
	var v float64
	f := Math["ET4L"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}

func BenchmarkNET4L(b *testing.B) {
	var v float64
	f := Math["NET4L"]
	for i := 0; i < b.N; i++ {
		v = f.Float64Function(2, 3, 4, 5)
	}
	result = v
}
