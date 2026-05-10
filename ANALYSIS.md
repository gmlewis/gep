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

### Phase A: Build the applied-design substrate

Add a new package tree, likely something like `design`, `design/eval`, or
`domains/shared`, with:

- a run manifest that records seed, config, simulator version, and artifact IDs
- evaluator workers for batched parallel simulation
- a constraint interface that can reject, repair, or penalize candidates
- multi-objective score aggregation
- novelty archive support
- checkpointing and replay support

### Phase B: Add artifact emitters and scenario sets

For each domain, define:

- the candidate representation used during evolution
- the emitted artifact format used by downstream tools
- the scenario set used for train/validation/test evaluation
- the acceptance criteria that must be satisfied before a design is considered
  real

### Phase C: Add domain pilot projects

Start with one pilot in each major domain:

- one analog or digital circuit problem
- one voxel or implicit-geometry structural problem
- one additional domain where simulator-backed search is clearly valuable

The goal is not to "cover everything" immediately. The goal is to prove the
full pipeline: evolve -> decode -> simulate -> constrain -> validate -> export.

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
