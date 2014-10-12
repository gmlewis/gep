// Package genome provides the basis for a single GEP genome.
// A genome consists of one or more genes.
package genome

import (
	"log"
	"math/rand"
	"strings"

	"github.com/gmlewis/gep/functions"
	mn "github.com/gmlewis/gep/functions/math_nodes"
	"github.com/gmlewis/gep/gene"
)

// Genome contains the genes that make up the genome.
// It also provides the linking function and the score that results
// from evaluating the genome against the fitness function.
type Genome struct {
	Genes    []*gene.Gene
	LinkFunc string
	Score    float64
}

// New creates a new genome from the given genes and linking function.
func New(genes []*gene.Gene, linkFunc string) *Genome {
	return &Genome{Genes: genes, LinkFunc: linkFunc}
}

// EvalBool evaluates the genome as a boolean expression and returns the result.
// in represents the boolean inputs available to the genome.
// fm is the map of available boolean functions to the genome.
func (g *Genome) EvalBool(in []bool, fm functions.FuncMap) bool {
	lf, ok := fm[g.LinkFunc]
	if !ok {
		log.Printf("Unable to find linking function: %v\n", g.LinkFunc)
		return false
	}
	result := g.Genes[0].EvalBool(in, fm)
	for i := 1; i < len(g.Genes); i++ {
		result = lf.BoolFunction(result, g.Genes[i].EvalBool(in, fm), false, false)
	}
	return result
}

// EvalMath evaluates the genome as a floating-point expression and returns the result.
// in represents the float64 inputs available to the genome.
func (g *Genome) EvalMath(in []float64) float64 {
	lf, ok := mn.Math[g.LinkFunc]
	if !ok {
		log.Printf("Unable to find linking function: %v\n", g.LinkFunc)
		return 0.0
	}
	result := g.Genes[0].EvalMath(in)
	for i := 1; i < len(g.Genes); i++ {
		result = lf.Float64Function(result, g.Genes[i].EvalMath(in), 0.0, 0.0)
	}
	return result
}

// String returns the Karva representation of the genome.
func (g Genome) String() string {
	result := []string{}
	for _, v := range g.Genes {
		result = append(result, v.String())
	}
	return strings.Join(result, "|"+g.LinkFunc+"|")
}

// Mutate mutates a genome by performing numMutations random symbol exchanges within the genome.
func (g *Genome) Mutate(numMutations int) {
	for i := 0; i < numMutations; i++ {
		n := rand.Intn(len(g.Genes))
		// fmt.Printf("\nMutating gene #%v, before:\n%v\n", n, g.Genes[n])
		g.Genes[n].Mutate()
		// fmt.Printf("after:\n%v\n", g.Genes[n])
	}
}

// Dup duplicates the genome into the provided destination genome.
func (g *Genome) Dup() *Genome {
	if g == nil {
		log.Printf("denome.Dup error: src and dst must be non-nil\n")
		return nil
	}
	dst := &Genome{
		Genes:    make([]*gene.Gene, len(g.Genes)),
		LinkFunc: g.LinkFunc,
		Score:    g.Score,
	}
	for i := range g.Genes {
		dst.Genes[i] = g.Genes[i].Dup()
	}
	return dst
}
