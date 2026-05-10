// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/gmlewis/gep/v2/grammars"
)

// ProgramFromSymbols creates a Program from pre-extracted per-gene Karva
// symbol-name slices and optional per-gene float64 constant slices. Pass nil
// for constsPerGene when the genome has no constants (e.g. boolean or integer
// genes).
func ProgramFromSymbols(symsPerGene [][]string, constsPerGene [][]float64, linkFunc string) Program {
	genes := make([]Expressor, len(symsPerGene))
	for i, syms := range symsPerGene {
		var consts []float64
		if i < len(constsPerGene) {
			consts = constsPerGene[i]
		}
		genes[i] = KarvaExpressor{Symbols: syms, Float64Constants: consts}
	}
	return Program{Genes: genes, LinkFunc: linkFunc}
}

// and numeric constants. It computes the argument-order tree from the grammar
// at render time so no typed node-arity information from the original genome
// is required. This allows experiment and example code to call codegen directly
// from a core.Genome[T] without importing the legacy gene or genome packages.
type KarvaExpressor struct {
	// Symbols contains the individual Karva symbol names (e.g. "And", "d0", "c0").
	Symbols []string
	// Float64Constants contains the constant values referenced by "c0", "c1", …
	// symbols. For boolean or integer genes with no constants this may be nil.
	Float64Constants []float64
}

// Expression builds the expression tree for this gene using the grammar and
// records any required helper functions into helpers.
func (e KarvaExpressor) Expression(grammar *grammars.Grammar, helpers grammars.HelperMap) (string, error) {
	argOrder := buildArgOrderFromGrammar(e.Symbols, grammar)
	numTerminals := computeNumTerminals(e.Symbols, e.Float64Constants)
	return Expression(e.Symbols, e.Float64Constants, numTerminals, argOrder, grammar, helpers)
}

// buildArgOrderFromGrammar computes the child-symbol index map by looking up
// function arities in the grammar's function map. This mirrors the logic in
// core.buildArgOrder but uses the grammar instead of typed node objects, keeping
// codegen free of any dependency on the core or gene packages.
func buildArgOrderFromGrammar(symbols []string, grammar *grammars.Grammar) [][]int {
	argOrder := make([][]int, len(symbols))
	argCount := 0
	for i, sym := range symbols {
		s, ok := grammar.Functions.FuncMap[sym]
		if !ok {
			continue
		}
		f, ok := s.(*grammars.Function)
		if !ok {
			continue
		}
		n := f.Terminals()
		if n <= 0 {
			continue
		}
		args := make([]int, n)
		for j := range n {
			argCount++
			args[j] = argCount
		}
		argOrder[i] = args
	}
	return argOrder
}

// computeNumTerminals returns the numTerminals value expected by Expression: the
// count of distinct input variables ("d0", "d1", …) plus the number of constants.
func computeNumTerminals(symbols []string, constants []float64) int {
	maxIdx := -1
	for _, sym := range symbols {
		if len(sym) > 1 && sym[0] == 'd' {
			if idx, err := strconv.Atoi(sym[1:]); err == nil && idx > maxIdx {
				maxIdx = idx
			}
		}
	}
	return (maxIdx + 1) + len(constants)
}

// Expressor renders a single gene expression for a grammar while collecting any
// helper functions required by that expression.
type Expressor interface {
	Expression(grammar *grammars.Grammar, helpers grammars.HelperMap) (string, error)
}

// Program describes the legacy genome data needed by the dedicated codegen
// subsystem.
type Program struct {
	Genes    []Expressor
	LinkFunc string
}

// Expression renders a single gene expression from legacy gene data.
func Expression(symbols []string, constants []float64, numTerminals int, argOrder [][]int, grammar *grammars.Grammar, helpers grammars.HelperMap) (string, error) {
	return buildExpression(symbols, constants, numTerminals, 0, argOrder, grammar, helpers)
}

func buildExpression(symbols []string, constants []float64, numTerminals, symbolIndex int, argOrder [][]int, grammar *grammars.Grammar, helpers grammars.HelperMap) (string, error) {
	if symbolIndex >= len(symbols) {
		return "", fmt.Errorf("bad symbolIndex %v for symbols: %v", symbolIndex, symbols)
	}

	sym := symbols[symbolIndex]
	if s, ok := grammar.Functions.FuncMap[sym]; ok {
		f, ok := s.(*grammars.Function)
		if !ok {
			return "", fmt.Errorf("unable to cast symbol %v to grammar function", sym)
		}
		exp := f.Chardata
		if _, ok := helpers[f.SymbolName]; !ok {
			if v, ok := grammar.Helpers.HelperMap[f.SymbolName]; ok {
				helpers[f.SymbolName] = v
			}
		}

		args := argOrder[symbolIndex]
		if len(args) < f.Terminals() {
			return "", fmt.Errorf("programming error: symbol %q args length mismatch: len(args)=%v, want %v; check FuncType", sym, len(args), f.Terminals())
		}

		for i := range f.Terminals() {
			e, err := buildExpression(symbols, constants, numTerminals, args[i], argOrder, grammar, helpers)
			if err != nil {
				return "", err
			}
			exp = strings.Replace(exp, "x"+strconv.Itoa(i), e, -1)
		}

		return exp, nil
	}

	if sym[0:1] == "d" {
		index, err := strconv.Atoi(sym[1:])
		if err != nil {
			return "", fmt.Errorf("unable to parse variable index: sym=%v", sym)
		}

		if n := numTerminals - len(constants); index >= n {
			return "", fmt.Errorf("terminal symbol name %q exceeds number of terminals (%v)", sym, n)
		}
		return fmt.Sprintf("d[%v]", index), nil
	}

	if sym[0:1] == "c" {
		index, err := strconv.Atoi(sym[1:])
		if err != nil {
			return "", fmt.Errorf("unable to parse constant index: sym=%v", sym)
		}

		if index >= len(constants) {
			return "", fmt.Errorf("constant symbol name %q exceeds length of constant slice (%v)", sym, len(constants))
		}
		return fmt.Sprintf("%v", constants[index]), nil
	}

	return "", fmt.Errorf("unable to render function: sym=%v for symbols %v", sym, symbols)
}

type generator struct {
	w       io.Writer
	grammar *grammars.Grammar
	program Program
	subs    map[string]string
}

// Generate renders a full program for the supplied grammar and returns the
// formatted source code bytes.
func Generate(program Program, grammar *grammars.Grammar) ([]byte, error) {
	g := &generator{
		grammar: grammar,
		program: program,
		subs: map[string]string{
			"CHARX": "X",
		},
	}
	return g.generate()
}

// Write renders a full program for the supplied grammar and writes it to w.
func Write(w io.Writer, program Program, grammar *grammars.Grammar) error {
	code, err := Generate(program, grammar)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s", code)
	return err
}

func (g *generator) generate() ([]byte, error) {
	var buf bytes.Buffer
	g.w = &buf
	g.write(g.grammar.Open)

	for _, h := range g.grammar.Headers {
		if h.Type != "default" {
			continue
		}
		g.write(h.Chardata)
		g.write(g.grammar.Endline)
	}

	for _, t := range g.grammar.Tempvars {
		if t.Type != "default" {
			continue
		}
		g.write(t.Chardata)
		g.subs["tempvarname"] = t.Varname
		g.write(g.grammar.Endline)
	}

	helpers := make(grammars.HelperMap)
	s, ok := g.grammar.Functions.FuncMap[g.program.LinkFunc]
	if !ok {
		return nil, fmt.Errorf("unable to find grammar linking function: %v", g.program.LinkFunc)
	}

	glf, ok := s.(*grammars.Function)
	if !ok {
		return nil, fmt.Errorf("error casting link function: %v", s.Symbol())
	}

	exps := []string{""}
	for i, gene := range g.program.Genes {
		exp, err := gene.Expression(g.grammar, helpers)
		if err != nil {
			return nil, err
		}

		if i > 0 {
			merge := strings.Replace(glf.Uniontype, "{tempvarname}", g.subs["tempvarname"], -1)
			merge = strings.Replace(merge, "{member}", exp, -1)
			merge = strings.Replace(merge, "{symbol}", glf.SymbolName, -1)
			exps = append(exps, merge)
		} else {
			exps = append(exps, g.subs["tempvarname"]+" = "+exp)
		}
	}

	exps = append(exps, "")
	fmt.Fprintln(g.w, strings.Join(exps, "\n"))

	for _, f := range g.grammar.Footers {
		if f.Type != "default" {
			continue
		}
		g.write(f.Chardata)
		g.write(g.grammar.Endline)
	}

	if len(helpers) > 0 {
		keys := make([]string, 0, len(helpers))
		for k := range helpers {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			g.write(g.grammar.Endline)
			g.write(helpers[k])
		}
	}

	clean, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.Bytes(), err
	}
	return clean, nil
}

func (g *generator) write(s string) {
	s = strings.Replace(s, "{CRLF}", "\n", -1)
	s = strings.Replace(s, "{TAB}", "\t", -1)
	for k, v := range g.subs {
		s = strings.Replace(s, fmt.Sprintf("{%v}", k), v, -1)
	}
	fmt.Fprint(g.w, s)
}
