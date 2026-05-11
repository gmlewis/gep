// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package intNodes

import (
	"fmt"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/functions"
)

// intNodeAdapter adapts a functions.FuncNode to core.Node[int].
type intNodeAdapter struct {
	fn functions.FuncNode
}

func (a intNodeAdapter) Symbol() string   { return a.fn.Symbol() }
func (a intNodeAdapter) Arity() int       { return a.fn.Terminals() }
func (a intNodeAdapter) Eval(x []int) int { return a.fn.IntFunction(x) }

// CatalogFrom creates a typed core.Catalog[int] from a legacy functions.FuncMap.
// All entries in fm are registered; an error is returned if any registration fails.
func CatalogFrom(fm functions.FuncMap) (*core.Catalog[int], error) {
	cat := core.NewCatalog[int]()
	for _, n := range fm {
		if err := cat.Register(intNodeAdapter{n}); err != nil {
			return nil, err
		}
	}
	return cat, nil
}

// CatalogFromNames creates a typed core.Catalog[int] containing only the named
// functions from Int. An error is returned if any name is not found or if any
// registration fails.
func CatalogFromNames(names []string) (*core.Catalog[int], error) {
	fm := make(functions.FuncMap, len(names))
	for _, sym := range names {
		fn, ok := Int[sym]
		if !ok {
			return nil, fmt.Errorf("intNodes.CatalogFromNames: unsupported int function %q", sym)
		}
		fm[sym] = fn
	}
	return CatalogFrom(fm)
}

// LinkFuncFrom returns a core.LinkFunc[int] for the named function in Int. The
// returned link operator folds the named binary integer operator across any
// number of gene outputs, so 1-gene and multi-gene genomes are both handled
// safely. An error is returned if the name is not found or if the named node is
// not binary.
func LinkFuncFrom(sym string) (core.LinkFunc[int], error) {
	fn, ok := Int[sym]
	if !ok {
		return core.LinkFunc[int]{}, fmt.Errorf("intNodes.LinkFuncFrom: unsupported int function %q", sym)
	}
	if fn.Terminals() != 2 {
		return core.LinkFunc[int]{}, fmt.Errorf("intNodes.LinkFuncFrom: function %q has arity %d, want binary operator", sym, fn.Terminals())
	}
	return core.NewLinkFunc[int](sym, func(v []int) int {
		switch len(v) {
		case 0:
			return 0
		case 1:
			return v[0]
		}
		result := v[0]
		for _, next := range v[1:] {
			result = fn.IntFunction([]int{result, next})
		}
		return result
	})
}
