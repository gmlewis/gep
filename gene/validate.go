package gene

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gmlewis/gep/v2/functions"
	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
	in "github.com/gmlewis/gep/v2/functions/int_nodes"
	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
	vin "github.com/gmlewis/gep/v2/functions/vector_int_nodes"
)

func validateSymbols(symbols []string, lookup functions.FuncMap, inputLen, constantsLen int, allowConstants bool) error {
	var errs []error
	for _, sym := range symbols {
		if sym == "" {
			errs = append(errs, errors.New("empty symbol"))
			continue
		}
		if _, ok := lookup[sym]; ok {
			continue
		}
		switch sym[0] {
		case 'd':
			index, err := strconv.Atoi(sym[1:])
			if err != nil {
				errs = append(errs, fmt.Errorf("unable to parse variable index: sym=%q: %w", sym, err))
				continue
			}
			if index >= inputLen {
				errs = append(errs, fmt.Errorf("error evaluating gene %q: index %v >= d length (%v)", sym, index, inputLen))
			}
		case 'c':
			if !allowConstants {
				errs = append(errs, fmt.Errorf("unknown gene symbol %q", sym))
				continue
			}
			index, err := strconv.Atoi(sym[1:])
			if err != nil {
				errs = append(errs, fmt.Errorf("unable to parse constant index: sym=%q: %w", sym, err))
				continue
			}
			if index >= constantsLen {
				errs = append(errs, fmt.Errorf("error evaluating gene %q: index %v >= c length (%v)", sym, index, constantsLen))
			}
		default:
			errs = append(errs, fmt.Errorf("unknown gene symbol %q", sym))
		}
	}
	return errors.Join(errs...)
}

func (g *Gene) validateIntSymbols(input []int) error {
	return validateSymbols(g.Symbols, in.Int, len(input), len(g.Constants), true)
}

func (g *Gene) validateMathSymbols(input []float64) error {
	return validateSymbols(g.Symbols, mn.Math, len(input), len(g.Constants), true)
}

func (g *Gene) validateBoolSymbols(input []bool) error {
	return validateSymbols(g.Symbols, bn.BoolAllGates, len(input), len(g.Constants), false)
}

func (g *Gene) validateVectorIntSymbols(input []VectorInt) error {
	return validateSymbols(g.Symbols, vin.VectorIntFuncs, len(input), len(g.Constants), true)
}
