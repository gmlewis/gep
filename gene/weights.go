package gene

import (
	"log"

	"github.com/gmlewis/gep/v2/functions"
	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
	in "github.com/gmlewis/gep/v2/functions/int_nodes"
	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
	vin "github.com/gmlewis/gep/v2/functions/vector_int_nodes"
)

// FuncWeight contains the symbol name and its weight to be used during
// a run of the GEP algorithm. A symbol with weight 5, for example, will
// be five times more likely to be used than a symbol with weight 1.
type FuncWeight struct {
	Symbol string
	Weight int
}

// AllSymbolsEqualWeights returns all symbols with equal weights
// for the given node type.
func AllSymbolsEqualWeights(funcType functions.FuncType) []FuncWeight {
	var lookup functions.FuncMap
	switch funcType {
	case functions.Bool:
		lookup = bn.BoolAllGates
	case functions.Int:
		lookup = in.Int
	case functions.Float64:
		lookup = mn.Math
	case functions.VectorInts:
		lookup = vin.VectorIntFuncs
	default:
		log.Fatalf("unknown funcType: %v", funcType)
	}

	result := make([]FuncWeight, 0, len(lookup))
	for key := range lookup {
		result = append(result, FuncWeight{
			Symbol: key,
			Weight: 1,
		})
	}
	return result
}
