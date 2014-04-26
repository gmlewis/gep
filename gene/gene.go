// Package gene provides the basis for a single gene in GEP.
package gene

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gmlewis/gep/functions"
	mn "github.com/gmlewis/gep/functions/math_nodes"
)

// FuncWeight contains the symbol name and its weight to be used during
// a run of the GEP algorithm. A symbol with weight 5, for example, will
// be five times more likely to be used than a symbol with weight 1.
type FuncWeight struct {
	Symbol string
	Weight int
}

// Gene contains all the information needed to represent a single gene
// in a GEP expression.
type Gene struct {
	// Symbols is the slice of strings being used in this gene's expression.
	Symbols []string
	// Constants is the slice of floats available for use by this gene.
	Constants    []float64
	bf           func([]bool) bool
	mf           func([]float64) float64
	headSize     int
	choiceSlice  []string
	numTerminals int
}

// New creates a new gene based on the Karva string representation.
func New(x string) Gene {
	parts := strings.Split(x, ".")
	return Gene{Symbols: parts}
}

// RandomNew generates a new, random gene for further manipulation by the GEP
// algorithm. The headSize, tailSize, and numTerminals determine the respective
// properties of the gene, and functions provide the available functions and
// their respective weights to be used in the creation of the gene.
func RandomNew(headSize, tailSize, numTerminals int, functions []FuncWeight) Gene {
	totalWeight := numTerminals
	for _, f := range functions {
		totalWeight += f.Weight
	}
	choiceSlice := make([]string, 0, totalWeight)
	for i := 0; i < numTerminals; i++ {
		choiceSlice = append(choiceSlice, fmt.Sprintf("d%v", i))
	}
	for _, f := range functions {
		for i := 0; i < f.Weight; i++ {
			choiceSlice = append(choiceSlice, f.Symbol)
		}
	}
	choices := rand.Perm(totalWeight)
	r := Gene{
		Symbols:      make([]string, 0, headSize+tailSize),
		headSize:     headSize,
		choiceSlice:  choiceSlice,
		numTerminals: numTerminals,
	}
	for i := 0; i < headSize; i++ {
		choice := choices[i%len(choices)]
		r.Symbols = append(r.Symbols, choiceSlice[choice])
	}
	for i := 0; i < tailSize; i++ {
		choice := choices[i%len(choices)]
		r.Symbols = append(r.Symbols, choiceSlice[choice%numTerminals])
	}
	return r
}

// String returns the Karva representation of the gene.
func (g Gene) String() string {
	return strings.Join(g.Symbols, ".")
}

func (g *Gene) getBoolArgOrder(nodes functions.FuncMap) [][]int {
	argOrder := make([][]int, len(g.Symbols))
	argCount := 0
	for i := 0; i < len(g.Symbols); i++ {
		sym := g.Symbols[i]
		if s, ok := nodes[sym]; ok {
			if s.Terminals() > 0 {
				args := make([]int, s.Terminals())
				for j := 0; j < s.Terminals(); j++ {
					argCount++
					args[j] = argCount
				}
				argOrder[i] = args
			}
		}
	}
	return argOrder
}

// EvalBool evaluates the gene as a boolean expression and returns the result.
// in represents the boolean inputs available to the gene.
// nodes is the map of available boolean functions to the gene.
func (g *Gene) EvalBool(in []bool, nodes functions.FuncMap) bool {
	if g.bf == nil {
		argOrder := g.getBoolArgOrder(nodes)
		g.bf = g.buildBoolTree(0, argOrder, nodes)
	}
	return g.bf(in)
}

func (g *Gene) buildBoolTree(symbolIndex int, argOrder [][]int, nodes functions.FuncMap) func([]bool) bool {
	// log.Infof("buildBoolTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex > len(g.Symbols) {
		log.Printf("bad symbolIndex %v for symbols: %v\n", symbolIndex, g.Symbols)
		return func(a []bool) bool { return false }
	}
	sym := g.Symbols[symbolIndex]
	if s, ok := nodes[sym]; ok {
		switch s.Terminals() {
		case 0:
			return func(in []bool) bool {
				return s.BoolFunction(false, false, false, false)
			}
		case 1:
			args := argOrder[symbolIndex]
			f := g.buildBoolTree(args[0], argOrder, nodes)
			return func(in []bool) bool {
				return s.BoolFunction(f(in), false, false, false)
			}
		case 2:
			args := argOrder[symbolIndex]
			left := g.buildBoolTree(args[0], argOrder, nodes)
			right := g.buildBoolTree(args[1], argOrder, nodes)
			return func(in []bool) bool {
				return s.BoolFunction(left(in), right(in), false, false)
			}
		case 3:
			args := argOrder[symbolIndex]
			f1 := g.buildBoolTree(args[0], argOrder, nodes)
			f2 := g.buildBoolTree(args[1], argOrder, nodes)
			f3 := g.buildBoolTree(args[2], argOrder, nodes)
			return func(in []bool) bool {
				return s.BoolFunction(f1(in), f2(in), f3(in), false)
			}
		case 4:
			args := argOrder[symbolIndex]
			f1 := g.buildBoolTree(args[0], argOrder, nodes)
			f2 := g.buildBoolTree(args[1], argOrder, nodes)
			f3 := g.buildBoolTree(args[2], argOrder, nodes)
			f4 := g.buildBoolTree(args[3], argOrder, nodes)
			return func(in []bool) bool {
				return s.BoolFunction(f1(in), f2(in), f3(in), f4(in))
			}
		}
	} else { // No named symbol found - look for d0, d1, ...
		if sym[0:1] == "d" {
			if index, err := strconv.ParseInt(sym[1:], 10, 32); err != nil {
				log.Printf("unable to parse variable index: sym=%v\n", sym)
			} else {
				return func(in []bool) bool {
					if index >= int64(len(in)) {
						log.Printf("error evaluating gene %v: index %v >= d length (%v)\n", sym, index, len(in))
						return false
					}
					return in[index]
				}
			}
		}
	}
	log.Printf("unable to return function: sym=%v\n", sym)
	return func(in []bool) bool { return false }
}

func (g *Gene) getMathArgOrder() [][]int {
	argOrder := make([][]int, len(g.Symbols))
	argCount := 0
	for i := 0; i < len(g.Symbols); i++ {
		sym := g.Symbols[i]
		if s, ok := mn.Math[sym]; ok {
			if s.Terminals() > 0 {
				args := make([]int, s.Terminals())
				for j := 0; j < s.Terminals(); j++ {
					argCount++
					args[j] = argCount
				}
				argOrder[i] = args
			}
		}
	}
	return argOrder
}

// EvalMath evaluates the gene as a floating-point expression and returns the result.
// in represents the float64 inputs available to the gene.
func (g *Gene) EvalMath(in []float64) float64 {
	if g.mf == nil {
		argOrder := g.getMathArgOrder()
		g.mf = g.buildMathTree(0, argOrder)
	}
	return g.mf(in)
}

func (g *Gene) buildMathTree(symbolIndex int, argOrder [][]int) func([]float64) float64 {
	// log.Infof("buildMathTree(%v, %#v, ...)", symbolIndex, argOrder)
	if symbolIndex > len(g.Symbols) {
		log.Printf("bad symbolIndex %v for symbols: %v\n", symbolIndex, g.Symbols)
		return func(a []float64) float64 { return 0.0 }
	}
	sym := g.Symbols[symbolIndex]
	if s, ok := mn.Math[sym]; ok {
		switch s.Terminals() {
		case 0:
			return func(in []float64) float64 {
				return s.Float64Function(0.0, 0.0, 0.0, 0.0)
			}
		case 1:
			args := argOrder[symbolIndex]
			f := g.buildMathTree(args[0], argOrder)
			return func(in []float64) float64 {
				return s.Float64Function(f(in), 0.0, 0.0, 0.0)
			}
		case 2:
			args := argOrder[symbolIndex]
			left := g.buildMathTree(args[0], argOrder)
			right := g.buildMathTree(args[1], argOrder)
			return func(in []float64) float64 {
				return s.Float64Function(left(in), right(in), 0.0, 0.0)
			}
		case 3:
			args := argOrder[symbolIndex]
			f1 := g.buildMathTree(args[0], argOrder)
			f2 := g.buildMathTree(args[1], argOrder)
			f3 := g.buildMathTree(args[2], argOrder)
			return func(in []float64) float64 {
				return s.Float64Function(f1(in), f2(in), f3(in), 0.0)
			}
		case 4:
			args := argOrder[symbolIndex]
			f1 := g.buildMathTree(args[0], argOrder)
			f2 := g.buildMathTree(args[1], argOrder)
			f3 := g.buildMathTree(args[2], argOrder)
			f4 := g.buildMathTree(args[3], argOrder)
			return func(in []float64) float64 {
				return s.Float64Function(f1(in), f2(in), f3(in), f4(in))
			}
		}
	} else { // No named symbol found - look for d0, d1, ...
		if sym[0:1] == "d" {
			if index, err := strconv.ParseInt(sym[1:], 10, 32); err != nil {
				log.Printf("unable to parse variable index: sym=%v\n", sym)
			} else {
				return func(in []float64) float64 {
					if index >= int64(len(in)) {
						log.Printf("error evaluating gene %v: index %v >= d length (%v)\n", sym, index, len(in))
						return 0.0
					}
					return in[index]
				}
			}
		} else if sym[0:1] == "c" {
			if index, err := strconv.ParseInt(sym[1:], 10, 32); err != nil {
				log.Printf("unable to parse constant index: sym=%v\n", sym)
			} else {
				return func(in []float64) float64 {
					if index >= int64(len(g.Constants)) {
						log.Printf("error evaluating gene %v: index %v >= c length (%v)\n", sym, index, len(g.Constants))
						return 0.0
					}
					return g.Constants[index]
				}
			}
		}
	}
	log.Printf("unable to return function: sym=%v\n", sym)
	return func(in []float64) float64 { return 0.0 }
}

// Mutate mutates a gene by performing a single random symbol exchange within the gene.
func (g *Gene) Mutate() {
	position := rand.Intn(len(g.Symbols))
	if position < g.headSize {
		symbol := g.choiceSlice[rand.Intn(len(g.choiceSlice))]
		g.Symbols[position] = symbol
	} else { // Must choose strictly from terminals
		terminal := g.choiceSlice[rand.Intn(g.numTerminals)]
		g.Symbols[position] = terminal
	}
	// Invalidate the cached function
	g.bf = nil
	g.mf = nil
}
