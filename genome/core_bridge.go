// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package genome

import (
	"fmt"

	"github.com/gmlewis/gep/v2/core"
	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
	in "github.com/gmlewis/gep/v2/functions/int_nodes"
	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
	"github.com/gmlewis/gep/v2/gene"
)

// CoreBool converts a legacy boolean genome into the typed core representation.
func (g *Genome) CoreBool() (core.Genome[bool], error) {
	if g == nil {
		return core.Genome[bool]{}, fmt.Errorf("genome.CoreBool: genome cannot be nil")
	}
	linkNode, ok := bn.BoolAllGates[g.LinkFunc]
	if !ok {
		return core.Genome[bool]{}, fmt.Errorf("genome.CoreBool: link function %q not found", g.LinkFunc)
	}
	link, err := core.NewLinkFunc[bool](g.LinkFunc, linkNode.BoolFunction)
	if err != nil {
		return core.Genome[bool]{}, fmt.Errorf("genome.CoreBool: %w", err)
	}
	genes := make([]core.Gene[bool], len(g.Genes))
	for i, legacyGene := range g.Genes {
		cg, err := legacyGene.CoreBool()
		if err != nil {
			return core.Genome[bool]{}, fmt.Errorf("genome.CoreBool: gene[%d]: %w", i, err)
		}
		genes[i] = cg
	}
	return core.Genome[bool]{Genes: genes, Link: link}, nil
}

// CoreInt converts a legacy integer genome into the typed core representation.
func (g *Genome) CoreInt() (core.Genome[int], error) {
	if g == nil {
		return core.Genome[int]{}, fmt.Errorf("genome.CoreInt: genome cannot be nil")
	}
	linkNode, ok := in.Int[g.LinkFunc]
	if !ok {
		return core.Genome[int]{}, fmt.Errorf("genome.CoreInt: link function %q not found", g.LinkFunc)
	}
	link, err := core.NewLinkFunc[int](g.LinkFunc, linkNode.IntFunction)
	if err != nil {
		return core.Genome[int]{}, fmt.Errorf("genome.CoreInt: %w", err)
	}
	genes := make([]core.Gene[int], len(g.Genes))
	for i, legacyGene := range g.Genes {
		cg, err := legacyGene.CoreInt()
		if err != nil {
			return core.Genome[int]{}, fmt.Errorf("genome.CoreInt: gene[%d]: %w", i, err)
		}
		genes[i] = cg
	}
	return core.Genome[int]{Genes: genes, Link: link}, nil
}

// CoreFloat64 converts a legacy floating-point genome into the typed core representation.
func (g *Genome) CoreFloat64() (core.Genome[float64], error) {
	if g == nil {
		return core.Genome[float64]{}, fmt.Errorf("genome.CoreFloat64: genome cannot be nil")
	}
	linkNode, ok := mn.Math[g.LinkFunc]
	if !ok {
		return core.Genome[float64]{}, fmt.Errorf("genome.CoreFloat64: link function %q not found", g.LinkFunc)
	}
	link, err := core.NewLinkFunc[float64](g.LinkFunc, linkNode.Float64Function)
	if err != nil {
		return core.Genome[float64]{}, fmt.Errorf("genome.CoreFloat64: %w", err)
	}
	genes := make([]core.Gene[float64], len(g.Genes))
	for i, legacyGene := range g.Genes {
		cg, err := legacyGene.CoreFloat64()
		if err != nil {
			return core.Genome[float64]{}, fmt.Errorf("genome.CoreFloat64: gene[%d]: %w", i, err)
		}
		genes[i] = cg
	}
	return core.Genome[float64]{Genes: genes, Link: link}, nil
}

// NewFromCoreBool converts a typed boolean genome into the legacy representation.
func NewFromCoreBool(g core.Genome[bool]) (*Genome, error) {
	genes := make([]*gene.Gene, len(g.Genes))
	for i, cg := range g.Genes {
		lg, err := gene.NewFromCoreBool(cg)
		if err != nil {
			return nil, fmt.Errorf("genome.NewFromCoreBool: gene[%d]: %w", i, err)
		}
		genes[i] = lg
	}
	if g.Link == nil {
		return nil, fmt.Errorf("genome.NewFromCoreBool: link operator cannot be nil")
	}
	return New(genes, g.Link.Symbol()), nil
}

// NewFromCoreInt converts a typed integer genome into the legacy representation.
func NewFromCoreInt(g core.Genome[int]) (*Genome, error) {
	genes := make([]*gene.Gene, len(g.Genes))
	for i, cg := range g.Genes {
		lg, err := gene.NewFromCoreInt(cg)
		if err != nil {
			return nil, fmt.Errorf("genome.NewFromCoreInt: gene[%d]: %w", i, err)
		}
		genes[i] = lg
	}
	if g.Link == nil {
		return nil, fmt.Errorf("genome.NewFromCoreInt: link operator cannot be nil")
	}
	return New(genes, g.Link.Symbol()), nil
}

// NewFromCoreFloat64 converts a typed floating-point genome into the legacy representation.
func NewFromCoreFloat64(g core.Genome[float64]) (*Genome, error) {
	genes := make([]*gene.Gene, len(g.Genes))
	for i, cg := range g.Genes {
		lg, err := gene.NewFromCoreFloat64(cg)
		if err != nil {
			return nil, fmt.Errorf("genome.NewFromCoreFloat64: gene[%d]: %w", i, err)
		}
		genes[i] = lg
	}
	if g.Link == nil {
		return nil, fmt.Errorf("genome.NewFromCoreFloat64: link operator cannot be nil")
	}
	return New(genes, g.Link.Symbol()), nil
}
