// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

//go:generate go run gen-benchmarks.go

// Package intNodes defines the integer function collection available for the GEP algorithm.
package intNodes

import (
	"log"
	"math"

	"github.com/gmlewis/gep/v2/functions"
)

// IntNode is an integer function used for the formation of GEP expressions.
type IntNode struct {
	index     int
	symbol    string
	terminals int
	function  func(x []int) int
}

// Symbol returns the Karva symbol for this integer function.
func (n IntNode) Symbol() string {
	return n.symbol
}

// Terminals returns the number of input terminals for this integer function.
func (n IntNode) Terminals() int {
	return n.terminals
}

// BoolFunction is unused in this package and returns an error.
func (n IntNode) BoolFunction([]bool) bool {
	log.Println("error calling BoolFunction on IntNode model.")
	return false
}

// IntFunction calls the integer function and returns the result.
func (n IntNode) IntFunction(x []int) int {
	return n.function(x)
}

// Float64Function is unused in this package and returns an error.
func (n IntNode) Float64Function([]float64) float64 {
	log.Println("error calling Float64Function on IntNode model.")
	return 0.0
}

// VectorIntFunction allows FuncMap to implement interace functions.FuncMap.
func (n IntNode) VectorIntFunction([]functions.VectorInt) functions.VectorInt {
	log.Println("error calling VectorIntFunction on IntNode model.")
	return functions.VectorInt{}
}

func safeDiv(a, b int) int {
	if b == 0 {
		return math.MaxInt
	}
	return a / b
}

// Int lists all the available integer functions for this package.
var Int = functions.FuncMap{
	// TODO(gmlewis): Change functions to operate on the entire length of the slice.
	"+":     IntNode{0, "+", 2, func(x []int) int { return (x[0] + x[1]) }},
	"-":     IntNode{1, "-", 2, func(x []int) int { return (x[0] - x[1]) }},
	"*":     IntNode{2, "*", 2, func(x []int) int { return (x[0] * x[1]) }},
	"/":     IntNode{3, "/", 2, func(x []int) int { return (safeDiv(x[0], x[1])) }},
	"Neg":   IntNode{17, "Neg", 1, func(x []int) int { return (-(x[0])) }},
	"Nop":   IntNode{16, "Nop", 1, func(x []int) int { return (x[0]) }},
	"Add3":  IntNode{84, "Add3", 3, func(x []int) int { return (x[0] + x[1] + x[2]) }},
	"Sub3":  IntNode{86, "Sub3", 3, func(x []int) int { return (x[0] - x[1] - x[2]) }},
	"Mul3":  IntNode{88, "Mul3", 3, func(x []int) int { return (x[0] * x[1] * x[2]) }},
	"Div3":  IntNode{90, "Div3", 3, func(x []int) int { return safeDiv(safeDiv(x[0], x[1]), x[2]) }},
	"Add4":  IntNode{85, "Add4", 4, func(x []int) int { return (x[0] + x[1] + x[2] + x[3]) }},
	"Sub4":  IntNode{87, "Sub4", 4, func(x []int) int { return (x[0] - x[1] - x[2] - x[3]) }},
	"Mul4":  IntNode{89, "Mul4", 4, func(x []int) int { return (x[0] * x[1] * x[2] * x[3]) }},
	"Div4":  IntNode{91, "Div4", 4, func(x []int) int { return safeDiv(safeDiv(safeDiv(x[0], x[1]), x[2]), x[3]) }},
	"Min2":  IntNode{92, "Min2", 2, func(x []int) int { return gepMin2(x[0], x[1]) }},
	"Min3":  IntNode{93, "Min3", 3, func(x []int) int { return gepMin3(x[0], x[1], x[2]) }},
	"Min4":  IntNode{94, "Min4", 4, func(x []int) int { return gepMin4(x[0], x[1], x[2], x[3]) }},
	"Max2":  IntNode{95, "Max2", 2, func(x []int) int { return gepMax2(x[0], x[1]) }},
	"Max3":  IntNode{96, "Max3", 3, func(x []int) int { return gepMax3(x[0], x[1], x[2]) }},
	"Max4":  IntNode{97, "Max4", 4, func(x []int) int { return gepMax4(x[0], x[1], x[2], x[3]) }},
	"Avg2":  IntNode{98, "Avg2", 2, func(x []int) int { return ((x[0] + x[1]) / 2) }},
	"Avg3":  IntNode{99, "Avg3", 3, func(x []int) int { return ((x[0] + x[1] + x[2]) / 3) }},
	"Avg4":  IntNode{100, "Avg4", 4, func(x []int) int { return ((x[0] + x[1] + x[2] + x[3]) / 4) }},
	"Zero":  IntNode{70, "Zero", 1, func(x []int) int { return 0 }},
	"One":   IntNode{71, "One", 1, func(x []int) int { return 1 }},
	"Zero2": IntNode{72, "Zero2", 2, func(x []int) int { return 0 }},
	"One2":  IntNode{73, "One2", 2, func(x []int) int { return 1 }},
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
	"LT2E":  IntNode{129, "LT2E", 2, func(x []int) int { return gepLT2E(x[0], x[1]) }},
	"GT2E":  IntNode{130, "GT2E", 2, func(x []int) int { return gepGT2E(x[0], x[1]) }},
	"LOE2E": IntNode{131, "LOE2E", 2, func(x []int) int { return gepLOE2E(x[0], x[1]) }},
	"GOE2E": IntNode{132, "GOE2E", 2, func(x []int) int { return gepGOE2E(x[0], x[1]) }},
	"ET2E":  IntNode{133, "ET2E", 2, func(x []int) int { return gepET2E(x[0], x[1]) }},
	"NET2E": IntNode{134, "NET2E", 2, func(x []int) int { return gepNET2E(x[0], x[1]) }},
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
}

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
