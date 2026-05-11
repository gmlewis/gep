// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/core"
	designscenarios "github.com/gmlewis/gep/v2/design/scenarios"
	boolNodes "github.com/gmlewis/gep/v2/functions/bool_nodes"
)

func TestHalfAdderPipelineEvaluatorDecoderAndArtifacts(t *testing.T) {
	train, err := loadTrainScenarios()
	if err != nil {
		t.Fatalf("loadTrainScenarios() error = %v", err)
	}

	genome := mustGenomeFromSymbols(t,
		[]string{"Xor", "d0", "d1"},
		[]string{"And", "d0", "d1"},
	)

	if got := scoreCandidate(genome, train); got != 1000 {
		t.Fatalf("scoreCandidate(perfect genome) = %v, want 1000", got)
	}

	program, gateCount, err := decodeCircuitProgram("test-candidate", genome)
	if err != nil {
		t.Fatalf("decodeCircuitProgram() error = %v", err)
	}
	if got, want := gateCount, 2; got != want {
		t.Fatalf("gateCount = %d, want %d", got, want)
	}

	outDir := t.TempDir()
	if err := exportArtifacts(program, outDir); err != nil {
		t.Fatalf("exportArtifacts() error = %v", err)
	}

	jsonData := readFile(t, filepath.Join(outDir, "candidate.json"))
	spiceData := readFile(t, filepath.Join(outDir, "candidate.spice"))
	verilogData := readFile(t, filepath.Join(outDir, "candidate.v"))

	if !strings.Contains(jsonData, "\"gate.xor\"") {
		t.Fatalf("candidate.json missing XOR gate: %q", jsonData)
	}
	if !strings.Contains(spiceData, "gate.xor") {
		t.Fatalf("candidate.spice missing XOR gate: %q", spiceData)
	}
	if !strings.Contains(verilogData, "gate_xor") {
		t.Fatalf("candidate.v missing XOR gate: %q", verilogData)
	}

	if err := exportArtifacts(program, outDir); err != nil {
		t.Fatalf("exportArtifacts() second call error = %v", err)
	}
	if got := readFile(t, filepath.Join(outDir, "candidate.json")); got != jsonData {
		t.Fatalf("candidate.json changed between deterministic exports")
	}
}

func TestScoreCandidateRespectsScenarioComponentBounds(t *testing.T) {
	train, err := loadTrainScenarios()
	if err != nil {
		t.Fatalf("loadTrainScenarios() error = %v", err)
	}
	if got, want := scenarioMaxComponents(train[0]), 3; got != want {
		t.Fatalf("scenarioMaxComponents(train[0]) = %d, want %d", got, want)
	}

	genome := mustGenomeFromSymbols(t,
		[]string{"Xor", "And", "Or", "d0", "d1", "d0", "d1"},
		[]string{"And", "d0", "d1"},
	)
	if got := scoreCandidate(genome, train); got != 500 {
		t.Fatalf("scoreCandidate(oversized genome) = %v, want 500", got)
	}
}

func mustGenomeFromSymbols(t *testing.T, genes ...[]string) core.Genome[bool] {
	t.Helper()
	cat, err := boolNodes.CatalogFromNames([]string{"Not", "And", "Or", "Xor"})
	if err != nil {
		t.Fatalf("CatalogFromNames() error = %v", err)
	}
	link, err := boolNodes.LinkFuncFrom("Or")
	if err != nil {
		t.Fatalf("LinkFuncFrom() error = %v", err)
	}

	parsedGenes := make([]core.Gene[bool], 0, len(genes))
	for i, symbols := range genes {
		parsed, err := core.ParseSymbols(symbols, cat)
		if err != nil {
			t.Fatalf("ParseSymbols(gene %d) error = %v", i, err)
		}
		parsedGenes = append(parsedGenes, core.Gene[bool]{Symbols: parsed})
	}
	return core.Genome[bool]{Genes: parsedGenes, Link: link}
}

func readFile(t *testing.T, filename string) string {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", filename, err)
	}
	return string(data)
}

func TestScenarioMaxComponents(t *testing.T) {
	if got := scenarioMaxComponents(designscenarios.Scenario{}); got != 0 {
		t.Fatalf("scenarioMaxComponents(empty) = %d, want 0", got)
	}
	if got := scenarioMaxComponents(designscenarios.Scenario{Params: map[string]any{"max_components": 4.0}}); got != 4 {
		t.Fatalf("scenarioMaxComponents(float64) = %d, want 4", got)
	}
	if got := scenarioMaxComponents(designscenarios.Scenario{Params: map[string]any{"max_components": "oops"}}); got != 0 {
		t.Fatalf("scenarioMaxComponents(invalid) = %d, want 0", got)
	}
}
