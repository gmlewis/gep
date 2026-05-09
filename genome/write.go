// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package genome

import (
	"io"

	"github.com/gmlewis/gep/v2/codegen"
	"github.com/gmlewis/gep/v2/grammars"
)

func (g *Genome) Write(w io.Writer, grammar *grammars.Grammar) error {
	genes := make([]codegen.Expressor, len(g.Genes))
	for i, gene := range g.Genes {
		genes[i] = gene
	}
	return codegen.Write(w, codegen.Program{
		Genes:    genes,
		LinkFunc: g.LinkFunc,
	}, grammar)
}
