// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package boolNodes

import (
	"fmt"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/functions"
)

// boolNodeAdapter adapts a functions.FuncNode to core.Node[bool].
type boolNodeAdapter struct {
	fn functions.FuncNode
}

func (a boolNodeAdapter) Symbol() string     { return a.fn.Symbol() }
func (a boolNodeAdapter) Arity() int         { return a.fn.Terminals() }
func (a boolNodeAdapter) Eval(x []bool) bool { return a.fn.BoolFunction(x) }

// CatalogFrom creates a typed core.Catalog[bool] from a legacy functions.FuncMap.
// All entries in fm are registered; an error is returned if any registration fails.
func CatalogFrom(fm functions.FuncMap) (*core.Catalog[bool], error) {
	cat := core.NewCatalog[bool]()
	for _, n := range fm {
		if err := cat.Register(boolNodeAdapter{n}); err != nil {
			return nil, err
		}
	}
	return cat, nil
}

// CatalogFromNames creates a typed core.Catalog[bool] containing only the named
// functions from BoolAllGates. An error is returned if any name is not found or
// if any registration fails.
func CatalogFromNames(names []string) (*core.Catalog[bool], error) {
	fm := make(functions.FuncMap, len(names))
	for _, sym := range names {
		fn, ok := BoolAllGates[sym]
		if !ok {
			return nil, fmt.Errorf("boolNodes.CatalogFromNames: unsupported boolean function %q", sym)
		}
		fm[sym] = fn
	}
	return CatalogFrom(fm)
}

// LinkFuncFrom returns a core.LinkFunc[bool] for the named function in
// BoolAllGates. An error is returned if the name is not found.
func LinkFuncFrom(sym string) (core.LinkFunc[bool], error) {
	fn, ok := BoolAllGates[sym]
	if !ok {
		return core.LinkFunc[bool]{}, fmt.Errorf("boolNodes.LinkFuncFrom: unsupported boolean function %q", sym)
	}
	return core.NewLinkFunc[bool](sym, fn.BoolFunction)
}
