// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package boolNodes defines the Boolean function collections available for the GEP algorithm.
package boolNodes

import (
	"log"
)

// BoolNode is a boolean function used for the formation of GEP expressions.
type BoolNode struct {
	index     int
	symbol    string
	terminals int
	function  func(x []bool) bool
}

// Symbol returns the Karva symbol for this boolean function.
func (n BoolNode) Symbol() string {
	return n.symbol
}

// Terminals returns the number of input terminals for this boolean function.
func (n BoolNode) Terminals() int {
	return n.terminals
}

// BoolFunction calls the boolean function and returns the result.
func (n BoolNode) BoolFunction(x []bool) bool {
	return n.function(x)
}

// IntFunction is unused in this package and returns an error.
func (n BoolNode) IntFunction([]int) int {
	log.Println("error calling IntFunction on BoolNode model.")
	return 0
}

// Float64Function is unused in this package and returns an error.
func (n BoolNode) Float64Function([]float64) float64 {
	log.Println("error calling Float64Function on BoolNode model.")
	return 0.0
}
