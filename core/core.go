// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package core defines the typed generic GEP surface.
package core

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

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

// Catalog is a typed registry that maps Karva symbol names to Node[T] implementations.
type Catalog[T any] struct {
	nodes map[string]Node[T]
}

// NewCatalog creates a new empty Catalog.
func NewCatalog[T any]() *Catalog[T] {
	return &Catalog[T]{nodes: make(map[string]Node[T])}
}

// Register adds node to the catalog.
// It returns an error if node is nil, if the node's symbol is empty, or if a
// node with the same symbol is already registered.
func (c *Catalog[T]) Register(node Node[T]) error {
	if node == nil {
		return fmt.Errorf("core.Catalog.Register: node cannot be nil")
	}
	sym := node.Symbol()
	if sym == "" {
		return fmt.Errorf("core.Catalog.Register: node symbol cannot be empty")
	}
	if _, ok := c.nodes[sym]; ok {
		return fmt.Errorf("core.Catalog.Register: symbol %q already registered", sym)
	}
	c.nodes[sym] = node
	return nil
}

// Lookup returns the Node for the given symbol name and whether it was found.
func (c *Catalog[T]) Lookup(name string) (Node[T], bool) {
	n, ok := c.nodes[name]
	return n, ok
}

// Symbols returns a sorted slice of all function symbol names registered in the catalog.
func (c *Catalog[T]) Symbols() []string {
	syms := make([]string, 0, len(c.nodes))
	for sym := range c.nodes {
		syms = append(syms, sym)
	}
	sort.Strings(syms)
	return syms
}

// MaxArity returns the maximum arity of any node registered in the catalog,
// or 0 if the catalog is empty.
func (c *Catalog[T]) MaxArity() int {
	max := 0
	for _, n := range c.nodes {
		if a := n.Arity(); a > max {
			max = a
		}
	}
	return max
}

// ParseSymbols converts a Karva symbol string slice into a typed []Symbol[T].
// Terminal symbols (d0, d1, ...) and constant symbols (c0, c1, ...) are
// recognised by their conventional prefixes. All other symbols are looked up in
// catalog; an error is returned for any symbol that is not found.
func ParseSymbols[T any](symbols []string, catalog *Catalog[T]) ([]Symbol[T], error) {
	if catalog == nil {
		return nil, fmt.Errorf("core.ParseSymbols: catalog cannot be nil")
	}
	result := make([]Symbol[T], 0, len(symbols))
	for i, s := range symbols {
		if s == "" {
			return nil, fmt.Errorf("core.ParseSymbols: symbol[%d] is empty", i)
		}
		// Terminal: d<N> where N is a non-negative integer.
		if len(s) > 1 && s[0] == 'd' {
			if index, err := strconv.Atoi(s[1:]); err == nil && index >= 0 {
				sym, err := NewTerminalSymbol[T](index)
				if err != nil {
					return nil, fmt.Errorf("core.ParseSymbols: symbol[%d] %q: %w", i, s, err)
				}
				result = append(result, sym)
				continue
			}
			// index < 0 or non-numeric suffix: not a terminal; fall through to catalog.
		}
		// Constant: c<N> where N is a non-negative integer.
		if len(s) > 1 && s[0] == 'c' {
			if index, err := strconv.Atoi(s[1:]); err == nil && index >= 0 {
				sym, err := NewConstantSymbol[T](index)
				if err != nil {
					return nil, fmt.Errorf("core.ParseSymbols: symbol[%d] %q: %w", i, s, err)
				}
				result = append(result, sym)
				continue
			}
			// index < 0 or non-numeric suffix: not a constant; fall through to catalog.
		}
		// Function: look up in catalog
		node, ok := catalog.Lookup(s)
		if !ok {
			return nil, fmt.Errorf("core.ParseSymbols: symbol[%d] %q not found in catalog", i, s)
		}
		sym, err := NewFunctionSymbol[T](node)
		if err != nil {
			return nil, fmt.Errorf("core.ParseSymbols: symbol[%d] %q: %w", i, s, err)
		}
		result = append(result, sym)
	}
	return result, nil
}

// buildArgOrder computes child-symbol indices for each function symbol in a
// Karva sequence. The entry at index i holds the argument positions for
// symbol i; terminal and constant symbols receive a nil entry.
func buildArgOrder[T any](symbols []Symbol[T]) [][]int {
	argOrder := make([][]int, len(symbols))
	argCount := 0
	for i, sym := range symbols {
		if sym.Kind != SymbolKindFunction {
			continue
		}
		n := sym.Node.Arity()
		if n <= 0 {
			continue
		}
		args := make([]int, n)
		for j := 0; j < n; j++ {
			argCount++
			args[j] = argCount
		}
		argOrder[i] = args
	}
	return argOrder
}

// evalTree recursively evaluates the expression subtree rooted at symbolIndex.
func evalTree[T any](symbols []Symbol[T], argOrder [][]int, symbolIndex int, terminals, constants []T) (T, error) {
	var zero T
	if symbolIndex >= len(symbols) {
		return zero, fmt.Errorf("core: symbolIndex %d out of range (len=%d)", symbolIndex, len(symbols))
	}
	sym := symbols[symbolIndex]
	switch sym.Kind {
	case SymbolKindFunction:
		childIndices := argOrder[symbolIndex]
		args := make([]T, len(childIndices))
		for i, idx := range childIndices {
			v, err := evalTree(symbols, argOrder, idx, terminals, constants)
			if err != nil {
				return zero, err
			}
			args[i] = v
		}
		return sym.Node.Eval(args), nil
	case SymbolKindTerminal:
		idx := sym.TerminalIndex
		if idx >= len(terminals) {
			return zero, fmt.Errorf("core: terminal index %d out of range (len=%d)", idx, len(terminals))
		}
		return terminals[idx], nil
	case SymbolKindConstant:
		idx := sym.ConstantIndex
		if idx >= len(constants) {
			return zero, fmt.Errorf("core: constant index %d out of range (len=%d)", idx, len(constants))
		}
		return constants[idx], nil
	default:
		return zero, fmt.Errorf("core: symbol[%d] has unknown kind %v", symbolIndex, sym.Kind)
	}
}

// Eval evaluates the gene's Karva expression for the given terminal inputs.
// Constants are read from g.Constants.
func (g Gene[T]) Eval(terminals []T) (T, error) {
	var zero T
	if err := g.Validate(); err != nil {
		return zero, fmt.Errorf("core.Gene.Eval: %w", err)
	}
	argOrder := buildArgOrder(g.Symbols)
	return evalTree(g.Symbols, argOrder, 0, terminals, g.Constants)
}

// Eval evaluates the genome by evaluating all genes and combining their outputs
// with the link operator.
func (g Genome[T]) Eval(terminals []T) (T, error) {
	var zero T
	if err := g.Validate(); err != nil {
		return zero, fmt.Errorf("core.Genome.Eval: %w", err)
	}
	geneOutputs := make([]T, len(g.Genes))
	for i, gene := range g.Genes {
		v, err := gene.Eval(terminals)
		if err != nil {
			return zero, fmt.Errorf("core.Genome.Eval: gene[%d]: %w", i, err)
		}
		geneOutputs[i] = v
	}
	return g.Link.Eval(geneOutputs), nil
}

// NewRandomGene creates a random Gene[T] from the given catalog.
//
// headSize is the number of positions in the gene head. The tail size is derived
// automatically as max(1, (maxArity-1)*headSize+1) where maxArity is the maximum
// arity of any node in cat; this guarantees a syntactically valid Karva expression.
//
// numTerminals is the count of available d-prefixed input variables (d0..d(n-1)).
// numConstants is the count of available c-prefixed constants (c0..c(n-1)).
// rng is the random source; if nil the global math/rand source is used.
//
// Returns an error if cat is nil, headSize < 0, numTerminals < 0, numConstants < 0,
// or numTerminals+numConstants == 0.
func NewRandomGene[T any](cat *Catalog[T], headSize, numTerminals, numConstants int, rng *rand.Rand) (Gene[T], error) {
	if cat == nil {
		return Gene[T]{}, fmt.Errorf("core.NewRandomGene: catalog cannot be nil")
	}
	if headSize < 0 {
		return Gene[T]{}, fmt.Errorf("core.NewRandomGene: headSize must be >= 0")
	}
	if numTerminals < 0 {
		return Gene[T]{}, fmt.Errorf("core.NewRandomGene: numTerminals must be >= 0")
	}
	if numConstants < 0 {
		return Gene[T]{}, fmt.Errorf("core.NewRandomGene: numConstants must be >= 0")
	}
	if numTerminals+numConstants == 0 {
		return Gene[T]{}, fmt.Errorf("core.NewRandomGene: numTerminals+numConstants must be > 0")
	}

	intn := func(n int) int {
		if rng != nil {
			return rng.Intn(n)
		}
		return rand.Intn(n) //nolint:gosec
	}

	// Calculate tail size from max arity so that any parse tree within the head
	// can always be completed by terminals in the tail.
	maxA := cat.MaxArity()
	tailSize := 1
	if maxA > 1 && headSize > 0 {
		tailSize = (maxA-1)*headSize + 1
	}

	// Build terminal choice slice: d0..d(numTerminals-1), c0..c(numConstants-1).
	termChoices := make([]Symbol[T], 0, numTerminals+numConstants)
	for i := 0; i < numTerminals; i++ {
		sym, _ := NewTerminalSymbol[T](i)
		termChoices = append(termChoices, sym)
	}
	for i := 0; i < numConstants; i++ {
		sym, _ := NewConstantSymbol[T](i)
		termChoices = append(termChoices, sym)
	}

	// Build function choice slice from catalog (sorted for determinism).
	funcNames := cat.Symbols()

	// Head: any function or terminal.
	headChoiceCount := len(funcNames) + len(termChoices)
	syms := make([]Symbol[T], 0, headSize+tailSize)
	for i := 0; i < headSize; i++ {
		choice := intn(headChoiceCount)
		if choice < len(funcNames) {
			node, _ := cat.Lookup(funcNames[choice])
			sym, _ := NewFunctionSymbol[T](node)
			syms = append(syms, sym)
		} else {
			syms = append(syms, termChoices[choice-len(funcNames)])
		}
	}

	// Tail: terminals only.
	for i := 0; i < tailSize; i++ {
		syms = append(syms, termChoices[intn(len(termChoices))])
	}

	return Gene[T]{Symbols: syms}, nil
}

// KarvaString returns the dot-separated Karva representation of the gene's symbols.
// For example: "+.d0.d1" or "*.c0.d1".
func (g Gene[T]) KarvaString() string {
	names := make([]string, len(g.Symbols))
	for i, sym := range g.Symbols {
		names[i] = sym.Name
	}
	return strings.Join(names, ".")
}

// Dup returns a deep copy of the gene.
func (g Gene[T]) Dup() Gene[T] {
	syms := make([]Symbol[T], len(g.Symbols))
	copy(syms, g.Symbols)
	consts := make([]T, len(g.Constants))
	copy(consts, g.Constants)
	return Gene[T]{Symbols: syms, Constants: consts}
}

// KarvaString returns the multi-gene Karva representation of the genome.
// Each gene is rendered with its own KarvaString and genes are joined by the
// link operator symbol, e.g. "+.d0.d1|+|*.d0.d1".
// Returns an empty string when the genome has no genes or a nil link operator.
func (g Genome[T]) KarvaString() string {
	if len(g.Genes) == 0 || g.Link == nil {
		return ""
	}
	parts := make([]string, len(g.Genes))
	for i, gene := range g.Genes {
		parts[i] = gene.KarvaString()
	}
	sep := "|" + g.Link.Symbol() + "|"
	return strings.Join(parts, sep)
}

// Dup returns a deep copy of the genome.  The link operator is shared (it is
// assumed to be stateless and immutable).
func (g Genome[T]) Dup() Genome[T] {
	genes := make([]Gene[T], len(g.Genes))
	for i := range g.Genes {
		genes[i] = g.Genes[i].Dup()
	}
	return Genome[T]{Genes: genes, Link: g.Link}
}

// NewRandomGenome creates a random Genome[T] containing numGenes random genes.
// Each gene is created with NewRandomGene using the provided parameters.
// link is the typed link operator used to combine gene outputs.
// rng is the random source; if nil the global math/rand source is used.
//
// Returns an error if numGenes <= 0, link is nil, or any gene creation fails.
func NewRandomGenome[T any](cat *Catalog[T], numGenes, headSize, numTerminals, numConstants int, link LinkOperator[T], rng *rand.Rand) (Genome[T], error) {
	if numGenes <= 0 {
		return Genome[T]{}, fmt.Errorf("core.NewRandomGenome: numGenes must be > 0")
	}
	if link == nil {
		return Genome[T]{}, fmt.Errorf("core.NewRandomGenome: link cannot be nil")
	}
	genes := make([]Gene[T], numGenes)
	for i := range genes {
		g, err := NewRandomGene(cat, headSize, numTerminals, numConstants, rng)
		if err != nil {
			return Genome[T]{}, fmt.Errorf("core.NewRandomGenome: gene[%d]: %w", i, err)
		}
		genes[i] = g
	}
	return Genome[T]{Genes: genes, Link: link}, nil
}
