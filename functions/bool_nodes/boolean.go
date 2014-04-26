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
	function  func(x0, x1, x2, x3 bool) bool
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
func (n BoolNode) BoolFunction(a, b, c, d bool) bool {
	return n.function(a, b, c, d)
}

// Float64Function is unused in this package and returns an error.
func (n BoolNode) Float64Function(a, b, c, d float64) float64 {
	log.Println("error calling Float64Function on BoolNode model.")
	return 0.0
}
