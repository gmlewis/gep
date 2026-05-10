# Gene Expression Programming (GEP) in Go

`github.com/gmlewis/gep/v2` is a typed Gene Expression Programming engine for
scientific and engineering search in Go.

The repository now has a clear default architecture:

- `core` defines typed genes, genomes, symbols, catalogs, and link operators
- `evolution` runs typed population search with configurable mutation,
  recombination, transposition, selection, statistics, and termination
- `problems` provides reusable typed scoring helpers for common boolean and
  floating-point tasks
- `codegen` renders evolved Karva programs through optional grammar backends
- `env` and `gymnasium` provide an exploratory environment/agent layer for
  discrete and tuple-space experimentation

## Status

The primary workflow is the typed stack:

- `core`
- `evolution`
- `problems`
- `codegen`

The `env` subsystem is usable for discrete and tuple-space agent experiments,
but it is still an exploratory RL adapter rather than a complete modern RL
framework.

Legacy `gene` and `genome` packages remain only as compatibility/reference
layers. New workflow code should not build on them.

## Package map

| Package | Role | Use it when |
| --- | --- | --- |
| `core` | Typed GEP representation and random genome construction | You need `Node[T]`, `Genome[T]`, `Catalog[T]`, or direct genome evaluation |
| `evolution` | Typed population search engine | You need seeded experiments, operators, stopping criteria, or per-generation statistics |
| `evolution/*` | Operator and evaluation subsystems | You are tuning mutation, recombination, selection, transposition, termination, or statistics behavior |
| `problems` | Reusable typed scoring seams | Your problem is a reusable boolean or regression task instead of a one-off experiment |
| `codegen` | Grammar-backed code generation | You want Go (or other grammar-backed) source emitted from evolved Karva expressions |
| `functions/*_nodes` | Ready-made node catalogs | You want to start from the built-in boolean, integer, float, or vector-int operators |
| `grammars` | Code-generation grammars | You want to render evolved programs into source code |
| `env` / `gymnasium` | Exploratory environment integration | You are experimenting with Gymnasium-style environments and discrete action/observation spaces |
| `experiments/*` | End-to-end examples | You want concrete entrypoints that exercise the typed stack |
| `gene`, `genome` | Legacy compatibility layers | You are maintaining compatibility code, not building new features |

## Quick start

The fastest path is:

1. build or reuse a typed catalog
2. define a typed scoring function over `core.Genome[T]`
3. create a seeded `evolution.Generation[T]`
4. evolve until the stop condition is met
5. optionally render the result with `codegen`

```go
package main

import (
	"fmt"
	"log"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/evolution"
	boolNodes "github.com/gmlewis/gep/v2/functions/bool_nodes"
)

var nandCases = []struct {
	in  []bool
	out bool
}{
	{[]bool{false, false}, true},
	{[]bool{false, true}, true},
	{[]bool{true, false}, true},
	{[]bool{true, true}, false},
}

func scoreNAND(g core.Genome[bool]) float64 {
	hits := 0
	for _, tc := range nandCases {
		got, err := g.Eval(tc.in)
		if err != nil {
			return 0
		}
		if got == tc.out {
			hits++
		}
	}
	return 1000.0 * float64(hits) / float64(len(nandCases))
}

func main() {
	cat, err := boolNodes.CatalogFromNames([]string{"Not", "And", "Or"})
	if err != nil {
		log.Fatal(err)
	}
	link, err := boolNodes.LinkFuncFrom("Or")
	if err != nil {
		log.Fatal(err)
	}

	gen, err := evolution.NewWithSeed(42, cat, 30, 7, 1, 2, 0, link, scoreNAND)
	if err != nil {
		log.Fatal(err)
	}

	best := gen.Evolve(250)
	fmt.Println(best.Score)
	fmt.Println(best.Genome.KarvaString())
}
```

## Optional code generation

If you want source code from an evolved genome, convert it to a `codegen.Program`
and render it with a grammar:

```go
prog := codegen.ProgramFromSymbols(
	best.Genome.SymbolNamesPerGene(),
	nil,
	best.Genome.Link.Symbol(),
)
grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
if err != nil {
	return err
}
return codegen.Write(os.Stdout, prog, grammar)
```

See:

- `experiments/nand`
- `experiments/symbolic_regression`

## Reproducible experiments

Use `evolution.NewWithSeed` whenever you care about deterministic replay.
For each run, record at least:

- seed
- package version / commit SHA
- catalog contents
- link operator
- population size and gene geometry
- mutation, recombination, and transposition configs
- stopping criteria
- scoring function definition and dataset/problem snapshot

If you emit code or downstream artifacts, store the final `KarvaString`,
`SymbolNamesPerGene`, constants, and rendered output together.

## Extending the engine

### Add a new typed domain

1. Define a Go type for your terminals and gene outputs.
2. Implement `core.Node[T]` for each function/operator in the domain.
3. Register those nodes in a `core.Catalog[T]`.
4. Define a typed link operator with `core.NewLinkFunc`.
5. Write a scoring function over `core.Genome[T]`.
6. Evolve with `evolution.New` or `evolution.NewWithSeed`.

### Add reusable problems

Put reusable scoring/problem definitions into `problems` or a sibling package
with typed seams. Keep one-off experiment scoring logic close to the experiment
entrypoint.

### Add code generation

If the output can be expressed through the grammar system, use `codegen` and
`grammars`. If not, treat the evolved genome as an intermediate representation
and write a domain-specific emitter.

### Add RL or simulator-backed workflows

Use `core` and `evolution` as the search engine, then place simulator calls,
reward aggregation, train/validation splits, and artifact generation in a
domain-specific package. The current `env` package is a useful reference for
agent orchestration, but advanced RL work will often want a richer typed layer.

## Included entrypoints

- `go run ./experiments/nand`
- `go run ./experiments/odd-3-parity`
- `go run ./experiments/odd-7-parity`
- `go run ./experiments/6-multiplexer`
- `go run ./experiments/symbolic_regression`
- `go run ./examples/gymnasium/toy_text/blackjack-go`

## Quality gates

Repo-level verification:

- `./scripts/test-all.sh`
- `./scripts/bench-all.sh`

GitHub Actions runs CI and benchmark workflows from `.github/workflows/`.

## License

Copyright 2014-2026 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0. See `LICENSE`.
