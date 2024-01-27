// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gmlewis/gep/v2/grammars"
)

func (g *Gene) buildExp(symbolIndex int, argOrder [][]int, grammar *grammars.Grammar, helpers grammars.HelperMap) (string, error) {
	if symbolIndex > len(g.Symbols) {
		return "", fmt.Errorf("bad symbolIndex %v for symbols: %v", symbolIndex, g.Symbols)
	}

	sym := g.Symbols[symbolIndex]
	if s, ok := grammar.Functions.FuncMap[sym]; ok {
		f, ok := s.(*grammars.Function)
		if !ok {
			return "", fmt.Errorf("unable to cast symbol %v to grammar function", sym)
		}
		exp := f.Chardata
		// Look up the SymbolName in the grammar's Helpers list to see if there is a replacement.
		if _, ok := helpers[f.SymbolName]; !ok {
			if v, ok := grammar.Helpers.HelperMap[f.SymbolName]; ok {
				helpers[f.SymbolName] = v
			}
		}

		args := argOrder[symbolIndex]
		if len(args) < f.Terminals() {
			log.Fatalf("programming error: symbol %q args length mismatch: len(args)=%v, want %v; check FuncType", sym, len(args), f.Terminals())
		}

		for i := 0; i < f.Terminals(); i++ {
			e, err := g.buildExp(args[i], argOrder, grammar, helpers)
			if err != nil {
				return "", err
			}
			exp = strings.Replace(exp, "x"+strconv.Itoa(i), e, -1)
		}

		return exp, nil
	}

	// No named symbol found - look for d0, d1, ... or constants c0, c1, ...
	if sym[0:1] == "d" {
		index, err := strconv.Atoi(sym[1:])
		if err != nil {
			return "", fmt.Errorf("unable to parse variable index: sym=%v", sym)
		}

		if n := g.numTerminals - len(g.Constants); index > n {
			log.Fatalf("programming error: terminal symbol name %q exceeds number of terminals (%v)", sym, n)
		}
		return fmt.Sprintf("d[%v]", index), nil
	}

	if sym[0:1] == "c" {
		index, err := strconv.Atoi(sym[1:])
		if err != nil {
			return "", fmt.Errorf("unable to parse constant index: sym=%v", sym)
		}

		if index > len(g.Constants) {
			log.Fatalf("programming error: constant symbol name %q exceeds length of constant slice (%v)", sym, len(g.Constants))
		}
		return fmt.Sprintf("%0.2f", g.Constants[index]), nil
	}

	return "", fmt.Errorf("unable to render function: sym=%v for gene %#v", sym, g)
}

// Expression builds up the expression tree and returns the resulting string.
// While building, it keeps track of any helper functions that are needed.
func (g *Gene) Expression(grammar *grammars.Grammar, helpers grammars.HelperMap) (string, error) {
	argOrder := g.getArgOrder()
	s, err := g.buildExp(0, argOrder, grammar, helpers)
	return s, err
}
