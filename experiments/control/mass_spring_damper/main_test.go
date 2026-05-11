// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/design"
	"github.com/gmlewis/gep/v2/design/checkpoint"
	"github.com/gmlewis/gep/v2/design/promotion"
	boolNodes "github.com/gmlewis/gep/v2/functions/bool_nodes"
)

func TestMassSpringDamperPipelineSimulatorAndArtifacts(t *testing.T) {
	train, err := loadTrainScenarios()
	if err != nil {
		t.Fatalf("loadTrainScenarios() error = %v", err)
	}
	genome := mustGenomeFromSymbols(t, []string{"Or", "d0", "Not", "d1"})
	if got := scoreCandidate(genome, train); got <= 0 {
		t.Fatalf("scoreCandidate() = %.2f, want > 0", got)
	}

	plant, err := scenarioPlant(train[0])
	if err != nil {
		t.Fatalf("scenarioPlant() error = %v", err)
	}
	sim, err := simulateController(genome, plant)
	if err != nil {
		t.Fatalf("simulateController() error = %v", err)
	}
	if sim.HardFailed {
		t.Fatal("simulateController() hard failed, want stable trajectory")
	}

	policy := decodeControllerPolicy("test-candidate", genome)
	policy.Scores = map[string]float64{"train_mean": 1, "validation_mean": 2, "test_mean": 3}
	outDir := t.TempDir()
	artifactRefs, err := exportArtifacts(policy, outDir)
	if err != nil {
		t.Fatalf("exportArtifacts() error = %v", err)
	}
	if got, want := len(artifactRefs), 2; got != want {
		t.Fatalf("len(artifactRefs) = %d, want %d", got, want)
	}

	jsonData := readFile(t, filepath.Join(outDir, controllerJSONArtifactName))
	summaryData := readFile(t, filepath.Join(outDir, controllerSummaryArtifact))
	if !strings.Contains(jsonData, `"candidate_id": "test-candidate"`) {
		t.Fatalf("controller_policy.json missing candidate ID: %q", jsonData)
	}
	if !strings.Contains(summaryData, "karva:") {
		t.Fatalf("controller_summary.txt missing karva line: %q", summaryData)
	}

	if _, err := exportArtifacts(policy, outDir); err != nil {
		t.Fatalf("exportArtifacts() second call error = %v", err)
	}
	if got := readFile(t, filepath.Join(outDir, controllerJSONArtifactName)); got != jsonData {
		t.Fatal("controller_policy.json changed between deterministic exports")
	}
}

func TestRunPilotDeterministicFixedSeed(t *testing.T) {
	cfg := runConfig{
		Seed:           20260511,
		PopulationSize: 72,
		Generations:    140,
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

	for _, name := range []string{controllerJSONArtifactName, controllerSummaryArtifact, promotionReportArtifactName} {
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

func TestRunPilotPromotionCheckpointAndManifest(t *testing.T) {
	outDir := t.TempDir()
	result, err := runPilot(runConfig{
		Seed:           20260511,
		PopulationSize: 72,
		Generations:    140,
		OutputDir:      outDir,
	})
	if err != nil {
		t.Fatalf("runPilot() error = %v", err)
	}
	if !result.Promoted {
		t.Fatalf("runPilot() promoted = false, want true")
	}

	reportPath := filepath.Join(outDir, promotionReportArtifactName)
	reportData := readFile(t, reportPath)
	var report promotion.PromotionReport
	if err := json.Unmarshal([]byte(reportData), &report); err != nil {
		t.Fatalf("json.Unmarshal(promotion report) error = %v", err)
	}
	if !report.Promoted {
		t.Fatalf("promotion report promoted = false, want true")
	}
	if got, want := report.CandidateID, result.CandidateID; got != want {
		t.Fatalf("promotion report candidate_id = %q, want %q", got, want)
	}
	if got, want := len(report.Decisions), 3; got != want {
		t.Fatalf("promotion report decisions len = %d, want %d", got, want)
	}
	for _, d := range report.Decisions {
		if !d.Passed {
			t.Fatalf("promotion decision for split %q failed: %s", d.Split, d.Reason)
		}
	}

	manifest, err := design.LoadRunManifestFile(filepath.Join(outDir, runManifestArtifactName))
	if err != nil {
		t.Fatalf("LoadRunManifestFile() error = %v", err)
	}
	if got, want := len(manifest.Artifacts), 4; got != want {
		t.Fatalf("manifest artifacts len = %d, want %d", got, want)
	}
	if got, want := len(manifest.ScenarioSplits), 3; got != want {
		t.Fatalf("manifest scenario splits len = %d, want %d", got, want)
	}
	for _, ref := range manifest.Artifacts {
		if ref.Path == "" {
			t.Fatalf("manifest artifact %q has empty path", ref.Name)
		}
		if _, err := os.Stat(ref.Path); err != nil {
			t.Fatalf("os.Stat(%q) error = %v", ref.Path, err)
		}
	}

	snap, err := checkpoint.LoadFile(filepath.Join(outDir, checkpointArtifactName))
	if err != nil {
		t.Fatalf("checkpoint.LoadFile() error = %v", err)
	}
	replayedManifest, err := checkpoint.ReplayManifest(snap)
	if err != nil {
		t.Fatalf("checkpoint.ReplayManifest() error = %v", err)
	}
	if !reflect.DeepEqual(replayedManifest, *manifest) {
		t.Fatalf("checkpoint manifest mismatch:\n got=%#v\nwant=%#v", replayedManifest, *manifest)
	}
	if len(snap.ArtifactRefs) < len(manifest.Artifacts) {
		t.Fatalf("checkpoint artifact refs len = %d, want at least %d", len(snap.ArtifactRefs), len(manifest.Artifacts))
	}
	if !reflect.DeepEqual(snap.ArtifactRefs[:len(manifest.Artifacts)], manifest.Artifacts) {
		t.Fatalf("checkpoint artifact refs do not preserve manifest artifact refs")
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
