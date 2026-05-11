// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package checkpoint_test

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/design"
	"github.com/gmlewis/gep/v2/design/checkpoint"
	"github.com/gmlewis/gep/v2/design/novelty"
	"github.com/gmlewis/gep/v2/design/objectives"
)

// --- helpers ---

func makeSnapshot() *checkpoint.Snapshot {
	return &checkpoint.Snapshot{
		Manifest: design.RunManifest{
			RunID: "run-001",
			RunConfig: design.RunConfig{
				Domain:         "test-domain",
				Experiment:     "test-exp",
				PopulationSize: 50,
				NumGenerations: 10,
			},
			Seeds: []design.SeedRecord{
				{Name: "rng", Value: 42, Purpose: "population init"},
			},
			Artifacts: []design.ArtifactRef{
				{Name: "best", Path: "/tmp/best.json", Format: "json", Kind: "elite"},
			},
		},
		Elites: []checkpoint.EliteRecord{
			{
				CandidateID:    "c-1",
				Generation:     3,
				AggregateScore: 0.95,
				Breakdown: objectives.ScoreBreakdown{
					Contributions: []objectives.WeightedContribution{
						{Name: "accuracy", RawScore: 0.95, Weight: 1.0, WeightedScore: 0.95},
					},
				},
				Behavior: novelty.BehaviorVector{0.1, 0.2, 0.3},
			},
			{
				CandidateID:    "c-2",
				Generation:     5,
				AggregateScore: 0.88,
				Breakdown: objectives.ScoreBreakdown{
					Contributions: []objectives.WeightedContribution{
						{Name: "accuracy", RawScore: 0.88, Weight: 1.0, WeightedScore: 0.88},
					},
				},
				Behavior: novelty.BehaviorVector{0.4, 0.5, 0.6},
			},
		},
		Aggregate: checkpoint.AggregateSnapshot{
			BestScore:                0.95,
			MeanScore:                0.75,
			TotalCandidatesEvaluated: 500,
		},
		Novelty: checkpoint.NoveltySnapshot{
			Config: novelty.ArchiveConfig{K: 3, MaxSize: 100},
			Entries: []novelty.ArchiveEntry{
				{Behavior: novelty.BehaviorVector{0.1, 0.2, 0.3}, Label: "c-1"},
				{Behavior: novelty.BehaviorVector{0.4, 0.5, 0.6}, Label: "c-2"},
			},
		},
		ArtifactRefs: []design.ArtifactRef{
			{Name: "checkpoint", Path: "/tmp/ckpt.json", Format: "json", Kind: "checkpoint"},
		},
	}
}

// --- Save / Load round-trip ---

func TestRoundTripBuffer(t *testing.T) {
	original := makeSnapshot()

	var buf bytes.Buffer
	if err := checkpoint.Save(&buf, original); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	restored, err := checkpoint.Load(&buf)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if !reflect.DeepEqual(original, restored) {
		t.Fatalf("round-trip mismatch:\n  original=%#v\n  restored=%#v", original, restored)
	}
}

func TestSchemaVersionSetOnSave(t *testing.T) {
	snap := makeSnapshot()
	snap.SchemaVersion = 0 // deliberately unset

	var buf bytes.Buffer
	if err := checkpoint.Save(&buf, snap); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	restored, err := checkpoint.Load(&buf)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if restored.SchemaVersion != 1 {
		t.Fatalf("SchemaVersion = %d, want 1", restored.SchemaVersion)
	}
}

// --- Error cases ---

func TestSaveNilWriter(t *testing.T) {
	if err := checkpoint.Save(nil, makeSnapshot()); err == nil {
		t.Fatal("Save(nil writer) should return error")
	}
}

func TestSaveNilSnapshot(t *testing.T) {
	var buf bytes.Buffer
	if err := checkpoint.Save(&buf, nil); err == nil {
		t.Fatal("Save(nil snapshot) should return error")
	}
}

func TestLoadNilReader(t *testing.T) {
	if _, err := checkpoint.Load(nil); err == nil {
		t.Fatal("Load(nil reader) should return error")
	}
}

func TestLoadUnknownVersion(t *testing.T) {
	json := `{"schema_version":99,"manifest":{},"aggregate":{},"novelty":{"config":{}}}`
	_, err := checkpoint.Load(strings.NewReader(json))
	if err == nil {
		t.Fatal("Load with unknown schema_version should return error")
	}
}

func TestLoadMultipleObjects(t *testing.T) {
	// Two JSON objects in the stream should be rejected.
	json := `{"schema_version":1,"manifest":{},"aggregate":{},"novelty":{"config":{}}}` +
		"\n" + `{"schema_version":1,"manifest":{},"aggregate":{},"novelty":{"config":{}}}`
	_, err := checkpoint.Load(strings.NewReader(json))
	if err == nil {
		t.Fatal("Load with two JSON objects should return error")
	}
}

// --- File round-trip ---

func TestRoundTripFile(t *testing.T) {
	original := makeSnapshot()

	dir := t.TempDir()
	path := dir + "/ckpt.json"

	if err := checkpoint.SaveFile(path, original); err != nil {
		t.Fatalf("SaveFile error: %v", err)
	}

	restored, err := checkpoint.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile error: %v", err)
	}

	if !reflect.DeepEqual(original, restored) {
		t.Fatalf("file round-trip mismatch:\n  original=%#v\n  restored=%#v", original, restored)
	}
}

func TestSaveFileInvalidPath(t *testing.T) {
	if err := checkpoint.SaveFile("/nonexistent/dir/ckpt.json", makeSnapshot()); err == nil {
		t.Fatal("SaveFile with invalid path should return error")
	}
}

func TestLoadFileInvalidPath(t *testing.T) {
	if _, err := checkpoint.LoadFile("/nonexistent/dir/ckpt.json"); err == nil {
		t.Fatal("LoadFile with invalid path should return error")
	}
}

// --- Replay helpers ---

func TestReplayManifest(t *testing.T) {
	original := makeSnapshot()
	var buf bytes.Buffer
	if err := checkpoint.Save(&buf, original); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	restored, err := checkpoint.Load(&buf)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	manifest, err := checkpoint.ReplayManifest(restored)
	if err != nil {
		t.Fatalf("ReplayManifest error: %v", err)
	}
	if !reflect.DeepEqual(manifest, original.Manifest) {
		t.Fatalf("ReplayManifest mismatch:\n  got=%#v\n  want=%#v", manifest, original.Manifest)
	}
}

func TestReplayManifestNilSnapshot(t *testing.T) {
	if _, err := checkpoint.ReplayManifest(nil); err == nil {
		t.Fatal("ReplayManifest(nil) should return error")
	}
}

func TestReplayElites(t *testing.T) {
	original := makeSnapshot()
	var buf bytes.Buffer
	if err := checkpoint.Save(&buf, original); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	restored, err := checkpoint.Load(&buf)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	elites, err := checkpoint.ReplayElites(restored)
	if err != nil {
		t.Fatalf("ReplayElites error: %v", err)
	}
	if !reflect.DeepEqual(elites, original.Elites) {
		t.Fatalf("ReplayElites mismatch:\n  got=%#v\n  want=%#v", elites, original.Elites)
	}
}

func TestReplayElitesNilSnapshot(t *testing.T) {
	if _, err := checkpoint.ReplayElites(nil); err == nil {
		t.Fatal("ReplayElites(nil) should return error")
	}
}

func TestReplayElitesReturnsCopy(t *testing.T) {
	snap := makeSnapshot()
	elites, err := checkpoint.ReplayElites(snap)
	if err != nil {
		t.Fatalf("ReplayElites error: %v", err)
	}
	elites[0].CandidateID = "mutated"
	if snap.Elites[0].CandidateID == "mutated" {
		t.Fatal("ReplayElites should return a copy, not a reference")
	}
}

func TestReplayNoveltyEntries(t *testing.T) {
	original := makeSnapshot()
	var buf bytes.Buffer
	if err := checkpoint.Save(&buf, original); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	restored, err := checkpoint.Load(&buf)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	cfg, entries, err := checkpoint.ReplayNoveltyEntries(restored)
	if err != nil {
		t.Fatalf("ReplayNoveltyEntries error: %v", err)
	}
	if !reflect.DeepEqual(cfg, original.Novelty.Config) {
		t.Fatalf("ReplayNoveltyEntries config mismatch:\n  got=%#v\n  want=%#v", cfg, original.Novelty.Config)
	}
	if !reflect.DeepEqual(entries, original.Novelty.Entries) {
		t.Fatalf("ReplayNoveltyEntries entries mismatch:\n  got=%#v\n  want=%#v", entries, original.Novelty.Entries)
	}
}

func TestReplayNoveltyEntriesNilSnapshot(t *testing.T) {
	if _, _, err := checkpoint.ReplayNoveltyEntries(nil); err == nil {
		t.Fatal("ReplayNoveltyEntries(nil) should return error")
	}
}

func TestReplayNoveltyEntriesReturnsCopy(t *testing.T) {
	snap := makeSnapshot()
	_, entries, err := checkpoint.ReplayNoveltyEntries(snap)
	if err != nil {
		t.Fatalf("ReplayNoveltyEntries error: %v", err)
	}
	entries[0].Label = "mutated"
	if snap.Novelty.Entries[0].Label == "mutated" {
		t.Fatal("ReplayNoveltyEntries should return a copy, not a reference")
	}
}

// --- Novelty archive re-seeding ---

// TestReplayCanSeedNoveltyArchive demonstrates that a replay can restore the
// novelty archive from a checkpoint and produce the same Score results as the
// original archive.
func TestReplayCanSeedNoveltyArchive(t *testing.T) {
	original := makeSnapshot()

	// Build an archive from the original snapshot entries.
	buildArchive := func(cfg novelty.ArchiveConfig, entries []novelty.ArchiveEntry) *novelty.Archive {
		a := novelty.NewArchive(cfg)
		for _, e := range entries {
			a.Add(e)
		}
		return a
	}

	origCfg, origEntries, err := checkpoint.ReplayNoveltyEntries(original)
	if err != nil {
		t.Fatalf("ReplayNoveltyEntries error: %v", err)
	}
	originalArchive := buildArchive(origCfg, origEntries)

	// Save and reload the snapshot.
	var buf bytes.Buffer
	if err := checkpoint.Save(&buf, original); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	restored, err := checkpoint.Load(&buf)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	replayCfg, replayEntries, err := checkpoint.ReplayNoveltyEntries(restored)
	if err != nil {
		t.Fatalf("ReplayNoveltyEntries (restored) error: %v", err)
	}
	restoredArchive := buildArchive(replayCfg, replayEntries)

	// Both archives should produce identical scores for the same query.
	query := novelty.BehaviorVector{0.25, 0.35, 0.45}
	origScore := originalArchive.Score(query)
	restoredScore := restoredArchive.Score(query)

	if !reflect.DeepEqual(origScore, restoredScore) {
		t.Fatalf("novelty score mismatch after replay:\n  original=%#v\n  restored=%#v",
			origScore, restoredScore)
	}
}

// --- Determinism ---

// TestSaveDeterministic ensures that saving the same snapshot twice produces
// identical byte output.
func TestSaveDeterministic(t *testing.T) {
	snap := makeSnapshot()

	var buf1, buf2 bytes.Buffer
	if err := checkpoint.Save(&buf1, snap); err != nil {
		t.Fatalf("Save 1 error: %v", err)
	}
	if err := checkpoint.Save(&buf2, snap); err != nil {
		t.Fatalf("Save 2 error: %v", err)
	}

	if !bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		t.Fatal("Save is not deterministic: two identical calls produced different output")
	}
}

// --- Empty elites / empty novelty ---

func TestRoundTripEmptyElitesAndNovelty(t *testing.T) {
	snap := &checkpoint.Snapshot{
		Manifest: design.RunManifest{
			RunID: "run-empty",
			RunConfig: design.RunConfig{
				Domain:         "test",
				PopulationSize: 10,
			},
		},
		Aggregate: checkpoint.AggregateSnapshot{BestScore: 0},
		Novelty:   checkpoint.NoveltySnapshot{Config: novelty.ArchiveConfig{K: 1}},
	}

	var buf bytes.Buffer
	if err := checkpoint.Save(&buf, snap); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	restored, err := checkpoint.Load(&buf)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if !reflect.DeepEqual(snap, restored) {
		t.Fatalf("empty snapshot round-trip mismatch:\n  original=%#v\n  restored=%#v", snap, restored)
	}
}
