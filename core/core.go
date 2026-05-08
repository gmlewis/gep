// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package core defines the typed generic GEP surface.
package core

import "fmt"

// Node defines a typed function node in a Karva expression.
type Node[T any] interface {
	// Symbol is the Karva string representation of the function.
	Symbol() string
	// Arity is the number of input terminals consumed by the function.
	Arity() int
	// Eval evaluates the node for the given arguments.
	Eval([]T) T
}

// LinkOperator defines a typed link function used by a genome to combine gene outputs.
type LinkOperator[T any] interface {
	// Symbol is the Karva representation of the link operator.
	Symbol() string
	// Eval combines gene outputs into a final result.
	Eval([]T) T
}

// SymbolKind identifies the typed symbol category used by a gene.
type SymbolKind uint8

const (
	// SymbolKindUnknown indicates an uninitialized symbol.
	SymbolKindUnknown SymbolKind = iota
	// SymbolKindFunction identifies function symbols backed by a Node.
	SymbolKindFunction
	// SymbolKindTerminal identifies data terminal symbols.
	SymbolKindTerminal
	// SymbolKindConstant identifies constant symbols.
	SymbolKindConstant
)

// Symbol is the typed symbol representation for a Karva sequence.
type Symbol[T any] struct {
	Kind SymbolKind
	Name string

	Node Node[T]

	TerminalIndex int
	ConstantIndex int
}

// NewFunctionSymbol creates a typed function symbol.
func NewFunctionSymbol[T any](node Node[T]) (Symbol[T], error) {
	if node == nil {
		return Symbol[T]{}, fmt.Errorf("core.NewFunctionSymbol: node cannot be nil")
	}
	name := node.Symbol()
	if name == "" {
		return Symbol[T]{}, fmt.Errorf("core.NewFunctionSymbol: node symbol cannot be empty")
	}
	if node.Arity() < 0 {
		return Symbol[T]{}, fmt.Errorf("core.NewFunctionSymbol: node arity must be >= 0")
	}
	return Symbol[T]{
		Kind: SymbolKindFunction,
		Name: name,
		Node: node,
	}, nil
}

// NewTerminalSymbol creates a typed terminal symbol.
func NewTerminalSymbol[T any](index int) (Symbol[T], error) {
	if index < 0 {
		return Symbol[T]{}, fmt.Errorf("core.NewTerminalSymbol: index must be >= 0")
	}
	return Symbol[T]{
		Kind:          SymbolKindTerminal,
		Name:          fmt.Sprintf("d%v", index),
		TerminalIndex: index,
	}, nil
}

// NewConstantSymbol creates a typed constant symbol.
func NewConstantSymbol[T any](index int) (Symbol[T], error) {
	if index < 0 {
		return Symbol[T]{}, fmt.Errorf("core.NewConstantSymbol: index must be >= 0")
	}
	return Symbol[T]{
		Kind:          SymbolKindConstant,
		Name:          fmt.Sprintf("c%v", index),
		ConstantIndex: index,
	}, nil
}

// Gene is a typed GEP gene.
type Gene[T any] struct {
	Symbols   []Symbol[T]
	Constants []T
}

// Validate validates a typed gene.
func (g Gene[T]) Validate() error {
	if len(g.Symbols) == 0 {
		return fmt.Errorf("core.Gene.Validate: gene must contain at least one symbol")
	}
	for i, sym := range g.Symbols {
		switch sym.Kind {
		case SymbolKindFunction:
			if sym.Node == nil {
				return fmt.Errorf("core.Gene.Validate: function symbol[%d] missing node", i)
			}
		case SymbolKindTerminal, SymbolKindConstant:
			// valid
		default:
			return fmt.Errorf("core.Gene.Validate: symbol[%d] has unknown kind %v", i, sym.Kind)
		}
	}
	return nil
}

// Genome is a typed GEP genome.
type Genome[T any] struct {
	Genes []Gene[T]
	Link  LinkOperator[T]
}

// Validate validates a typed genome.
func (g Genome[T]) Validate() error {
	if len(g.Genes) == 0 {
		return fmt.Errorf("core.Genome.Validate: genome must contain at least one gene")
	}
	if g.Link == nil {
		return fmt.Errorf("core.Genome.Validate: link operator cannot be nil")
	}
	for i := range g.Genes {
		if err := g.Genes[i].Validate(); err != nil {
			return fmt.Errorf("core.Genome.Validate: gene[%d]: %w", i, err)
		}
	}
	return nil
}

// LinkFunc adapts a typed function to a LinkOperator.
type LinkFunc[T any] struct {
	name string
	fn   func([]T) T
}

// NewLinkFunc creates a LinkFunc adapter.
func NewLinkFunc[T any](name string, fn func([]T) T) (LinkFunc[T], error) {
	if name == "" {
		return LinkFunc[T]{}, fmt.Errorf("core.NewLinkFunc: name cannot be empty")
	}
	if fn == nil {
		return LinkFunc[T]{}, fmt.Errorf("core.NewLinkFunc: fn cannot be nil")
	}
	return LinkFunc[T]{name: name, fn: fn}, nil
}

// Symbol returns the Karva representation of the link operator.
func (l LinkFunc[T]) Symbol() string { return l.name }

// Eval evaluates the link operator.
func (l LinkFunc[T]) Eval(values []T) T { return l.fn(values) }
