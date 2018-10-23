// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

//go:generate go run gen-benchmarks.go

// Package vectorIntNodes defines the vector of integers function collection available for the GEP algorithm.
package vectorIntNodes

import (
	"log"

	"github.com/gmlewis/gep/v2/functions"
)

type VectorInt = functions.VectorInt

// VectorIntNode is a vector of integers function used for the formation of GEP expressions.
type VectorIntNode struct {
	index     int
	symbol    string
	terminals int
	function  func(x []VectorInt) VectorInt
}

// Symbol returns the Karva symbol for this vector of integers function.
func (n VectorIntNode) Symbol() string {
	return n.symbol
}

// Terminals returns the number of input terminals for this vector of integers function.
func (n VectorIntNode) Terminals() int {
	return n.terminals
}

// BoolFunction is unused in this package and returns an error.
func (n VectorIntNode) BoolFunction([]bool) bool {
	log.Println("error calling BoolFunction on VectorIntNode model.")
	return false
}

// IntFunction calls the vector of integers function and returns the result.
func (n VectorIntNode) IntFunction([]int) int {
	return 0
}

// Float64Function is unused in this package and returns an error.
func (n VectorIntNode) Float64Function([]float64) float64 {
	log.Println("error calling Float64Function on VectorIntNode model.")
	return 0.0
}

// VectorIntFunction allows FuncMap to implement interace functions.FuncMap.
func (n VectorIntNode) VectorIntFunction(x []VectorInt) VectorInt {
	return n.function(x)
}

// sliceOp is a function used on each index of a VectorInt.
type sliceOp func(in []int) int

// ProcessVector is a helper function that processes each index of an array of VectorInts.
func ProcessVector(x []VectorInt, op sliceOp) (result VectorInt) {
	if len(x) == 0 {
		return result
	}
	result = make([]int, len(x[0]))
	for i := 0; i < len(x[0]); i++ {
		args := make([]int, len(x))
		for j := 0; j < len(x); j++ {
			args[j] = x[j][i]
		}
		result[i] = op(args)
	}
	return result
}

// Int lists all the available vector of integers functions for this package.
var VectorIntFuncs = functions.FuncMap{
	// TODO(gmlewis): Change functions to operate on variable-length slices.
	"+": VectorIntNode{0, "+", 2, func(x []VectorInt) VectorInt {
		op := func(in []int) (result int) {
			for i := 0; i < 2; /* len(in) */ i++ {
				result += in[i]
			}
			return result
		}
		return ProcessVector(x, op)
	}},
	"-": VectorIntNode{1, "-", 2, func(x []VectorInt) VectorInt {
		op := func(in []int) (result int) {
			result = in[0]
			for i := 1; i < 2; /* len(in) */ i++ {
				result -= in[i]
			}
			return result
		}
		return ProcessVector(x, op)
	}},
	"*": VectorIntNode{2, "*", 2, func(x []VectorInt) VectorInt {
		op := func(in []int) (result int) {
			result = in[0]
			for i := 1; i < 2; /* len(in) */ i++ {
				result *= in[i]
			}
			return result
		}
		return ProcessVector(x, op)
	}},
	"/": VectorIntNode{3, "/", 2, func(x []VectorInt) VectorInt {
		op := func(in []int) (result int) {
			result = in[0]
			for i := 1; i < 2; /* len(in) */ i++ {
				if in[0] == 0 {
					result = 0
				} else {
					result /= in[i]
				}
			}
			return result
		}
		return ProcessVector(x, op)
	}},
}
