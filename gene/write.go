// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"github.com/gmlewis/gep/v2/codegen"
	"github.com/gmlewis/gep/v2/grammars"
)

// Expression builds up the expression tree and returns the resulting string.
// While building, it keeps track of any helper functions that are needed.
func (g *Gene) Expression(grammar *grammars.Grammar, helpers grammars.HelperMap) (string, error) {
	argOrder, err := g.getArgOrder()
	if err != nil {
		return "", err
	}
	return codegen.Expression(g.Symbols, g.Constants, g.numTerminals, argOrder, grammar, helpers)
}
