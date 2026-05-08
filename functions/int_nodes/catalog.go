// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package intNodes

import (
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
