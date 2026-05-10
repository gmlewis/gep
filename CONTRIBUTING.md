# Contributing

Thanks for contributing to `github.com/gmlewis/gep/v2`.

The default goal is to keep strengthening the typed GEP platform so it
is useful for scientific, engineering, and simulator-backed search
workflows.

## Project direction

New work should prefer the typed architecture:

- `core` for typed genes, genomes, catalogs, and link operators
- `evolution` for typed population search and operator orchestration
- `problems` for reusable typed scoring/problem definitions
- `codegen` for grammar-backed rendering

The `env` subsystem is an exploratory RL adapter. It is a good place for
environment orchestration work, but not a reason to add new workflow code to the
legacy `gene` or `genome` packages.

Legacy `gene` and `genome` packages are compatibility/reference layers. Avoid
building new features on them.

## Architecture guardrails

### 1. Prefer typed seams

New features should be built on `core.Genome[T]`, `core.Node[T]`, and
`evolution.Generation[T]`.

### 2. Keep package responsibilities narrow

- `core`: representation and evaluation primitives
- `evolution`: search loop and operator orchestration
- `problems`: reusable problem definitions
- `codegen`: program rendering
- `env`: environment/agent orchestration
- `experiments`: demonstrations and regression fixtures

### 3. Treat docs as part of the feature

Exported APIs should have strong godoc-style comments. If a change
alters the intended public workflow, update the root docs as part of
the same change.

### 4. Preserve reproducibility

Prefer seeded execution paths and explicit configuration surfaces. When adding a
new experiment or benchmark, make it straightforward to rerun the exact setup.

## Common contribution types

### Add a new function set

Preferred path:

1. define typed `core.Node[T]` implementations for the new domain
2. register them in a `core.Catalog[T]`
3. expose helpers similar to `CatalogFromNames` / `LinkFuncFrom` if they improve usability
4. add focused tests for node behavior and catalog construction

### Add a new problem/domain package

Reusable scoring or dataset logic belongs in `problems` or a closely related
typed package. Keep the seam in terms of `func(core.Genome[T]) float64` so it
plugs directly into `evolution.New`.

### Add a new code generation backend

If the grammar system is a good fit, extend `grammars` and keep rendering logic
inside `codegen`. If a domain needs a custom emitter, keep that emitter outside
`core` and `evolution`.

### Add a new environment or RL integration

Keep environment metadata and transport separate from search/training logic.
When possible, design new environment-facing code around typed state/action
models instead of widening the legacy int-centric path.

## Verification

Run the repo quality gates before sending a change:

```bash
./scripts/test-all.sh
./scripts/bench-all.sh
```

Use additional targeted commands when appropriate, for example:

```bash
go test ./core ./evolution ./codegen ./problems ./env
go test ./experiments/... ./examples/...
```

## Pull requests

Please include:

- the problem being solved
- the architectural surface being changed
- any new public API or behavior
- any docs that were updated alongside the change
- any follow-up work that remains intentionally out of scope

All submissions, including submissions by project members, require review via
GitHub pull requests.

## Contributor license agreement

Before we can use your code, you must sign the
[Google Individual Contributor License Agreement](https://cla.developers.google.com/about/google-individual).
The CLA is necessary mainly because you own the copyright to your changes, even
after your contribution becomes part of the codebase, so we need your
permission to use and distribute your code. You do not need to sign the CLA
until after a contribution has been reviewed and approved, but it must be in
place before the change can be merged.

Contributions made by corporations are covered by the
[Software Grant and Corporate Contributor License Agreement](https://cla.developers.google.com/about/google-corporate).
