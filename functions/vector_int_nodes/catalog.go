// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package vectorIntNodes

import (
	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/functions"
)

// vectorIntNodeAdapter adapts a functions.FuncNode to core.Node[functions.VectorInt].
type vectorIntNodeAdapter struct {
	fn functions.FuncNode
}

func (a vectorIntNodeAdapter) Symbol() string { return a.fn.Symbol() }
func (a vectorIntNodeAdapter) Arity() int     { return a.fn.Terminals() }
func (a vectorIntNodeAdapter) Eval(x []functions.VectorInt) functions.VectorInt {
	return a.fn.VectorIntFunction(x)
}

// CatalogFrom creates a typed core.Catalog[functions.VectorInt] from a legacy functions.FuncMap.
// All entries in fm are registered; an error is returned if any registration fails.
func CatalogFrom(fm functions.FuncMap) (*core.Catalog[functions.VectorInt], error) {
	cat := core.NewCatalog[functions.VectorInt]()
	for _, n := range fm {
		if err := cat.Register(vectorIntNodeAdapter{n}); err != nil {
			return nil, err
		}
	}
	return cat, nil
}
