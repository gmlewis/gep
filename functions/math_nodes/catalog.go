// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package mathNodes

import (
	"fmt"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/functions"
)

// float64NodeAdapter adapts a functions.FuncNode to core.Node[float64].
type float64NodeAdapter struct {
	fn functions.FuncNode
}

func (a float64NodeAdapter) Symbol() string           { return a.fn.Symbol() }
func (a float64NodeAdapter) Arity() int               { return a.fn.Terminals() }
func (a float64NodeAdapter) Eval(x []float64) float64 { return a.fn.Float64Function(x) }

// CatalogFrom creates a typed core.Catalog[float64] from a legacy functions.FuncMap.
// All entries in fm are registered; an error is returned if any registration fails.
func CatalogFrom(fm functions.FuncMap) (*core.Catalog[float64], error) {
	cat := core.NewCatalog[float64]()
	for _, n := range fm {
		if err := cat.Register(float64NodeAdapter{n}); err != nil {
			return nil, err
		}
	}
	return cat, nil
}

// CatalogFromNames creates a typed core.Catalog[float64] containing only the
// named functions from Math. An error is returned if any name is not found or
// if any registration fails.
func CatalogFromNames(names []string) (*core.Catalog[float64], error) {
	fm := make(functions.FuncMap, len(names))
	for _, sym := range names {
		fn, ok := Math[sym]
		if !ok {
			return nil, fmt.Errorf("mathNodes.CatalogFromNames: unsupported math function %q", sym)
		}
		fm[sym] = fn
	}
	return CatalogFrom(fm)
}

// LinkFuncFrom returns a core.LinkFunc[float64] for the named function in Math.
// The returned link operator folds the named binary float64 operator across any
// number of gene outputs, so 1-gene and multi-gene genomes are both handled
// safely. An error is returned if the name is not found or if the named node is
// not binary.
func LinkFuncFrom(sym string) (core.LinkFunc[float64], error) {
	fn, ok := Math[sym]
	if !ok {
		return core.LinkFunc[float64]{}, fmt.Errorf("mathNodes.LinkFuncFrom: unsupported math function %q", sym)
	}
	if fn.Terminals() != 2 {
		return core.LinkFunc[float64]{}, fmt.Errorf("mathNodes.LinkFuncFrom: function %q has arity %d, want binary operator", sym, fn.Terminals())
	}
	return core.NewLinkFunc[float64](sym, func(v []float64) float64 {
		switch len(v) {
		case 0:
			return 0
		case 1:
			return v[0]
		}
		result := v[0]
		for _, next := range v[1:] {
			result = fn.Float64Function([]float64{result, next})
		}
		return result
	})
}
