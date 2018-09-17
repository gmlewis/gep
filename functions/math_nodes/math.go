// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package mathNodes defines the floating-point function collection available for the GEP algorithm.
package mathNodes

import (
	"log"
	"math"

	"github.com/gmlewis/gep/v2/functions"
)

// MathNode is a floating-point function used for the formation of GEP expressions.
type MathNode struct {
	index     int
	symbol    string
	terminals int
	function  func(x []float64) float64
}

// Symbol returns the Karva symbol for this floating-point function.
func (n MathNode) Symbol() string {
	return n.symbol
}

// Terminals returns the number of input terminals for this floating-point function.
func (n MathNode) Terminals() int {
	return n.terminals
}

// BoolFunction is unused in this package and returns an error.
func (n MathNode) BoolFunction([]bool) bool {
	log.Println("error calling BoolFunction on MathNode model.")
	return false
}

// IntFunction is unused in this package and returns an error.
func (n MathNode) IntFunction([]int) int {
	log.Println("error calling IntFunction on MathNode model.")
	return 0
}

// Float64Function calls the floating-point function and returns the result.
func (n MathNode) Float64Function(x []float64) float64 {
	return n.function(x)
}

// Math lists all the available floating-point functions for this package.
var Math = functions.FuncMap{
	// TODO(gmlewis): Change functions to operate on the entire length of the slice.
	"+":     MathNode{0, "+", 2, func(x []float64) float64 { return (x[0] + x[1]) }},
	"-":     MathNode{1, "-", 2, func(x []float64) float64 { return (x[0] - x[1]) }},
	"*":     MathNode{2, "*", 2, func(x []float64) float64 { return (x[0] * x[1]) }},
	"/":     MathNode{3, "/", 2, func(x []float64) float64 { return (x[0] / x[1]) }},
	"Mod":   MathNode{4, "Mod", 2, func(x []float64) float64 { return gepMod(x[0], x[1]) }},
	"Pow":   MathNode{5, "Pow", 2, func(x []float64) float64 { return math.Pow(x[0], x[1]) }},
	"Sqrt":  MathNode{6, "Sqrt", 1, func(x []float64) float64 { return math.Sqrt(x[0]) }},
	"Exp":   MathNode{7, "Exp", 1, func(x []float64) float64 { return math.Exp(x[0]) }},
	"Pow10": MathNode{8, "Pow10", 1, func(x []float64) float64 { return math.Pow(10.0, x[0]) }},
	"Ln":    MathNode{9, "Ln", 1, func(x []float64) float64 { return math.Log(x[0]) }},
	"Log":   MathNode{10, "Log", 1, func(x []float64) float64 { return math.Log10(x[0]) }},
	"Log2":  MathNode{83, "Log2", 2, func(x []float64) float64 { return gepLog2(x[0], x[1]) }},
	"Floor": MathNode{12, "Floor", 1, func(x []float64) float64 { return math.Floor(x[0]) }},
	"Ceil":  MathNode{13, "Ceil", 1, func(x []float64) float64 { return math.Ceil(x[0]) }},
	"Abs":   MathNode{14, "Abs", 1, func(x []float64) float64 { return math.Abs(x[0]) }},
	"Inv":   MathNode{15, "Inv", 1, func(x []float64) float64 { return (1.0 / (x[0])) }},
	"Neg":   MathNode{17, "Neg", 1, func(x []float64) float64 { return (-(x[0])) }},
	"Nop":   MathNode{16, "Nop", 1, func(x []float64) float64 { return (x[0]) }},
	"X2":    MathNode{76, "X2", 1, func(x []float64) float64 { return math.Pow(x[0], 2.0) }},
	"X3":    MathNode{77, "X3", 1, func(x []float64) float64 { return math.Pow(x[0], 3.0) }},
	"X4":    MathNode{78, "X4", 1, func(x []float64) float64 { return math.Pow(x[0], 4.0) }},
	"X5":    MathNode{79, "X5", 1, func(x []float64) float64 { return math.Pow(x[0], 5.0) }},
	"3Rt":   MathNode{80, "3Rt", 1, func(x []float64) float64 { return gep3Rt(x[0]) }},
	"4Rt":   MathNode{81, "4Rt", 1, func(x []float64) float64 { return math.Pow(x[0], (1.0 / 4.0)) }},
	"5Rt":   MathNode{82, "5Rt", 1, func(x []float64) float64 { return gep5Rt(x[0]) }},
	"Add3":  MathNode{84, "Add3", 3, func(x []float64) float64 { return (x[0] + x[1] + x[2]) }},
	"Sub3":  MathNode{86, "Sub3", 3, func(x []float64) float64 { return (x[0] - x[1] - x[2]) }},
	"Mul3":  MathNode{88, "Mul3", 3, func(x []float64) float64 { return (x[0] * x[1] * x[2]) }},
	"Div3":  MathNode{90, "Div3", 3, func(x []float64) float64 { return (x[0] / x[1] / x[2]) }},
	"Add4":  MathNode{85, "Add4", 4, func(x []float64) float64 { return (x[0] + x[1] + x[2] + x[3]) }},
	"Sub4":  MathNode{87, "Sub4", 4, func(x []float64) float64 { return (x[0] - x[1] - x[2] - x[3]) }},
	"Mul4":  MathNode{89, "Mul4", 4, func(x []float64) float64 { return (x[0] * x[1] * x[2] * x[3]) }},
	"Div4":  MathNode{91, "Div4", 4, func(x []float64) float64 { return (x[0] / x[1] / x[2] / x[3]) }},
	"Min2":  MathNode{92, "Min2", 2, func(x []float64) float64 { return gepMin2(x[0], x[1]) }},
	"Min3":  MathNode{93, "Min3", 3, func(x []float64) float64 { return gepMin3(x[0], x[1], x[2]) }},
	"Min4":  MathNode{94, "Min4", 4, func(x []float64) float64 { return gepMin4(x[0], x[1], x[2], x[3]) }},
	"Max2":  MathNode{95, "Max2", 2, func(x []float64) float64 { return gepMax2(x[0], x[1]) }},
	"Max3":  MathNode{96, "Max3", 3, func(x []float64) float64 { return gepMax3(x[0], x[1], x[2]) }},
	"Max4":  MathNode{97, "Max4", 4, func(x []float64) float64 { return gepMax4(x[0], x[1], x[2], x[3]) }},
	"Avg2":  MathNode{98, "Avg2", 2, func(x []float64) float64 { return ((x[0] + x[1]) / 2.0) }},
	"Avg3":  MathNode{99, "Avg3", 3, func(x []float64) float64 { return ((x[0] + x[1] + x[2]) / 3.0) }},
	"Avg4":  MathNode{100, "Avg4", 4, func(x []float64) float64 { return ((x[0] + x[1] + x[2] + x[3]) / 4.0) }},
	"Logi":  MathNode{11, "Logi", 1, func(x []float64) float64 { return gepLogi(x[0]) }},
	"Logi2": MathNode{101, "Logi2", 2, func(x []float64) float64 { return gepLogi2(x[0], x[1]) }},
	"Logi3": MathNode{102, "Logi3", 3, func(x []float64) float64 { return gepLogi3(x[0], x[1], x[2]) }},
	"Logi4": MathNode{103, "Logi4", 4, func(x []float64) float64 { return gepLogi4(x[0], x[1], x[2], x[3]) }},
	"Gau":   MathNode{104, "Gau", 1, func(x []float64) float64 { return gepGau(x[0]) }},
	"Gau2":  MathNode{105, "Gau2", 2, func(x []float64) float64 { return gepGau2(x[0], x[1]) }},
	"Gau3":  MathNode{106, "Gau3", 3, func(x []float64) float64 { return gepGau3(x[0], x[1], x[2]) }},
	"Gau4":  MathNode{107, "Gau4", 4, func(x []float64) float64 { return gepGau4(x[0], x[1], x[2], x[3]) }},
	"Zero":  MathNode{70, "Zero", 1, func(x []float64) float64 { return (0.0) }},
	"One":   MathNode{71, "One", 1, func(x []float64) float64 { return (1.0) }},
	"Zero2": MathNode{72, "Zero2", 2, func(x []float64) float64 { return (0.0) }},
	"One2":  MathNode{73, "One2", 2, func(x []float64) float64 { return (1.0) }},
	"Pi":    MathNode{74, "Pi", 1, func(x []float64) float64 { return (math.Pi) }},
	"E":     MathNode{75, "E", 1, func(x []float64) float64 { return (math.E) }},
	"Sin":   MathNode{18, "Sin", 1, func(x []float64) float64 { return math.Sin(x[0]) }},
	"Cos":   MathNode{19, "Cos", 1, func(x []float64) float64 { return math.Cos(x[0]) }},
	"Tan":   MathNode{20, "Tan", 1, func(x []float64) float64 { return math.Tan(x[0]) }},
	"Csc":   MathNode{21, "Csc", 1, func(x []float64) float64 { return (1.0 / math.Sin(x[0])) }},
	"Sec":   MathNode{22, "Sec", 1, func(x []float64) float64 { return (1.0 / math.Cos(x[0])) }},
	"Cot":   MathNode{23, "Cot", 1, func(x []float64) float64 { return (1.0 / math.Tan(x[0])) }},
	"Asin":  MathNode{24, "Asin", 1, func(x []float64) float64 { return math.Asin(x[0]) }},
	"Acos":  MathNode{25, "Acos", 1, func(x []float64) float64 { return math.Acos(x[0]) }},
	"Atan":  MathNode{26, "Atan", 1, func(x []float64) float64 { return math.Atan(x[0]) }},
	"Acsc":  MathNode{27, "Acsc", 1, func(x []float64) float64 { return gepAcsc(x[0]) }},
	"Asec":  MathNode{28, "Asec", 1, func(x []float64) float64 { return gepAsec(x[0]) }},
	"Acot":  MathNode{29, "Acot", 1, func(x []float64) float64 { return gepAcot(x[0]) }},
	"Sinh":  MathNode{30, "Sinh", 1, func(x []float64) float64 { return math.Sinh(x[0]) }},
	"Cosh":  MathNode{31, "Cosh", 1, func(x []float64) float64 { return math.Cosh(x[0]) }},
	"Tanh":  MathNode{32, "Tanh", 1, func(x []float64) float64 { return math.Tanh(x[0]) }},
	"Csch":  MathNode{33, "Csch", 1, func(x []float64) float64 { return (1.0 / math.Sinh(x[0])) }},
	"Sech":  MathNode{34, "Sech", 1, func(x []float64) float64 { return (1.0 / math.Cosh(x[0])) }},
	"Coth":  MathNode{35, "Coth", 1, func(x []float64) float64 { return (1.0 / math.Tanh(x[0])) }},
	"Asinh": MathNode{36, "Asinh", 1, func(x []float64) float64 { return gepAsinh(x[0]) }},
	"Acosh": MathNode{37, "Acosh", 1, func(x []float64) float64 { return gepAcosh(x[0]) }},
	"Atanh": MathNode{38, "Atanh", 1, func(x []float64) float64 { return gepAtanh(x[0]) }},
	"Acsch": MathNode{39, "Acsch", 1, func(x []float64) float64 { return gepAcsch(x[0]) }},
	"Asech": MathNode{40, "Asech", 1, func(x []float64) float64 { return gepAsech(x[0]) }},
	"Acoth": MathNode{41, "Acoth", 1, func(x []float64) float64 { return gepAcoth(x[0]) }},
	"NOT":   MathNode{108, "NOT", 1, func(x []float64) float64 { return (1.0 - x[0]) }},
	"OR1":   MathNode{42, "OR1", 2, func(x []float64) float64 { return gepOR1(x[0], x[1]) }},
	"OR2":   MathNode{43, "OR2", 2, func(x []float64) float64 { return gepOR2(x[0], x[1]) }},
	"OR3":   MathNode{109, "OR3", 2, func(x []float64) float64 { return gepOR3(x[0], x[1]) }},
	"OR4":   MathNode{110, "OR4", 2, func(x []float64) float64 { return gepOR4(x[0], x[1]) }},
	"OR5":   MathNode{111, "OR5", 2, func(x []float64) float64 { return gepOR5(x[0], x[1]) }},
	"OR6":   MathNode{112, "OR6", 2, func(x []float64) float64 { return gepOR6(x[0], x[1]) }},
	"AND1":  MathNode{44, "AND1", 2, func(x []float64) float64 { return gepAND1(x[0], x[1]) }},
	"AND2":  MathNode{45, "AND2", 2, func(x []float64) float64 { return gepAND2(x[0], x[1]) }},
	"AND3":  MathNode{113, "AND3", 2, func(x []float64) float64 { return gepAND3(x[0], x[1]) }},
	"AND4":  MathNode{114, "AND4", 2, func(x []float64) float64 { return gepAND4(x[0], x[1]) }},
	"AND5":  MathNode{115, "AND5", 2, func(x []float64) float64 { return gepAND5(x[0], x[1]) }},
	"AND6":  MathNode{116, "AND6", 2, func(x []float64) float64 { return gepAND6(x[0], x[1]) }},
	"LT2A":  MathNode{46, "LT2A", 2, func(x []float64) float64 { return gepLT2A(x[0], x[1]) }},
	"GT2A":  MathNode{47, "GT2A", 2, func(x []float64) float64 { return gepGT2A(x[0], x[1]) }},
	"LOE2A": MathNode{48, "LOE2A", 2, func(x []float64) float64 { return gepLOE2A(x[0], x[1]) }},
	"GOE2A": MathNode{49, "GOE2A", 2, func(x []float64) float64 { return gepGOE2A(x[0], x[1]) }},
	"ET2A":  MathNode{50, "ET2A", 2, func(x []float64) float64 { return gepET2A(x[0], x[1]) }},
	"NET2A": MathNode{51, "NET2A", 2, func(x []float64) float64 { return gepNET2A(x[0], x[1]) }},
	"LT2B":  MathNode{52, "LT2B", 2, func(x []float64) float64 { return gepLT2B(x[0], x[1]) }},
	"GT2B":  MathNode{53, "GT2B", 2, func(x []float64) float64 { return gepGT2B(x[0], x[1]) }},
	"LOE2B": MathNode{54, "LOE2B", 2, func(x []float64) float64 { return gepLOE2B(x[0], x[1]) }},
	"GOE2B": MathNode{55, "GOE2B", 2, func(x []float64) float64 { return gepGOE2B(x[0], x[1]) }},
	"ET2B":  MathNode{56, "ET2B", 2, func(x []float64) float64 { return gepET2B(x[0], x[1]) }},
	"NET2B": MathNode{57, "NET2B", 2, func(x []float64) float64 { return gepNET2B(x[0], x[1]) }},
	"LT2C":  MathNode{117, "LT2C", 2, func(x []float64) float64 { return gepLT2C(x[0], x[1]) }},
	"GT2C":  MathNode{118, "GT2C", 2, func(x []float64) float64 { return gepGT2C(x[0], x[1]) }},
	"LOE2C": MathNode{119, "LOE2C", 2, func(x []float64) float64 { return gepLOE2C(x[0], x[1]) }},
	"GOE2C": MathNode{120, "GOE2C", 2, func(x []float64) float64 { return gepGOE2C(x[0], x[1]) }},
	"ET2C":  MathNode{121, "ET2C", 2, func(x []float64) float64 { return gepET2C(x[0], x[1]) }},
	"NET2C": MathNode{122, "NET2C", 2, func(x []float64) float64 { return gepNET2C(x[0], x[1]) }},
	"LT2D":  MathNode{123, "LT2D", 2, func(x []float64) float64 { return gepLT2D(x[0], x[1]) }},
	"GT2D":  MathNode{124, "GT2D", 2, func(x []float64) float64 { return gepGT2D(x[0], x[1]) }},
	"LOE2D": MathNode{125, "LOE2D", 2, func(x []float64) float64 { return gepLOE2D(x[0], x[1]) }},
	"GOE2D": MathNode{126, "GOE2D", 2, func(x []float64) float64 { return gepGOE2D(x[0], x[1]) }},
	"ET2D":  MathNode{127, "ET2D", 2, func(x []float64) float64 { return gepET2D(x[0], x[1]) }},
	"NET2D": MathNode{128, "NET2D", 2, func(x []float64) float64 { return gepNET2D(x[0], x[1]) }},
	"LT2E":  MathNode{129, "LT2E", 2, func(x []float64) float64 { return gepLT2E(x[0], x[1]) }},
	"GT2E":  MathNode{130, "GT2E", 2, func(x []float64) float64 { return gepGT2E(x[0], x[1]) }},
	"LOE2E": MathNode{131, "LOE2E", 2, func(x []float64) float64 { return gepLOE2E(x[0], x[1]) }},
	"GOE2E": MathNode{132, "GOE2E", 2, func(x []float64) float64 { return gepGOE2E(x[0], x[1]) }},
	"ET2E":  MathNode{133, "ET2E", 2, func(x []float64) float64 { return gepET2E(x[0], x[1]) }},
	"NET2E": MathNode{134, "NET2E", 2, func(x []float64) float64 { return gepNET2E(x[0], x[1]) }},
	"LT2F":  MathNode{135, "LT2F", 2, func(x []float64) float64 { return gepLT2F(x[0], x[1]) }},
	"GT2F":  MathNode{136, "GT2F", 2, func(x []float64) float64 { return gepGT2F(x[0], x[1]) }},
	"LOE2F": MathNode{137, "LOE2F", 2, func(x []float64) float64 { return gepLOE2F(x[0], x[1]) }},
	"GOE2F": MathNode{138, "GOE2F", 2, func(x []float64) float64 { return gepGOE2F(x[0], x[1]) }},
	"ET2F":  MathNode{139, "ET2F", 2, func(x []float64) float64 { return gepET2F(x[0], x[1]) }},
	"NET2F": MathNode{140, "NET2F", 2, func(x []float64) float64 { return gepNET2F(x[0], x[1]) }},
	"LT2G":  MathNode{141, "LT2G", 2, func(x []float64) float64 { return gepLT2G(x[0], x[1]) }},
	"GT2G":  MathNode{142, "GT2G", 2, func(x []float64) float64 { return gepGT2G(x[0], x[1]) }},
	"LOE2G": MathNode{143, "LOE2G", 2, func(x []float64) float64 { return gepLOE2G(x[0], x[1]) }},
	"GOE2G": MathNode{144, "GOE2G", 2, func(x []float64) float64 { return gepGOE2G(x[0], x[1]) }},
	"ET2G":  MathNode{145, "ET2G", 2, func(x []float64) float64 { return gepET2G(x[0], x[1]) }},
	"NET2G": MathNode{146, "NET2G", 2, func(x []float64) float64 { return gepNET2G(x[0], x[1]) }},
	"LT3A":  MathNode{58, "LT3A", 3, func(x []float64) float64 { return gepLT3A(x[0], x[1], x[2]) }},
	"GT3A":  MathNode{59, "GT3A", 3, func(x []float64) float64 { return gepGT3A(x[0], x[1], x[2]) }},
	"LOE3A": MathNode{60, "LOE3A", 3, func(x []float64) float64 { return gepLOE3A(x[0], x[1], x[2]) }},
	"GOE3A": MathNode{61, "GOE3A", 3, func(x []float64) float64 { return gepGOE3A(x[0], x[1], x[2]) }},
	"ET3A":  MathNode{62, "ET3A", 3, func(x []float64) float64 { return gepET3A(x[0], x[1], x[2]) }},
	"NET3A": MathNode{63, "NET3A", 3, func(x []float64) float64 { return gepNET3A(x[0], x[1], x[2]) }},
	"LT3B":  MathNode{147, "LT3B", 3, func(x []float64) float64 { return gepLT3B(x[0], x[1], x[2]) }},
	"GT3B":  MathNode{148, "GT3B", 3, func(x []float64) float64 { return gepGT3B(x[0], x[1], x[2]) }},
	"LOE3B": MathNode{149, "LOE3B", 3, func(x []float64) float64 { return gepLOE3B(x[0], x[1], x[2]) }},
	"GOE3B": MathNode{150, "GOE3B", 3, func(x []float64) float64 { return gepGOE3B(x[0], x[1], x[2]) }},
	"ET3B":  MathNode{151, "ET3B", 3, func(x []float64) float64 { return gepET3B(x[0], x[1], x[2]) }},
	"NET3B": MathNode{152, "NET3B", 3, func(x []float64) float64 { return gepNET3B(x[0], x[1], x[2]) }},
	"LT3C":  MathNode{153, "LT3C", 3, func(x []float64) float64 { return gepLT3C(x[0], x[1], x[2]) }},
	"GT3C":  MathNode{154, "GT3C", 3, func(x []float64) float64 { return gepGT3C(x[0], x[1], x[2]) }},
	"LOE3C": MathNode{155, "LOE3C", 3, func(x []float64) float64 { return gepLOE3C(x[0], x[1], x[2]) }},
	"GOE3C": MathNode{156, "GOE3C", 3, func(x []float64) float64 { return gepGOE3C(x[0], x[1], x[2]) }},
	"ET3C":  MathNode{157, "ET3C", 3, func(x []float64) float64 { return gepET3C(x[0], x[1], x[2]) }},
	"NET3C": MathNode{158, "NET3C", 3, func(x []float64) float64 { return gepNET3C(x[0], x[1], x[2]) }},
	"LT3D":  MathNode{159, "LT3D", 3, func(x []float64) float64 { return gepLT3D(x[0], x[1], x[2]) }},
	"GT3D":  MathNode{160, "GT3D", 3, func(x []float64) float64 { return gepGT3D(x[0], x[1], x[2]) }},
	"LOE3D": MathNode{161, "LOE3D", 3, func(x []float64) float64 { return gepLOE3D(x[0], x[1], x[2]) }},
	"GOE3D": MathNode{162, "GOE3D", 3, func(x []float64) float64 { return gepGOE3D(x[0], x[1], x[2]) }},
	"ET3D":  MathNode{163, "ET3D", 3, func(x []float64) float64 { return gepET3D(x[0], x[1], x[2]) }},
	"NET3D": MathNode{164, "NET3D", 3, func(x []float64) float64 { return gepNET3D(x[0], x[1], x[2]) }},
	"LT3E":  MathNode{165, "LT3E", 3, func(x []float64) float64 { return gepLT3E(x[0], x[1], x[2]) }},
	"GT3E":  MathNode{166, "GT3E", 3, func(x []float64) float64 { return gepGT3E(x[0], x[1], x[2]) }},
	"LOE3E": MathNode{167, "LOE3E", 3, func(x []float64) float64 { return gepLOE3E(x[0], x[1], x[2]) }},
	"GOE3E": MathNode{168, "GOE3E", 3, func(x []float64) float64 { return gepGOE3E(x[0], x[1], x[2]) }},
	"ET3E":  MathNode{169, "ET3E", 3, func(x []float64) float64 { return gepET3E(x[0], x[1], x[2]) }},
	"NET3E": MathNode{170, "NET3E", 3, func(x []float64) float64 { return gepNET3E(x[0], x[1], x[2]) }},
	"LT3F":  MathNode{171, "LT3F", 3, func(x []float64) float64 { return gepLT3F(x[0], x[1], x[2]) }},
	"GT3F":  MathNode{172, "GT3F", 3, func(x []float64) float64 { return gepGT3F(x[0], x[1], x[2]) }},
	"LOE3F": MathNode{173, "LOE3F", 3, func(x []float64) float64 { return gepLOE3F(x[0], x[1], x[2]) }},
	"GOE3F": MathNode{174, "GOE3F", 3, func(x []float64) float64 { return gepGOE3F(x[0], x[1], x[2]) }},
	"ET3F":  MathNode{175, "ET3F", 3, func(x []float64) float64 { return gepET3F(x[0], x[1], x[2]) }},
	"NET3F": MathNode{176, "NET3F", 3, func(x []float64) float64 { return gepNET3F(x[0], x[1], x[2]) }},
	"LT3G":  MathNode{177, "LT3G", 3, func(x []float64) float64 { return gepLT3G(x[0], x[1], x[2]) }},
	"GT3G":  MathNode{178, "GT3G", 3, func(x []float64) float64 { return gepGT3G(x[0], x[1], x[2]) }},
	"LOE3G": MathNode{179, "LOE3G", 3, func(x []float64) float64 { return gepLOE3G(x[0], x[1], x[2]) }},
	"GOE3G": MathNode{180, "GOE3G", 3, func(x []float64) float64 { return gepGOE3G(x[0], x[1], x[2]) }},
	"ET3G":  MathNode{181, "ET3G", 3, func(x []float64) float64 { return gepET3G(x[0], x[1], x[2]) }},
	"NET3G": MathNode{182, "NET3G", 3, func(x []float64) float64 { return gepNET3G(x[0], x[1], x[2]) }},
	"LT3H":  MathNode{183, "LT3H", 3, func(x []float64) float64 { return gepLT3H(x[0], x[1], x[2]) }},
	"GT3H":  MathNode{184, "GT3H", 3, func(x []float64) float64 { return gepGT3H(x[0], x[1], x[2]) }},
	"LOE3H": MathNode{185, "LOE3H", 3, func(x []float64) float64 { return gepLOE3H(x[0], x[1], x[2]) }},
	"GOE3H": MathNode{186, "GOE3H", 3, func(x []float64) float64 { return gepGOE3H(x[0], x[1], x[2]) }},
	"ET3H":  MathNode{187, "ET3H", 3, func(x []float64) float64 { return gepET3H(x[0], x[1], x[2]) }},
	"NET3H": MathNode{188, "NET3H", 3, func(x []float64) float64 { return gepNET3H(x[0], x[1], x[2]) }},
	"LT3I":  MathNode{189, "LT3I", 3, func(x []float64) float64 { return gepLT3I(x[0], x[1], x[2]) }},
	"GT3I":  MathNode{190, "GT3I", 3, func(x []float64) float64 { return gepGT3I(x[0], x[1], x[2]) }},
	"LOE3I": MathNode{191, "LOE3I", 3, func(x []float64) float64 { return gepLOE3I(x[0], x[1], x[2]) }},
	"GOE3I": MathNode{192, "GOE3I", 3, func(x []float64) float64 { return gepGOE3I(x[0], x[1], x[2]) }},
	"ET3I":  MathNode{193, "ET3I", 3, func(x []float64) float64 { return gepET3I(x[0], x[1], x[2]) }},
	"NET3I": MathNode{194, "NET3I", 3, func(x []float64) float64 { return gepNET3I(x[0], x[1], x[2]) }},
	"LT3J":  MathNode{195, "LT3J", 3, func(x []float64) float64 { return gepLT3J(x[0], x[1], x[2]) }},
	"GT3J":  MathNode{196, "GT3J", 3, func(x []float64) float64 { return gepGT3J(x[0], x[1], x[2]) }},
	"LOE3J": MathNode{197, "LOE3J", 3, func(x []float64) float64 { return gepLOE3J(x[0], x[1], x[2]) }},
	"GOE3J": MathNode{198, "GOE3J", 3, func(x []float64) float64 { return gepGOE3J(x[0], x[1], x[2]) }},
	"ET3J":  MathNode{199, "ET3J", 3, func(x []float64) float64 { return gepET3J(x[0], x[1], x[2]) }},
	"NET3J": MathNode{200, "NET3J", 3, func(x []float64) float64 { return gepNET3J(x[0], x[1], x[2]) }},
	"LT3K":  MathNode{201, "LT3K", 3, func(x []float64) float64 { return gepLT3K(x[0], x[1], x[2]) }},
	"GT3K":  MathNode{202, "GT3K", 3, func(x []float64) float64 { return gepGT3K(x[0], x[1], x[2]) }},
	"LOE3K": MathNode{203, "LOE3K", 3, func(x []float64) float64 { return gepLOE3K(x[0], x[1], x[2]) }},
	"GOE3K": MathNode{204, "GOE3K", 3, func(x []float64) float64 { return gepGOE3K(x[0], x[1], x[2]) }},
	"ET3K":  MathNode{205, "ET3K", 3, func(x []float64) float64 { return gepET3K(x[0], x[1], x[2]) }},
	"NET3K": MathNode{206, "NET3K", 3, func(x []float64) float64 { return gepNET3K(x[0], x[1], x[2]) }},
	"LT3L":  MathNode{207, "LT3L", 3, func(x []float64) float64 { return gepLT3L(x[0], x[1], x[2]) }},
	"GT3L":  MathNode{208, "GT3L", 3, func(x []float64) float64 { return gepGT3L(x[0], x[1], x[2]) }},
	"LOE3L": MathNode{209, "LOE3L", 3, func(x []float64) float64 { return gepLOE3L(x[0], x[1], x[2]) }},
	"GOE3L": MathNode{210, "GOE3L", 3, func(x []float64) float64 { return gepGOE3L(x[0], x[1], x[2]) }},
	"ET3L":  MathNode{211, "ET3L", 3, func(x []float64) float64 { return gepET3L(x[0], x[1], x[2]) }},
	"NET3L": MathNode{212, "NET3L", 3, func(x []float64) float64 { return gepNET3L(x[0], x[1], x[2]) }},
	"LT4A":  MathNode{64, "LT4A", 4, func(x []float64) float64 { return gepLT4A(x[0], x[1], x[2], x[3]) }},
	"GT4A":  MathNode{65, "GT4A", 4, func(x []float64) float64 { return gepGT4A(x[0], x[1], x[2], x[3]) }},
	"LOE4A": MathNode{66, "LOE4A", 4, func(x []float64) float64 { return gepLOE4A(x[0], x[1], x[2], x[3]) }},
	"GOE4A": MathNode{67, "GOE4A", 4, func(x []float64) float64 { return gepGOE4A(x[0], x[1], x[2], x[3]) }},
	"ET4A":  MathNode{68, "ET4A", 4, func(x []float64) float64 { return gepET4A(x[0], x[1], x[2], x[3]) }},
	"NET4A": MathNode{69, "NET4A", 4, func(x []float64) float64 { return gepNET4A(x[0], x[1], x[2], x[3]) }},
	"LT4B":  MathNode{213, "LT4B", 4, func(x []float64) float64 { return gepLT4B(x[0], x[1], x[2], x[3]) }},
	"GT4B":  MathNode{214, "GT4B", 4, func(x []float64) float64 { return gepGT4B(x[0], x[1], x[2], x[3]) }},
	"LOE4B": MathNode{215, "LOE4B", 4, func(x []float64) float64 { return gepLOE4B(x[0], x[1], x[2], x[3]) }},
	"GOE4B": MathNode{216, "GOE4B", 4, func(x []float64) float64 { return gepGOE4B(x[0], x[1], x[2], x[3]) }},
	"ET4B":  MathNode{217, "ET4B", 4, func(x []float64) float64 { return gepET4B(x[0], x[1], x[2], x[3]) }},
	"NET4B": MathNode{218, "NET4B", 4, func(x []float64) float64 { return gepNET4B(x[0], x[1], x[2], x[3]) }},
	"LT4C":  MathNode{219, "LT4C", 4, func(x []float64) float64 { return gepLT4C(x[0], x[1], x[2], x[3]) }},
	"GT4C":  MathNode{220, "GT4C", 4, func(x []float64) float64 { return gepGT4C(x[0], x[1], x[2], x[3]) }},
	"LOE4C": MathNode{221, "LOE4C", 4, func(x []float64) float64 { return gepLOE4C(x[0], x[1], x[2], x[3]) }},
	"GOE4C": MathNode{222, "GOE4C", 4, func(x []float64) float64 { return gepGOE4C(x[0], x[1], x[2], x[3]) }},
	"ET4C":  MathNode{223, "ET4C", 4, func(x []float64) float64 { return gepET4C(x[0], x[1], x[2], x[3]) }},
	"NET4C": MathNode{224, "NET4C", 4, func(x []float64) float64 { return gepNET4C(x[0], x[1], x[2], x[3]) }},
	"LT4D":  MathNode{225, "LT4D", 4, func(x []float64) float64 { return gepLT4D(x[0], x[1], x[2], x[3]) }},
	"GT4D":  MathNode{226, "GT4D", 4, func(x []float64) float64 { return gepGT4D(x[0], x[1], x[2], x[3]) }},
	"LOE4D": MathNode{227, "LOE4D", 4, func(x []float64) float64 { return gepLOE4D(x[0], x[1], x[2], x[3]) }},
	"GOE4D": MathNode{228, "GOE4D", 4, func(x []float64) float64 { return gepGOE4D(x[0], x[1], x[2], x[3]) }},
	"ET4D":  MathNode{229, "ET4D", 4, func(x []float64) float64 { return gepET4D(x[0], x[1], x[2], x[3]) }},
	"NET4D": MathNode{230, "NET4D", 4, func(x []float64) float64 { return gepNET4D(x[0], x[1], x[2], x[3]) }},
	"LT4E":  MathNode{231, "LT4E", 4, func(x []float64) float64 { return gepLT4E(x[0], x[1], x[2], x[3]) }},
	"GT4E":  MathNode{232, "GT4E", 4, func(x []float64) float64 { return gepGT4E(x[0], x[1], x[2], x[3]) }},
	"LOE4E": MathNode{233, "LOE4E", 4, func(x []float64) float64 { return gepLOE4E(x[0], x[1], x[2], x[3]) }},
	"GOE4E": MathNode{234, "GOE4E", 4, func(x []float64) float64 { return gepGOE4E(x[0], x[1], x[2], x[3]) }},
	"ET4E":  MathNode{235, "ET4E", 4, func(x []float64) float64 { return gepET4E(x[0], x[1], x[2], x[3]) }},
	"NET4E": MathNode{236, "NET4E", 4, func(x []float64) float64 { return gepNET4E(x[0], x[1], x[2], x[3]) }},
	"LT4F":  MathNode{237, "LT4F", 4, func(x []float64) float64 { return gepLT4F(x[0], x[1], x[2], x[3]) }},
	"GT4F":  MathNode{238, "GT4F", 4, func(x []float64) float64 { return gepGT4F(x[0], x[1], x[2], x[3]) }},
	"LOE4F": MathNode{239, "LOE4F", 4, func(x []float64) float64 { return gepLOE4F(x[0], x[1], x[2], x[3]) }},
	"GOE4F": MathNode{240, "GOE4F", 4, func(x []float64) float64 { return gepGOE4F(x[0], x[1], x[2], x[3]) }},
	"ET4F":  MathNode{241, "ET4F", 4, func(x []float64) float64 { return gepET4F(x[0], x[1], x[2], x[3]) }},
	"NET4F": MathNode{242, "NET4F", 4, func(x []float64) float64 { return gepNET4F(x[0], x[1], x[2], x[3]) }},
	"LT4G":  MathNode{243, "LT4G", 4, func(x []float64) float64 { return gepLT4G(x[0], x[1], x[2], x[3]) }},
	"GT4G":  MathNode{244, "GT4G", 4, func(x []float64) float64 { return gepGT4G(x[0], x[1], x[2], x[3]) }},
	"LOE4G": MathNode{245, "LOE4G", 4, func(x []float64) float64 { return gepLOE4G(x[0], x[1], x[2], x[3]) }},
	"GOE4G": MathNode{246, "GOE4G", 4, func(x []float64) float64 { return gepGOE4G(x[0], x[1], x[2], x[3]) }},
	"ET4G":  MathNode{247, "ET4G", 4, func(x []float64) float64 { return gepET4G(x[0], x[1], x[2], x[3]) }},
	"NET4G": MathNode{248, "NET4G", 4, func(x []float64) float64 { return gepNET4G(x[0], x[1], x[2], x[3]) }},
	"LT4H":  MathNode{249, "LT4H", 4, func(x []float64) float64 { return gepLT4H(x[0], x[1], x[2], x[3]) }},
	"GT4H":  MathNode{250, "GT4H", 4, func(x []float64) float64 { return gepGT4H(x[0], x[1], x[2], x[3]) }},
	"LOE4H": MathNode{251, "LOE4H", 4, func(x []float64) float64 { return gepLOE4H(x[0], x[1], x[2], x[3]) }},
	"GOE4H": MathNode{252, "GOE4H", 4, func(x []float64) float64 { return gepGOE4H(x[0], x[1], x[2], x[3]) }},
	"ET4H":  MathNode{253, "ET4H", 4, func(x []float64) float64 { return gepET4H(x[0], x[1], x[2], x[3]) }},
	"NET4H": MathNode{254, "NET4H", 4, func(x []float64) float64 { return gepNET4H(x[0], x[1], x[2], x[3]) }},
	"LT4I":  MathNode{255, "LT4I", 4, func(x []float64) float64 { return gepLT4I(x[0], x[1], x[2], x[3]) }},
	"GT4I":  MathNode{256, "GT4I", 4, func(x []float64) float64 { return gepGT4I(x[0], x[1], x[2], x[3]) }},
	"LOE4I": MathNode{257, "LOE4I", 4, func(x []float64) float64 { return gepLOE4I(x[0], x[1], x[2], x[3]) }},
	"GOE4I": MathNode{258, "GOE4I", 4, func(x []float64) float64 { return gepGOE4I(x[0], x[1], x[2], x[3]) }},
	"ET4I":  MathNode{259, "ET4I", 4, func(x []float64) float64 { return gepET4I(x[0], x[1], x[2], x[3]) }},
	"NET4I": MathNode{260, "NET4I", 4, func(x []float64) float64 { return gepNET4I(x[0], x[1], x[2], x[3]) }},
	"LT4J":  MathNode{261, "LT4J", 4, func(x []float64) float64 { return gepLT4J(x[0], x[1], x[2], x[3]) }},
	"GT4J":  MathNode{262, "GT4J", 4, func(x []float64) float64 { return gepGT4J(x[0], x[1], x[2], x[3]) }},
	"LOE4J": MathNode{263, "LOE4J", 4, func(x []float64) float64 { return gepLOE4J(x[0], x[1], x[2], x[3]) }},
	"GOE4J": MathNode{264, "GOE4J", 4, func(x []float64) float64 { return gepGOE4J(x[0], x[1], x[2], x[3]) }},
	"ET4J":  MathNode{265, "ET4J", 4, func(x []float64) float64 { return gepET4J(x[0], x[1], x[2], x[3]) }},
	"NET4J": MathNode{266, "NET4J", 4, func(x []float64) float64 { return gepNET4J(x[0], x[1], x[2], x[3]) }},
	"LT4K":  MathNode{267, "LT4K", 4, func(x []float64) float64 { return gepLT4K(x[0], x[1], x[2], x[3]) }},
	"GT4K":  MathNode{268, "GT4K", 4, func(x []float64) float64 { return gepGT4K(x[0], x[1], x[2], x[3]) }},
	"LOE4K": MathNode{269, "LOE4K", 4, func(x []float64) float64 { return gepLOE4K(x[0], x[1], x[2], x[3]) }},
	"GOE4K": MathNode{270, "GOE4K", 4, func(x []float64) float64 { return gepGOE4K(x[0], x[1], x[2], x[3]) }},
	"ET4K":  MathNode{271, "ET4K", 4, func(x []float64) float64 { return gepET4K(x[0], x[1], x[2], x[3]) }},
	"NET4K": MathNode{272, "NET4K", 4, func(x []float64) float64 { return gepNET4K(x[0], x[1], x[2], x[3]) }},
	"LT4L":  MathNode{273, "LT4L", 4, func(x []float64) float64 { return gepLT4L(x[0], x[1], x[2], x[3]) }},
	"GT4L":  MathNode{274, "GT4L", 4, func(x []float64) float64 { return gepGT4L(x[0], x[1], x[2], x[3]) }},
	"LOE4L": MathNode{275, "LOE4L", 4, func(x []float64) float64 { return gepLOE4L(x[0], x[1], x[2], x[3]) }},
	"GOE4L": MathNode{276, "GOE4L", 4, func(x []float64) float64 { return gepGOE4L(x[0], x[1], x[2], x[3]) }},
	"ET4L":  MathNode{277, "ET4L", 4, func(x []float64) float64 { return gepET4L(x[0], x[1], x[2], x[3]) }},
	"NET4L": MathNode{278, "NET4L", 4, func(x []float64) float64 { return gepNET4L(x[0], x[1], x[2], x[3]) }},
}

func gep3Rt(x float64) float64 {
	if x < 0.0 {
		return -math.Pow(-x, (1.0 / 3.0))
	}
	return math.Pow(x, (1.0 / 3.0))
}

func gep5Rt(x float64) float64 {
	if x < 0.0 {
		return -math.Pow(-x, (1.0 / 5.0))
	}
	return math.Pow(x, (1.0 / 5.0))
}

func gepLog2(x, y float64) float64 {
	if y == 0.0 {
		return 0.0
	}
	return math.Log(x) / math.Log(y)
}

func gepMod(x, y float64) float64 {
	// The built-in function is incorrect for cases such as -1.0 and 0.2.
	return ((x / y) - float64(int(x/y))) * y
}

func gepLogi(x float64) float64 {
	if math.Abs(x) > 709.0 {
		return 1.0 / (1.0 + math.Exp(math.Abs(x)/x*709.0))
	}
	return 1.0 / (1.0 + math.Exp(-x))
}

func gepLogi2(x, y float64) float64 {
	if math.Abs(x+y) > 709.0 {
		return 1.0 / (1.0 + math.Exp(math.Abs(x+y)/(x+y)*709.0))
	}
	return 1.0 / (1.0 + math.Exp(-(x + y)))
}

func gepLogi3(x, y, z float64) float64 {
	if math.Abs(x+y+z) > 709.0 {
		return 1.0 / (1.0 + math.Exp(math.Abs(x+y+z)/(x+y+z)*709.0))
	}
	return 1.0 / (1.0 + math.Exp(-(x + y + z)))
}

func gepLogi4(a, b, c, d float64) float64 {
	if math.Abs(a+b+c+d) > 709.0 {
		return 1.0 / (1.0 + math.Exp(math.Abs(a+b+c+d)/(a+b+c+d)*709.0))
	}
	return 1.0 / (1.0 + math.Exp(-(a + b + c + d)))
}

func gepGau(x float64) float64 {
	return math.Exp(-math.Pow(x, 2.0))
}

func gepGau2(x, y float64) float64 {
	return math.Exp(-math.Pow((x + y), 2.0))
}

func gepGau3(x, y, z float64) float64 {
	return math.Exp(-math.Pow((x + y + z), 2.0))
}

func gepGau4(a, b, c, d float64) float64 {
	return math.Exp(-math.Pow((a + b + c + d), 2.0))
}

func gepAcsc(x float64) float64 {
	varSign := 0.0
	if x < 0.0 {
		varSign = -1.0
	} else {
		if x > 0.0 {
			varSign = 1.0
		} else {
			varSign = 0.0
		}
	}
	return math.Atan(varSign / math.Sqrt(x*x-1.0))
}

func gepAsec(x float64) float64 {
	varSign := 0.0
	if x < 0.0 {
		varSign = -1.0
	} else {
		if x > 0.0 {
			varSign = 1.0
		} else {
			varSign = 0.0
		}
	}

	if math.Abs(x) == 1.0 {
		if x == -1.0 {
			return 4.0 * math.Atan(1.0)
		}
		return 0.0
	}
	return 2.0*math.Atan(1.0) - math.Atan(varSign/math.Sqrt(x*x-1.0))
}

func gepAcot(x float64) float64 {
	return math.Atan(1.0 / x)
}

func gepAsinh(x float64) float64 {
	return math.Log(x + math.Sqrt(x*x+1.0))
}

func gepAcosh(x float64) float64 {
	return math.Log(x + math.Sqrt(x*x-1.0))
}

func gepAtanh(x float64) float64 {
	return math.Log((1.0+x)/(1.0-x)) / 2.0
}

func gepAcsch(x float64) float64 {
	varSign := 0.0
	if x < 0.0 {
		varSign = -1.0
	} else {
		if x > 0.0 {
			varSign = 1.0
		} else {
			varSign = 0.0
		}
	}
	return math.Log((varSign*math.Sqrt(x*x+1.0) + 1.0) / x)
}

func gepAsech(x float64) float64 {
	return math.Log((math.Sqrt(-x*x+1.0) + 1.0) / x)
}

func gepAcoth(x float64) float64 {
	return math.Log((x+1.0)/(x-1.0)) / 2.0
}

func gepMin2(x, y float64) float64 {
	if x > y {
		return y
	}
	return x
}

func gepMin3(x, y, z float64) float64 {
	varTemp := x
	if varTemp > y {
		varTemp = y
	}
	if varTemp > z {
		varTemp = z
	}
	return varTemp
}

func gepMin4(a, b, c, d float64) float64 {
	varTemp := a
	if varTemp > b {
		varTemp = b
	}
	if varTemp > c {
		varTemp = c
	}
	if varTemp > d {
		varTemp = d
	}
	return varTemp
}

func gepMax2(x, y float64) float64 {
	if x < y {
		return y
	}
	return x
}

func gepMax3(x, y, z float64) float64 {
	varTemp := x
	if varTemp < y {
		varTemp = y
	}
	if varTemp < z {
		varTemp = z
	}
	return varTemp
}

func gepMax4(a, b, c, d float64) float64 {
	varTemp := a
	if varTemp < b {
		varTemp = b
	}
	if varTemp < c {
		varTemp = c
	}
	if varTemp < d {
		varTemp = d
	}
	return varTemp
}

func gepOR1(x, y float64) float64 {
	if (x < 0.0) || (y < 0.0) {
		return 1.0
	}
	return 0.0
}

func gepOR2(x, y float64) float64 {
	if (x >= 0.0) || (y >= 0.0) {
		return 1.0
	}
	return 0.0
}

func gepOR3(x, y float64) float64 {
	if (x <= 0.0) || (y <= 0.0) {
		return 1.0
	}
	return 0.0
}

func gepOR4(x, y float64) float64 {
	if (x < 1.0) || (y < 1.0) {
		return 1.0
	}
	return 0.0
}

func gepOR5(x, y float64) float64 {
	if (x >= 1.0) || (y >= 1.0) {
		return 1.0
	}
	return 0.0
}

func gepOR6(x, y float64) float64 {
	if (x <= 1.0) || (y <= 1.0) {
		return 1.0
	}
	return 0.0
}

func gepAND1(x, y float64) float64 {
	if (x < 0.0) && (y < 0.0) {
		return 1.0
	}
	return 0.0
}

func gepAND2(x, y float64) float64 {
	if (x >= 0.0) && (y >= 0.0) {
		return 1.0
	}
	return 0.0
}

func gepAND3(x, y float64) float64 {
	if (x <= 0.0) && (y <= 0.0) {
		return 1.0
	}
	return 0.0
}

func gepAND4(x, y float64) float64 {
	if (x < 1.0) && (y < 1.0) {
		return 1.0
	}
	return 0.0
}

func gepAND5(x, y float64) float64 {
	if (x >= 1.0) && (y >= 1.0) {
		return 1.0
	}
	return 0.0
}

func gepAND6(x, y float64) float64 {
	if (x <= 1.0) && (y <= 1.0) {
		return 1.0
	}
	return 0.0
}

func gepLT2A(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

func gepGT2A(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func gepLOE2A(x, y float64) float64 {
	if x <= y {
		return x
	}
	return y
}

func gepGOE2A(x, y float64) float64 {
	if x >= y {
		return x
	}
	return y
}

func gepET2A(x, y float64) float64 {
	if x == y {
		return x
	}
	return y
}

func gepNET2A(x, y float64) float64 {
	if x != y {
		return x
	}
	return y
}

func gepLT2B(x, y float64) float64 {
	if x < y {
		return 1.0
	}
	return 0.0
}

func gepGT2B(x, y float64) float64 {
	if x > y {
		return 1.0
	}
	return 0.0
}

func gepLOE2B(x, y float64) float64 {
	if x <= y {
		return 1.0
	}
	return 0.0
}

func gepGOE2B(x, y float64) float64 {
	if x >= y {
		return 1.0
	}
	return 0.0
}

func gepET2B(x, y float64) float64 {
	if x == y {
		return 1.0
	}
	return 0.0
}

func gepNET2B(x, y float64) float64 {
	if x != y {
		return 1.0
	}
	return 0.0
}

func gepLT2C(x, y float64) float64 {
	if x < y {
		return (x + y)
	}
	return (x - y)
}

func gepGT2C(x, y float64) float64 {
	if x > y {
		return (x + y)
	}
	return (x - y)
}

func gepLOE2C(x, y float64) float64 {
	if x <= y {
		return (x + y)
	}
	return (x - y)
}

func gepGOE2C(x, y float64) float64 {
	if x >= y {
		return (x + y)
	}
	return (x - y)
}

func gepET2C(x, y float64) float64 {
	if x == y {
		return (x + y)
	}
	return (x - y)
}

func gepNET2C(x, y float64) float64 {
	if x != y {
		return (x + y)
	}
	return (x - y)
}

func gepLT2D(x, y float64) float64 {
	if x < y {
		return (x * y)
	}
	return (x / y)
}

func gepGT2D(x, y float64) float64 {
	if x > y {
		return (x * y)
	}
	return (x / y)
}

func gepLOE2D(x, y float64) float64 {
	if x <= y {
		return (x * y)
	}
	return (x / y)
}

func gepGOE2D(x, y float64) float64 {
	if x >= y {
		return (x * y)
	}
	return (x / y)
}

func gepET2D(x, y float64) float64 {
	if x == y {
		return (x * y)
	}
	return (x / y)
}

func gepNET2D(x, y float64) float64 {
	if x != y {
		return (x * y)
	}
	return (x / y)
}

func gepLT2E(x, y float64) float64 {
	if x < y {
		return (x + y)
	}
	return (x * y)
}

func gepGT2E(x, y float64) float64 {
	if x > y {
		return (x + y)
	}
	return (x * y)
}

func gepLOE2E(x, y float64) float64 {
	if x <= y {
		return (x + y)
	}
	return (x * y)
}

func gepGOE2E(x, y float64) float64 {
	if x >= y {
		return (x + y)
	}
	return (x * y)
}

func gepET2E(x, y float64) float64 {
	if x == y {
		return (x + y)
	}
	return (x * y)
}

func gepNET2E(x, y float64) float64 {
	if x != y {
		return (x + y)
	}
	return (x * y)
}

func gepLT2F(x, y float64) float64 {
	if x < y {
		return (x + y)
	}
	return math.Sin(x * y)
}

func gepGT2F(x, y float64) float64 {
	if x > y {
		return (x + y)
	}
	return math.Sin(x * y)
}

func gepLOE2F(x, y float64) float64 {
	if x <= y {
		return (x + y)
	}
	return math.Sin(x * y)
}

func gepGOE2F(x, y float64) float64 {
	if x >= y {
		return (x + y)
	}
	return math.Sin(x * y)
}

func gepET2F(x, y float64) float64 {
	if x == y {
		return (x + y)
	}
	return math.Sin(x * y)
}

func gepNET2F(x, y float64) float64 {
	if x != y {
		return (x + y)
	}
	return math.Sin(x * y)
}

func gepLT2G(x, y float64) float64 {
	if x < y {
		return (x + y)
	}
	return math.Atan(x * y)
}

func gepGT2G(x, y float64) float64 {
	if x > y {
		return (x + y)
	}
	return math.Atan(x * y)
}

func gepLOE2G(x, y float64) float64 {
	if x <= y {
		return (x + y)
	}
	return math.Atan(x * y)
}

func gepGOE2G(x, y float64) float64 {
	if x >= y {
		return (x + y)
	}
	return math.Atan(x * y)
}

func gepET2G(x, y float64) float64 {
	if x == y {
		return (x + y)
	}
	return math.Atan(x * y)
}

func gepNET2G(x, y float64) float64 {
	if x != y {
		return (x + y)
	}
	return math.Atan(x * y)
}

func gepLT3A(x, y, z float64) float64 {
	if x < 0.0 {
		return y
	}
	return z
}

func gepGT3A(x, y, z float64) float64 {
	if x > 0.0 {
		return y
	}
	return z
}

func gepLOE3A(x, y, z float64) float64 {
	if x <= 0.0 {
		return y
	}
	return z
}

func gepGOE3A(x, y, z float64) float64 {
	if x >= 0.0 {
		return y
	}
	return z
}

func gepET3A(x, y, z float64) float64 {
	if x == 0.0 {
		return y
	}
	return z
}

func gepNET3A(x, y, z float64) float64 {
	if x != 0.0 {
		return y
	}
	return z
}

func gepLT3B(x, y, z float64) float64 {
	if (x + y) < z {
		return (x + y)
	}
	return z
}

func gepGT3B(x, y, z float64) float64 {
	if (x + y) > z {
		return (x + y)
	}
	return z
}

func gepLOE3B(x, y, z float64) float64 {
	if (x + y) <= z {
		return (x + y)
	}
	return z
}

func gepGOE3B(x, y, z float64) float64 {
	if (x + y) >= z {
		return (x + y)
	}
	return z
}

func gepET3B(x, y, z float64) float64 {
	if (x + y) == z {
		return (x + y)
	}
	return z
}

func gepNET3B(x, y, z float64) float64 {
	if (x + y) != z {
		return (x + y)
	}
	return z
}

func gepLT3C(x, y, z float64) float64 {
	if (x + y) < z {
		return (x + y)
	}
	return (x + z)
}

func gepGT3C(x, y, z float64) float64 {
	if (x + y) > z {
		return (x + y)
	}
	return (x + z)
}

func gepLOE3C(x, y, z float64) float64 {
	if (x + y) <= z {
		return (x + y)
	}
	return (x + z)
}

func gepGOE3C(x, y, z float64) float64 {
	if (x + y) >= z {
		return (x + y)
	}
	return (x + z)
}

func gepET3C(x, y, z float64) float64 {
	if (x + y) == z {
		return (x + y)
	}
	return (x + z)
}

func gepNET3C(x, y, z float64) float64 {
	if (x + y) != z {
		return (x + y)
	}
	return (x + z)
}

func gepLT3D(x, y, z float64) float64 {
	if (x + y) < z {
		return (x + y)
	}
	return (x - z)
}

func gepGT3D(x, y, z float64) float64 {
	if (x + y) > z {
		return (x + y)
	}
	return (x - z)
}

func gepLOE3D(x, y, z float64) float64 {
	if (x + y) <= z {
		return (x + y)
	}
	return (x - z)
}

func gepGOE3D(x, y, z float64) float64 {
	if (x + y) >= z {
		return (x + y)
	}
	return (x - z)
}

func gepET3D(x, y, z float64) float64 {
	if (x + y) == z {
		return (x + y)
	}
	return (x - z)
}

func gepNET3D(x, y, z float64) float64 {
	if (x + y) != z {
		return (x + y)
	}
	return (x - z)
}

func gepLT3E(x, y, z float64) float64 {
	if (x + y) < z {
		return (x + y)
	}
	return (x * z)
}

func gepGT3E(x, y, z float64) float64 {
	if (x + y) > z {
		return (x + y)
	}
	return (x * z)
}

func gepLOE3E(x, y, z float64) float64 {
	if (x + y) <= z {
		return (x + y)
	}
	return (x * z)
}

func gepGOE3E(x, y, z float64) float64 {
	if (x + y) >= z {
		return (x + y)
	}
	return (x * z)
}

func gepET3E(x, y, z float64) float64 {
	if (x + y) == z {
		return (x + y)
	}
	return (x * z)
}

func gepNET3E(x, y, z float64) float64 {
	if (x + y) != z {
		return (x + y)
	}
	return (x * z)
}

func gepLT3F(x, y, z float64) float64 {
	if (x + y) < z {
		return (x + y)
	}
	return (x / z)
}

func gepGT3F(x, y, z float64) float64 {
	if (x + y) > z {
		return (x + y)
	}
	return (x / z)
}

func gepLOE3F(x, y, z float64) float64 {
	if (x + y) <= z {
		return (x + y)
	}
	return (x / z)
}

func gepGOE3F(x, y, z float64) float64 {
	if (x + y) >= z {
		return (x + y)
	}
	return (x / z)
}

func gepET3F(x, y, z float64) float64 {
	if (x + y) == z {
		return (x + y)
	}
	return (x / z)
}

func gepNET3F(x, y, z float64) float64 {
	if (x + y) != z {
		return (x + y)
	}
	return (x / z)
}

func gepLT3G(x, y, z float64) float64 {
	if (x + y) < z {
		return (x * y)
	}
	return (x + z)
}

func gepGT3G(x, y, z float64) float64 {
	if (x + y) > z {
		return (x * y)
	}
	return (x + z)
}

func gepLOE3G(x, y, z float64) float64 {
	if (x + y) <= z {
		return (x * y)
	}
	return (x + z)
}

func gepGOE3G(x, y, z float64) float64 {
	if (x + y) >= z {
		return (x * y)
	}
	return (x + z)
}

func gepET3G(x, y, z float64) float64 {
	if (x + y) == z {
		return (x * y)
	}
	return (x + z)
}

func gepNET3G(x, y, z float64) float64 {
	if (x + y) != z {
		return (x * y)
	}
	return (x + z)
}

func gepLT3H(x, y, z float64) float64 {
	if (x + y) < z {
		return (x * y)
	}
	return (x - z)
}

func gepGT3H(x, y, z float64) float64 {
	if (x + y) > z {
		return (x * y)
	}
	return (x - z)
}

func gepLOE3H(x, y, z float64) float64 {
	if (x + y) <= z {
		return (x * y)
	}
	return (x - z)
}

func gepGOE3H(x, y, z float64) float64 {
	if (x + y) >= z {
		return (x * y)
	}
	return (x - z)
}

func gepET3H(x, y, z float64) float64 {
	if (x + y) == z {
		return (x * y)
	}
	return (x - z)
}

func gepNET3H(x, y, z float64) float64 {
	if (x + y) != z {
		return (x * y)
	}
	return (x - z)
}

func gepLT3I(x, y, z float64) float64 {
	if (x + y) < z {
		return (x * y)
	}
	return (x * z)
}

func gepGT3I(x, y, z float64) float64 {
	if (x + y) > z {
		return (x * y)
	}
	return (x * z)
}

func gepLOE3I(x, y, z float64) float64 {
	if (x + y) <= z {
		return (x * y)
	}
	return (x * z)
}

func gepGOE3I(x, y, z float64) float64 {
	if (x + y) >= z {
		return (x * y)
	}
	return (x * z)
}

func gepET3I(x, y, z float64) float64 {
	if (x + y) == z {
		return (x * y)
	}
	return (x * z)
}

func gepNET3I(x, y, z float64) float64 {
	if (x + y) != z {
		return (x * y)
	}
	return (x * z)
}

func gepLT3J(x, y, z float64) float64 {
	if (x + y) < z {
		return (x * y)
	}
	return (x / z)
}

func gepGT3J(x, y, z float64) float64 {
	if (x + y) > z {
		return (x * y)
	}
	return (x / z)
}

func gepLOE3J(x, y, z float64) float64 {
	if (x + y) <= z {
		return (x * y)
	}
	return (x / z)
}

func gepGOE3J(x, y, z float64) float64 {
	if (x + y) >= z {
		return (x * y)
	}
	return (x / z)
}

func gepET3J(x, y, z float64) float64 {
	if (x + y) == z {
		return (x * y)
	}
	return (x / z)
}

func gepNET3J(x, y, z float64) float64 {
	if (x + y) != z {
		return (x * y)
	}
	return (x / z)
}

func gepLT3K(x, y, z float64) float64 {
	if (x + y) < z {
		return (x + y + z)
	}
	return math.Sin(x * y * z)
}

func gepGT3K(x, y, z float64) float64 {
	if (x + y) > z {
		return (x + y + z)
	}
	return math.Sin(x * y * z)
}

func gepLOE3K(x, y, z float64) float64 {
	if (x + y) <= z {
		return (x + y + z)
	}
	return math.Sin(x * y * z)
}

func gepGOE3K(x, y, z float64) float64 {
	if (x + y) >= z {
		return (x + y + z)
	}
	return math.Sin(x * y * z)
}

func gepET3K(x, y, z float64) float64 {
	if (x + y) == z {
		return (x + y + z)
	}
	return math.Sin(x * y * z)
}

func gepNET3K(x, y, z float64) float64 {
	if (x + y) != z {
		return (x + y + z)
	}
	return math.Sin(x * y * z)
}

func gepLT3L(x, y, z float64) float64 {
	if (x + y) < z {
		return (x + y + z)
	}
	return math.Atan(x * y * z)
}

func gepGT3L(x, y, z float64) float64 {
	if (x + y) > z {
		return (x + y + z)
	}
	return math.Atan(x * y * z)
}

func gepLOE3L(x, y, z float64) float64 {
	if (x + y) <= z {
		return (x + y + z)
	}
	return math.Atan(x * y * z)
}

func gepGOE3L(x, y, z float64) float64 {
	if (x + y) >= z {
		return (x + y + z)
	}
	return math.Atan(x * y * z)
}

func gepET3L(x, y, z float64) float64 {
	if (x + y) == z {
		return (x + y + z)
	}
	return math.Atan(x * y * z)
}

func gepNET3L(x, y, z float64) float64 {
	if (x + y) != z {
		return (x + y + z)
	}
	return math.Atan(x * y * z)
}

func gepLT4A(a, b, c, d float64) float64 {
	if a < b {
		return c
	}
	return d
}

func gepGT4A(a, b, c, d float64) float64 {
	if a > b {
		return c
	}
	return d
}

func gepLOE4A(a, b, c, d float64) float64 {
	if a <= b {
		return c
	}
	return d
}

func gepGOE4A(a, b, c, d float64) float64 {
	if a >= b {
		return c
	}
	return d
}

func gepET4A(a, b, c, d float64) float64 {
	if a == b {
		return c
	}
	return d
}

func gepNET4A(a, b, c, d float64) float64 {
	if a != b {
		return c
	}
	return d
}

func gepLT4B(a, b, c, d float64) float64 {
	if (a + b) < (c + d) {
		return c
	}
	return d
}

func gepGT4B(a, b, c, d float64) float64 {
	if (a + b) > (c + d) {
		return c
	}
	return d
}

func gepLOE4B(a, b, c, d float64) float64 {
	if (a + b) <= (c + d) {
		return c
	}
	return d
}

func gepGOE4B(a, b, c, d float64) float64 {
	if (a + b) >= (c + d) {
		return c
	}
	return d
}

func gepET4B(a, b, c, d float64) float64 {
	if (a + b) == (c + d) {
		return c
	}
	return d
}

func gepNET4B(a, b, c, d float64) float64 {
	if (a + b) != (c + d) {
		return c
	}
	return d
}

func gepLT4C(a, b, c, d float64) float64 {
	if (a + b) < (c + d) {
		return (a + b)
	}
	return (c + d)
}

func gepGT4C(a, b, c, d float64) float64 {
	if (a + b) > (c + d) {
		return (a + b)
	}
	return (c + d)
}

func gepLOE4C(a, b, c, d float64) float64 {
	if (a + b) <= (c + d) {
		return (a + b)
	}
	return (c + d)
}

func gepGOE4C(a, b, c, d float64) float64 {
	if (a + b) >= (c + d) {
		return (a + b)
	}
	return (c + d)
}

func gepET4C(a, b, c, d float64) float64 {
	if (a + b) == (c + d) {
		return (a + b)
	}
	return (c + d)
}

func gepNET4C(a, b, c, d float64) float64 {
	if (a + b) != (c + d) {
		return (a + b)
	}
	return (c + d)
}

func gepLT4D(a, b, c, d float64) float64 {
	if (a + b) < (c + d) {
		return (a + b)
	}
	return (c - d)
}

func gepGT4D(a, b, c, d float64) float64 {
	if (a + b) > (c + d) {
		return (a + b)
	}
	return (c - d)
}

func gepLOE4D(a, b, c, d float64) float64 {
	if (a + b) <= (c + d) {
		return (a + b)
	}
	return (c - d)
}

func gepGOE4D(a, b, c, d float64) float64 {
	if (a + b) >= (c + d) {
		return (a + b)
	}
	return (c - d)
}

func gepET4D(a, b, c, d float64) float64 {
	if (a + b) == (c + d) {
		return (a + b)
	}
	return (c - d)
}

func gepNET4D(a, b, c, d float64) float64 {
	if (a + b) != (c + d) {
		return (a + b)
	}
	return (c - d)
}

func gepLT4E(a, b, c, d float64) float64 {
	if (a + b) < (c + d) {
		return (a + b)
	}
	return (c * d)
}

func gepGT4E(a, b, c, d float64) float64 {
	if (a + b) > (c + d) {
		return (a + b)
	}
	return (c * d)
}

func gepLOE4E(a, b, c, d float64) float64 {
	if (a + b) <= (c + d) {
		return (a + b)
	}
	return (c * d)
}

func gepGOE4E(a, b, c, d float64) float64 {
	if (a + b) >= (c + d) {
		return (a + b)
	}
	return (c * d)
}

func gepET4E(a, b, c, d float64) float64 {
	if (a + b) == (c + d) {
		return (a + b)
	}
	return (c * d)
}

func gepNET4E(a, b, c, d float64) float64 {
	if (a + b) != (c + d) {
		return (a + b)
	}
	return (c * d)
}

func gepLT4F(a, b, c, d float64) float64 {
	if (a + b) < (c + d) {
		return (a + b)
	}
	return (c / d)
}

func gepGT4F(a, b, c, d float64) float64 {
	if (a + b) > (c + d) {
		return (a + b)
	}
	return (c / d)
}

func gepLOE4F(a, b, c, d float64) float64 {
	if (a + b) <= (c + d) {
		return (a + b)
	}
	return (c / d)
}

func gepGOE4F(a, b, c, d float64) float64 {
	if (a + b) >= (c + d) {
		return (a + b)
	}
	return (c / d)
}

func gepET4F(a, b, c, d float64) float64 {
	if (a + b) == (c + d) {
		return (a + b)
	}
	return (c / d)
}

func gepNET4F(a, b, c, d float64) float64 {
	if (a + b) != (c + d) {
		return (a + b)
	}
	return (c / d)
}

func gepLT4G(a, b, c, d float64) float64 {
	if (a + b) < (c + d) {
		return (a * b)
	}
	return (c + d)
}

func gepGT4G(a, b, c, d float64) float64 {
	if (a + b) > (c + d) {
		return (a * b)
	}
	return (c + d)
}

func gepLOE4G(a, b, c, d float64) float64 {
	if (a + b) <= (c + d) {
		return (a * b)
	}
	return (c + d)
}

func gepGOE4G(a, b, c, d float64) float64 {
	if (a + b) >= (c + d) {
		return (a * b)
	}
	return (c + d)
}

func gepET4G(a, b, c, d float64) float64 {
	if (a + b) == (c + d) {
		return (a * b)
	}
	return (c + d)
}

func gepNET4G(a, b, c, d float64) float64 {
	if (a + b) != (c + d) {
		return (a * b)
	}
	return (c + d)
}

func gepLT4H(a, b, c, d float64) float64 {
	if (a + b) < (c + d) {
		return (a * b)
	}
	return (c - d)
}

func gepGT4H(a, b, c, d float64) float64 {
	if (a + b) > (c + d) {
		return (a * b)
	}
	return (c - d)
}

func gepLOE4H(a, b, c, d float64) float64 {
	if (a + b) <= (c + d) {
		return (a * b)
	}
	return (c - d)
}

func gepGOE4H(a, b, c, d float64) float64 {
	if (a + b) >= (c + d) {
		return (a * b)
	}
	return (c - d)
}

func gepET4H(a, b, c, d float64) float64 {
	if (a + b) == (c + d) {
		return (a * b)
	}
	return (c - d)
}

func gepNET4H(a, b, c, d float64) float64 {
	if (a + b) != (c + d) {
		return (a * b)
	}
	return (c - d)
}

func gepLT4I(a, b, c, d float64) float64 {
	if (a + b) < (c + d) {
		return (a * b)
	}
	return (c * d)
}

func gepGT4I(a, b, c, d float64) float64 {
	if (a + b) > (c + d) {
		return (a * b)
	}
	return (c * d)
}

func gepLOE4I(a, b, c, d float64) float64 {
	if (a + b) <= (c + d) {
		return (a * b)
	}
	return (c * d)
}

func gepGOE4I(a, b, c, d float64) float64 {
	if (a + b) >= (c + d) {
		return (a * b)
	}
	return (c * d)
}

func gepET4I(a, b, c, d float64) float64 {
	if (a + b) == (c + d) {
		return (a * b)
	}
	return (c * d)
}

func gepNET4I(a, b, c, d float64) float64 {
	if (a + b) != (c + d) {
		return (a * b)
	}
	return (c * d)
}

func gepLT4J(a, b, c, d float64) float64 {
	if (a + b) < (c + d) {
		return (a * b)
	}
	return (c / d)
}

func gepGT4J(a, b, c, d float64) float64 {
	if (a + b) > (c + d) {
		return (a * b)
	}
	return (c / d)
}

func gepLOE4J(a, b, c, d float64) float64 {
	if (a + b) <= (c + d) {
		return (a * b)
	}
	return (c / d)
}

func gepGOE4J(a, b, c, d float64) float64 {
	if (a + b) >= (c + d) {
		return (a * b)
	}
	return (c / d)
}

func gepET4J(a, b, c, d float64) float64 {
	if (a + b) == (c + d) {
		return (a * b)
	}
	return (c / d)
}

func gepNET4J(a, b, c, d float64) float64 {
	if (a + b) != (c + d) {
		return (a * b)
	}
	return (c / d)
}

func gepLT4K(a, b, c, d float64) float64 {
	if (a + b) < (c + d) {
		return math.Sin(a * b)
	}
	return math.Sin(c * d)
}

func gepGT4K(a, b, c, d float64) float64 {
	if (a + b) > (c + d) {
		return math.Sin(a * b)
	}
	return math.Sin(c * d)
}

func gepLOE4K(a, b, c, d float64) float64 {
	if (a + b) <= (c + d) {
		return math.Sin(a * b)
	}
	return math.Sin(c * d)
}

func gepGOE4K(a, b, c, d float64) float64 {
	if (a + b) >= (c + d) {
		return math.Sin(a * b)
	}
	return math.Sin(c * d)
}

func gepET4K(a, b, c, d float64) float64 {
	if (a + b) == (c + d) {
		return math.Sin(a * b)
	}
	return math.Sin(c * d)
}

func gepNET4K(a, b, c, d float64) float64 {
	if (a + b) != (c + d) {
		return math.Sin(a * b)
	}
	return math.Sin(c * d)
}

func gepLT4L(a, b, c, d float64) float64 {
	if (a + b) < (c + d) {
		return math.Atan(a * b)
	}
	return math.Atan(c * d)
}

func gepGT4L(a, b, c, d float64) float64 {
	if (a + b) > (c + d) {
		return math.Atan(a * b)
	}
	return math.Atan(c * d)
}

func gepLOE4L(a, b, c, d float64) float64 {
	if (a + b) <= (c + d) {
		return math.Atan(a * b)
	}
	return math.Atan(c * d)
}

func gepGOE4L(a, b, c, d float64) float64 {
	if (a + b) >= (c + d) {
		return math.Atan(a * b)
	}
	return math.Atan(c * d)
}

func gepET4L(a, b, c, d float64) float64 {
	if (a + b) == (c + d) {
		return math.Atan(a * b)
	}
	return math.Atan(c * d)
}

func gepNET4L(a, b, c, d float64) float64 {
	if (a + b) != (c + d) {
		return math.Atan(a * b)
	}
	return math.Atan(c * d)
}
