// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/core"
	designscenarios "github.com/gmlewis/gep/v2/design/scenarios"
	boolNodes "github.com/gmlewis/gep/v2/functions/bool_nodes"
)

func TestBracketPipelineEvaluatorDecoderAndArtifacts(t *testing.T) {
	train, err := loadTrainScenarios()
	if err != nil {
		t.Fatalf("loadTrainScenarios() error = %v", err)
	}

	genome := mustGenomeFromSymbols(t, []string{"And", "d0", "d3"})
	if got := scoreCandidate(genome, train); got <= 0 {
		t.Fatalf("scoreCandidate() = %.2f, want > 0", got)
	}

	program, err := decodeVoxelProgram("test-candidate", genome, train[0])
	if err != nil {
		t.Fatalf("decodeVoxelProgram() error = %v", err)
	}
	if got := len(program.Design.Occupied); got == 0 {
		t.Fatalf("len(program.Design.Occupied) = %d, want > 0", got)
	}

	outDir := t.TempDir()
	if err := exportArtifacts(program, outDir); err != nil {
		t.Fatalf("exportArtifacts() error = %v", err)
	}

	jsonData := readFile(t, filepath.Join(outDir, candidateJSONArtifactName))
	objData := readFile(t, filepath.Join(outDir, candidateOBJArtifactName))
	summaryData := readFile(t, filepath.Join(outDir, candidateSummaryArtifactName))

	if !strings.Contains(jsonData, `"candidate_id": "test-candidate"`) {
		t.Fatalf("candidate.json missing candidate ID: %q", jsonData)
	}
	if !strings.Contains(objData, "o bracket") {
		t.Fatalf("candidate.obj missing object name: %q", objData)
	}
	if !strings.Contains(summaryData, "occupied:") {
		t.Fatalf("candidate.txt missing occupancy summary: %q", summaryData)
	}

	if err := exportArtifacts(program, outDir); err != nil {
		t.Fatalf("exportArtifacts() second call error = %v", err)
	}
	if got := readFile(t, filepath.Join(outDir, candidateJSONArtifactName)); got != jsonData {
		t.Fatal("candidate.json changed between deterministic exports")
	}
}

func TestScoreScenarioRejectsOverMaxCells(t *testing.T) {
	scenario := designscenarios.Scenario{
		ID:    "tiny",
		Split: designscenarios.Train,
		Params: map[string]any{
			"volume": map[string]any{
				"size_x": 4.0,
				"size_y": 3.0,
				"size_z": 2.0,
			},
			"max_cells": 1.0,
		},
	}

	alwaysFill := mustGenomeFromSymbols(t, []string{"Or", "d0", "Not", "d0"})
	if got := scoreScenario(alwaysFill, scenario); got != 0 {
		t.Fatalf("scoreScenario(alwaysFill) = %.2f, want 0", got)
	}
}

func TestRunPilotDeterministicFixedSeed(t *testing.T) {
	cfg := runConfig{
		Seed:           20260511,
		PopulationSize: 64,
		Generations:    100,
	}
	outDir1 := t.TempDir()
	cfg.OutputDir = outDir1
	got1, err := runPilot(cfg)
	if err != nil {
		t.Fatalf("runPilot() #1 error = %v", err)
	}

	outDir2 := t.TempDir()
	cfg.OutputDir = outDir2
	got2, err := runPilot(cfg)
	if err != nil {
		t.Fatalf("runPilot() #2 error = %v", err)
	}

	if !reflect.DeepEqual(got1, got2) {
		t.Fatalf("runPilot() nondeterministic:\nrun1=%+v\nrun2=%+v", got1, got2)
	}

	for _, name := range []string{
		candidateJSONArtifactName,
		candidateOBJArtifactName,
		candidateSummaryArtifactName,
	} {
		p1 := filepath.Join(outDir1, name)
		p2 := filepath.Join(outDir2, name)
		if _, err := os.Stat(p1); err != nil {
			t.Fatalf("os.Stat(%q) error = %v", p1, err)
		}
		if _, err := os.Stat(p2); err != nil {
			t.Fatalf("os.Stat(%q) error = %v", p2, err)
		}
		if got, want := readFile(t, p1), readFile(t, p2); got != want {
			t.Fatalf("%s differs across deterministic runs", name)
		}
	}
}

func mustGenomeFromSymbols(t *testing.T, geneSymbols ...[]string) core.Genome[bool] {
	t.Helper()
	cat, err := boolNodes.CatalogFromNames([]string{"Not", "And", "Or", "Xor"})
	if err != nil {
		t.Fatalf("CatalogFromNames() error = %v", err)
	}
	link, err := boolNodes.LinkFuncFrom("Or")
	if err != nil {
		t.Fatalf("LinkFuncFrom() error = %v", err)
	}

	genes := make([]core.Gene[bool], 0, len(geneSymbols))
	for i, symbols := range geneSymbols {
		parsed, err := core.ParseSymbols(symbols, cat)
		if err != nil {
			t.Fatalf("ParseSymbols(gene %d) error = %v", i, err)
		}
		genes = append(genes, core.Gene[bool]{Symbols: parsed})
	}
	return core.Genome[bool]{Genes: genes, Link: link}
}

func readFile(t *testing.T, filename string) string {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", filename, err)
	}
	return string(data)
}
