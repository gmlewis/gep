// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package boolNodes

import (
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
