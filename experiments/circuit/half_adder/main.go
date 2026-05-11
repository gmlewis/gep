// -*- compile-command: "go test ./experiments/circuit/half_adder/..."; -*-
// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/design"
	"github.com/gmlewis/gep/v2/design/checkpoint"
	"github.com/gmlewis/gep/v2/design/objectives"
	"github.com/gmlewis/gep/v2/design/promotion"
	designscenarios "github.com/gmlewis/gep/v2/design/scenarios"
	"github.com/gmlewis/gep/v2/domains/circuit"
	circuitartifacts "github.com/gmlewis/gep/v2/domains/circuit/artifacts"
	circuitscenarios "github.com/gmlewis/gep/v2/domains/circuit/scenarios"
	"github.com/gmlewis/gep/v2/evolution"
	boolNodes "github.com/gmlewis/gep/v2/functions/bool_nodes"
)

var halfAdderTruthTable = []struct {
	A, B  bool
	Sum   bool
	Carry bool
}{
	{A: false, B: false, Sum: false, Carry: false},
	{A: false, B: true, Sum: true, Carry: false},
	{A: true, B: false, Sum: true, Carry: false},
	{A: true, B: true, Sum: false, Carry: true},
}

type runConfig struct {
	Seed           int64
	PopulationSize int
	Generations    int
	OutputDir      string
}

type runResult struct {
	CandidateID string
	Score       float64
	Karva       string
	GateCount   int
	Promoted    bool
}

const (
	candidateJSONArtifactName    = "candidate.json"
	candidateSPICEArtifactName   = "candidate.spice"
	candidateVerilogArtifactName = "candidate.v"
	promotionReportArtifactName  = "promotion_report.json"
	runManifestArtifactName      = "run_manifest.json"
	checkpointArtifactName       = "checkpoint.json"
)

var halfAdderObjectiveDefs = []objectives.ObjectiveDef{{
	Name:   "truth_table_accuracy",
	Weight: 1,
	Kind:   objectives.Soft,
}}

func main() {
	cfg := runConfig{}
	flag.Int64Var(&cfg.Seed, "seed", 20260511, "deterministic evolution seed")
	flag.IntVar(&cfg.PopulationSize, "population", 60, "evolution population size")
	flag.IntVar(&cfg.Generations, "generations", 120, "maximum evolution generations")
	flag.StringVar(&cfg.OutputDir, "out", filepath.Join("artifacts", "circuit", "half_adder"), "output directory for emitted artifacts")
	flag.Parse()

	result, err := runPilot(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "half_adder: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("half_adder complete: candidate=%s score=%.2f gates=%d promoted=%t karva=%s\n", result.CandidateID, result.Score, result.GateCount, result.Promoted, result.Karva)
	fmt.Printf("artifacts written to %s\n", cfg.OutputDir)
}

func runPilot(cfg runConfig) (runResult, error) {
	if cfg.PopulationSize <= 0 {
		return runResult{}, errors.New("population must be > 0")
	}
	if cfg.Generations <= 0 {
		return runResult{}, errors.New("generations must be > 0")
	}

	set, registry, err := loadScenarioRegistry()
	if err != nil {
		return runResult{}, err
	}
	trainScenarios := registry.BySplit(designscenarios.Train)
	validationScenarios := registry.BySplit(designscenarios.Validation)
	if len(trainScenarios) == 0 {
		return runResult{}, errors.New("circuit fixture scenarios contain no train split")
	}
	if len(validationScenarios) == 0 {
		return runResult{}, errors.New("circuit fixture scenarios contain no validation split")
	}

	candidateID := fmt.Sprintf("half-adder-seed-%d", cfg.Seed)
	manifest := design.RunManifest{
		RunID: candidateID,
		RunConfig: design.RunConfig{
			Domain:         "circuit",
			Experiment:     "half_adder",
			PopulationSize: cfg.PopulationSize,
			NumGenerations: cfg.Generations,
		},
		Seeds: []design.SeedRecord{{
			Name:    "evolution_seed",
			Value:   cfg.Seed,
			Purpose: "deterministic evolution seed",
		}},
		ScenarioSplits: []design.ScenarioSplitSummary{
			{SplitName: string(designscenarios.Train), ScenarioCount: len(trainScenarios), Source: set.Source},
			{SplitName: string(designscenarios.Validation), ScenarioCount: len(validationScenarios), Source: set.Source},
		},
	}

	catalog, err := boolNodes.CatalogFromNames([]string{"Not", "And", "Or", "Xor"})
	if err != nil {
		return runResult{}, fmt.Errorf("build boolean function catalog: %w", err)
	}
	link, err := boolNodes.LinkFuncFrom("Or")
	if err != nil {
		return runResult{}, fmt.Errorf("build link function: %w", err)
	}

	population, err := evolution.NewWithSeed(cfg.Seed, catalog, cfg.PopulationSize, 4, 2, 2, 0, link, func(g core.Genome[bool]) float64 {
		return scoreCandidate(g, trainScenarios)
	})
	if err != nil {
		return runResult{}, fmt.Errorf("create seeded population: %w", err)
	}
	population.StopFunc = func(best evolution.Individual[bool]) bool {
		return best.Score >= 1000
	}

	best := population.Evolve(cfg.Generations)
	program, gateCount, err := decodeCircuitProgram(candidateID, best.Genome)
	if err != nil {
		return runResult{}, err
	}
	artifactRefs, err := exportArtifacts(program, cfg.OutputDir)
	if err != nil {
		return runResult{}, err
	}
	manifest.Artifacts = append(manifest.Artifacts, artifactRefs...)

	trainResults, err := evaluateScenarios(best.Genome, gateCount, trainScenarios)
	if err != nil {
		return runResult{}, err
	}
	validationResults, err := evaluateScenarios(best.Genome, gateCount, validationScenarios)
	if err != nil {
		return runResult{}, err
	}
	report := promotion.Evaluate(candidateID, []promotion.AcceptanceCriterion{
		{Split: designscenarios.Train, MinAggregateScore: 400},
		{Split: designscenarios.Validation, MinAggregateScore: 400},
	}, []promotion.SplitEvalSummary{
		promotion.SummarizeSplit(designscenarios.Train, trainResults),
		promotion.SummarizeSplit(designscenarios.Validation, validationResults),
	})

	promotionReportPath := filepath.Join(cfg.OutputDir, promotionReportArtifactName)
	if err := writeJSONFile(promotionReportPath, report); err != nil {
		return runResult{}, fmt.Errorf("write promotion report: %w", err)
	}
	manifest.Artifacts = append(manifest.Artifacts, design.ArtifactRef{
		Name:   "promotion_report",
		Path:   promotionReportPath,
		Format: "json",
		Kind:   "promotion.report",
	})

	manifestPath := filepath.Join(cfg.OutputDir, runManifestArtifactName)
	if err := design.WriteRunManifestFile(manifestPath, &manifest); err != nil {
		return runResult{}, fmt.Errorf("write run manifest: %w", err)
	}
	manifest.Artifacts = append(manifest.Artifacts, design.ArtifactRef{
		Name:   "run_manifest",
		Path:   manifestPath,
		Format: "json",
		Kind:   "run.manifest",
	})
	if err := design.WriteRunManifestFile(manifestPath, &manifest); err != nil {
		return runResult{}, fmt.Errorf("rewrite run manifest with full artifact refs: %w", err)
	}

	checkpointPath := filepath.Join(cfg.OutputDir, checkpointArtifactName)
	checkpointRefs := append([]design.ArtifactRef(nil), manifest.Artifacts...)
	checkpointRefs = append(checkpointRefs, design.ArtifactRef{
		Name:   "checkpoint",
		Path:   checkpointPath,
		Format: "json",
		Kind:   "run.checkpoint",
	})
	snapshot := &checkpoint.Snapshot{
		Manifest: manifest,
		Elites: []checkpoint.EliteRecord{{
			CandidateID:    candidateID,
			Generation:     cfg.Generations,
			AggregateScore: report.Summaries[len(report.Summaries)-1].MeanAggregateScore,
			Breakdown:      validationResults[0].Breakdown,
		}},
		Aggregate: checkpoint.AggregateSnapshot{
			BestScore:                best.Score,
			MeanScore:                report.Summaries[len(report.Summaries)-1].MeanAggregateScore,
			TotalCandidatesEvaluated: cfg.PopulationSize * cfg.Generations,
		},
		ArtifactRefs: checkpointRefs,
	}
	if err := checkpoint.SaveFile(checkpointPath, snapshot); err != nil {
		return runResult{}, fmt.Errorf("write checkpoint snapshot: %w", err)
	}
	if _, err := checkpoint.LoadFile(checkpointPath); err != nil {
		return runResult{}, fmt.Errorf("reload checkpoint snapshot: %w", err)
	}

	return runResult{
		CandidateID: candidateID,
		Score:       best.Score,
		Karva:       best.Genome.KarvaString(),
		GateCount:   gateCount,
		Promoted:    report.Promoted,
	}, nil
}

func loadScenarioRegistry() (*designscenarios.ScenarioSet, *designscenarios.ScenarioRegistry, error) {
	set, err := circuitscenarios.LoadFixtureSet()
	if err != nil {
		return nil, nil, fmt.Errorf("load circuit fixture scenarios: %w", err)
	}
	registry := &designscenarios.ScenarioRegistry{Sets: []designscenarios.ScenarioSet{*set}}
	if err := registry.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validate circuit fixture scenarios: %w", err)
	}
	return set, registry, nil
}

func loadTrainScenarios() ([]designscenarios.Scenario, error) {
	_, registry, err := loadScenarioRegistry()
	if err != nil {
		return nil, err
	}
	train := registry.BySplit(designscenarios.Train)
	if len(train) == 0 {
		return nil, errors.New("circuit fixture scenarios contain no train split")
	}
	return train, nil
}

func scoreCandidate(genome core.Genome[bool], scenarios []designscenarios.Scenario) float64 {
	if len(genome.Genes) < 2 {
		return 0
	}

	baseScore, evalErr := halfAdderAccuracyScore(genome)
	if evalErr != nil {
		return 0
	}

	_, gateCount, err := decodeCircuitProgram("score", genome)
	if err != nil {
		return 0
	}

	if len(scenarios) == 0 {
		return baseScore
	}

	var totalScore float64
	for _, sc := range scenarios {
		maxComponents := scenarioMaxComponents(sc)
		if maxComponents > 0 && gateCount > maxComponents {
			continue
		}
		totalScore += baseScore
	}
	return totalScore / float64(len(scenarios))
}

func halfAdderAccuracyScore(genome core.Genome[bool]) (float64, error) {
	correct, total, evalErr := evaluateHalfAdder(genome)
	if evalErr != nil {
		return 0, evalErr
	}
	if total == 0 {
		return 0, errors.New("evaluate half adder: no outputs produced")
	}
	return 1000.0 * float64(correct) / float64(total), nil
}

func evaluateScenarios(genome core.Genome[bool], gateCount int, scenarios []designscenarios.Scenario) ([]objectives.AggregateResult, error) {
	baseScore, err := halfAdderAccuracyScore(genome)
	if err != nil {
		return nil, fmt.Errorf("evaluate truth-table score: %w", err)
	}
	results := make([]objectives.AggregateResult, 0, len(scenarios))
	for _, sc := range scenarios {
		maxComponents := scenarioMaxComponents(sc)
		rejected := maxComponents > 0 && gateCount > maxComponents
		results = append(results, objectives.Score(halfAdderObjectiveDefs, map[string]float64{
			"truth_table_accuracy": baseScore,
		}, rejected, 0))
	}
	return results, nil
}

func evaluateHalfAdder(genome core.Genome[bool]) (correct int, total int, err error) {
	if len(genome.Genes) < 2 {
		return 0, 0, errors.New("genome must contain at least two genes (sum and carry)")
	}

	for _, tc := range halfAdderTruthTable {
		inputs := []bool{tc.A, tc.B}
		sumOut, sumErr := genome.Genes[0].Eval(inputs)
		if sumErr != nil {
			return 0, 0, fmt.Errorf("evaluate sum gene: %w", sumErr)
		}
		carryOut, carryErr := genome.Genes[1].Eval(inputs)
		if carryErr != nil {
			return 0, 0, fmt.Errorf("evaluate carry gene: %w", carryErr)
		}

		if sumOut == tc.Sum {
			correct++
		}
		total++
		if carryOut == tc.Carry {
			correct++
		}
		total++
	}

	return correct, total, nil
}

func decodeCircuitProgram(candidateID string, genome core.Genome[bool]) (circuit.CircuitProgram, int, error) {
	if len(genome.Genes) < 2 {
		return circuit.CircuitProgram{}, 0, errors.New("decode circuit: genome must contain at least two genes")
	}

	components := []circuit.Component{
		{
			NodeID: "in_d0",
			Name:   "in_d0",
			Type:   "source.input",
			Outputs: []circuit.Port{{
				Node: "in_d0",
				Name: "out",
			}},
		},
		{
			NodeID: "in_d1",
			Name:   "in_d1",
			Type:   "source.input",
			Outputs: []circuit.Port{{
				Node: "in_d1",
				Name: "out",
			}},
		},
	}

	gateCount := 0
	for geneIndex := 0; geneIndex < 2; geneIndex++ {
		geneComponents, rootPort, count, err := decodeGeneToComponents(geneIndex, genome.Genes[geneIndex])
		if err != nil {
			return circuit.CircuitProgram{}, 0, err
		}
		gateCount += count
		components = append(components, geneComponents...)

		outputNodeID := circuit.NodeID(fmt.Sprintf("out_%d", geneIndex))
		outputName := "sum"
		if geneIndex == 1 {
			outputName = "carry"
		}
		components = append(components, circuit.Component{
			NodeID: outputNodeID,
			Name:   outputName,
			Type:   "probe.output",
			Inputs: []circuit.Port{rootPort},
		})
	}

	program := circuit.CircuitProgram{
		CandidateID: candidateID,
		Graph:       circuit.CircuitGraph{Components: components},
		Spec: circuit.CircuitSpec{
			Name:   "half_adder",
			Domain: "circuit",
		},
	}
	if err := program.Validate(); err != nil {
		return circuit.CircuitProgram{}, 0, fmt.Errorf("decode circuit: %w", err)
	}
	return program, gateCount, nil
}

func decodeGeneToComponents(geneIndex int, gene core.Gene[bool]) ([]circuit.Component, circuit.Port, int, error) {
	if err := gene.Validate(); err != nil {
		return nil, circuit.Port{}, 0, fmt.Errorf("decode gene %d: %w", geneIndex, err)
	}

	argOrder, err := buildArgOrder(gene.Symbols)
	if err != nil {
		return nil, circuit.Port{}, 0, fmt.Errorf("decode gene %d: %w", geneIndex, err)
	}
	used := make(map[int]struct{})
	collectUsedSymbols(0, argOrder, used)

	componentBySymbol := map[int]circuit.NodeID{}
	gateCount := 0
	for i, sym := range gene.Symbols {
		if _, ok := used[i]; !ok {
			continue
		}
		if sym.Kind != core.SymbolKindFunction {
			continue
		}

		nodeID := circuit.NodeID(fmt.Sprintf("g%d_n%d", geneIndex, i))
		componentBySymbol[i] = nodeID
		gateCount++
	}

	components := make([]circuit.Component, 0, gateCount)
	for i, sym := range gene.Symbols {
		if _, ok := used[i]; !ok {
			continue
		}
		if sym.Kind != core.SymbolKindFunction {
			continue
		}
		nodeID := componentBySymbol[i]

		inputs := make([]circuit.Port, 0, len(argOrder[i]))
		for _, childIndex := range argOrder[i] {
			portRef, portErr := symbolPortRef(gene.Symbols[childIndex], childIndex, componentBySymbol)
			if portErr != nil {
				return nil, circuit.Port{}, 0, fmt.Errorf("decode gene %d symbol[%d]: %w", geneIndex, i, portErr)
			}
			inputs = append(inputs, portRef)
		}

		components = append(components, circuit.Component{
			NodeID:  nodeID,
			Name:    strings.ToLower(sym.Name),
			Type:    "gate." + strings.ToLower(sym.Name),
			Inputs:  inputs,
			Outputs: []circuit.Port{{Node: nodeID, Name: "out"}},
		})
	}

	rootPort, err := symbolPortRef(gene.Symbols[0], 0, componentBySymbol)
	if err != nil {
		return nil, circuit.Port{}, 0, fmt.Errorf("decode gene %d root: %w", geneIndex, err)
	}
	return components, rootPort, gateCount, nil
}

func buildArgOrder(symbols []core.Symbol[bool]) ([][]int, error) {
	argOrder := make([][]int, len(symbols))
	argIndex := 0
	for i, sym := range symbols {
		if sym.Kind != core.SymbolKindFunction {
			continue
		}
		arity := sym.Node.Arity()
		if arity < 0 {
			return nil, fmt.Errorf("symbol[%d] has negative arity", i)
		}
		if arity == 0 {
			continue
		}
		args := make([]int, arity)
		for j := 0; j < arity; j++ {
			argIndex++
			if argIndex >= len(symbols) {
				return nil, errors.New("invalid karva expression: missing symbol arguments")
			}
			args[j] = argIndex
		}
		argOrder[i] = args
	}
	return argOrder, nil
}

func collectUsedSymbols(index int, argOrder [][]int, used map[int]struct{}) {
	if _, ok := used[index]; ok {
		return
	}
	used[index] = struct{}{}
	for _, child := range argOrder[index] {
		collectUsedSymbols(child, argOrder, used)
	}
}

func symbolPortRef(sym core.Symbol[bool], symbolIndex int, componentBySymbol map[int]circuit.NodeID) (circuit.Port, error) {
	switch sym.Kind {
	case core.SymbolKindFunction:
		nodeID, ok := componentBySymbol[symbolIndex]
		if !ok {
			return circuit.Port{}, fmt.Errorf("missing component mapping for symbol[%d]", symbolIndex)
		}
		return circuit.Port{Node: nodeID, Name: "out"}, nil
	case core.SymbolKindTerminal:
		return circuit.Port{Node: circuit.NodeID("in_d" + strconv.Itoa(sym.TerminalIndex)), Name: "out"}, nil
	default:
		return circuit.Port{}, fmt.Errorf("unsupported symbol kind %v", sym.Kind)
	}
}

func scenarioMaxComponents(sc designscenarios.Scenario) int {
	raw, ok := sc.Params["max_components"]
	if !ok {
		return 0
	}
	switch v := raw.(type) {
	case float64:
		return int(v)
	case int:
		return v
	default:
		return 0
	}
}

func exportArtifacts(program circuit.CircuitProgram, outputDir string) ([]design.ArtifactRef, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("create output directory %q: %w", outputDir, err)
	}

	jsonBytes, err := circuitartifacts.JSON(program)
	if err != nil {
		return nil, fmt.Errorf("emit JSON artifact: %w", err)
	}
	jsonPath := filepath.Join(outputDir, candidateJSONArtifactName)
	if err := os.WriteFile(jsonPath, jsonBytes, 0o644); err != nil {
		return nil, fmt.Errorf("write JSON artifact: %w", err)
	}

	spiceText, err := circuitartifacts.SPICE(program)
	if err != nil {
		return nil, fmt.Errorf("emit SPICE artifact: %w", err)
	}
	spicePath := filepath.Join(outputDir, candidateSPICEArtifactName)
	if err := os.WriteFile(spicePath, []byte(spiceText), 0o644); err != nil {
		return nil, fmt.Errorf("write SPICE artifact: %w", err)
	}

	verilogText, err := circuitartifacts.Verilog(program)
	if err != nil {
		return nil, fmt.Errorf("emit Verilog artifact: %w", err)
	}
	verilogPath := filepath.Join(outputDir, candidateVerilogArtifactName)
	if err := os.WriteFile(verilogPath, []byte(verilogText), 0o644); err != nil {
		return nil, fmt.Errorf("write Verilog artifact: %w", err)
	}

	return []design.ArtifactRef{
		{Name: "promoted_circuit_json", Path: jsonPath, Format: "json", Kind: "circuit.program"},
		{Name: "promoted_circuit_spice", Path: spicePath, Format: "spice", Kind: "circuit.netlist"},
		{Name: "promoted_circuit_verilog", Path: verilogPath, Format: "verilog", Kind: "circuit.verilog"},
	}, nil
}

func writeJSONFile(filename string, value any) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create %q: %w", filename, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		return fmt.Errorf("encode %q: %w", filename, err)
	}
	return nil
}
