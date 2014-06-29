package gene

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gmlewis/gep/grammars"
)

func (g *Gene) buildExp(symbolIndex int, argOrder [][]int, grammar *grammars.Grammar) (string, error) {
	if symbolIndex > len(g.Symbols) {
		return "", fmt.Errorf("bad symbolIndex %v for symbols: %v\n", symbolIndex, g.Symbols)
	}
	sym := g.Symbols[symbolIndex]
	if s, ok := grammar.Functions.FuncMap[sym]; ok {
		f, ok := s.(*grammars.Function)
		if !ok {
			return "", fmt.Errorf("unable to cast symbol %v to grammar function", sym)
		}
		exp := f.Chardata
		args := argOrder[symbolIndex]
		for i := 0; i < f.Terminals(); i++ {
			e, err := g.buildExp(args[i], argOrder, grammar)
			if err != nil {
				return "", err
			}
			exp = strings.Replace(exp, "x"+strconv.Itoa(i), e, -1)
		}
		return exp, nil
	} else { // No named symbol found - look for d0, d1, ...
		if sym[0:1] == "d" {
			if index, err := strconv.ParseInt(sym[1:], 10, 32); err != nil {
				return "", fmt.Errorf("unable to parse variable index: sym=%v\n", sym)
			} else {
				return fmt.Sprintf("d[%v]", index), nil
			}
		}
	}
	return "", fmt.Errorf("unable to render function: sym=%v for gene %#v\n", sym, g)
}

// Expression builds up the expression tree and returns the resulting string.
func (g *Gene) Expression(grammar *grammars.Grammar) (string, error) {
	argOrder := g.getBoolArgOrder(grammar.Functions.FuncMap)
	s, err := g.buildExp(0, argOrder, grammar)
	return s, err
}
