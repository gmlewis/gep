# ANALYSIS

## Bottom line

The architectural rewrite is complete. This repository is now best understood as
a **typed Gene Expression Programming engine** with optional code generation and
an exploratory environment layer, not as a one-off collection of experiments.

That means the next chapter should not be another internal cleanup phase. The
next chapter should be **applied scientific design**:

1. evolve novel electronic circuits against simulator-backed specifications
2. evolve voxel and implicit-geometry designs against computational engineering
   constraints
3. use the same search core for other invention problems where simulator-backed
   RL and evolutionary search work well together

This document describes the concrete steps needed to do that.

## What the repo already gives us

The core platform pieces are now in place:

- `core` supplies typed `Node[T]`, `Gene[T]`, `Genome[T]`, `Catalog[T]`, and
  `LinkOperator[T]`
- `evolution` supplies seeded, typed population search with configurable
  mutation, recombination, transposition, selection, statistics, and
  termination
- `problems` supplies reusable typed scoring seams for common boolean and float
  tasks
- `codegen` can render evolved Karva expressions into textual artifacts through
  grammars
- `env` demonstrates environment/agent orchestration, but is still oriented
  around the current discrete/tuple Gymnasium path rather than a full modern RL
  stack

So the right mental model is:

> **Use `core` + `evolution` as the search engine. Build domain-specific
> simulator, constraint, artifact, and evaluation layers around it.**

## What still has to be built for practical discovery work

The current engine is strong enough to be the search core, but practical
scientific design work still needs a new application-facing layer.

### 1. Add a domain-facing design loop, not just experiments

For real design tasks, create a new package family around concepts like:

- `Candidate` or `Artifact`
- `Decoder` from genome to domain object
- `Constraint`
- `Evaluator`
- `Scenario` or `Spec`
- `RunManifest`

Do **not** force advanced scientific design through the current
`env.GymnasiumAgents` path. That package is a useful reference, but serious
applied work needs richer typed state, artifact, and simulator interfaces.

### 2. Treat modern RL best practices as mandatory

Whether the outer loop is called RL, evolutionary RL, black-box policy search,
or simulator-backed optimization, the practical workflow should include:

- seeded and replayable runs
- vectorized or batched evaluation
- train/validation/test splits over scenarios, operating corners, or load cases
- reward decomposition into hard constraints plus soft objectives
- curriculum from easier to harder design targets
- novelty/diversity pressure so the search does not collapse too early
- checkpointing and experiment manifests
- held-out verification before promoting a design
- surrogate modeling when the simulator is expensive
- human review gates before fabrication or lab execution

### 3. Make output artifacts first-class

For real use, the engine should not stop at a fitness score. Each domain package
should emit artifacts that another tool can consume:

- circuit netlists, schematics, or Verilog
- voxel grids, implicit geometry scripts, or meshes
- controller code, parameter files, or manufacturing manifests

The current `codegen` package is excellent for textual outputs. When the target
artifact is geometric or graph-structured, use a domain-specific emitter in Go
and treat `codegen` as optional.

## Recommended cross-cutting roadmap

Before tackling any one scientific domain, add these platform pieces.

Each milestone below is intentionally written so it can be used as the source
prompt for a standalone `/delegate` PR. The milestones are designed for
**sequential execution only**. Do not run them in parallel.

### Phase A: Build the applied-design substrate

#### `PA-01`: Create the shared `design` package and run-manifest schema

Goal:

- establish the package tree and durable metadata model that all later applied
  design milestones will build on

Dependencies:

- none

Required outcome:

- add a new top-level `design` package with package-level godoc explaining its
  role as the shared applied-design substrate
- define serializable types for at least:
  - `RunManifest`
  - `RunConfig`
  - `ArtifactRef`
  - `ScenarioSplitSummary`
  - `SeedRecord`
- add helpers for JSON round-trip loading/writing of run manifests
- add `testdata/` manifest fixtures that exercise both minimal and populated
  manifests
- do **not** add evaluator, constraint, novelty, or checkpoint logic yet

Mechanically verifiable completion:

- `go test ./design/...` passes
- tests prove manifest JSON round-trips without losing seed, config, artifact,
  or scenario-split metadata

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added: `design` (`design/doc.go`, `design/manifest.go`)
- required schema types added:
  `RunManifest`, `RunConfig`, `ArtifactRef`, `ScenarioSplitSummary`,
  `SeedRecord`
- JSON manifest helpers added:
  `LoadRunManifest`, `LoadRunManifestFile`, `WriteRunManifest`,
  `WriteRunManifestFile`
- fixtures added:
  `design/testdata/manifest_minimal.json`,
  `design/testdata/manifest_populated.json`
- verification tests added:
  `design/manifest_test.go` (in-memory and file round-trip coverage for
  minimal and populated manifests, including seed/config/artifact/scenario
  metadata retention)
- command proof:
  `go test ./design/...`

#### `PA-02`: Add batched evaluation worker abstractions in `design/eval`

Goal:

- create the concurrency-safe evaluation substrate used by future circuit,
  voxel, and control pilots

Dependencies:

- depends on `PA-01`

Required outcome:

- add a new `design/eval` package with package-level godoc
- define typed request/result abstractions for batched evaluation, including at
  least:
  - candidate identity
  - scenario identity
  - batch request
  - batch result
  - evaluator interface
  - runner/dispatcher abstraction
- support context cancellation and bounded worker counts
- include a deterministic fake evaluator in tests so concurrency behavior is
  verifiable without external simulators
- do **not** integrate with `core`/`evolution` yet beyond generic type-safe
  request/result plumbing

Mechanically verifiable completion:

- `go test ./design/eval/...` passes
- tests prove batched execution preserves all requests, respects worker-count
  limits, and returns deterministic results for a seeded fake evaluator

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added: `design/eval` (`design/eval/doc.go`, `design/eval/eval.go`)
- required abstractions added:
  `CandidateID`, `ScenarioID`, `BatchRequestItem`, `BatchResultItem`,
  `BatchRequest`, `BatchResult`, `Evaluator`, `Runner`, `WorkerRunner`
- dispatcher behavior implemented:
  context cancellation support, bounded worker pool execution, request/result
  identity propagation
- deterministic fake evaluator coverage added:
  `design/eval/eval_test.go`
- verification tests added:
  request preservation and deterministic results for seeded fake evaluator,
  worker-limit enforcement, context-cancellation handling, and worker-count
  validation
- command proof:
  `go test ./design/eval/...`

#### `PA-03`: Add constraint and validation reporting in `design/constraints`

Goal:

- make constraints explicit, composable, and testable before any domain pilots
  are introduced

Dependencies:

- depends on `PA-01`

Required outcome:

- add a new `design/constraints` package with package-level godoc
- define typed constraint interfaces and result types that support the three
  required behaviors from this roadmap:
  - reject a candidate
  - repair a candidate
  - penalize a candidate
- define a `ValidationReport`-style aggregate that records:
  - which constraints ran
  - what decisions they made
  - what repairs or penalties were applied
- add deterministic tests for:
  - stop-on-reject behavior
  - repair chaining
  - penalty accumulation

Mechanically verifiable completion:

- `go test ./design/constraints/...` passes
- tests prove the aggregate report is stable and deterministic for a fixed input

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added: `design/constraints` (`design/constraints/doc.go`,
  `design/constraints/constraints.go`)
- required types and interfaces added:
  `Decision` (Pass/Repair/Penalize/Reject), `ConstraintResult`,
  `Constraint[T]` interface, `ValidationReport[T]`
- `Validate[T]` function implements stop-on-reject, repair chaining, and
  penalty accumulation
- deterministic tests added:
  `design/constraints/constraints_test.go` (stop-on-reject, repair chaining,
  penalty accumulation, reject-skips-pending-penalties, repair-then-reject,
  empty constraints, determinism for fixed input, Decision.String coverage)
- command proof:
  `go test ./design/constraints/...`

#### `PA-04`: Add shared multi-objective scoring in `design/objectives`

Goal:

- provide a reusable way to combine hard constraints and soft objectives into a
  single score breakdown that later pilots can share

Dependencies:

- depends on `PA-03`

Required outcome:

- add a new `design/objectives` package with package-level godoc
- define serializable types for at least:
  - objective definition
  - weighted score contribution
  - score breakdown
  - aggregate score result
- support at least:
  - weighted aggregation for soft objectives
  - explicit hard-fail gating from validation results
  - deterministic score ordering for ties
- add tests that show a candidate with the same raw objective values always gets
  the same aggregate score and breakdown ordering

Mechanically verifiable completion:

- `go test ./design/objectives/...` passes
- tests prove hard failures dominate soft-objective aggregation

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added: `design/objectives` (`design/objectives/doc.go`,
  `design/objectives/objectives.go`)
- required types added:
  `ObjectiveKind` (Soft/Hard), `ObjectiveDef`, `WeightedContribution`,
  `ScoreBreakdown`, `AggregateResult`
- `Score` function implements weighted soft aggregation, hard-fail gating from
  `rejected` flag (e.g. `constraints.ValidationReport.Rejected`), and
  hard-fail gating from Hard-kind objectives with zero/negative raw scores
- `Less` helper provides deterministic tie-breaking comparator for sorting
  slices of `AggregateResult`; contributions are recorded in definition order
  so identical inputs always yield identical breakdowns
- deterministic tests added:
  `design/objectives/objectives_test.go` (weighted aggregation, hard-fail from
  rejected, hard-fail from zero/negative Hard objective, hard-fail dominates
  soft objectives, penalty reduces aggregate, determinism for fixed input,
  contributions in definition order, empty defs, missing raw score defaults to
  zero, Less ordering, Less tie-breaking, Less equal candidates, Less usable
  for sort, ObjectiveKind.String coverage)
- command proof:
  `go test ./design/objectives/...`

#### `PA-05`: Add novelty archive support in `design/novelty`

Goal:

- introduce reusable novelty/diversity infrastructure before domain pilots start
  optimizing toward a single objective

Dependencies:

- depends on `PA-04`

Required outcome:

- add a new `design/novelty` package with package-level godoc
- define types for at least:
  - behavior vector
  - novelty archive entry
  - distance function contract
  - archive configuration
  - novelty score result
- implement a deterministic k-nearest-neighbor novelty calculation over stored
  behavior vectors
- include tests covering:
  - empty archive behavior
  - archive insertion
  - stable novelty scores for a fixed archive and fixed query vector

Mechanically verifiable completion:

- `go test ./design/novelty/...` passes
- tests prove novelty scores are deterministic and archive growth is correct

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added: `design/novelty` (`design/novelty/doc.go`,
  `design/novelty/novelty.go`)
- required types added:
  `BehaviorVector`, `ArchiveEntry`, `DistanceFunc`, `ArchiveConfig`,
  `NoveltyResult`, `Archive`
- built-in distance functions added: `SquaredEuclidean` (default),
  `EuclideanDistance`; `ArchiveConfig.Distance` accepts any custom
  `DistanceFunc`
- `Archive` implements `Add` (with optional `MaxSize` cap), `Len`, `Entries`
  (returns a defensive copy), and `Score` (pure k-NN mean-distance function)
- `Score` is deterministic: same archive contents + same query always produce
  identical `NoveltyResult`; neighbor distances are returned in ascending order
- tests added: `design/novelty/novelty_test.go`
  (SquaredEuclidean same vector / known value / symmetric / mismatched lengths,
  EuclideanDistance, NewArchive defaults / default K, empty archive score,
  Add increases Len, Add respects MaxSize, MaxSize zero means unbounded,
  Entries returns copy, k=1 score, k-nearest mean, fewer-than-K entries,
  neighbor distances ascending, deterministic for fixed archive, Score does not
  mutate archive, repeated Score calls stable, custom distance func, archive
  growth correctness)
- command proof:
  `go test ./design/novelty/...`

#### `PA-06`: Add checkpoint and replay support in `design/checkpoint`

Goal:

- persist enough state to pause, inspect, and replay applied-design runs without
  re-inventing ad hoc storage in each future domain

Dependencies:

- depends on `PA-01`
- depends on `PA-02`
- depends on `PA-04`
- depends on `PA-05`

Required outcome:

- add a new `design/checkpoint` package with package-level godoc
- define serializable snapshot types that combine:
  - the `RunManifest`
  - elite candidate metadata
  - aggregate scores and breakdowns
  - novelty/archive state
  - artifact references
- implement save/load helpers with versioned on-disk JSON snapshots
- add replay helpers that can restore checkpoint metadata for a later run even
  if the evaluator itself is stubbed or fake in tests
- add integration tests that save a checkpoint, reload it, and compare the
  restored snapshot to the original

Mechanically verifiable completion:

- `go test ./design/checkpoint/... ./design/...` passes
- tests prove checkpoint round-trips preserve manifest, elite metadata, scores,
  and novelty state

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added: `design/checkpoint` (`design/checkpoint/doc.go`,
  `design/checkpoint/checkpoint.go`)
- required types added:
  `EliteRecord`, `AggregateSnapshot`, `NoveltySnapshot`, `Snapshot`
- `Snapshot` combines: `design.RunManifest`, `[]EliteRecord`,
  `AggregateSnapshot`, `NoveltySnapshot`, `[]design.ArtifactRef`
- save/load helpers added: `Save`/`SaveFile` (write versioned indented JSON),
  `Load`/`LoadFile` (decode and validate schema version)
- schema versioning: `SchemaVersion` field stamped on every save; unknown
  versions are rejected with a descriptive error
- replay helpers added: `ReplayManifest`, `ReplayElites`,
  `ReplayNoveltyEntries`; each returns a defensive copy and accepts a nil
  snapshot with a clear error
- integration tests added: `design/checkpoint/checkpoint_test.go`
  (Save/Load buffer round-trip, schema version stamped on Save,
  nil-writer / nil-snapshot / nil-reader error cases, unknown version error,
  two-object stream rejection, file round-trip, invalid file path errors,
  ReplayManifest / ReplayElites / ReplayNoveltyEntries with nil and copy
  checks, novelty archive re-seeding from replay, deterministic Save output,
  empty-elites and empty-novelty round-trip)
- command proof:
  `go test ./design/checkpoint/... ./design/...`

### Phase B: Add artifact emitters and scenario sets

#### `PB-01`: Add shared scenario-set support in `design/scenarios`

Goal:

- standardize train/validation/test scenario handling before any domain pilot
  invents its own incompatible scenario format

Dependencies:

- depends on `PA-01`

Required outcome:

- add a new `design/scenarios` package with package-level godoc
- define serializable types for:
  - `ScenarioID`
  - `ScenarioSet`
  - `ScenarioSplit`
  - `ScenarioRegistry`
- support fixture loading from `testdata/`
- validate that a scenario cannot belong to multiple splits at once
- add deterministic tests for fixture loading, split validation, and ordering

Mechanically verifiable completion:

- `go test ./design/scenarios/...` passes
- tests prove train/validation/test fixtures load and validate deterministically

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added: `design/scenarios` (`design/scenarios/doc.go`,
  `design/scenarios/scenarios.go`)
- required serializable types added:
  `ScenarioID`, `ScenarioSplit` (with `Train`/`Validation`/`Test` constants),
  `Scenario` (ID + split + optional tags + optional params), `ScenarioSet`
  (named ordered collection with optional source), `ScenarioRegistry`
  (aggregates one or more sets)
- fixture loading helpers added:
  `LoadScenarioSet` (JSON reader), `LoadScenarioSetFile` (file path); both
  use `DisallowUnknownFields` and reject multi-object streams
- split validation added: `ScenarioRegistry.Validate` returns a descriptive
  error for the first conflict — cross-split duplicate ID or same-split
  repeated ID — scanning in declaration order for determinism
- query helpers added: `ScenarioRegistry.ByID` (first match in declaration
  order, bool ok), `ScenarioRegistry.BySplit` (ordered slice, new allocation)
- testdata fixtures added:
  `design/scenarios/testdata/set_smoke.json` (2 train + 1 validation + 2 test
  scenarios, tags exercised), `design/scenarios/testdata/set_params.json`
  (3 scenarios with JSON params map)
- deterministic tests added: `design/scenarios/scenarios_test.go`
  (nil reader, malformed JSON, two-object stream, unknown field, empty set,
  missing file, smoke fixture load, params fixture load + params round-trip,
  smoke split assignments, smoke ordering, smoke validates clean, empty
  registry, valid registry, cross-split duplicate, same-split duplicate,
  duplicate across sets, Validate determinism, ByID found/missing/first-wins,
  BySplit empty/returns-copy/declaration-order, multi-set registry, fixture
  loading determinism, BySplit determinism)
- command proof:
  `go test ./design/scenarios/...`

#### `PB-02`: Add shared promotion and acceptance criteria in `design/promotion`

Goal:

- formalize what it means for a candidate to be "good enough" to promote after
  train/validation evaluation

Dependencies:

- depends on `PA-04`
- depends on `PB-01`

Required outcome:

- add a new `design/promotion` package with package-level godoc
- define serializable types for:
  - acceptance criterion
  - promotion decision
  - promotion report
  - per-split evaluation summary
- implement helpers that decide promotion eligibility from:
  - aggregate score breakdowns
  - validation results
  - split-specific thresholds
- add tests for pass, fail, and edge-threshold cases

Mechanically verifiable completion:

- `go test ./design/promotion/...` passes
- tests prove promotion decisions are deterministic and threshold-driven

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added: `design/promotion` (`design/promotion/doc.go`,
  `design/promotion/promotion.go`)
- required serializable types added:
  `AcceptanceCriterion` (split + MinAggregateScore threshold +
  RequireNoHardFail gate), `SplitEvalSummary` (split + count + mean/min/max
  aggregate score + hard-fail count), `PromotionDecision` (split + passed flag
  + reason string), `PromotionReport` (candidateID + summaries + decisions +
  promoted flag)
- helpers added:
  `SummarizeSplit` (computes SplitEvalSummary from a slice of
  objectives.AggregateResult values), `Decide` (evaluates one
  AcceptanceCriterion against a set of summaries, returning a
  PromotionDecision with a human-readable reason), `Evaluate` (runs all
  criteria, records decisions in criteria order, sets Promoted only when every
  criterion passes)
- promotion eligibility rules implemented:
  threshold-driven (MeanAggregateScore >= MinAggregateScore), hard-fail
  gating (RequireNoHardFail + HardFailCount > 0), missing-summary failure,
  empty-count failure
- deterministic tests added: `design/promotion/promotion_test.go`
  (SummarizeSplit: empty/single/multiple/hard-fail count/determinism/no-inf,
  Decide: pass above threshold/pass at exact threshold/fail below
  threshold/fail no matching summary/fail empty count/fail hard-fail
  required/pass hard-fail not required/determinism, Evaluate: no
  criteria/all pass/one fails/decisions in criteria order/summaries
  preserved/determinism/hard-fail blocks promotion/edge threshold at
  boundary/just below boundary, JSON round-trips for AcceptanceCriterion,
  PromotionReport, SplitEvalSummary)
- command proof:
  `go test ./design/promotion/...`

#### `PB-03`: Define the circuit domain core model in `domains/circuit`

Goal:

- create the reusable circuit candidate representation that future pilots and
  emitters will target

Dependencies:

- depends on `PA-01`

Required outcome:

- add a new `domains/circuit` package with package-level godoc
- define serializable types for at least:
  - `NodeID`
  - `Port`
  - `Component`
  - `CircuitGraph`
  - `CircuitProgram`
  - `CircuitSpec`
  - `CircuitConstraint`
- add validation helpers for:
  - duplicate node IDs
  - illegal port references
  - missing component names/types
- add JSON round-trip tests and validation tests
- do **not** add simulator execution or full pilot logic yet

Mechanically verifiable completion:

- `go test ./domains/circuit/...` passes
- tests prove invalid graphs fail validation and valid graphs round-trip through
  JSON

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added: `domains/circuit` (`domains/circuit/doc.go`,
  `domains/circuit/circuit.go`)
- required serializable types added:
  `NodeID`, `Port`, `Component`, `CircuitGraph`, `CircuitProgram`,
  `CircuitSpec`, `CircuitConstraint`
- validation helpers added:
  `CircuitGraph.Validate` (duplicate `node_id` detection, missing component
  name/type checks, illegal port-reference checks against known nodes),
  `CircuitProgram.Validate` (delegates to graph validation)
- deterministic tests added: `domains/circuit/circuit_test.go`
  (valid graph passes; duplicate node ID fails; missing component name/type
  fail; illegal port reference fails; validation determinism; program validate
  delegation; CircuitProgram JSON round-trip + post-round-trip validation)
- command proof:
  `go test ./domains/circuit/...`

#### `PB-04`: Add circuit artifact emitters and reusable circuit scenario fixtures

Goal:

- make the circuit domain capable of producing durable artifacts and reusable
  train/validation/test fixtures before the first circuit pilot is attempted

Dependencies:

- depends on `PB-01`
- depends on `PB-03`

Required outcome:

- add a `domains/circuit/artifacts` package with emitters for at least:
  - canonical JSON graph output
  - SPICE-style netlist text
  - structural-Verilog-style text
- add a `domains/circuit/scenarios` package or `testdata/` fixture set with:
  - a training split
  - a validation split
  - a test split
- keep the scenarios small and deterministic; they should be reusable by later
  pilots without requiring external tools
- add tests that verify emitted artifacts are stable for a fixed graph and that
  scenario fixtures load via `design/scenarios`

Mechanically verifiable completion:

- `go test ./domains/circuit/... ./design/scenarios/...` passes
- tests prove artifact emitters are deterministic and scenario fixtures are
  valid

Status: ✅ Completed (2026-05-11)

Completion evidence:

- packages added:
  `domains/circuit/artifacts` (`domains/circuit/artifacts/doc.go`,
  `domains/circuit/artifacts/artifacts.go`), `domains/circuit/scenarios`
  (`domains/circuit/scenarios/doc.go`, `domains/circuit/scenarios/scenarios.go`)
- artifact emitters added:
  `artifacts.JSON` (canonical indented JSON for `circuit.CircuitProgram`),
  `artifacts.SPICE` (deterministic SPICE-style netlist text),
  `artifacts.Verilog` (deterministic structural-Verilog-style text); all
  validate the underlying graph before emission
- reusable fixture set added:
  `domains/circuit/scenarios/testdata/set_smoke.json` with deterministic
  train/validation/test circuit scenarios, loaded through
  `design/scenarios.LoadScenarioSet` by `scenarios.LoadFixtureSet`
- deterministic tests added:
  `domains/circuit/artifacts/artifacts_test.go` (exact JSON/SPICE/Verilog
  golden output, repeated-call determinism, invalid-graph rejection),
  `domains/circuit/scenarios/scenarios_test.go` (fixture load, validation via
  `design/scenarios`, split counts, repeated-load determinism)
- command proof:
  `go test ./domains/circuit/... ./design/scenarios/...`

#### `PB-05`: Define the voxel domain core model in `domains/voxel`

Goal:

- create the reusable voxel/implicit-geometry candidate representation that
  future engineering pilots will target

Dependencies:

- depends on `PA-01`

Required outcome:

- add a new `domains/voxel` package with package-level godoc
- define serializable types for at least:
  - `VoxelProgram`
  - `VoxelDesign`
  - `DesignVolume`
  - `Material`
  - `LoadCase`
  - `InterfaceRegion`
  - `VoxelSpec`
- add validation helpers for:
  - out-of-bounds occupied cells
  - overlapping forbidden/interface regions
  - empty or malformed design volumes
- add JSON round-trip tests and validation tests
- do **not** add external geometry kernels or pilot-specific evaluation logic yet

Mechanically verifiable completion:

- `go test ./domains/voxel/...` passes
- tests prove invalid voxel designs fail validation and valid designs round-trip
  through JSON

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added: `domains/voxel` (`domains/voxel/doc.go`,
  `domains/voxel/voxel.go`)
- required serializable core types added:
  `VoxelProgram`, `VoxelDesign`, `DesignVolume`, `Material`, `LoadCase`,
  `InterfaceRegion`, and `VoxelSpec`
- deterministic validation helpers added:
  `DesignVolume.Validate` rejects empty/malformed volumes, malformed or
  out-of-bounds forbidden/interface regions, and overlap between forbidden and
  interface regions; `VoxelDesign.Validate` rejects occupied cells outside the
  design volume; `VoxelProgram.Validate` delegates to design validation
- deterministic tests added: `domains/voxel/voxel_test.go`
  (valid volume/program validation, empty/malformed volume rejection,
  overlapping forbidden/interface region rejection, out-of-bounds occupied cell
  rejection, repeated-call determinism, and JSON round-trip validation)
- command proof:
  `go test ./domains/voxel/...`

#### `PB-06`: Add voxel artifact emitters and reusable voxel scenario fixtures

Goal:

- make the voxel domain capable of producing durable artifacts and reusable
  engineering scenario fixtures before the first voxel pilot is attempted

Dependencies:

- depends on `PB-01`
- depends on `PB-05`

Required outcome:

- add a `domains/voxel/artifacts` package with emitters for at least:
  - canonical JSON design output
  - one mesh-like export format (`.obj` or `.stl`)
  - one simple human-readable summary format
- add a `domains/voxel/scenarios` package or `testdata/` fixture set with:
  - a training split
  - a validation split
  - a test split
- keep the scenarios deterministic and pure-Go; do not require an external
  geometry kernel at this stage
- add tests that verify stable artifact output and valid scenario loading

Mechanically verifiable completion:

- `go test ./domains/voxel/... ./design/scenarios/...` passes
- tests prove artifact emitters are deterministic and scenario fixtures are
  valid

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added: `domains/voxel/artifacts` (`domains/voxel/artifacts/doc.go`,
  `domains/voxel/artifacts/artifacts.go`)
- emitters added:
  - `JSON`: canonical indented JSON output for `VoxelProgram`
  - `OBJ`: Wavefront OBJ mesh with one unit cube per occupied voxel cell
  - `Summary`: concise human-readable plain-text overview
- package added: `domains/voxel/scenarios` (`domains/voxel/scenarios/doc.go`,
  `domains/voxel/scenarios/scenarios.go`)
- embedded fixture set added:
  `domains/voxel/scenarios/testdata/set_smoke.json` with 4 scenarios:
  2 training splits, 1 validation split, 1 test split
- deterministic tests added:
  `domains/voxel/artifacts/artifacts_test.go` (exact JSON/OBJ/Summary
  output verification, determinism across two calls, rejection of invalid
  programs) and `domains/voxel/scenarios/scenarios_test.go` (fixture
  loading, split counts, determinism across two loads)
- command proof:
  `go test ./domains/voxel/... ./design/scenarios/...`

### Phase C: Add domain pilot projects

#### `PC-01`: Add the first circuit pilot in `experiments/circuit/half_adder`

Goal:

- prove the end-to-end circuit path on a small, fully deterministic, pure-Go
  pilot before introducing external circuit simulators

Dependencies:

- depends on `PA-02`
- depends on `PA-03`
- depends on `PA-04`
- depends on `PB-04`

Required outcome:

- add a new pilot entrypoint at `experiments/circuit/half_adder`
- use the shared circuit domain types and reusable scenario fixtures
- implement a pure-Go truth-table evaluator for a bounded half-adder search
  problem
- evolve candidates, decode them to `CircuitGraph`, validate constraints, and
  emit artifacts using the shared circuit emitters
- add tests that prove the evaluator, decoder, and exported artifacts are wired
  together correctly
- do **not** require ngspice, Xyce, or any other external simulator yet

Mechanically verifiable completion:

- `go test ./experiments/circuit/half_adder/... ./domains/circuit/...` passes
- `go run ./experiments/circuit/half_adder` completes and emits deterministic
  artifacts for a fixed seed

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added:
  `experiments/circuit/half_adder` (`experiments/circuit/half_adder/main.go`)
- deterministic pilot entrypoint added:
  `runPilot` loads shared circuit fixture scenarios, performs seeded typed
  evolution for a bounded two-gene half-adder search, decodes the best genome
  to `domains/circuit.CircuitProgram`, validates the graph, and exports JSON,
  SPICE-style, and Verilog artifacts via `domains/circuit/artifacts`
- pure-Go bounded evaluator added:
  `evaluateHalfAdder` scores SUM/CARRY truth-table correctness without external
  simulators; `scoreCandidate` enforces fixture-provided `max_components`
  limits for deterministic bounded search behavior
- decoder + exporter wiring added:
  `decodeCircuitProgram` / `decodeGeneToComponents` map active Karva symbols to
  a structural `CircuitGraph`; `exportArtifacts` writes deterministic
  `candidate.json`, `candidate.spice`, and `candidate.v`
- deterministic integration-style tests added:
  `experiments/circuit/half_adder/main_test.go` proves evaluator + decoder +
  artifact export wiring, validates deterministic re-export stability, and
  verifies scenario component-bound enforcement
- command proof:
  `go test ./experiments/circuit/half_adder/... ./domains/circuit/...`
  `go run ./experiments/circuit/half_adder --seed 20260511 --out /tmp/half_adder_artifacts`

#### `PC-02`: Add circuit promotion, checkpoint, and held-out validation wiring

Goal:

- finish the full circuit pilot loop so it proves
  evolve -> decode -> constrain -> validate -> promote -> export -> checkpoint

Dependencies:

- depends on `PA-06`
- depends on `PB-02`
- depends on `PC-01`

Required outcome:

- wire the half-adder pilot through:
  - train/validation split handling
  - promotion decisions
  - checkpoint save/load
  - run-manifest generation
- add an integration test that runs the pilot with a fixed seed and verifies the
  promotion report plus checkpoint artifacts are created and reloadable
- ensure promoted outputs reference their emitted circuit artifacts through
  `ArtifactRef` metadata

Mechanically verifiable completion:

- `go test ./experiments/circuit/half_adder/... ./design/...` passes
- integration tests prove the promoted candidate survives held-out validation
  and round-trips through checkpoint restore

Status: ✅ Completed (2026-05-11)

Completion evidence:

- pilot loop extended in `experiments/circuit/half_adder/main.go` to execute:
  evolve -> decode -> constrain -> validate -> promote -> export -> checkpoint
- train/validation split wiring added via shared circuit fixture registry
  (`loadScenarioRegistry`) and deterministic split summaries recorded in
  `design.RunManifest`
- promotion wiring added using `design/promotion`:
  per-split aggregate results are summarized with
  `promotion.SummarizeSplit` and evaluated by `promotion.Evaluate` against
  explicit train/validation acceptance criteria
- run-manifest and checkpoint wiring added:
  `run_manifest.json` is emitted with deterministic run config, seed records,
  scenario-split metadata, and artifact references; `checkpoint.json` is saved
  and immediately reloaded with `design/checkpoint` to prove replayability
- promoted-output artifact metadata added:
  emitted circuit artifacts (`candidate.json`, `candidate.spice`, `candidate.v`)
  are recorded as `design.ArtifactRef` entries and propagated through manifest
  and checkpoint metadata
- deterministic integration test added:
  `TestRunPilotPromotionCheckpointAndManifest` in
  `experiments/circuit/half_adder/main_test.go` runs the fixed-seed pilot and
  verifies promotion success, manifest artifact references, and checkpoint
  manifest round-trip parity
- command proof:
  `go test ./experiments/circuit/half_adder/... ./design/...`
  `go run ./experiments/circuit/half_adder --seed 20260511 --out /tmp/half_adder_pc02_artifacts`

#### `PC-03`: Add the first voxel pilot in `experiments/voxel/bracket`

Goal:

- prove the end-to-end voxel path on a deterministic, pure-Go structural pilot
  before any external geometry or FEA integration is attempted

Dependencies:

- depends on `PA-02`
- depends on `PA-03`
- depends on `PA-04`
- depends on `PB-06`

Required outcome:

- add a new pilot entrypoint at `experiments/voxel/bracket`
- use the shared voxel domain types and reusable scenario fixtures
- implement a pure-Go occupancy-grid or lattice-style bracket evaluator with a
  deterministic heuristic score
- evolve candidates, decode them to `VoxelDesign`, validate constraints, and
  emit JSON plus mesh-like artifacts through the shared voxel emitters
- add tests that prove the evaluator, decoder, and artifact emitters are wired
  together correctly
- do **not** require PicoGK, noroyon, or external FEA tooling yet

Mechanically verifiable completion:

- `go test ./experiments/voxel/bracket/... ./domains/voxel/...` passes
- `go run ./experiments/voxel/bracket` completes and emits deterministic
  artifacts for a fixed seed

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added: `experiments/voxel/bracket`
  (`experiments/voxel/bracket/main.go`,
  `experiments/voxel/bracket/main_test.go`)
- pilot entrypoint implemented:
  deterministic pure-Go voxel bracket loop evolves boolean genomes, decodes
  candidates into `domains/voxel.VoxelDesign`, validates constraints, and emits
  shared voxel artifacts (`candidate.json`, `candidate.obj`, `candidate.txt`)
- shared domain wiring implemented:
  fixture scenarios are loaded from `domains/voxel/scenarios` and parsed into
  bounded design volumes plus max-cell constraints for deterministic
  train-split evaluation
- evaluator and decoder wiring implemented:
  occupancy-grid decoding uses deterministic lattice features and pure-Go
  structural heuristics (left/right interface coverage, path connectivity, and
  bounded density fit) with deterministic scoring for fixed seeds
- deterministic tests added:
  `TestBracketPipelineEvaluatorDecoderAndArtifacts`,
  `TestScoreScenarioRejectsOverMaxCells`, and
  `TestRunPilotDeterministicFixedSeed` prove evaluator, decoder, artifact
  emitters, and fixed-seed run determinism are wired together correctly
- command proof:
  `go test ./experiments/voxel/bracket/... ./domains/voxel/...`
  `go run ./experiments/voxel/bracket --seed 20260511 --out /tmp/bracket_pc03_artifacts`

#### `PC-04`: Add voxel promotion, checkpoint, and held-out validation wiring

Goal:

- finish the full voxel pilot loop so it proves
  evolve -> decode -> constrain -> validate -> promote -> export -> checkpoint

Dependencies:

- depends on `PA-06`
- depends on `PB-02`
- depends on `PC-03`

Required outcome:

- wire the bracket pilot through:
  - train/validation split handling
  - promotion decisions
  - checkpoint save/load
  - run-manifest generation
- add an integration test that runs the pilot with a fixed seed and verifies the
  promotion report plus checkpoint artifacts are created and reloadable
- ensure promoted outputs reference their emitted voxel artifacts through
  `ArtifactRef` metadata

Mechanically verifiable completion:

- `go test ./experiments/voxel/bracket/... ./design/...` passes
- integration tests prove the promoted candidate survives held-out validation
  and round-trips through checkpoint restore

Status: ✅ Completed (2026-05-11)

Completion evidence:

- voxel bracket pilot loop in `experiments/voxel/bracket/main.go` now executes:
  evolve -> decode -> constrain -> validate -> promote -> export -> checkpoint
- train/validation split wiring added via shared voxel fixture registry
  (`loadScenarioRegistry`), with deterministic split summaries recorded in
  `design.RunManifest`
- promotion wiring added with `design/promotion`:
  per-split aggregate results are summarized using
  `promotion.SummarizeSplit` and evaluated by `promotion.Evaluate` against
  explicit train/validation acceptance thresholds
- run-manifest and checkpoint wiring added:
  `run_manifest.json` is emitted with deterministic run configuration, seed
  records, scenario split metadata, and artifact references; `checkpoint.json`
  is saved and immediately reloaded with `design/checkpoint`
- promoted-output artifact metadata added:
  `exportArtifacts` now returns `[]design.ArtifactRef` for emitted voxel JSON,
  OBJ, and summary artifacts (`promoted_voxel_*`) and they are attached to the
  manifest/checkpoint refs
- deterministic integration proof added:
  `TestRunPilotPromotionCheckpointAndManifest` verifies fixed-seed promotion,
  promotion report creation, run-manifest artifact references, checkpoint
  reload, and checkpoint manifest replay equivalence
- command proof:
  `go test ./experiments/voxel/bracket/... ./design/...`
  `go run ./experiments/voxel/bracket --seed 20260511 --out /tmp/bracket_pc04_artifacts`

#### `PC-05`: Add a third pilot in `experiments/control/mass_spring_damper`

Goal:

- demonstrate that the applied-design substrate is not limited to circuits and
  geometry by shipping one additional pure-Go simulator-backed discovery domain

Dependencies:

- depends on `PA-02`
- depends on `PA-03`
- depends on `PA-04`
- depends on `PA-06`
- depends on `PB-01`
- depends on `PB-02`

Required outcome:

- add a new pilot entrypoint at `experiments/control/mass_spring_damper`
- implement a pure-Go plant simulator with:
  - train/validation/test disturbance or initial-condition splits
  - deterministic seeded evaluation
  - exported controller/policy artifact output
- use the shared manifest, promotion, and checkpoint infrastructure
- add tests proving the pilot can evolve, validate, promote, and export a
  controller artifact without external dependencies

Mechanically verifiable completion:

- `go test ./experiments/control/mass_spring_damper/... ./design/...` passes
- `go run ./experiments/control/mass_spring_damper` completes and emits a
  deterministic promoted controller artifact for a fixed seed

Status: ✅ Completed (2026-05-11)

Completion evidence:

- package added:
  `experiments/control/mass_spring_damper`
  (`experiments/control/mass_spring_damper/main.go`)
- deterministic control pilot entrypoint added:
  `runPilot` loads deterministic embedded train/validation/test control
  scenarios, performs seeded typed evolution, evaluates candidates with a
  pure-Go mass-spring-damper simulator, runs split-based promotion, exports
  controller artifacts, emits a run manifest, and saves/reloads a checkpoint
- pure-Go simulator-backed evaluation added:
  `simulateController` integrates plant dynamics with deterministic fixed-step
  Euler updates; `scoreScenarioAggregate` scores tracking quality, terminal
  error, and control effort while hard-failing unstable trajectories
- promoted controller artifact wiring added:
  `decodeControllerPolicy`/`exportArtifacts` emit deterministic
  `controller_policy.json` and `controller_summary.txt` artifacts, then attach
  artifact refs to manifest/checkpoint metadata
- deterministic integration proof added:
  `experiments/control/mass_spring_damper/main_test.go` verifies simulator and
  artifact export wiring, fixed-seed deterministic runs, promotion report
  creation across train/validation/test splits, run-manifest artifact refs, and
  checkpoint replay equivalence
- command proof:
  `go test ./experiments/control/mass_spring_damper/... ./design/...`
  `go run ./experiments/control/mass_spring_damper --seed 20260511 --out /tmp/mass_spring_damper_pc05_artifacts`

#### `PC-06`: Add a cross-domain discovery regression suite and update docs

Goal:

- make the new applied-design pipeline a stable, repo-visible contract rather
  than a set of isolated experiments

Dependencies:

- depends on `PC-02`
- depends on `PC-04`
- depends on `PC-05`

Required outcome:

- add cross-domain regression coverage that exercises the promoted pipeline for:
  - the circuit pilot
  - the voxel pilot
  - the control pilot
- ensure the regression suite verifies the full flow:
  evolve -> decode -> constrain -> validate -> promote -> export -> checkpoint
- update the root docs to list the new pilot entrypoints and the applied-design
  package map once the regression suite is in place
- keep the regression suite pure-Go and CI-friendly

Mechanically verifiable completion:

- `go test ./...` passes with the new pilots and regression suite included
- docs reference the new applied-design packages and pilot entrypoints

The goal is still not to "cover everything" immediately. The goal is to prove
the full pipeline repeatedly and mechanically:

> **evolve -> decode -> simulate -> constrain -> validate -> promote -> export -> checkpoint**

## 1. Novel electronic circuit design

This is a strong fit for the engine because the search space is symbolic,
combinatorial, and naturally simulator-backed.

### What should be designed first

Pick a circuit class with clear objectives and reliable simulation:

- logic block or arithmetic cell
- oscillator
- analog filter
- low-noise amplifier
- power stage or regulator subcircuit
- sensor front-end

Do not start with a full chip. Start with a narrow block that has:

- a bounded component library
- a small number of operating conditions
- measurable objectives
- a well-defined export format

### Tool collection

| Need | Recommended tool class |
| --- | --- |
| Search engine | `core` + `evolution` from this repo |
| Candidate decoding | New Go package such as `domains/circuit` |
| Circuit simulation | A SPICE-class simulator such as ngspice or Xyce |
| Artifact export | SPICE netlist emitter; optionally schematic/Verilog emitters |
| Constraint checking | Go validators for topology, parameter bounds, and manufacturability rules |
| Experiment tracking | Run manifest files plus optional external tracking system |
| Surrogate modeling | Optional Python/JAX/PyTorch sidecar when simulation cost becomes dominant |

### Required inputs

At minimum, the circuit workflow needs:

- a component library:
  - allowed devices
  - parameter ranges
  - any topology templates you want to permit
- a spec:
  - gain, bandwidth, delay, power, noise, area, stability, yield, or other
    objectives
- operating scenarios:
  - input stimuli
  - loads
  - voltage rails
  - temperature/process corners
- hard constraints:
  - illegal topologies
  - forbidden connections
  - max current, max voltage, area limits, etc.
- evaluation policy:
  - reward weights
  - pass/fail thresholds
  - train/validation/test scenario split

### What the evolved artifact should look like

Use one of these output formats:

1. **SPICE netlist** for immediate simulation and validation
2. **JSON circuit graph** for easier downstream tooling, then compile that to a
   netlist
3. **Structural Verilog** for digital logic-style designs

In practice, the cleanest path is:

- evolve an internal circuit graph or DSL in Go
- export that graph to SPICE for simulation
- optionally export the winning result to schematic or HDL form

### How the engine should be used

#### Step 1: Define a circuit DSL in Go

Create a `domains/circuit` package with types like:

- `NodeID`
- `Component`
- `Port`
- `CircuitGraph`
- `CircuitProgram`
- `CircuitSpec`
- `CircuitConstraint`
- `NetlistArtifact`

Then decide what a genome means.

Two practical choices:

1. **Topology program**: the genome builds a circuit graph by composing
   primitives such as `Series`, `Parallel`, `Mirror`, `DiffPair`, `RC`,
   `Buffer`, `SelectWidth`, `SelectLength`, etc.
2. **Template parameterizer**: the genome fills in parameters of a constrained
   template instead of inventing arbitrary topology from scratch

Start with template parameterization if fabrication realism matters more than
topological novelty. Start with topology programs if novelty is the point.

#### Step 2: Decode genomes to circuit artifacts

Add a decoder:

- `func Decode(genome core.Genome[CircuitProgram]) (CircuitGraph, error)`

Then add an emitter:

- `func (g CircuitGraph) SPICE() ([]byte, error)`

This decoder is where you enforce:

- valid pin counts
- legal component wiring
- node naming
- bounded parameter selection

#### Step 3: Simulate in batches

For every candidate:

1. export the circuit to SPICE
2. run the simulator across the training scenario batch
3. extract measurements
4. compute reward and penalties

Batching matters. A serious circuit workflow should evaluate many candidates
across many corners concurrently instead of one-at-a-time shelling out in an
experiment `main.go`.

#### Step 4: Use modern RL-style reward structure

Reward should be a composition of:

- **hard-gate constraints**: invalid or dangerous circuits get rejected or
  heavily penalized
- **primary task reward**: how well the design meets the spec
- **secondary objectives**: power, area, noise, robustness, component count
- **novelty bonus**: reward structurally different but still valid candidates

Do not rely on a single scalar metric from one nominal operating point.

#### Step 5: Validate on held-out corners

After each generation, or at least for the elite set:

- rerun the best designs on validation corners that were not used for direct
  reward
- only promote designs that stay good under held-out conditions

This is the circuit analogue of train/validation/test evaluation in modern RL.

#### Step 6: Export and review

The final promoted artifact should include:

- circuit graph JSON
- SPICE netlist
- simulation measurements
- run manifest
- rendered human-readable summary

That package set is what turns "interesting evolved expression" into "design we
can inspect, rerun, and potentially fabricate."

### Best initial repo additions for circuits

Add packages and entrypoints roughly like:

- `domains/circuit`
- `domains/circuit/spice`
- `domains/circuit/constraints`
- `domains/circuit/metrics`
- `experiments/circuit/filter`
- `experiments/circuit/oscillator`

## 2. Novel voxel-based design with computational engineering tools

This is another excellent use case, especially for lattice structures, thermal
parts, fixtures, fluid manifolds, brackets, porous media, or topology-optimized
components.

The engine should evolve **geometry-building programs**, not just scalars.

### Tool collection

| Need | Recommended tool class |
| --- | --- |
| Search engine | `core` + `evolution` from this repo |
| Geometry kernel | A voxel or implicit-geometry toolchain, e.g. a PicoGK- or noroyon-style stack |
| Domain package | New Go package such as `domains/voxel` |
| Simulation | FEA, CFD, thermal, acoustics, or custom engineering solver |
| Artifact export | Voxel field, mesh (`.stl`, `.3mf`, `.obj`), or a reproducible Go builder |
| Constraint checking | Go validators for geometry, interfaces, printability, load paths, and forbidden regions |
| Visualization | Mesh/voxel viewer for inspection of elites |

### Required inputs

At minimum, the voxel workflow needs:

- a design volume or workspace
- material options
- boundary conditions:
  - loads
  - supports
  - flow ports
  - thermal contacts
- manufacturing limits:
  - minimum wall thickness
  - maximum overhang
  - minimum feature size
  - keep-out zones
- required interfaces:
  - bolt holes
  - connectors
  - mating surfaces
- objective set:
  - stiffness-to-weight
  - pressure drop
  - heat transfer
  - resonance
  - mass
  - compliance

### What the evolved artifact should look like

Choose one primary internal representation, then export others:

1. **voxel occupancy grid**
2. **signed-distance or implicit-geometry program**
3. **CSG-like construction script in Go**

Then export winners to:

- mesh files such as STL or 3MF
- sparse voxel or field data
- a Go builder script that deterministically rebuilds the geometry

For this domain, a Go builder or geometry AST is often better than the current
grammar-based textual codegen path.

### How Go should be used to model constraints

Go should own the constraint system.

Define types like:

- `VoxelDesign`
- `VoxelProgram`
- `DesignVolume`
- `Material`
- `LoadCase`
- `InterfaceRegion`
- `Constraint`
- `ValidationReport`

Then implement constraints as ordinary Go code, for example:

- `MinWallThickness`
- `MaxOverhang`
- `ConnectedLoadPath`
- `KeepOutRegion`
- `RequiredSymmetry`
- `RequiredMountingFaces`
- `MassBudget`
- `MaxDeflection`

Each constraint should support one of three outcomes:

1. reject the candidate immediately
2. repair the candidate into a valid form
3. apply a penalty and continue evaluation

That makes the constraints explicit, testable, and reusable.

### How the engine should be used

#### Step 1: Define geometry-building primitives

Create `core.Node[T]` implementations or a decoder DSL around operations like:

- union
- subtract
- intersect
- extrude
- shell
- lattice fill
- thicken
- mirror
- place connector
- carve channel

The output should be a geometry program or AST, not just a final scalar score.

#### Step 2: Decode the genome into geometry

Add a decoder:

- `func Decode(genome core.Genome[VoxelProgram]) (VoxelDesign, error)`

Then compile that design into your geometry kernel and export intermediate and
final artifacts.

#### Step 3: Run engineering simulation

Evaluate each candidate across a batch of load cases or flow/thermal scenarios.
Examples:

- FEA for structural stiffness
- CFD for pressure drop
- thermal solver for cooling efficiency
- modal analysis for vibration behavior

#### Step 4: Use constrained, multi-objective reward

Reward should combine:

- validity/manufacturability
- task performance
- robustness over multiple scenarios
- simplicity or fabrication cost
- novelty/diversity

As with circuits, held-out validation scenarios matter.

#### Step 5: Export fabrication-ready artifacts

A winning design should produce:

- geometry source representation
- mesh export
- simulation summary
- validation report
- run manifest

### Best initial repo additions for voxel work

Add packages and entrypoints roughly like:

- `domains/voxel`
- `domains/voxel/constraints`
- `domains/voxel/artifacts`
- `domains/voxel/metrics`
- `experiments/voxel/bracket`
- `experiments/voxel/manifold`

## 3. Other practical scientific domains

Once the circuit and voxel stacks exist, the same engine can be applied much
more broadly.

| Domain | Candidate artifact | Typical simulator or feedback source | Likely outputs |
| --- | --- | --- | --- |
| Antenna and RF structure design | topology or geometry program | EM solver | netlists, meshes, fabrication geometry |
| Microfluidic device design | channel/layout graph | CFD and pressure/flow models | masks, meshes, printable geometry |
| Materials or metamaterial design | lattice/unit-cell program | FEA, thermal, acoustic, or surrogate models | geometry, property tables, process settings |
| Robot morphology or end-effector design | body/fixture geometry plus controller hooks | physics simulator | CAD/mesh plus control config |
| Control-law discovery | symbolic controller or policy | simulator or plant model | Go/C code, config files, validation traces |
| Sensor placement and experiment design | placement graph or protocol | simulator or recorded datasets | deployment plan, protocol, justification |
| Chemical process recipe tuning | symbolic recipe or schedule | plant simulator or lab feedback | structured recipe file and validation report |

## How to approach any new discovery problem

No matter the scientific area, use this sequence.

### 1. Choose the artifact first

Decide what the evolved thing really is:

- a text program
- a graph
- a netlist
- a geometry AST
- a controller
- a recipe

If that is vague, the project will drift.

### 2. Define a bounded design language

Do not begin with "arbitrary invention." Begin with a bounded DSL and grow it.
The function catalog should encode what kinds of novelty are allowed.

### 3. Make constraints explicit in Go

Constraints should be code, not comments. They should be unit-testable and
reusable across experiments.

### 4. Separate train, validation, and promotion

Search on the train scenarios, compare elites on validation scenarios, and only
promote designs that survive held-out checks.

### 5. Keep artifacts and manifests

Every promoted design should come with:

- the genome or symbolic representation
- emitted artifact(s)
- simulator outputs
- scoring breakdown
- seed and run config

### 6. Add a surrogate only after the exact loop works

Surrogate models are useful when evaluation is expensive, but they should speed
up a trustworthy pipeline, not replace the first trustworthy pipeline.

## Concrete next steps for this repository

If the goal is to turn `gep` into a practical discovery platform, the next work
should be:

1. add a shared applied-design substrate for manifests, constraints, and batched
   evaluator workers
2. add `domains/circuit` with SPICE export and at least one simulator-backed
   pilot problem
3. add `domains/voxel` with geometry constraints and at least one structural or
   thermal pilot problem
4. add held-out validation, novelty tracking, and artifact persistence to both
   pilot domains
5. document the promoted workflow with reproducible examples the same way the
   current experiments document the typed architecture

That is the path from "modern typed GEP engine" to "useful scientific discovery
platform."
