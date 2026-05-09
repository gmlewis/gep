# ANALYSIS

## Bottom line

This repository is a **valuable prototype and exploration vehicle**, but it is **not yet close to the architecture needed for the goals you described**: a general-purpose, importable Gene Expression Programming platform that researchers could use for reinforcement learning, scientific discovery, engineering design, chemistry, physics, or other domain-heavy search problems.

Today it is best described as:

- a working **GEP proof of concept**
- a set of **narrow experiments**
- a partially modernized codebase that still carries **strong pre-generics design constraints**
- a repository with **real useful ideas**, but also **major architectural ceilings, missing abstractions, and several correctness/maintainability risks**

The biggest issue is not that the repo is "broken". The tests pass, and several experiments work. The real issue is that the codebase is still organized like a **research spike from an earlier Go era**, not like a **modern reusable scientific computing library**.

## What exists today

### 1. Core GEP implementation

There is a clear core stack:

- `functions/` defines available node catalogs and `FuncNode`/`FuncType` (`functions/functions.go`)
- `gene/` holds Karva symbols, constants, mutation, evaluation, and code rendering (`gene/gene.go`)
- `genome/` links genes, stores score, evaluates genomes, and writes generated code (`genome/genome.go`, `genome/write.go`)
- `model/` manages populations and evolution (`model/model.go`)
- `fitness/` contains reusable bool/float fitness helpers (`fitness/bool/fitness.go`, `fitness/float/fitness.go`)

That layering is a good starting point and should be preserved conceptually.

### 2. Grammar / code generation system

The repo can render evolved expressions into Go code through XML grammars (`grammars/grammars.go`, `gene/write.go`, `genome/write.go`).

This is an interesting capability and is worth keeping, but it is currently too tightly coupled to the core representation and too Go/XML-centric to serve as the backbone of a broad platform.

### 3. RL / environment exploration

There is an exploratory Gymnasium-inspired layer:

- a generic `Environment` interface (`gymnasium/gymnasium.go`)
- `common.Space` and `common.Obs` abstractions (`common/common.go`)
- one pure-Go environment: Blackjack (`gymnasium/envs/toy_text/blackjack/blackjack.go`)
- a GEP-based agent wrapper (`model/gymnasium-agents.go`)
- example programs under `examples/gymnasium/`

This is promising as a direction, but currently it is still toy-scale.

### 4. Evidence the repo does real work

The repo is not empty or purely aspirational:

- symbolic regression works (`experiments/symbolic_regression/main.go`)
- boolean problems work (`experiments/nand`, `experiments/odd-3-parity`, `experiments/odd-7-parity`)
- tests cover many core packages (`gene`, `genome`, `functions/*`, `fitness/*`, `model`)

So the right framing is **"promising prototype that needs major redesign"**, not **"abandonware with nothing usable"**.

## The biggest architectural constraints

## 1. The core is still built around pre-generics type splitting

This is the single biggest structural limitation.

The core abstraction is not `Gene[T]`, `Genome[T]`, `Node[T]`, or a typed IR. Instead it uses:

- `FuncType` enums (`functions/functions.go`)
- one `FuncNode` interface with methods for **all** supported types
- duplicated per-type execution paths for bool, int, float64, and vector-int

Examples:

- `functions.FuncNode` requires `BoolFunction`, `IntFunction`, `Float64Function`, and `VectorIntFunction` on the same interface (`functions/functions.go`)
- `gene/` has separate evaluators for bool/int/math/vector-int (`gene/bools.go`, `gene/ints.go`, `gene/math.go`, `gene/vector_ints.go`)
- `genome/` repeats the same pattern (`genome/bools.go`, `genome/ints.go`, `genome/math.go`)

This causes several platform-level problems:

- **duplication** across every supported value type
- **friction** adding new domains or data types
- **runtime type selection** where compile-time typing should exist
- APIs that are harder to understand, harder to validate, and easier to misuse

For the repo to become a serious platform, this part needs a **major overhaul**, likely into a generic typed core.

## 2. The internal representation is too stringly-typed

The core representation uses raw strings for almost everything:

- `Gene.Symbols []string` (`gene/gene.go`)
- terminals represented as `"d0"`, `"d1"`, ...
- constants represented as `"c0"`, `"c1"`, ...
- `Genome.LinkFunc string` (`genome/genome.go`)

This was a reasonable prototype choice, but it becomes a scaling problem when the library is meant to support many fields and many users:

- it weakens validation
- it pushes errors to runtime
- it makes refactoring and extension fragile
- it couples parsing, execution, mutation, and rendering around ad hoc string conventions

A platform for scientific use should move toward a **typed symbol/IR layer** with explicit node kinds, terminals, constants, and linking operators.

## 3. The public API is positional and prototype-shaped

The central constructor is still:

`model.New(fs, funcType, numIndividuals, headSize, numGenesPerGenome, numTerminals, numConstants, linkFunc, sf, debug)` (`model/model.go`)

That API is difficult to discover and easy to misuse:

- many ordered parameters
- mixed concerns in one call
- weak validation
- string-based `linkFunc`
- runtime `funcType`
- no builder/options/config object

A library meant to be "trivial to import" needs a very different shape:

- explicit config structs
- typed builders
- defaults
- validation with returned errors
- package-level concepts for experiments, domains, operators, and evaluation

## 4. Library code still uses fatal exits in places that should return errors

There are many `log.Fatalf` calls in library code, including parsing, evaluation, and model logic:

- `gene/gene.go`
- `gene/write.go`
- `model/model.go`

That is a major platform blocker.

A reusable library must not terminate the host process because a user provided a bad symbol, a mismatched function catalog, or an unsupported type. Those need to become typed errors.

## The evolution/search engine needs major work

## 1. The implemented evolutionary loop is incomplete

`Generation.Evolve` currently performs:

- best-genome preservation
- replication
- mutation

while several classic GEP operators are commented out (`model/model.go`):

- IS transposition
- RIS transposition
- gene transposition
- one-point recombination
- two-point recombination
- gene recombination

So the engine is not a full, robust, configurable GEP framework yet. It is a partial implementation.

For a serious platform, evolution operators should be:

- explicit
- configurable
- independently testable
- measurable
- interchangeable

## 2. Important behavior is hard-coded

Examples from `model/model.go`:

- stopping condition assumes success at `Score >= 1000`
- `getBest` initializes `bestScore := 0.0`
- mutation counts and crossover counts are randomly chosen with fixed formulas

Why this matters:

- a universal "perfect score is 1000" assumption does not generalize across domains
- `bestScore := 0.0` is wrong for fully negative fitness landscapes
- selection pressure and operator schedules are not first-class configuration

For research use, the engine should support:

- arbitrary optimization objectives
- pluggable stopping criteria
- minimization and maximization
- multi-objective search
- constraints and penalties
- experiment-specific operator schedules

## 3. Reproducibility is not designed in

Randomness is drawn directly from the global `math/rand` package throughout the core and RL layers:

- `model/model.go`
- `gene/gene.go`
- `gymnasium/envs/toy_text/blackjack/blackjack.go`

There is no clear experiment-level seed management or RNG injection.

That is a major issue for scientific and engineering work. A platform intended for discovery must make reproducibility easy:

- seeded runs
- reproducible operators
- logged configs
- versioned experiment metadata
- replayable runs

## 4. There are cache / duplication correctness risks

Several parts of the core cache compiled closures and symbol maps, but invalidation/copy semantics are unsafe.

### `Gene.Dup` copies cached compiled functions

`Gene.Dup` copies `bf`, `intF`, and `mf` from the source gene (`gene/gene.go`).

That is dangerous because those closures were built from the original gene state. A duplicate should either:

- rebuild lazily from copied symbols/constants
- or copy a typed immutable compiled form that is explicitly safe to share

As written, this creates a risk of duplicated genes evaluating through stale or source-bound closures.

### Mutation and crossover do not fully invalidate caches

- `Gene.Mutate` clears compiled funcs, but does not clear `SymbolMap` (`gene/gene.go`)
- `Genome.Mutate` does not clear `Genome.SymbolMap` (`genome/genome.go`)
- `singleCrossover` rewrites symbol slices directly and does not invalidate compiled caches or symbol maps (`model/model.go`)

This is a strong sign that the current implementation has outgrown its prototype architecture.

## 5. There are concrete correctness bugs / sharp edges

Examples:

- `getBest` mishandles all-negative scores because `bestScore` starts at `0.0` (`model/model.go`)
- `gene/math.go` and `gene/vector_ints.go` use `if symbolIndex > len(g.Symbols)` before indexing, which should be `>=`
- `gene/write.go` uses off-by-one bounds checks for terminals/constants
- `genome/write.go` dereferences `s.Symbol()` in the failure branch where `ok == false`

None of these mean the whole project is unusable, but together they show the engine is not yet hardened enough for broad reuse.

## The RL / environment stack is still a narrow experiment

## 1. Only one pure-Go environment exists

`gymnasium.Make` supports only:

- `"Blackjack-v1"` (`gymnasium/gymnasium.go`)

That is useful as a proof of concept, but it is nowhere near the breadth needed for the repo's long-term goals.

## 2. Space support exists on paper more than in practice

`common.Space` includes fields for:

- `Discrete`
- `Tuple`
- `MultiBinary`
- `MultiDiscrete`
- `Box`

But `model.NewGymnasiumAgents` actually implements only a narrow subset:

- action spaces: `Discrete`, `Tuple`
- observation spaces: `Discrete`, `Tuple`

and explicitly returns "not yet implemented" for other space types (`model/gymnasium-agents.go`).

That is a major mismatch between abstraction and implementation.

## 3. The RL policy representation is still int-centric and environment-specific

`GymnasiumAgents` creates integer-based genomes from observation space shape (`model/gymnasium-agents.go`) and `Genome.Evaluate` only supports populating `*int` or `*[]int` actions (`genome/genome.go`).

This blocks many realistic RL settings:

- continuous control
- structured actions
- probabilistic policies
- value functions
- actor-critic methods
- recurrent/stateful agents
- partially observed environments

Even within the current abstraction, `common.Obs.Unmarshal(dst interface{})` is too weak and too manual for a large platform (`common/common.go`).

## 4. The training/evolution loop is still toy-grade

In `GymnasiumAgents.Evolve`:

- replication is disabled with a comment noting diversity problems
- evolution uses only mutation + crossover + elitism
- reward is simply assigned onto `Genome.Score`
- all scores are reset to zero after each evolve

That may be fine for exploration, but it is not yet a serious RL research framework.

Missing pieces include:

- episodic statistics and logging
- evaluation/train splits
- curriculum support
- batched/vectorized environments
- proper experiment tracking
- policy/value abstractions
- alternate selection strategies
- reward normalization strategies
- safety around noisy fitness

## 5. The non-Go RL example is only a transport demo

The NATS example (`examples/gymnasium/toy_text/blackjack-nats/main.go`) is not an actual learning agent. It effectively hardcodes action publication and demonstrates message transport, not a reusable RL stack.

That means the RL story is still mostly aspirational.

## The repo is not yet packaged like a reusable scientific platform

## 1. Documentation is still experiment-first, not platform-first

The root README mostly presents:

- a short project description
- the NAND and symbolic regression demos
- a note that the repo is still experimental

It does **not** yet provide the material a downstream researcher would need:

- architectural overview
- package map
- API entry points
- extension guide
- domain modeling guide
- experiment reproducibility guide
- examples for embedding the library into another project
- roadmap / stability guarantees

For the stated ambition, documentation needs a full rewrite.

## 2. CI / repo automation is outdated

The repo still has:

- `.travis.yml` targeting Go `1.10.x`, `1.11.x`, and `tip`
- no GitHub Actions workflows under `.github/workflows/`

That is out of step with the current module (`go 1.24.0` in `go.mod`) and weakens confidence for outside users.

## 3. Contribution scaffolding is outdated

`CONTRIBUTING.md` still reads like an old Google open-source template with CLA guidance, but it does not explain:

- project direction
- review expectations
- coding conventions specific to this repo
- testing expectations
- how to add a new domain / function set / environment

That is fine for a dormant experiment, but not for a platform trying to attract serious contributions.

## 4. The repo is dominated by experiments and entrypoints

There are many `package main` programs across `experiments/`, `examples/`, benchmark generators, and historic slides/examples. In contrast, there is relatively little platform-level documentation around the reusable library surface.

This gives the repo the feel of:

- a lab notebook
- a set of demos
- a historical artifact

more than a polished library.

## 5. Test coverage is meaningful but not platform-grade

There are 17 `_test.go` files, and many core packages are tested. That is good.

But the tests are still mostly package-local behavior tests, not the sort of coverage needed for a platform:

- determinism/reproducibility tests
- API stability tests
- operator conformance tests
- large integration benchmarks
- domain extension tests
- failure-mode tests
- performance regression tests

## What should be preserved

These are the strongest foundations worth keeping.

### Preserve conceptually

- the **gene / genome / population** layering
- the Karva-expression-based representation
- the idea of **optional code generation**
- the existing boolean / symbolic-regression experiments as regression fixtures
- the basic environment abstraction direction

### Preserve as assets

- tests that prove current behavior
- example problem sets
- function catalogs as raw domain knowledge
- the pure-Go Blackjack environment as a reference example

## What should be overhauled

## 1. The typed core API

Move from:

- `FuncType`
- `FuncNode` with all-type methods
- stringly symbols

to something like:

- `Node[T]`
- `Gene[T]`
- `Genome[T]`
- typed terminals/constants/symbol definitions

## 2. The experiment / evolution configuration surface

Replace long positional constructors with:

- config structs
- options
- explicit operator registries
- explicit stopping criteria
- injected RNG
- explicit evaluation mode

## 3. Error handling

Replace `log.Fatalf`-style library exits with returned errors and validation.

## 4. The RL abstraction layer

Keep the direction, but redesign it around:

- typed observations/actions
- broader space support
- environment seeding
- batched/vectorized execution
- clear policy/evaluator interfaces

## 5. Documentation and examples

Rewrite the root docs around:

- what the library is
- what is stable
- how to model a new problem
- how to extend operators/functions
- how to run reproducible experiments

## What likely needs full rewrite

## 1. The duplicated per-type evaluator stack

The bool/int/float/vector-int split across `functions`, `gene`, and `genome` is the clearest place where a rewrite will pay off.

## 2. The evolution engine

The current implementation is too hard-coded, partially implemented, and cache-fragile to serve as the long-term core.

It is a good reference for behavior, but probably not the right long-term implementation.

## 3. The grammar boundary

The current XML grammar model is powerful but oversized relative to what the generator actually uses. It should be reconsidered as an **optional backend** behind a smaller typed expression/codegen IR.

## Missing capabilities required for the long-term vision

If the goal is "researchers can import this and use GEP for discovery across many fields", the repo is currently missing at least the following classes of capability:

### Scientific platform requirements

- reproducible experiment definitions
- metrics / run logging / checkpoints
- benchmark suites
- comparative operator studies
- configuration serialization
- dataset/problem registries
- deterministic replay

### Library ergonomics

- stable public API
- typed extension points
- examples for embedding in another application
- better package boundaries
- clearer versioning/stability story

### Search/evolution features

- complete GEP operator support
- pluggable selection strategies
- constraints
- multi-objective optimization
- novelty/diversity support
- island/population strategies
- parallel/distributed evaluation

### Domain modeling features

- richer terminal/constant abstractions
- typed domain primitives
- graph/tree/sequence support beyond scalar primitives
- structured outputs
- problem-specific constraints and repair

### RL features

- more environments
- proper continuous spaces
- policy/value abstractions
- vectorized environments
- training diagnostics
- hybrid evolutionary/RL strategies

### Production/research tooling

- modern CI
- profiling/benchmark harnesses
- compatibility guarantees
- stronger failure handling
- clearer observability

## Recommended modernization path

## Phase 1: Stabilize the current core

Before any major redesign:

- fix the concrete correctness bugs
- remove fatal exits from library code
- fix cache invalidation / duplication semantics
- make score handling and stopping criteria configurable
- add deterministic seeding

This would make the current prototype safer to evolve.

## Phase 2: Introduce a typed generic core

Build a new core around:

- typed nodes
- typed terminals/constants
- typed genomes
- typed link operators

Keep the current code as a compatibility/reference layer if needed, but do not keep extending the old `FuncType` architecture.

## Phase 3: Rebuild the evolution engine as a first-class subsystem

Design explicit packages for:

- selection
- mutation
- recombination
- transposition
- termination
- evaluation
- statistics

Each should be configurable and testable.

## Phase 4: Split platform concerns cleanly

Separate the project into clearer subsystems:

- `core` / typed expression engine
- `evolution`
- `problems` or `domains`
- `codegen`
- `env` / RL
- `experiments`

Right now those concerns are too interwoven.

## Phase 5: Reposition the repo as a platform

Once the architecture is modernized:

- rewrite the README
- add a package map and extension guide
- document the stable surface
- provide "import this library" examples
- add modern CI and benchmark automation

## Phase 1-4 implementation audit (2026-05-09)

`ANALYSIS.md` is the source of truth for the intended meaning of Phases 1-4.
The notes below evaluate the current repository against that original intent
without reinterpreting the phases.

### Phase 1 status: complete (after P1-A and P1-B)

What appears complete:

- deterministic seeding exists in both the typed and legacy evolution paths
- score orientation and stopping behavior are configurable
- cache invalidation / duplication semantics were materially improved

P1-A completion evidence:

- direct library `log.Printf` / `fmt.Printf` emissions were removed from legacy
  `gene`, `genome`, and `model` packages
- explicit error-return surfaces were added for legacy invalid states:
  - `gene`: `New`, `String`, `SymbolCount`, `Mutate`, `Dup`, `EvalInt`,
    `EvalMath`, `EvalBool`, `EvalVectorInt`, `AllSymbolsEqualWeights`
  - `genome`: `EvalInt`, `EvalMath`, `EvalBool`, `EvalIntTuple`, `Dup`,
    `EvaluateWithScore`
  - `model`: `Evolve`, `singleCrossover`, `maxArity`
- legacy codegen now has an explicit error-return surface:
  `(*genome.Genome).Write(...) error`
- focused tests now prove no direct stdout/stderr emission on representative
  legacy error/stop paths:
  - `gene.TestNew_InvalidSymbolIndexesDoNotWriteStderr`
  - `genome.TestWrite_MissingLinkFunctionDoesNotWriteStdout`
  - `model.TestEvolve_DoesNotWriteStdoutOnStop`
- legacy index-based `for` loops were modernized to `range` where semantics are
  equivalent across the legacy packages (`gene`, `genome`, `model`)
- full repo validation passed via `scripts/test-all.sh`

### Phase 2 status: complete (after P2-A)

What appears complete:

- a typed generic `core` package exists
- a typed generic `evolution` engine exists and is used as the default
  execution surface for the experiment entrypoints
- compatibility conversion code is isolated behind explicitly named adapter
  surfaces (`gene/core_bridge.go`, `genome/core_bridge.go`)
- legacy `gene`, `genome`, and `model` packages are explicitly marked as
  deprecated compatibility/reference layers with migration guidance

Representative locations:

- `core/core.go`
- `gene/core_bridge.go`
- `genome/core_bridge.go`
- `model/doc.go`
- `experiments/*`
- `examples/gymnasium/*`

### Phase 3 status: complete (after P3-A)

What appears complete:

- explicit subsystems exist for selection, mutation, recombination,
  transposition, termination, evaluation, and statistics
- those subsystems are integrated into the typed `evolution` engine

Representative locations:

- `evolution/evolution.go`
- `evolution/mutation/mutation.go`
- `evolution/transposition/transposition.go`

### Phase 4 status: partially complete (after P4-A, P4-B, P4-C)

The current repository partially satisfies the original Phase 4 intent.
The work completed under "Phase 4" has progressed from pure import-boundary
enforcement to actual subsystem separation.

What has been completed:

- dedicated `codegen` subsystem separates code generation from representation
  (P4-A)
- dedicated `env` / RL subsystem separates Gymnasium agent orchestration from
  legacy `model` (P4-B)
- dedicated `problems` subsystem promotes reusable fitness/problem definitions
  into typed seams over `core` + `evolution` (P4-C)

What is still missing relative to the original plan:

- experiments and examples still depend primarily on legacy representation and
  legacy model packages rather than the typed `core` + `evolution` stack
- legacy packages are still active workflow packages, not merely
  compatibility/reference layers

Representative locations:

- `experiments/*`
- `examples/gymnasium/*`
- `evolution/package_map_test.go`

## Corrective milestone backlog for Phase 1-4 deviations

These milestones are intended to correct the current deviations from the
original plan. They are deliberately written so each milestone can be used as
the source prompt for a separate `/delegate` PR.

### P1-A: Finish legacy stabilization and explicit error surfaces

Goal:

- complete the remaining stabilization work on the legacy path

Required outcome:

- replace library `log.Printf` / `fmt.Printf` error handling in legacy
  `gene`, `genome`, `model`, and legacy codegen paths with explicit error
  returns or narrow compatibility wrappers
- eliminate direct stdout/stderr emission from library packages
- preserve existing behavior only where it is intentionally CLI/example-facing

Status: ✅ 100% complete (2026-05-09)

Completion proof:

- no direct runtime `log.Printf` / `fmt.Printf` emissions remain in legacy
  library code for:
  - `gene/*.go`
  - `genome/*.go`
  - `model/*.go`
- `genome/write.go` now exposes `Write(...) error` as the explicit
  codegen error surface
- explicit error-return APIs were added across legacy packages:
  - `gene`: `New`, `String`, `SymbolCount`, `Mutate`, `Dup`, `Eval*`,
    `AllSymbolsEqualWeights`
  - `genome`: `Eval*`, `EvalIntTuple`, `Dup`, `EvaluateWithScore`
  - `model`: `Evolve(...) (*genome.Genome, error)`, `singleCrossover`,
    `maxArity`
- added regression tests:
  - `gene/gene_test.go`: `TestNew_InvalidSymbolIndexesDoNotWriteStderr`
  - `gene/gene_test.go`: `TestNew_InvalidSymbolIndexesReturnsError`,
    `TestAllSymbolsEqualWeights_UnknownFuncTypeReturnsError`,
    `TestEvalInt_UnknownSymbolReturnsError`,
    `TestMutate_TooFewChoicesReturnsError`
  - `genome/genome_test.go`:
    `TestEvaluateWithScore_NilScoringFuncReturnsError`,
    `TestEvalInt_MissingLinkFunctionReturnsError`
  - `genome/write_test.go`:
    `TestWrite_MissingLinkFunctionReturnsError`,
    `TestWrite_MissingLinkFunctionDoesNotWriteStdout`
  - `model/model_test.go`:
    `TestEvolve_DoesNotWriteStdoutOnStop`,
    `TestMaxArity_UnknownFuncTypeReturnsError`,
    `TestGetBest_NilScoringFuncReturnsError`,
    `TestSingleCrossover_MismatchedGenesReturnsError`
- verification commands passed:
  - `go test ./gene ./genome ./model`
  - `./scripts/test-all.sh`

### P1-B: Modernize legacy-style `for` loops to `range` where equivalent

Goal:

- improve readability by replacing legacy index-based loop forms with modern Go
  `range` loops where behavior remains equivalent

Required outcome:

- audit repository packages for legacy-style `for` loops that can be safely
  rewritten as `range` loops without changing semantics
- convert those loops to `range` syntax in a behavior-preserving way
- keep explicit index/counter loops only where they are required for semantics
  (for example, custom stepping, reverse traversal, or mutation during index
  iteration)

Status: ✅ 100% complete (2026-05-09)

Completion proof:

- completed a legacy package sweep (`gene`, `genome`, `model`) and converted
  equivalent index-based loops to idiomatic `range` loops in:
  - `gene/gene.go`
  - `gene/write.go`
  - `genome/bools.go`
  - `genome/genome.go`
  - `genome/ints.go`
  - `genome/math.go`
  - `model/model.go`
- retained explicit index/counter loops only where semantics are intentionally
  index-driven (for example benchmark counters and random/index-position
  mutation/crossover logic)
- verification commands passed:
  - `go test ./gene ./genome ./model`
  - `./scripts/test-all.sh`

### P2-A: Make typed core/evolution the default execution surface

Goal:

- make `core` + `evolution` the primary public execution path and reduce the
  legacy stack to compatibility/reference status

Required outcome:

- add typed entrypoints where the repo still depends on legacy `model`
- move compatibility conversion code behind clearly named adapter surfaces
- stop adding new workflow code that depends directly on `functions.FuncType`,
  legacy `gene`, legacy `genome`, or legacy `model`
- mark any remaining legacy packages as `deprecated` to the package-level documentation and document how the new packages should be used in their places

This milestone corrects the remaining Phase 2 deviation.

Status: ✅ 100% complete (2026-05-09)

Completion proof:

- migrated the experiment workflow entrypoints from legacy `model.New` to typed
  `evolution.New`:
  - `experiments/nand/main.go`
  - `experiments/odd-3-parity/main.go`
  - `experiments/odd-7-parity/main.go`
  - `experiments/6-multiplexer/main.go`
  - `experiments/symbolic_regression/main.go`
- updated experiment scoring paths to evaluate typed `core.Genome[T]` directly
  and use explicit compatibility adapters only at the codegen boundary
  (`genome.NewFromCoreBool` / `genome.NewFromCoreFloat64`) when writing legacy
  grammar output
- preserved compatibility conversion behind clearly named adapter surfaces:
  - `gene/core_bridge.go`
  - `genome/core_bridge.go`
- marked remaining legacy packages as deprecated with package-level migration
  guidance:
  - `gene/doc.go`
  - `genome/doc.go`
  - `model/doc.go`
- updated typed benchmark coverage in
  `experiments/symbolic_regression/main_test.go` to construct populations via
  `evolution.New`
- verification commands passed:
  - `go test ./experiments/symbolic_regression`
  - `go test ./experiments/nand ./experiments/odd-3-parity ./experiments/odd-7-parity ./experiments/6-multiplexer`
  - `./scripts/test-all.sh`

### P3-A: Extract transposition into a first-class evolution subsystem

Goal:

- create an explicit `evolution/transposition` subsystem

Required outcome:

- move IS, RIS, and gene transposition orchestration out of
  `evolution/mutation`
- give `evolution.Generation` a distinct transposition configuration/invocation
  stage
- keep operator tests focused on behavior and validity of genomes
- continue prioritizing excellent end-user godoc-style documentation for all exported structs and functions

This milestone corrects the remaining Phase 3 deviation.

Status: ✅ 100% complete (2026-05-09)

Completion proof:

- extracted typed transposition orchestration into a dedicated
  `evolution/transposition` subsystem with package-level godoc and exported
  `Config` / `Apply` API:
  - `evolution/transposition/transposition.go`
- removed IS, RIS, and gene transposition orchestration from
  `evolution/mutation`, leaving that package focused on point mutation and
  inversion:
  - `evolution/mutation/mutation.go`
- gave `evolution.Generation` a distinct transposition configuration and
  invocation stage via `TranspositionConfig` and `Transpose`, and updated the
  typed evolution loop to run `select -> recombine -> mutate -> transpose`:
  - `evolution/evolution.go`
- kept operator tests focused on behavior and genome validity by moving
  transposition-specific operator coverage into the new subsystem tests and
  leaving `evolution` tests focused on stage wiring/integration:
  - `evolution/transposition/transposition_test.go`
  - `evolution/mutation/mutation_test.go`
  - `evolution/evolution_test.go`
- verification commands passed:
  - `go test ./evolution/...`
  - `./scripts/test-all.sh`

### P4-A: Create a dedicated codegen subsystem

Goal:

- separate code generation from representation packages

Required outcome:

- move code generation orchestration out of `gene/write.go`,
  `genome/write.go`, and ad hoc grammar coupling into a dedicated `codegen`
  package (or small `codegen/*` package tree)
- keep legacy `gene` / `genome` write methods as thin adapters only, if they
  remain at all and remove them if they no longer remain
- continue prioritizing excellent end-user godoc-style documentation for all exported structs and functions

Status: ✅ 100% complete (2026-05-09)

Completion proof:

- added a dedicated `codegen` subsystem with package-level godoc and exported
  `Expressor`, `Program`, `Expression`, `Generate`, and `Write` APIs:
  - `codegen/doc.go`
  - `codegen/codegen.go`
- moved legacy gene-expression rendering orchestration out of `gene/write.go`
  into `codegen.Expression`, leaving `(*gene.Gene).Expression` as a thin
  adapter that supplies legacy gene data and arg-order metadata:
  - `gene/write.go`
- moved legacy full-program rendering orchestration out of `genome/write.go`
  into `codegen.Generate` / `codegen.Write`, leaving `(*genome.Genome).Write`
  as a thin adapter over the dedicated subsystem:
  - `genome/write.go`
- added dedicated subsystem regression coverage and preserved legacy adapter
  coverage:
  - `codegen/codegen_test.go`
  - `gene/write_test.go`
  - `genome/write_test.go`
- updated Phase 4 package-boundary coverage so the new dedicated `codegen`
  subsystem is part of the inspected architecture and legacy representation
  packages are allowed to depend on it as thin adapters:
  - `evolution/package_map_test.go`
- verification commands passed:
  - `go test ./codegen ./gene ./genome ./evolution`
  - `./scripts/test-all.sh`

### P4-B: Create a dedicated env / RL subsystem

Goal:

- separate environment definitions from agent/training orchestration

Required outcome:

- keep `gymnasium` focused on environments and environment metadata
- move Gymnasium agent evolution/training logic out of
  `model/gymnasium-agents.go` into a dedicated `env` / RL package
- ensure RL flows no longer depend on legacy `model` as their primary engine
- continue prioritizing excellent end-user godoc-style documentation for all exported structs and functions

Status: ✅ 100% complete (2026-05-09)

Completion proof:

- added a dedicated `env` subsystem with package-level godoc and exported
  `GymnasiumAgents`, `GymnasiumAgentsOption`, `NewGymnasiumAgents`,
  `WithAppendEpisodeSteps`, `WithDebug`, `WithHeadSize`, `WithNumConstants`,
  and `WithNumIndividuals` APIs:
  - `env/doc.go`
  - `env/env.go`
- moved all Gymnasium agent evolution/training orchestration out of
  `model/gymnasium-agents.go` into `env/env.go`; the `env` package imports
  only `common`, `functions`, `gene`, `genome`, and standard library —
  **it does not import `model`**
- added `env/env_test.go` with full coverage of `NewGymnasiumAgents`,
  `EvaluateAgent`, `RewardAgent`, `Evolve`, and `processObservations`
- added deprecation notice to `model.GymnasiumAgents` directing consumers to
  the new `env` package
- migrated `examples/gymnasium/toy_text/blackjack-go/main.go` from
  `model.GymnasiumAgents` to `env.GymnasiumAgents`
- added `TestPhase4MilestoneB_EnvPackageBoundaries` to
  `evolution/package_map_test.go` that enforces the env package boundary
  (no `model` import allowed)
- `gymnasium` package remains unchanged and focused solely on environment
  definitions and metadata

Verification commands:

  ```
  go test ./env/...
  go test ./evolution/...   # TestPhase4MilestoneB_EnvPackageBoundaries passes
  go build ./examples/gymnasium/toy_text/blackjack-go/...
  ./scripts/test-all.sh
  ```

### P4-C: Create a dedicated problems / domains subsystem

Goal:

- separate reusable problem definitions from ad hoc experiments

Required outcome:

- promote reusable fitness/problem definitions into a dedicated `problems`
  or `domains` package tree
- keep `fitness` only as a narrow helper layer, or fold it into the new
  subsystem where appropriate
- define typed problem-facing seams over `core` + `evolution`
- continue prioritizing excellent end-user godoc-style documentation for all exported structs and functions

Status: ✅ 100% complete (2026-05-09)

Completion proof:

- added a dedicated `problems` subsystem with package-level godoc and exported
  `Case`, `BoolProblem`, `FloatProblem`, `NewBoolProblem`, `NewFloatProblem`,
  `BoolProblem.NumHitsScoringFunc`, `FloatProblem.NumHitsAbsScoringFunc`,
  `FloatProblem.MeanSquaredErrorAbsScoringFunc`, and
  `FloatProblem.RSquareScoringFunc` APIs:
  - `problems/doc.go`
  - `problems/problems.go`
- `problems` imports only `core` and `fitness/*`; it does not import any legacy
  `gene`, `genome`, or `model` packages, satisfying the typed seam requirement
- `fitness/bool` and `fitness/float` remain unchanged as narrow helper layers
  providing raw slice-based fitness math; `problems` bridges from typed
  `core.Genome[T]` evaluation to those fitness helpers
- scoring functions return `func(core.Genome[T]) float64`, the exact type
  accepted by `evolution.New` as its `scoringFunc` argument — no additional
  adaptation is required at call sites
- package boundary enforcement added to `evolution/package_map_test.go`:
  - `TestPhase4MilestoneC_ProblemsPackageBoundaries` passes
  - `allowedImportPrefixes` extended with a `problems` case
- focused tests added in `problems/problems_test.go` covering constructor
  validation, all three boolean scoring paths, and all three float scoring paths
- verification commands passed:
  ```
  go test ./problems/...
  go test ./evolution/...   # TestPhase4MilestoneC_ProblemsPackageBoundaries passes
  ./scripts/test-all.sh
  ```

### P4-D: Rewire experiments and examples to the separated subsystems

Goal:

- make experiments/examples consume the intended architecture rather than the
  legacy prototype architecture

Required outcome:

- migrate experiments/examples to import `core`, `evolution`, `codegen`,
  `env`, and/or `problems`
- stop using legacy `model`, `gene`, and `genome` as the default workflow
  path for new or migrated examples and completely delete all unused legacy packages
- continue prioritizing excellent end-user godoc-style documentation for all exported structs and functions

Dependencies:

- depends on `P2-A`, `P3-A`, `P4-A`, `P4-B`, and `P4-C`

### P4-E: Retire Phase 4 import-boundary enforcement as the definition of success

Goal:

- stop treating import-boundary tests as the implementation of Phase 4

Required outcome:

- remove or greatly reduce `evolution/package_map_test.go` assertions that do
  not prove user-visible architectural behavior
- replace them with integration tests that exercise the real separated seams:
  typed evolution, codegen, env/RL, and problem/domain wiring
- continue prioritizing excellent end-user godoc-style documentation for all exported structs and functions

Dependencies:

- depends on `P4-A`, `P4-B`, `P4-C`, and `P4-D`

## Final assessment

This repo contains **real intellectual value** and should not be dismissed. The core ideas, the problem experiments, and the codegen direction are all worth preserving.

But to achieve the goal of becoming a broadly useful GEP platform for scientific and engineering discovery, the repo needs:

- **substantial architectural overhaul**
- **selective rewrites of major subsystems**
- **a modern typed API**
- **a rebuilt evolution engine**
- **a much broader and better-designed RL/domain layer**

In short:

> **Keep the ideas, the tests, the example problems, and the overall decomposition. Rewrite the platform around a modern typed core and treat the current codebase as the prototype/specimen, not the final foundation.**
