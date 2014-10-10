// NAND is a simple experiment to run the GEP algorithm using the Boolean logic package.
// Given a set of input functions (Not, And, and Or), this solves how to create a NAND gate
// from those basic building blocks. This experiment usually converges to a solution within
// the first generation of evolution.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	bn "github.com/gmlewis/gep/functions/bool_nodes"
	"github.com/gmlewis/gep/gene"
	"github.com/gmlewis/gep/genome"
	"github.com/gmlewis/gep/grammars"
	"github.com/gmlewis/gep/model"
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

func validateNand(g *genome.Genome) float64 {
	correct := 0
	for _, n := range nandTests {
		r := g.EvalBool(n.in, bn.BoolAllGates)
		if r == n.out {
			correct++
		}
	}
	return 1000.0 * float64(correct) / float64(len(nandTests))
}

func main() {
	funcs := []gene.FuncWeight{
		{"Not", 1},
		{"And", 5},
		{"Or", 5},
	}
	numIn := len(nandTests[0].in)
	e := model.New(funcs, bn.BoolAllGates, 30, 7, 1, numIn, "Or", validateNand)
	s := e.Evolve(1000)

	// Write out the Go source code for the solution.
	gr, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		log.Printf("unable to load Boolean grammar: %v\n", err)
	}
	fmt.Printf("\n// gepModel is auto-generated Go source code for the\n")
	fmt.Printf("// nand solution karva expression:\n// %q, score=%v\n", s, validateNand(s))
	s.Write(os.Stdout, gr)
}
