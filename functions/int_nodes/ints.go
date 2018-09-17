// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package intNodes defines the floating-point function collection available for the GEP algorithm.
package intNodes

import (
	"log"

	"github.com/gmlewis/gep/v2/functions"
)

// IntNode is a floating-point function used for the formation of GEP expressions.
type IntNode struct {
	index     int
	symbol    string
	terminals int
	function  func(x []int) int
}

// Symbol returns the Karva symbol for this floating-point function.
func (n IntNode) Symbol() string {
	return n.symbol
}

// Terminals returns the number of input terminals for this floating-point function.
func (n IntNode) Terminals() int {
	return n.terminals
}

// BoolFunction is unused in this package and returns an error.
func (n IntNode) BoolFunction([]bool) bool {
	log.Println("error calling BoolFunction on IntNode model.")
	return false
}

// IntFunction calls the floating-point function and returns the result.
func (n IntNode) IntFunction(x []int) int {
	return n.function(x)
}

// Float64Function is unused in this package and returns an error.
func (n IntNode) Float64Function([]float64) float64 {
	log.Println("error calling Float64Function on IntNode model.")
	return 0.0
}

// Int lists all the available floating-point functions for this package.
var Int = functions.FuncMap{
	// TODO(gmlewis): Change functions to operate on the entire length of the slice.
	"+": IntNode{0, "+", 2, func(x []int) int { return (x[0] + x[1]) }},
	"-": IntNode{1, "-", 2, func(x []int) int { return (x[0] - x[1]) }},
	"*": IntNode{2, "*", 2, func(x []int) int { return (x[0] * x[1]) }},
	"/": IntNode{3, "/", 2, func(x []int) int { return (x[0] / x[1]) }},
	// "Mod": IntNode{4, "Mod", 2, func(x []int) int { return gepMod(x[0], x[1]) }},
	// "Pow":   IntNode{5, "Pow", 2, func(x []int) int { return math.Pow(x[0], x[1]) }},
	// "Sqrt":  IntNode{6, "Sqrt", 1, func(x []int) int { return math.Sqrt(x[0]) }},
	// "Exp":   IntNode{7, "Exp", 1, func(x []int) int { return math.Exp(x[0]) }},
	// "Pow10": IntNode{8, "Pow10", 1, func(x []int) int { return math.Pow(10.0, x[0]) }},
	// "Ln":    IntNode{9, "Ln", 1, func(x []int) int { return math.Log(x[0]) }},
	// "Log":   IntNode{10, "Log", 1, func(x []int) int { return math.Log10(x[0]) }},
	// "Log2": IntNode{83, "Log2", 2, func(x []int) int { return gepLog2(x[0], x[1]) }},
	// "Floor": IntNode{12, "Floor", 1, func(x []int) int { return math.Floor(x[0]) }},
	// "Ceil":  IntNode{13, "Ceil", 1, func(x []int) int { return math.Ceil(x[0]) }},
	// "Abs":   IntNode{14, "Abs", 1, func(x []int) int { return math.Abs(x[0]) }},
	"Inv": IntNode{15, "Inv", 1, func(x []int) int { return (1.0 / (x[0])) }},
	"Neg": IntNode{17, "Neg", 1, func(x []int) int { return (-(x[0])) }},
	"Nop": IntNode{16, "Nop", 1, func(x []int) int { return (x[0]) }},
	// "X2":    IntNode{76, "X2", 1, func(x []int) int { return math.Pow(x[0], 2.0) }},
	// "X3":    IntNode{77, "X3", 1, func(x []int) int { return math.Pow(x[0], 3.0) }},
	// "X4":    IntNode{78, "X4", 1, func(x []int) int { return math.Pow(x[0], 4.0) }},
	// "X5":    IntNode{79, "X5", 1, func(x []int) int { return math.Pow(x[0], 5.0) }},
	// "3Rt": IntNode{80, "3Rt", 1, func(x []int) int { return gep3Rt(x[0]) }},
	// "4Rt":   IntNode{81, "4Rt", 1, func(x []int) int { return math.Pow(x[0], (1.0 / 4.0)) }},
	// "5Rt":   IntNode{82, "5Rt", 1, func(x []int) int { return gep5Rt(x[0]) }},
	"Add3": IntNode{84, "Add3", 3, func(x []int) int { return (x[0] + x[1] + x[2]) }},
	"Sub3": IntNode{86, "Sub3", 3, func(x []int) int { return (x[0] - x[1] - x[2]) }},
	"Mul3": IntNode{88, "Mul3", 3, func(x []int) int { return (x[0] * x[1] * x[2]) }},
	"Div3": IntNode{90, "Div3", 3, func(x []int) int { return (x[0] / x[1] / x[2]) }},
	"Add4": IntNode{85, "Add4", 4, func(x []int) int { return (x[0] + x[1] + x[2] + x[3]) }},
	"Sub4": IntNode{87, "Sub4", 4, func(x []int) int { return (x[0] - x[1] - x[2] - x[3]) }},
	"Mul4": IntNode{89, "Mul4", 4, func(x []int) int { return (x[0] * x[1] * x[2] * x[3]) }},
	"Div4": IntNode{91, "Div4", 4, func(x []int) int { return (x[0] / x[1] / x[2] / x[3]) }},
	"Min2": IntNode{92, "Min2", 2, func(x []int) int { return gepMin2(x[0], x[1]) }},
	"Min3": IntNode{93, "Min3", 3, func(x []int) int { return gepMin3(x[0], x[1], x[2]) }},
	"Min4": IntNode{94, "Min4", 4, func(x []int) int { return gepMin4(x[0], x[1], x[2], x[3]) }},
	"Max2": IntNode{95, "Max2", 2, func(x []int) int { return gepMax2(x[0], x[1]) }},
	"Max3": IntNode{96, "Max3", 3, func(x []int) int { return gepMax3(x[0], x[1], x[2]) }},
	"Max4": IntNode{97, "Max4", 4, func(x []int) int { return gepMax4(x[0], x[1], x[2], x[3]) }},
	"Avg2": IntNode{98, "Avg2", 2, func(x []int) int { return ((x[0] + x[1]) / 2.0) }},
	"Avg3": IntNode{99, "Avg3", 3, func(x []int) int { return ((x[0] + x[1] + x[2]) / 3.0) }},
	"Avg4": IntNode{100, "Avg4", 4, func(x []int) int { return ((x[0] + x[1] + x[2] + x[3]) / 4.0) }},
	// "Logi":  IntNode{11, "Logi", 1, func(x []int) int { return gepLogi(x[0]) }},
	// "Logi2": IntNode{101, "Logi2", 2, func(x []int) int { return gepLogi2(x[0], x[1]) }},
	// "Logi3": IntNode{102, "Logi3", 3, func(x []int) int { return gepLogi3(x[0], x[1], x[2]) }},
	// "Logi4": IntNode{103, "Logi4", 4, func(x []int) int { return gepLogi4(x[0], x[1], x[2], x[3]) }},
	// "Gau":   IntNode{104, "Gau", 1, func(x []int) int { return gepGau(x[0]) }},
	// "Gau2":  IntNode{105, "Gau2", 2, func(x []int) int { return gepGau2(x[0], x[1]) }},
	// "Gau3":  IntNode{106, "Gau3", 3, func(x []int) int { return gepGau3(x[0], x[1], x[2]) }},
	// "Gau4":  IntNode{107, "Gau4", 4, func(x []int) int { return gepGau4(x[0], x[1], x[2], x[3]) }},
	"Zero":  IntNode{70, "Zero", 1, func(x []int) int { return 0 }},
	"One":   IntNode{71, "One", 1, func(x []int) int { return 1 }},
	"Zero2": IntNode{72, "Zero2", 2, func(x []int) int { return 0 }},
	"One2":  IntNode{73, "One2", 2, func(x []int) int { return 1 }},
	// "Pi":    IntNode{74, "Pi", 1, func(x []int) int { return (math.Pi) }},
	// "E":     IntNode{75, "E", 1, func(x []int) int { return (math.E) }},
	// "Sin":   IntNode{18, "Sin", 1, func(x []int) int { return math.Sin(x[0]) }},
	// "Cos":   IntNode{19, "Cos", 1, func(x []int) int { return math.Cos(x[0]) }},
	// "Tan":   IntNode{20, "Tan", 1, func(x []int) int { return math.Tan(x[0]) }},
	// "Csc":   IntNode{21, "Csc", 1, func(x []int) int { return (1.0 / math.Sin(x[0])) }},
	// "Sec":   IntNode{22, "Sec", 1, func(x []int) int { return (1.0 / math.Cos(x[0])) }},
	// "Cot":   IntNode{23, "Cot", 1, func(x []int) int { return (1.0 / math.Tan(x[0])) }},
	// "Asin":  IntNode{24, "Asin", 1, func(x []int) int { return math.Asin(x[0]) }},
	// "Acos":  IntNode{25, "Acos", 1, func(x []int) int { return math.Acos(x[0]) }},
	// "Atan":  IntNode{26, "Atan", 1, func(x []int) int { return math.Atan(x[0]) }},
	// "Acsc": IntNode{27, "Acsc", 1, func(x []int) int { return gepAcsc(x[0]) }},
	// "Asec": IntNode{28, "Asec", 1, func(x []int) int { return gepAsec(x[0]) }},
	// "Acot": IntNode{29, "Acot", 1, func(x []int) int { return gepAcot(x[0]) }},
	// "Sinh":  IntNode{30, "Sinh", 1, func(x []int) int { return math.Sinh(x[0]) }},
	// "Cosh":  IntNode{31, "Cosh", 1, func(x []int) int { return math.Cosh(x[0]) }},
	// "Tanh":  IntNode{32, "Tanh", 1, func(x []int) int { return math.Tanh(x[0]) }},
	// "Csch":  IntNode{33, "Csch", 1, func(x []int) int { return (1.0 / math.Sinh(x[0])) }},
	// "Sech":  IntNode{34, "Sech", 1, func(x []int) int { return (1.0 / math.Cosh(x[0])) }},
	// "Coth":  IntNode{35, "Coth", 1, func(x []int) int { return (1.0 / math.Tanh(x[0])) }},
	// "Asinh": IntNode{36, "Asinh", 1, func(x []int) int { return gepAsinh(x[0]) }},
	// "Acosh": IntNode{37, "Acosh", 1, func(x []int) int { return gepAcosh(x[0]) }},
	// "Atanh": IntNode{38, "Atanh", 1, func(x []int) int { return gepAtanh(x[0]) }},
	// "Acsch": IntNode{39, "Acsch", 1, func(x []int) int { return gepAcsch(x[0]) }},
	// "Asech": IntNode{40, "Asech", 1, func(x []int) int { return gepAsech(x[0]) }},
	// "Acoth": IntNode{41, "Acoth", 1, func(x []int) int { return gepAcoth(x[0]) }},
	// "NOT":   IntNode{108, "NOT", 1, func(x []int) int { return (1.0 - x[0]) }},
	// "OR1":   IntNode{42, "OR1", 2, func(x []int) int { return gepOR1(x[0], x[1]) }},
	// "OR2":   IntNode{43, "OR2", 2, func(x []int) int { return gepOR2(x[0], x[1]) }},
	// "OR3":   IntNode{109, "OR3", 2, func(x []int) int { return gepOR3(x[0], x[1]) }},
	// "OR4":   IntNode{110, "OR4", 2, func(x []int) int { return gepOR4(x[0], x[1]) }},
	// "OR5":   IntNode{111, "OR5", 2, func(x []int) int { return gepOR5(x[0], x[1]) }},
	// "OR6":   IntNode{112, "OR6", 2, func(x []int) int { return gepOR6(x[0], x[1]) }},
	// "AND1":  IntNode{44, "AND1", 2, func(x []int) int { return gepAND1(x[0], x[1]) }},
	// "AND2":  IntNode{45, "AND2", 2, func(x []int) int { return gepAND2(x[0], x[1]) }},
	// "AND3":  IntNode{113, "AND3", 2, func(x []int) int { return gepAND3(x[0], x[1]) }},
	// "AND4":  IntNode{114, "AND4", 2, func(x []int) int { return gepAND4(x[0], x[1]) }},
	// "AND5":  IntNode{115, "AND5", 2, func(x []int) int { return gepAND5(x[0], x[1]) }},
	// "AND6":  IntNode{116, "AND6", 2, func(x []int) int { return gepAND6(x[0], x[1]) }},
	"LT2A":  IntNode{46, "LT2A", 2, func(x []int) int { return gepLT2A(x[0], x[1]) }},
	"GT2A":  IntNode{47, "GT2A", 2, func(x []int) int { return gepGT2A(x[0], x[1]) }},
	"LOE2A": IntNode{48, "LOE2A", 2, func(x []int) int { return gepLOE2A(x[0], x[1]) }},
	"GOE2A": IntNode{49, "GOE2A", 2, func(x []int) int { return gepGOE2A(x[0], x[1]) }},
	"ET2A":  IntNode{50, "ET2A", 2, func(x []int) int { return gepET2A(x[0], x[1]) }},
	"NET2A": IntNode{51, "NET2A", 2, func(x []int) int { return gepNET2A(x[0], x[1]) }},
	"LT2B":  IntNode{52, "LT2B", 2, func(x []int) int { return gepLT2B(x[0], x[1]) }},
	"GT2B":  IntNode{53, "GT2B", 2, func(x []int) int { return gepGT2B(x[0], x[1]) }},
	"LOE2B": IntNode{54, "LOE2B", 2, func(x []int) int { return gepLOE2B(x[0], x[1]) }},
	"GOE2B": IntNode{55, "GOE2B", 2, func(x []int) int { return gepGOE2B(x[0], x[1]) }},
	"ET2B":  IntNode{56, "ET2B", 2, func(x []int) int { return gepET2B(x[0], x[1]) }},
	"NET2B": IntNode{57, "NET2B", 2, func(x []int) int { return gepNET2B(x[0], x[1]) }},
	"LT2C":  IntNode{117, "LT2C", 2, func(x []int) int { return gepLT2C(x[0], x[1]) }},
	"GT2C":  IntNode{118, "GT2C", 2, func(x []int) int { return gepGT2C(x[0], x[1]) }},
	"LOE2C": IntNode{119, "LOE2C", 2, func(x []int) int { return gepLOE2C(x[0], x[1]) }},
	"GOE2C": IntNode{120, "GOE2C", 2, func(x []int) int { return gepGOE2C(x[0], x[1]) }},
	"ET2C":  IntNode{121, "ET2C", 2, func(x []int) int { return gepET2C(x[0], x[1]) }},
	"NET2C": IntNode{122, "NET2C", 2, func(x []int) int { return gepNET2C(x[0], x[1]) }},
	// "LT2D":  IntNode{123, "LT2D", 2, func(x []int) int { return gepLT2D(x[0], x[1]) }},
	// "GT2D":  IntNode{124, "GT2D", 2, func(x []int) int { return gepGT2D(x[0], x[1]) }},
	// "LOE2D": IntNode{125, "LOE2D", 2, func(x []int) int { return gepLOE2D(x[0], x[1]) }},
	// "GOE2D": IntNode{126, "GOE2D", 2, func(x []int) int { return gepGOE2D(x[0], x[1]) }},
	// "ET2D":  IntNode{127, "ET2D", 2, func(x []int) int { return gepET2D(x[0], x[1]) }},
	// "NET2D": IntNode{128, "NET2D", 2, func(x []int) int { return gepNET2D(x[0], x[1]) }},
	"LT2E":  IntNode{129, "LT2E", 2, func(x []int) int { return gepLT2E(x[0], x[1]) }},
	"GT2E":  IntNode{130, "GT2E", 2, func(x []int) int { return gepGT2E(x[0], x[1]) }},
	"LOE2E": IntNode{131, "LOE2E", 2, func(x []int) int { return gepLOE2E(x[0], x[1]) }},
	"GOE2E": IntNode{132, "GOE2E", 2, func(x []int) int { return gepGOE2E(x[0], x[1]) }},
	"ET2E":  IntNode{133, "ET2E", 2, func(x []int) int { return gepET2E(x[0], x[1]) }},
	"NET2E": IntNode{134, "NET2E", 2, func(x []int) int { return gepNET2E(x[0], x[1]) }},
	// "LT2F":  IntNode{135, "LT2F", 2, func(x []int) int { return gepLT2F(x[0], x[1]) }},
	// "GT2F":  IntNode{136, "GT2F", 2, func(x []int) int { return gepGT2F(x[0], x[1]) }},
	// "LOE2F": IntNode{137, "LOE2F", 2, func(x []int) int { return gepLOE2F(x[0], x[1]) }},
	// "GOE2F": IntNode{138, "GOE2F", 2, func(x []int) int { return gepGOE2F(x[0], x[1]) }},
	// "ET2F":  IntNode{139, "ET2F", 2, func(x []int) int { return gepET2F(x[0], x[1]) }},
	// "NET2F": IntNode{140, "NET2F", 2, func(x []int) int { return gepNET2F(x[0], x[1]) }},
	// "LT2G":  IntNode{141, "LT2G", 2, func(x []int) int { return gepLT2G(x[0], x[1]) }},
	// "GT2G":  IntNode{142, "GT2G", 2, func(x []int) int { return gepGT2G(x[0], x[1]) }},
	// "LOE2G": IntNode{143, "LOE2G", 2, func(x []int) int { return gepLOE2G(x[0], x[1]) }},
	// "GOE2G": IntNode{144, "GOE2G", 2, func(x []int) int { return gepGOE2G(x[0], x[1]) }},
	// "ET2G":  IntNode{145, "ET2G", 2, func(x []int) int { return gepET2G(x[0], x[1]) }},
	// "NET2G": IntNode{146, "NET2G", 2, func(x []int) int { return gepNET2G(x[0], x[1]) }},
	"LT3A":  IntNode{58, "LT3A", 3, func(x []int) int { return gepLT3A(x[0], x[1], x[2]) }},
	"GT3A":  IntNode{59, "GT3A", 3, func(x []int) int { return gepGT3A(x[0], x[1], x[2]) }},
	"LOE3A": IntNode{60, "LOE3A", 3, func(x []int) int { return gepLOE3A(x[0], x[1], x[2]) }},
	"GOE3A": IntNode{61, "GOE3A", 3, func(x []int) int { return gepGOE3A(x[0], x[1], x[2]) }},
	"ET3A":  IntNode{62, "ET3A", 3, func(x []int) int { return gepET3A(x[0], x[1], x[2]) }},
	"NET3A": IntNode{63, "NET3A", 3, func(x []int) int { return gepNET3A(x[0], x[1], x[2]) }},
	"LT3B":  IntNode{147, "LT3B", 3, func(x []int) int { return gepLT3B(x[0], x[1], x[2]) }},
	"GT3B":  IntNode{148, "GT3B", 3, func(x []int) int { return gepGT3B(x[0], x[1], x[2]) }},
	"LOE3B": IntNode{149, "LOE3B", 3, func(x []int) int { return gepLOE3B(x[0], x[1], x[2]) }},
	"GOE3B": IntNode{150, "GOE3B", 3, func(x []int) int { return gepGOE3B(x[0], x[1], x[2]) }},
	"ET3B":  IntNode{151, "ET3B", 3, func(x []int) int { return gepET3B(x[0], x[1], x[2]) }},
	"NET3B": IntNode{152, "NET3B", 3, func(x []int) int { return gepNET3B(x[0], x[1], x[2]) }},
	"LT3C":  IntNode{153, "LT3C", 3, func(x []int) int { return gepLT3C(x[0], x[1], x[2]) }},
	"GT3C":  IntNode{154, "GT3C", 3, func(x []int) int { return gepGT3C(x[0], x[1], x[2]) }},
	"LOE3C": IntNode{155, "LOE3C", 3, func(x []int) int { return gepLOE3C(x[0], x[1], x[2]) }},
	"GOE3C": IntNode{156, "GOE3C", 3, func(x []int) int { return gepGOE3C(x[0], x[1], x[2]) }},
	"ET3C":  IntNode{157, "ET3C", 3, func(x []int) int { return gepET3C(x[0], x[1], x[2]) }},
	"NET3C": IntNode{158, "NET3C", 3, func(x []int) int { return gepNET3C(x[0], x[1], x[2]) }},
	"LT3D":  IntNode{159, "LT3D", 3, func(x []int) int { return gepLT3D(x[0], x[1], x[2]) }},
	"GT3D":  IntNode{160, "GT3D", 3, func(x []int) int { return gepGT3D(x[0], x[1], x[2]) }},
	"LOE3D": IntNode{161, "LOE3D", 3, func(x []int) int { return gepLOE3D(x[0], x[1], x[2]) }},
	"GOE3D": IntNode{162, "GOE3D", 3, func(x []int) int { return gepGOE3D(x[0], x[1], x[2]) }},
	"ET3D":  IntNode{163, "ET3D", 3, func(x []int) int { return gepET3D(x[0], x[1], x[2]) }},
	"NET3D": IntNode{164, "NET3D", 3, func(x []int) int { return gepNET3D(x[0], x[1], x[2]) }},
	"LT3E":  IntNode{165, "LT3E", 3, func(x []int) int { return gepLT3E(x[0], x[1], x[2]) }},
	"GT3E":  IntNode{166, "GT3E", 3, func(x []int) int { return gepGT3E(x[0], x[1], x[2]) }},
	"LOE3E": IntNode{167, "LOE3E", 3, func(x []int) int { return gepLOE3E(x[0], x[1], x[2]) }},
	"GOE3E": IntNode{168, "GOE3E", 3, func(x []int) int { return gepGOE3E(x[0], x[1], x[2]) }},
	"ET3E":  IntNode{169, "ET3E", 3, func(x []int) int { return gepET3E(x[0], x[1], x[2]) }},
	"NET3E": IntNode{170, "NET3E", 3, func(x []int) int { return gepNET3E(x[0], x[1], x[2]) }},
	// "LT3F":  IntNode{171, "LT3F", 3, func(x []int) int { return gepLT3F(x[0], x[1], x[2]) }},
	// "GT3F":  IntNode{172, "GT3F", 3, func(x []int) int { return gepGT3F(x[0], x[1], x[2]) }},
	// "LOE3F": IntNode{173, "LOE3F", 3, func(x []int) int { return gepLOE3F(x[0], x[1], x[2]) }},
	// "GOE3F": IntNode{174, "GOE3F", 3, func(x []int) int { return gepGOE3F(x[0], x[1], x[2]) }},
	// "ET3F":  IntNode{175, "ET3F", 3, func(x []int) int { return gepET3F(x[0], x[1], x[2]) }},
	// "NET3F": IntNode{176, "NET3F", 3, func(x []int) int { return gepNET3F(x[0], x[1], x[2]) }},
	"LT3G":  IntNode{177, "LT3G", 3, func(x []int) int { return gepLT3G(x[0], x[1], x[2]) }},
	"GT3G":  IntNode{178, "GT3G", 3, func(x []int) int { return gepGT3G(x[0], x[1], x[2]) }},
	"LOE3G": IntNode{179, "LOE3G", 3, func(x []int) int { return gepLOE3G(x[0], x[1], x[2]) }},
	"GOE3G": IntNode{180, "GOE3G", 3, func(x []int) int { return gepGOE3G(x[0], x[1], x[2]) }},
	"ET3G":  IntNode{181, "ET3G", 3, func(x []int) int { return gepET3G(x[0], x[1], x[2]) }},
	"NET3G": IntNode{182, "NET3G", 3, func(x []int) int { return gepNET3G(x[0], x[1], x[2]) }},
	"LT3H":  IntNode{183, "LT3H", 3, func(x []int) int { return gepLT3H(x[0], x[1], x[2]) }},
	"GT3H":  IntNode{184, "GT3H", 3, func(x []int) int { return gepGT3H(x[0], x[1], x[2]) }},
	"LOE3H": IntNode{185, "LOE3H", 3, func(x []int) int { return gepLOE3H(x[0], x[1], x[2]) }},
	"GOE3H": IntNode{186, "GOE3H", 3, func(x []int) int { return gepGOE3H(x[0], x[1], x[2]) }},
	"ET3H":  IntNode{187, "ET3H", 3, func(x []int) int { return gepET3H(x[0], x[1], x[2]) }},
	"NET3H": IntNode{188, "NET3H", 3, func(x []int) int { return gepNET3H(x[0], x[1], x[2]) }},
	"LT3I":  IntNode{189, "LT3I", 3, func(x []int) int { return gepLT3I(x[0], x[1], x[2]) }},
	"GT3I":  IntNode{190, "GT3I", 3, func(x []int) int { return gepGT3I(x[0], x[1], x[2]) }},
	"LOE3I": IntNode{191, "LOE3I", 3, func(x []int) int { return gepLOE3I(x[0], x[1], x[2]) }},
	"GOE3I": IntNode{192, "GOE3I", 3, func(x []int) int { return gepGOE3I(x[0], x[1], x[2]) }},
	"ET3I":  IntNode{193, "ET3I", 3, func(x []int) int { return gepET3I(x[0], x[1], x[2]) }},
	"NET3I": IntNode{194, "NET3I", 3, func(x []int) int { return gepNET3I(x[0], x[1], x[2]) }},
	// "LT3J":  IntNode{195, "LT3J", 3, func(x []int) int { return gepLT3J(x[0], x[1], x[2]) }},
	// "GT3J":  IntNode{196, "GT3J", 3, func(x []int) int { return gepGT3J(x[0], x[1], x[2]) }},
	// "LOE3J": IntNode{197, "LOE3J", 3, func(x []int) int { return gepLOE3J(x[0], x[1], x[2]) }},
	// "GOE3J": IntNode{198, "GOE3J", 3, func(x []int) int { return gepGOE3J(x[0], x[1], x[2]) }},
	// "ET3J":  IntNode{199, "ET3J", 3, func(x []int) int { return gepET3J(x[0], x[1], x[2]) }},
	// "NET3J": IntNode{200, "NET3J", 3, func(x []int) int { return gepNET3J(x[0], x[1], x[2]) }},
	// "LT3K":  IntNode{201, "LT3K", 3, func(x []int) int { return gepLT3K(x[0], x[1], x[2]) }},
	// "GT3K":  IntNode{202, "GT3K", 3, func(x []int) int { return gepGT3K(x[0], x[1], x[2]) }},
	// "LOE3K": IntNode{203, "LOE3K", 3, func(x []int) int { return gepLOE3K(x[0], x[1], x[2]) }},
	// "GOE3K": IntNode{204, "GOE3K", 3, func(x []int) int { return gepGOE3K(x[0], x[1], x[2]) }},
	// "ET3K":  IntNode{205, "ET3K", 3, func(x []int) int { return gepET3K(x[0], x[1], x[2]) }},
	// "NET3K": IntNode{206, "NET3K", 3, func(x []int) int { return gepNET3K(x[0], x[1], x[2]) }},
	// "LT3L":  IntNode{207, "LT3L", 3, func(x []int) int { return gepLT3L(x[0], x[1], x[2]) }},
	// "GT3L":  IntNode{208, "GT3L", 3, func(x []int) int { return gepGT3L(x[0], x[1], x[2]) }},
	// "LOE3L": IntNode{209, "LOE3L", 3, func(x []int) int { return gepLOE3L(x[0], x[1], x[2]) }},
	// "GOE3L": IntNode{210, "GOE3L", 3, func(x []int) int { return gepGOE3L(x[0], x[1], x[2]) }},
	// "ET3L":  IntNode{211, "ET3L", 3, func(x []int) int { return gepET3L(x[0], x[1], x[2]) }},
	// "NET3L": IntNode{212, "NET3L", 3, func(x []int) int { return gepNET3L(x[0], x[1], x[2]) }},
	"LT4A":  IntNode{64, "LT4A", 4, func(x []int) int { return gepLT4A(x[0], x[1], x[2], x[3]) }},
	"GT4A":  IntNode{65, "GT4A", 4, func(x []int) int { return gepGT4A(x[0], x[1], x[2], x[3]) }},
	"LOE4A": IntNode{66, "LOE4A", 4, func(x []int) int { return gepLOE4A(x[0], x[1], x[2], x[3]) }},
	"GOE4A": IntNode{67, "GOE4A", 4, func(x []int) int { return gepGOE4A(x[0], x[1], x[2], x[3]) }},
	"ET4A":  IntNode{68, "ET4A", 4, func(x []int) int { return gepET4A(x[0], x[1], x[2], x[3]) }},
	"NET4A": IntNode{69, "NET4A", 4, func(x []int) int { return gepNET4A(x[0], x[1], x[2], x[3]) }},
	"LT4B":  IntNode{213, "LT4B", 4, func(x []int) int { return gepLT4B(x[0], x[1], x[2], x[3]) }},
	"GT4B":  IntNode{214, "GT4B", 4, func(x []int) int { return gepGT4B(x[0], x[1], x[2], x[3]) }},
	"LOE4B": IntNode{215, "LOE4B", 4, func(x []int) int { return gepLOE4B(x[0], x[1], x[2], x[3]) }},
	"GOE4B": IntNode{216, "GOE4B", 4, func(x []int) int { return gepGOE4B(x[0], x[1], x[2], x[3]) }},
	"ET4B":  IntNode{217, "ET4B", 4, func(x []int) int { return gepET4B(x[0], x[1], x[2], x[3]) }},
	"NET4B": IntNode{218, "NET4B", 4, func(x []int) int { return gepNET4B(x[0], x[1], x[2], x[3]) }},
	"LT4C":  IntNode{219, "LT4C", 4, func(x []int) int { return gepLT4C(x[0], x[1], x[2], x[3]) }},
	"GT4C":  IntNode{220, "GT4C", 4, func(x []int) int { return gepGT4C(x[0], x[1], x[2], x[3]) }},
	"LOE4C": IntNode{221, "LOE4C", 4, func(x []int) int { return gepLOE4C(x[0], x[1], x[2], x[3]) }},
	"GOE4C": IntNode{222, "GOE4C", 4, func(x []int) int { return gepGOE4C(x[0], x[1], x[2], x[3]) }},
	"ET4C":  IntNode{223, "ET4C", 4, func(x []int) int { return gepET4C(x[0], x[1], x[2], x[3]) }},
	"NET4C": IntNode{224, "NET4C", 4, func(x []int) int { return gepNET4C(x[0], x[1], x[2], x[3]) }},
	"LT4D":  IntNode{225, "LT4D", 4, func(x []int) int { return gepLT4D(x[0], x[1], x[2], x[3]) }},
	"GT4D":  IntNode{226, "GT4D", 4, func(x []int) int { return gepGT4D(x[0], x[1], x[2], x[3]) }},
	"LOE4D": IntNode{227, "LOE4D", 4, func(x []int) int { return gepLOE4D(x[0], x[1], x[2], x[3]) }},
	"GOE4D": IntNode{228, "GOE4D", 4, func(x []int) int { return gepGOE4D(x[0], x[1], x[2], x[3]) }},
	"ET4D":  IntNode{229, "ET4D", 4, func(x []int) int { return gepET4D(x[0], x[1], x[2], x[3]) }},
	"NET4D": IntNode{230, "NET4D", 4, func(x []int) int { return gepNET4D(x[0], x[1], x[2], x[3]) }},
	"LT4E":  IntNode{231, "LT4E", 4, func(x []int) int { return gepLT4E(x[0], x[1], x[2], x[3]) }},
	"GT4E":  IntNode{232, "GT4E", 4, func(x []int) int { return gepGT4E(x[0], x[1], x[2], x[3]) }},
	"LOE4E": IntNode{233, "LOE4E", 4, func(x []int) int { return gepLOE4E(x[0], x[1], x[2], x[3]) }},
	"GOE4E": IntNode{234, "GOE4E", 4, func(x []int) int { return gepGOE4E(x[0], x[1], x[2], x[3]) }},
	"ET4E":  IntNode{235, "ET4E", 4, func(x []int) int { return gepET4E(x[0], x[1], x[2], x[3]) }},
	"NET4E": IntNode{236, "NET4E", 4, func(x []int) int { return gepNET4E(x[0], x[1], x[2], x[3]) }},
	// "LT4F":  IntNode{237, "LT4F", 4, func(x []int) int { return gepLT4F(x[0], x[1], x[2], x[3]) }},
	// "GT4F":  IntNode{238, "GT4F", 4, func(x []int) int { return gepGT4F(x[0], x[1], x[2], x[3]) }},
	// "LOE4F": IntNode{239, "LOE4F", 4, func(x []int) int { return gepLOE4F(x[0], x[1], x[2], x[3]) }},
	// "GOE4F": IntNode{240, "GOE4F", 4, func(x []int) int { return gepGOE4F(x[0], x[1], x[2], x[3]) }},
	// "ET4F":  IntNode{241, "ET4F", 4, func(x []int) int { return gepET4F(x[0], x[1], x[2], x[3]) }},
	// "NET4F": IntNode{242, "NET4F", 4, func(x []int) int { return gepNET4F(x[0], x[1], x[2], x[3]) }},
	"LT4G":  IntNode{243, "LT4G", 4, func(x []int) int { return gepLT4G(x[0], x[1], x[2], x[3]) }},
	"GT4G":  IntNode{244, "GT4G", 4, func(x []int) int { return gepGT4G(x[0], x[1], x[2], x[3]) }},
	"LOE4G": IntNode{245, "LOE4G", 4, func(x []int) int { return gepLOE4G(x[0], x[1], x[2], x[3]) }},
	"GOE4G": IntNode{246, "GOE4G", 4, func(x []int) int { return gepGOE4G(x[0], x[1], x[2], x[3]) }},
	"ET4G":  IntNode{247, "ET4G", 4, func(x []int) int { return gepET4G(x[0], x[1], x[2], x[3]) }},
	"NET4G": IntNode{248, "NET4G", 4, func(x []int) int { return gepNET4G(x[0], x[1], x[2], x[3]) }},
	"LT4H":  IntNode{249, "LT4H", 4, func(x []int) int { return gepLT4H(x[0], x[1], x[2], x[3]) }},
	"GT4H":  IntNode{250, "GT4H", 4, func(x []int) int { return gepGT4H(x[0], x[1], x[2], x[3]) }},
	"LOE4H": IntNode{251, "LOE4H", 4, func(x []int) int { return gepLOE4H(x[0], x[1], x[2], x[3]) }},
	"GOE4H": IntNode{252, "GOE4H", 4, func(x []int) int { return gepGOE4H(x[0], x[1], x[2], x[3]) }},
	"ET4H":  IntNode{253, "ET4H", 4, func(x []int) int { return gepET4H(x[0], x[1], x[2], x[3]) }},
	"NET4H": IntNode{254, "NET4H", 4, func(x []int) int { return gepNET4H(x[0], x[1], x[2], x[3]) }},
	"LT4I":  IntNode{255, "LT4I", 4, func(x []int) int { return gepLT4I(x[0], x[1], x[2], x[3]) }},
	"GT4I":  IntNode{256, "GT4I", 4, func(x []int) int { return gepGT4I(x[0], x[1], x[2], x[3]) }},
	"LOE4I": IntNode{257, "LOE4I", 4, func(x []int) int { return gepLOE4I(x[0], x[1], x[2], x[3]) }},
	"GOE4I": IntNode{258, "GOE4I", 4, func(x []int) int { return gepGOE4I(x[0], x[1], x[2], x[3]) }},
	"ET4I":  IntNode{259, "ET4I", 4, func(x []int) int { return gepET4I(x[0], x[1], x[2], x[3]) }},
	"NET4I": IntNode{260, "NET4I", 4, func(x []int) int { return gepNET4I(x[0], x[1], x[2], x[3]) }},
	// "LT4J":  IntNode{261, "LT4J", 4, func(x []int) int { return gepLT4J(x[0], x[1], x[2], x[3]) }},
	// "GT4J":  IntNode{262, "GT4J", 4, func(x []int) int { return gepGT4J(x[0], x[1], x[2], x[3]) }},
	// "LOE4J": IntNode{263, "LOE4J", 4, func(x []int) int { return gepLOE4J(x[0], x[1], x[2], x[3]) }},
	// "GOE4J": IntNode{264, "GOE4J", 4, func(x []int) int { return gepGOE4J(x[0], x[1], x[2], x[3]) }},
	// "ET4J":  IntNode{265, "ET4J", 4, func(x []int) int { return gepET4J(x[0], x[1], x[2], x[3]) }},
	// "NET4J": IntNode{266, "NET4J", 4, func(x []int) int { return gepNET4J(x[0], x[1], x[2], x[3]) }},
	// "LT4K":  IntNode{267, "LT4K", 4, func(x []int) int { return gepLT4K(x[0], x[1], x[2], x[3]) }},
	// "GT4K":  IntNode{268, "GT4K", 4, func(x []int) int { return gepGT4K(x[0], x[1], x[2], x[3]) }},
	// "LOE4K": IntNode{269, "LOE4K", 4, func(x []int) int { return gepLOE4K(x[0], x[1], x[2], x[3]) }},
	// "GOE4K": IntNode{270, "GOE4K", 4, func(x []int) int { return gepGOE4K(x[0], x[1], x[2], x[3]) }},
	// "ET4K":  IntNode{271, "ET4K", 4, func(x []int) int { return gepET4K(x[0], x[1], x[2], x[3]) }},
	// "NET4K": IntNode{272, "NET4K", 4, func(x []int) int { return gepNET4K(x[0], x[1], x[2], x[3]) }},
	// "LT4L":  IntNode{273, "LT4L", 4, func(x []int) int { return gepLT4L(x[0], x[1], x[2], x[3]) }},
	// "GT4L":  IntNode{274, "GT4L", 4, func(x []int) int { return gepGT4L(x[0], x[1], x[2], x[3]) }},
	// "LOE4L": IntNode{275, "LOE4L", 4, func(x []int) int { return gepLOE4L(x[0], x[1], x[2], x[3]) }},
	// "GOE4L": IntNode{276, "GOE4L", 4, func(x []int) int { return gepGOE4L(x[0], x[1], x[2], x[3]) }},
	// "ET4L":  IntNode{277, "ET4L", 4, func(x []int) int { return gepET4L(x[0], x[1], x[2], x[3]) }},
	// "NET4L": IntNode{278, "NET4L", 4, func(x []int) int { return gepNET4L(x[0], x[1], x[2], x[3]) }},
}

/*
func gep3Rt(x int) int {
	if x < 0.0 {
		return -math.Pow(-x, (1.0 / 3.0))
	}
	return math.Pow(x, (1.0 / 3.0))
}

func gep5Rt(x int) int {
	if x < 0.0 {
		return -math.Pow(-x, (1.0 / 5.0))
	}
	return math.Pow(x, (1.0 / 5.0))
}

func gepLog2(x, y int) int {
	if y == 0.0 {
		return 0.0
	}
	return math.Log(x) / math.Log(y)
}

func gepMod(x, y int) int {
	// The built-in function is incorrect for cases such as -1.0 and 0.2.
	return ((x / y) - int(int(x/y))) * y
}

func gepLogi(x int) int {
	if math.Abs(x) > 709.0 {
		return 1.0 / (1.0 + math.Exp(math.Abs(x)/x*709.0))
	}
	return 1.0 / (1.0 + math.Exp(-x))
}

func gepLogi2(x, y int) int {
	if math.Abs(x+y) > 709.0 {
		return 1.0 / (1.0 + math.Exp(math.Abs(x+y)/(x+y)*709.0))
	}
	return 1.0 / (1.0 + math.Exp(-(x + y)))
}

func gepLogi3(x, y, z int) int {
	if math.Abs(x+y+z) > 709.0 {
		return 1.0 / (1.0 + math.Exp(math.Abs(x+y+z)/(x+y+z)*709.0))
	}
	return 1.0 / (1.0 + math.Exp(-(x + y + z)))
}

func gepLogi4(a, b, c, d int) int {
	if math.Abs(a+b+c+d) > 709.0 {
		return 1.0 / (1.0 + math.Exp(math.Abs(a+b+c+d)/(a+b+c+d)*709.0))
	}
	return 1.0 / (1.0 + math.Exp(-(a + b + c + d)))
}

func gepGau(x int) int {
	return math.Exp(-math.Pow(x, 2.0))
}

func gepGau2(x, y int) int {
	return math.Exp(-math.Pow((x + y), 2.0))
}

func gepGau3(x, y, z int) int {
	return math.Exp(-math.Pow((x + y + z), 2.0))
}

func gepGau4(a, b, c, d int) int {
	return math.Exp(-math.Pow((a + b + c + d), 2.0))
}

func gepAcsc(x int) int {
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

func gepAsec(x int) int {
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

func gepAcot(x int) int {
	return math.Atan(1.0 / x)
}

func gepAsinh(x int) int {
	return math.Log(x + math.Sqrt(x*x+1.0))
}

func gepAcosh(x int) int {
	return math.Log(x + math.Sqrt(x*x-1.0))
}

func gepAtanh(x int) int {
	return math.Log((1.0+x)/(1.0-x)) / 2.0
}

func gepAcsch(x int) int {
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

func gepAsech(x int) int {
	return math.Log((math.Sqrt(-x*x+1.0) + 1.0) / x)
}

func gepAcoth(x int) int {
	return math.Log((x+1.0)/(x-1.0)) / 2.0
}
*/

func gepMin2(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func gepMin3(x, y, z int) int {
	varTemp := x
	if varTemp > y {
		varTemp = y
	}
	if varTemp > z {
		varTemp = z
	}
	return varTemp
}

func gepMin4(a, b, c, d int) int {
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

func gepMax2(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func gepMax3(x, y, z int) int {
	varTemp := x
	if varTemp < y {
		varTemp = y
	}
	if varTemp < z {
		varTemp = z
	}
	return varTemp
}

func gepMax4(a, b, c, d int) int {
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

/*
func gepOR1(x, y int) int {
	if (x < 0.0) || (y < 0.0) {
		return 1.0
	}
	return 0.0
}

func gepOR2(x, y int) int {
	if (x >= 0.0) || (y >= 0.0) {
		return 1.0
	}
	return 0.0
}

func gepOR3(x, y int) int {
	if (x <= 0.0) || (y <= 0.0) {
		return 1.0
	}
	return 0.0
}

func gepOR4(x, y int) int {
	if (x < 1.0) || (y < 1.0) {
		return 1.0
	}
	return 0.0
}

func gepOR5(x, y int) int {
	if (x >= 1.0) || (y >= 1.0) {
		return 1.0
	}
	return 0.0
}

func gepOR6(x, y int) int {
	if (x <= 1.0) || (y <= 1.0) {
		return 1.0
	}
	return 0.0
}

func gepAND1(x, y int) int {
	if (x < 0.0) && (y < 0.0) {
		return 1.0
	}
	return 0.0
}

func gepAND2(x, y int) int {
	if (x >= 0.0) && (y >= 0.0) {
		return 1.0
	}
	return 0.0
}

func gepAND3(x, y int) int {
	if (x <= 0.0) && (y <= 0.0) {
		return 1.0
	}
	return 0.0
}

func gepAND4(x, y int) int {
	if (x < 1.0) && (y < 1.0) {
		return 1.0
	}
	return 0.0
}

func gepAND5(x, y int) int {
	if (x >= 1.0) && (y >= 1.0) {
		return 1.0
	}
	return 0.0
}

func gepAND6(x, y int) int {
	if (x <= 1.0) && (y <= 1.0) {
		return 1.0
	}
	return 0.0
}
*/

func gepLT2A(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func gepGT2A(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func gepLOE2A(x, y int) int {
	if x <= y {
		return x
	}
	return y
}

func gepGOE2A(x, y int) int {
	if x >= y {
		return x
	}
	return y
}

func gepET2A(x, y int) int {
	if x == y {
		return x
	}
	return y
}

func gepNET2A(x, y int) int {
	if x != y {
		return x
	}
	return y
}

func gepLT2B(x, y int) int {
	if x < y {
		return 1.0
	}
	return 0.0
}

func gepGT2B(x, y int) int {
	if x > y {
		return 1.0
	}
	return 0.0
}

func gepLOE2B(x, y int) int {
	if x <= y {
		return 1.0
	}
	return 0.0
}

func gepGOE2B(x, y int) int {
	if x >= y {
		return 1.0
	}
	return 0.0
}

func gepET2B(x, y int) int {
	if x == y {
		return 1.0
	}
	return 0.0
}

func gepNET2B(x, y int) int {
	if x != y {
		return 1.0
	}
	return 0.0
}

func gepLT2C(x, y int) int {
	if x < y {
		return (x + y)
	}
	return (x - y)
}

func gepGT2C(x, y int) int {
	if x > y {
		return (x + y)
	}
	return (x - y)
}

func gepLOE2C(x, y int) int {
	if x <= y {
		return (x + y)
	}
	return (x - y)
}

func gepGOE2C(x, y int) int {
	if x >= y {
		return (x + y)
	}
	return (x - y)
}

func gepET2C(x, y int) int {
	if x == y {
		return (x + y)
	}
	return (x - y)
}

func gepNET2C(x, y int) int {
	if x != y {
		return (x + y)
	}
	return (x - y)
}

/*
func gepLT2D(x, y int) int {
	if x < y {
		return (x * y)
	}
	return (x / y)
}

func gepGT2D(x, y int) int {
	if x > y {
		return (x * y)
	}
	return (x / y)
}

func gepLOE2D(x, y int) int {
	if x <= y {
		return (x * y)
	}
	return (x / y)
}

func gepGOE2D(x, y int) int {
	if x >= y {
		return (x * y)
	}
	return (x / y)
}

func gepET2D(x, y int) int {
	if x == y {
		return (x * y)
	}
	return (x / y)
}

func gepNET2D(x, y int) int {
	if x != y {
		return (x * y)
	}
	return (x / y)
}
*/

func gepLT2E(x, y int) int {
	if x < y {
		return (x + y)
	}
	return (x * y)
}

func gepGT2E(x, y int) int {
	if x > y {
		return (x + y)
	}
	return (x * y)
}

func gepLOE2E(x, y int) int {
	if x <= y {
		return (x + y)
	}
	return (x * y)
}

func gepGOE2E(x, y int) int {
	if x >= y {
		return (x + y)
	}
	return (x * y)
}

func gepET2E(x, y int) int {
	if x == y {
		return (x + y)
	}
	return (x * y)
}

func gepNET2E(x, y int) int {
	if x != y {
		return (x + y)
	}
	return (x * y)
}

/*
func gepLT2F(x, y int) int {
	if x < y {
		return (x + y)
	}
	return math.Sin(x * y)
}

func gepGT2F(x, y int) int {
	if x > y {
		return (x + y)
	}
	return math.Sin(x * y)
}

func gepLOE2F(x, y int) int {
	if x <= y {
		return (x + y)
	}
	return math.Sin(x * y)
}

func gepGOE2F(x, y int) int {
	if x >= y {
		return (x + y)
	}
	return math.Sin(x * y)
}

func gepET2F(x, y int) int {
	if x == y {
		return (x + y)
	}
	return math.Sin(x * y)
}

func gepNET2F(x, y int) int {
	if x != y {
		return (x + y)
	}
	return math.Sin(x * y)
}

func gepLT2G(x, y int) int {
	if x < y {
		return (x + y)
	}
	return math.Atan(x * y)
}

func gepGT2G(x, y int) int {
	if x > y {
		return (x + y)
	}
	return math.Atan(x * y)
}

func gepLOE2G(x, y int) int {
	if x <= y {
		return (x + y)
	}
	return math.Atan(x * y)
}

func gepGOE2G(x, y int) int {
	if x >= y {
		return (x + y)
	}
	return math.Atan(x * y)
}

func gepET2G(x, y int) int {
	if x == y {
		return (x + y)
	}
	return math.Atan(x * y)
}

func gepNET2G(x, y int) int {
	if x != y {
		return (x + y)
	}
	return math.Atan(x * y)
}
*/

func gepLT3A(x, y, z int) int {
	if x < 0.0 {
		return y
	}
	return z
}

func gepGT3A(x, y, z int) int {
	if x > 0.0 {
		return y
	}
	return z
}

func gepLOE3A(x, y, z int) int {
	if x <= 0.0 {
		return y
	}
	return z
}

func gepGOE3A(x, y, z int) int {
	if x >= 0.0 {
		return y
	}
	return z
}

func gepET3A(x, y, z int) int {
	if x == 0.0 {
		return y
	}
	return z
}

func gepNET3A(x, y, z int) int {
	if x != 0.0 {
		return y
	}
	return z
}

func gepLT3B(x, y, z int) int {
	if (x + y) < z {
		return (x + y)
	}
	return z
}

func gepGT3B(x, y, z int) int {
	if (x + y) > z {
		return (x + y)
	}
	return z
}

func gepLOE3B(x, y, z int) int {
	if (x + y) <= z {
		return (x + y)
	}
	return z
}

func gepGOE3B(x, y, z int) int {
	if (x + y) >= z {
		return (x + y)
	}
	return z
}

func gepET3B(x, y, z int) int {
	if (x + y) == z {
		return (x + y)
	}
	return z
}

func gepNET3B(x, y, z int) int {
	if (x + y) != z {
		return (x + y)
	}
	return z
}

func gepLT3C(x, y, z int) int {
	if (x + y) < z {
		return (x + y)
	}
	return (x + z)
}

func gepGT3C(x, y, z int) int {
	if (x + y) > z {
		return (x + y)
	}
	return (x + z)
}

func gepLOE3C(x, y, z int) int {
	if (x + y) <= z {
		return (x + y)
	}
	return (x + z)
}

func gepGOE3C(x, y, z int) int {
	if (x + y) >= z {
		return (x + y)
	}
	return (x + z)
}

func gepET3C(x, y, z int) int {
	if (x + y) == z {
		return (x + y)
	}
	return (x + z)
}

func gepNET3C(x, y, z int) int {
	if (x + y) != z {
		return (x + y)
	}
	return (x + z)
}

func gepLT3D(x, y, z int) int {
	if (x + y) < z {
		return (x + y)
	}
	return (x - z)
}

func gepGT3D(x, y, z int) int {
	if (x + y) > z {
		return (x + y)
	}
	return (x - z)
}

func gepLOE3D(x, y, z int) int {
	if (x + y) <= z {
		return (x + y)
	}
	return (x - z)
}

func gepGOE3D(x, y, z int) int {
	if (x + y) >= z {
		return (x + y)
	}
	return (x - z)
}

func gepET3D(x, y, z int) int {
	if (x + y) == z {
		return (x + y)
	}
	return (x - z)
}

func gepNET3D(x, y, z int) int {
	if (x + y) != z {
		return (x + y)
	}
	return (x - z)
}

func gepLT3E(x, y, z int) int {
	if (x + y) < z {
		return (x + y)
	}
	return (x * z)
}

func gepGT3E(x, y, z int) int {
	if (x + y) > z {
		return (x + y)
	}
	return (x * z)
}

func gepLOE3E(x, y, z int) int {
	if (x + y) <= z {
		return (x + y)
	}
	return (x * z)
}

func gepGOE3E(x, y, z int) int {
	if (x + y) >= z {
		return (x + y)
	}
	return (x * z)
}

func gepET3E(x, y, z int) int {
	if (x + y) == z {
		return (x + y)
	}
	return (x * z)
}

func gepNET3E(x, y, z int) int {
	if (x + y) != z {
		return (x + y)
	}
	return (x * z)
}

/*
func gepLT3F(x, y, z int) int {
	if (x + y) < z {
		return (x + y)
	}
	return (x / z)
}

func gepGT3F(x, y, z int) int {
	if (x + y) > z {
		return (x + y)
	}
	return (x / z)
}

func gepLOE3F(x, y, z int) int {
	if (x + y) <= z {
		return (x + y)
	}
	return (x / z)
}

func gepGOE3F(x, y, z int) int {
	if (x + y) >= z {
		return (x + y)
	}
	return (x / z)
}

func gepET3F(x, y, z int) int {
	if (x + y) == z {
		return (x + y)
	}
	return (x / z)
}

func gepNET3F(x, y, z int) int {
	if (x + y) != z {
		return (x + y)
	}
	return (x / z)
}
*/

func gepLT3G(x, y, z int) int {
	if (x + y) < z {
		return (x * y)
	}
	return (x + z)
}

func gepGT3G(x, y, z int) int {
	if (x + y) > z {
		return (x * y)
	}
	return (x + z)
}

func gepLOE3G(x, y, z int) int {
	if (x + y) <= z {
		return (x * y)
	}
	return (x + z)
}

func gepGOE3G(x, y, z int) int {
	if (x + y) >= z {
		return (x * y)
	}
	return (x + z)
}

func gepET3G(x, y, z int) int {
	if (x + y) == z {
		return (x * y)
	}
	return (x + z)
}

func gepNET3G(x, y, z int) int {
	if (x + y) != z {
		return (x * y)
	}
	return (x + z)
}

func gepLT3H(x, y, z int) int {
	if (x + y) < z {
		return (x * y)
	}
	return (x - z)
}

func gepGT3H(x, y, z int) int {
	if (x + y) > z {
		return (x * y)
	}
	return (x - z)
}

func gepLOE3H(x, y, z int) int {
	if (x + y) <= z {
		return (x * y)
	}
	return (x - z)
}

func gepGOE3H(x, y, z int) int {
	if (x + y) >= z {
		return (x * y)
	}
	return (x - z)
}

func gepET3H(x, y, z int) int {
	if (x + y) == z {
		return (x * y)
	}
	return (x - z)
}

func gepNET3H(x, y, z int) int {
	if (x + y) != z {
		return (x * y)
	}
	return (x - z)
}

func gepLT3I(x, y, z int) int {
	if (x + y) < z {
		return (x * y)
	}
	return (x * z)
}

func gepGT3I(x, y, z int) int {
	if (x + y) > z {
		return (x * y)
	}
	return (x * z)
}

func gepLOE3I(x, y, z int) int {
	if (x + y) <= z {
		return (x * y)
	}
	return (x * z)
}

func gepGOE3I(x, y, z int) int {
	if (x + y) >= z {
		return (x * y)
	}
	return (x * z)
}

func gepET3I(x, y, z int) int {
	if (x + y) == z {
		return (x * y)
	}
	return (x * z)
}

func gepNET3I(x, y, z int) int {
	if (x + y) != z {
		return (x * y)
	}
	return (x * z)
}

/*
func gepLT3J(x, y, z int) int {
	if (x + y) < z {
		return (x * y)
	}
	return (x / z)
}

func gepGT3J(x, y, z int) int {
	if (x + y) > z {
		return (x * y)
	}
	return (x / z)
}

func gepLOE3J(x, y, z int) int {
	if (x + y) <= z {
		return (x * y)
	}
	return (x / z)
}

func gepGOE3J(x, y, z int) int {
	if (x + y) >= z {
		return (x * y)
	}
	return (x / z)
}

func gepET3J(x, y, z int) int {
	if (x + y) == z {
		return (x * y)
	}
	return (x / z)
}

func gepNET3J(x, y, z int) int {
	if (x + y) != z {
		return (x * y)
	}
	return (x / z)
}

func gepLT3K(x, y, z int) int {
	if (x + y) < z {
		return (x + y + z)
	}
	return math.Sin(x * y * z)
}

func gepGT3K(x, y, z int) int {
	if (x + y) > z {
		return (x + y + z)
	}
	return math.Sin(x * y * z)
}

func gepLOE3K(x, y, z int) int {
	if (x + y) <= z {
		return (x + y + z)
	}
	return math.Sin(x * y * z)
}

func gepGOE3K(x, y, z int) int {
	if (x + y) >= z {
		return (x + y + z)
	}
	return math.Sin(x * y * z)
}

func gepET3K(x, y, z int) int {
	if (x + y) == z {
		return (x + y + z)
	}
	return math.Sin(x * y * z)
}

func gepNET3K(x, y, z int) int {
	if (x + y) != z {
		return (x + y + z)
	}
	return math.Sin(x * y * z)
}

func gepLT3L(x, y, z int) int {
	if (x + y) < z {
		return (x + y + z)
	}
	return math.Atan(x * y * z)
}

func gepGT3L(x, y, z int) int {
	if (x + y) > z {
		return (x + y + z)
	}
	return math.Atan(x * y * z)
}

func gepLOE3L(x, y, z int) int {
	if (x + y) <= z {
		return (x + y + z)
	}
	return math.Atan(x * y * z)
}

func gepGOE3L(x, y, z int) int {
	if (x + y) >= z {
		return (x + y + z)
	}
	return math.Atan(x * y * z)
}

func gepET3L(x, y, z int) int {
	if (x + y) == z {
		return (x + y + z)
	}
	return math.Atan(x * y * z)
}

func gepNET3L(x, y, z int) int {
	if (x + y) != z {
		return (x + y + z)
	}
	return math.Atan(x * y * z)
}
*/

func gepLT4A(a, b, c, d int) int {
	if a < b {
		return c
	}
	return d
}

func gepGT4A(a, b, c, d int) int {
	if a > b {
		return c
	}
	return d
}

func gepLOE4A(a, b, c, d int) int {
	if a <= b {
		return c
	}
	return d
}

func gepGOE4A(a, b, c, d int) int {
	if a >= b {
		return c
	}
	return d
}

func gepET4A(a, b, c, d int) int {
	if a == b {
		return c
	}
	return d
}

func gepNET4A(a, b, c, d int) int {
	if a != b {
		return c
	}
	return d
}

func gepLT4B(a, b, c, d int) int {
	if (a + b) < (c + d) {
		return c
	}
	return d
}

func gepGT4B(a, b, c, d int) int {
	if (a + b) > (c + d) {
		return c
	}
	return d
}

func gepLOE4B(a, b, c, d int) int {
	if (a + b) <= (c + d) {
		return c
	}
	return d
}

func gepGOE4B(a, b, c, d int) int {
	if (a + b) >= (c + d) {
		return c
	}
	return d
}

func gepET4B(a, b, c, d int) int {
	if (a + b) == (c + d) {
		return c
	}
	return d
}

func gepNET4B(a, b, c, d int) int {
	if (a + b) != (c + d) {
		return c
	}
	return d
}

func gepLT4C(a, b, c, d int) int {
	if (a + b) < (c + d) {
		return (a + b)
	}
	return (c + d)
}

func gepGT4C(a, b, c, d int) int {
	if (a + b) > (c + d) {
		return (a + b)
	}
	return (c + d)
}

func gepLOE4C(a, b, c, d int) int {
	if (a + b) <= (c + d) {
		return (a + b)
	}
	return (c + d)
}

func gepGOE4C(a, b, c, d int) int {
	if (a + b) >= (c + d) {
		return (a + b)
	}
	return (c + d)
}

func gepET4C(a, b, c, d int) int {
	if (a + b) == (c + d) {
		return (a + b)
	}
	return (c + d)
}

func gepNET4C(a, b, c, d int) int {
	if (a + b) != (c + d) {
		return (a + b)
	}
	return (c + d)
}

func gepLT4D(a, b, c, d int) int {
	if (a + b) < (c + d) {
		return (a + b)
	}
	return (c - d)
}

func gepGT4D(a, b, c, d int) int {
	if (a + b) > (c + d) {
		return (a + b)
	}
	return (c - d)
}

func gepLOE4D(a, b, c, d int) int {
	if (a + b) <= (c + d) {
		return (a + b)
	}
	return (c - d)
}

func gepGOE4D(a, b, c, d int) int {
	if (a + b) >= (c + d) {
		return (a + b)
	}
	return (c - d)
}

func gepET4D(a, b, c, d int) int {
	if (a + b) == (c + d) {
		return (a + b)
	}
	return (c - d)
}

func gepNET4D(a, b, c, d int) int {
	if (a + b) != (c + d) {
		return (a + b)
	}
	return (c - d)
}

func gepLT4E(a, b, c, d int) int {
	if (a + b) < (c + d) {
		return (a + b)
	}
	return (c * d)
}

func gepGT4E(a, b, c, d int) int {
	if (a + b) > (c + d) {
		return (a + b)
	}
	return (c * d)
}

func gepLOE4E(a, b, c, d int) int {
	if (a + b) <= (c + d) {
		return (a + b)
	}
	return (c * d)
}

func gepGOE4E(a, b, c, d int) int {
	if (a + b) >= (c + d) {
		return (a + b)
	}
	return (c * d)
}

func gepET4E(a, b, c, d int) int {
	if (a + b) == (c + d) {
		return (a + b)
	}
	return (c * d)
}

func gepNET4E(a, b, c, d int) int {
	if (a + b) != (c + d) {
		return (a + b)
	}
	return (c * d)
}

/*
func gepLT4F(a, b, c, d int) int {
	if (a + b) < (c + d) {
		return (a + b)
	}
	return (c / d)
}

func gepGT4F(a, b, c, d int) int {
	if (a + b) > (c + d) {
		return (a + b)
	}
	return (c / d)
}

func gepLOE4F(a, b, c, d int) int {
	if (a + b) <= (c + d) {
		return (a + b)
	}
	return (c / d)
}

func gepGOE4F(a, b, c, d int) int {
	if (a + b) >= (c + d) {
		return (a + b)
	}
	return (c / d)
}

func gepET4F(a, b, c, d int) int {
	if (a + b) == (c + d) {
		return (a + b)
	}
	return (c / d)
}

func gepNET4F(a, b, c, d int) int {
	if (a + b) != (c + d) {
		return (a + b)
	}
	return (c / d)
}
*/

func gepLT4G(a, b, c, d int) int {
	if (a + b) < (c + d) {
		return (a * b)
	}
	return (c + d)
}

func gepGT4G(a, b, c, d int) int {
	if (a + b) > (c + d) {
		return (a * b)
	}
	return (c + d)
}

func gepLOE4G(a, b, c, d int) int {
	if (a + b) <= (c + d) {
		return (a * b)
	}
	return (c + d)
}

func gepGOE4G(a, b, c, d int) int {
	if (a + b) >= (c + d) {
		return (a * b)
	}
	return (c + d)
}

func gepET4G(a, b, c, d int) int {
	if (a + b) == (c + d) {
		return (a * b)
	}
	return (c + d)
}

func gepNET4G(a, b, c, d int) int {
	if (a + b) != (c + d) {
		return (a * b)
	}
	return (c + d)
}

func gepLT4H(a, b, c, d int) int {
	if (a + b) < (c + d) {
		return (a * b)
	}
	return (c - d)
}

func gepGT4H(a, b, c, d int) int {
	if (a + b) > (c + d) {
		return (a * b)
	}
	return (c - d)
}

func gepLOE4H(a, b, c, d int) int {
	if (a + b) <= (c + d) {
		return (a * b)
	}
	return (c - d)
}

func gepGOE4H(a, b, c, d int) int {
	if (a + b) >= (c + d) {
		return (a * b)
	}
	return (c - d)
}

func gepET4H(a, b, c, d int) int {
	if (a + b) == (c + d) {
		return (a * b)
	}
	return (c - d)
}

func gepNET4H(a, b, c, d int) int {
	if (a + b) != (c + d) {
		return (a * b)
	}
	return (c - d)
}

func gepLT4I(a, b, c, d int) int {
	if (a + b) < (c + d) {
		return (a * b)
	}
	return (c * d)
}

func gepGT4I(a, b, c, d int) int {
	if (a + b) > (c + d) {
		return (a * b)
	}
	return (c * d)
}

func gepLOE4I(a, b, c, d int) int {
	if (a + b) <= (c + d) {
		return (a * b)
	}
	return (c * d)
}

func gepGOE4I(a, b, c, d int) int {
	if (a + b) >= (c + d) {
		return (a * b)
	}
	return (c * d)
}

func gepET4I(a, b, c, d int) int {
	if (a + b) == (c + d) {
		return (a * b)
	}
	return (c * d)
}

func gepNET4I(a, b, c, d int) int {
	if (a + b) != (c + d) {
		return (a * b)
	}
	return (c * d)
}

/*
func gepLT4J(a, b, c, d int) int {
	if (a + b) < (c + d) {
		return (a * b)
	}
	return (c / d)
}

func gepGT4J(a, b, c, d int) int {
	if (a + b) > (c + d) {
		return (a * b)
	}
	return (c / d)
}

func gepLOE4J(a, b, c, d int) int {
	if (a + b) <= (c + d) {
		return (a * b)
	}
	return (c / d)
}

func gepGOE4J(a, b, c, d int) int {
	if (a + b) >= (c + d) {
		return (a * b)
	}
	return (c / d)
}

func gepET4J(a, b, c, d int) int {
	if (a + b) == (c + d) {
		return (a * b)
	}
	return (c / d)
}

func gepNET4J(a, b, c, d int) int {
	if (a + b) != (c + d) {
		return (a * b)
	}
	return (c / d)
}

func gepLT4K(a, b, c, d int) int {
	if (a + b) < (c + d) {
		return math.Sin(a * b)
	}
	return math.Sin(c * d)
}

func gepGT4K(a, b, c, d int) int {
	if (a + b) > (c + d) {
		return math.Sin(a * b)
	}
	return math.Sin(c * d)
}

func gepLOE4K(a, b, c, d int) int {
	if (a + b) <= (c + d) {
		return math.Sin(a * b)
	}
	return math.Sin(c * d)
}

func gepGOE4K(a, b, c, d int) int {
	if (a + b) >= (c + d) {
		return math.Sin(a * b)
	}
	return math.Sin(c * d)
}

func gepET4K(a, b, c, d int) int {
	if (a + b) == (c + d) {
		return math.Sin(a * b)
	}
	return math.Sin(c * d)
}

func gepNET4K(a, b, c, d int) int {
	if (a + b) != (c + d) {
		return math.Sin(a * b)
	}
	return math.Sin(c * d)
}

func gepLT4L(a, b, c, d int) int {
	if (a + b) < (c + d) {
		return math.Atan(a * b)
	}
	return math.Atan(c * d)
}

func gepGT4L(a, b, c, d int) int {
	if (a + b) > (c + d) {
		return math.Atan(a * b)
	}
	return math.Atan(c * d)
}

func gepLOE4L(a, b, c, d int) int {
	if (a + b) <= (c + d) {
		return math.Atan(a * b)
	}
	return math.Atan(c * d)
}

func gepGOE4L(a, b, c, d int) int {
	if (a + b) >= (c + d) {
		return math.Atan(a * b)
	}
	return math.Atan(c * d)
}

func gepET4L(a, b, c, d int) int {
	if (a + b) == (c + d) {
		return math.Atan(a * b)
	}
	return math.Atan(c * d)
}

func gepNET4L(a, b, c, d int) int {
	if (a + b) != (c + d) {
		return math.Atan(a * b)
	}
	return math.Atan(c * d)
}
*/
