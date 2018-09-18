// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package functions provides a map of available functions for the GEP algorithm.
package functions

// VectorInt is a vector of integers.
type VectorInt []int

// FuncMap is a map from the symbol name of a function to its defining FuncNode.
type FuncMap map[string]FuncNode

// FuncNode defines an available function for the GEP algorithm.
type FuncNode interface {
	// Symbol is the Karva string representation of the function.
	Symbol() string
	// Terminals is the number of input terminals for the function.
	Terminals() int

	// BoolFunction represents a general boolean function.
	BoolFunction([]bool) bool
	// IntFunction represents a general integer function.
	IntFunction([]int) int
	// Float64Function represents a general floating-point function.
	Float64Function([]float64) float64
	// VectorIntFunction represents a general vector of integers function.
	VectorIntFunction([]VectorInt) VectorInt
}

// FuncType represents the underlying function type of the model, genome, and gene.
// This is needed in place of generics.
type FuncType int

const (
	unknown FuncType = iota
	Bool
	Int
	Float64
	VectorInts
)
