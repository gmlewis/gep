// -*- compile-command: "go test ./experiments/circuit/half_adder/..."; -*-
// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gmlewis/gep/v2/core"
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
}

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
	fmt.Printf("half_adder complete: candidate=%s score=%.2f gates=%d karva=%s\n", result.CandidateID, result.Score, result.GateCount, result.Karva)
	fmt.Printf("artifacts written to %s\n", cfg.OutputDir)
}

func runPilot(cfg runConfig) (runResult, error) {
	if cfg.PopulationSize <= 0 {
		return runResult{}, errors.New("population must be > 0")
	}
	if cfg.Generations <= 0 {
		return runResult{}, errors.New("generations must be > 0")
	}

	trainScenarios, err := loadTrainScenarios()
	if err != nil {
		return runResult{}, err
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
	candidateID := fmt.Sprintf("half-adder-seed-%d", cfg.Seed)
	program, gateCount, err := decodeCircuitProgram(candidateID, best.Genome)
	if err != nil {
		return runResult{}, err
	}
	if err := exportArtifacts(program, cfg.OutputDir); err != nil {
		return runResult{}, err
	}

	return runResult{
		CandidateID: candidateID,
		Score:       best.Score,
		Karva:       best.Genome.KarvaString(),
		GateCount:   gateCount,
	}, nil
}

func loadTrainScenarios() ([]designscenarios.Scenario, error) {
	set, err := circuitscenarios.LoadFixtureSet()
	if err != nil {
		return nil, fmt.Errorf("load circuit fixture scenarios: %w", err)
	}
	registry := &designscenarios.ScenarioRegistry{Sets: []designscenarios.ScenarioSet{*set}}
	if err := registry.Validate(); err != nil {
		return nil, fmt.Errorf("validate circuit fixture scenarios: %w", err)
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

	correct, total, evalErr := evaluateHalfAdder(genome)
	if evalErr != nil || total == 0 {
		return 0
	}
	baseScore := 1000.0 * float64(correct) / float64(total)

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

func exportArtifacts(program circuit.CircuitProgram, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory %q: %w", outputDir, err)
	}

	jsonBytes, err := circuitartifacts.JSON(program)
	if err != nil {
		return fmt.Errorf("emit JSON artifact: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "candidate.json"), jsonBytes, 0o644); err != nil {
		return fmt.Errorf("write JSON artifact: %w", err)
	}

	spiceText, err := circuitartifacts.SPICE(program)
	if err != nil {
		return fmt.Errorf("emit SPICE artifact: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "candidate.spice"), []byte(spiceText), 0o644); err != nil {
		return fmt.Errorf("write SPICE artifact: %w", err)
	}

	verilogText, err := circuitartifacts.Verilog(program)
	if err != nil {
		return fmt.Errorf("emit Verilog artifact: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "candidate.v"), []byte(verilogText), 0o644); err != nil {
		return fmt.Errorf("write Verilog artifact: %w", err)
	}

	return nil
}
