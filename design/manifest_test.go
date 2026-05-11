// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package design

import (
	"bytes"
	"path/filepath"
	"reflect"
	"testing"
)

func mustLoadFixture(t *testing.T, fixture string) *RunManifest {
	t.Helper()

	m, err := LoadRunManifestFile(filepath.Join("testdata", fixture))
	if err != nil {
		t.Fatalf("LoadRunManifestFile(%q): %v", fixture, err)
	}
	return m
}

func roundTripManifest(t *testing.T, m *RunManifest) *RunManifest {
	t.Helper()

	var buf bytes.Buffer
	if err := WriteRunManifest(&buf, m); err != nil {
		t.Fatalf("WriteRunManifest: %v", err)
	}
	roundTrip, err := LoadRunManifest(&buf)
	if err != nil {
		t.Fatalf("LoadRunManifest: %v", err)
	}
	return roundTrip
}

func TestRunManifestRoundTrip_MinimalFixture(t *testing.T) {
	orig := mustLoadFixture(t, "manifest_minimal.json")
	roundTrip := roundTripManifest(t, orig)
	if !reflect.DeepEqual(roundTrip, orig) {
		t.Fatalf("round-trip mismatch:\n got: %#v\nwant: %#v", roundTrip, orig)
	}
}

func TestRunManifestRoundTrip_PopulatedFixture(t *testing.T) {
	orig := mustLoadFixture(t, "manifest_populated.json")
	roundTrip := roundTripManifest(t, orig)
	if !reflect.DeepEqual(roundTrip, orig) {
		t.Fatalf("round-trip mismatch:\n got: %#v\nwant: %#v", roundTrip, orig)
	}
	if !reflect.DeepEqual(roundTrip.RunConfig, orig.RunConfig) {
		t.Fatalf("run config mismatch: got %#v want %#v", roundTrip.RunConfig, orig.RunConfig)
	}
	if !reflect.DeepEqual(roundTrip.Seeds, orig.Seeds) {
		t.Fatalf("seed metadata mismatch: got %#v want %#v", roundTrip.Seeds, orig.Seeds)
	}
	if !reflect.DeepEqual(roundTrip.Artifacts, orig.Artifacts) {
		t.Fatalf("artifact metadata mismatch: got %#v want %#v", roundTrip.Artifacts, orig.Artifacts)
	}
	if !reflect.DeepEqual(roundTrip.ScenarioSplits, orig.ScenarioSplits) {
		t.Fatalf("scenario split metadata mismatch: got %#v want %#v", roundTrip.ScenarioSplits, orig.ScenarioSplits)
	}
}

func TestRunManifestFileRoundTrip_PopulatedFixture(t *testing.T) {
	orig := mustLoadFixture(t, "manifest_populated.json")

	filename := filepath.Join(t.TempDir(), "manifest.json")
	if err := WriteRunManifestFile(filename, orig); err != nil {
		t.Fatalf("WriteRunManifestFile: %v", err)
	}

	roundTrip, err := LoadRunManifestFile(filename)
	if err != nil {
		t.Fatalf("LoadRunManifestFile: %v", err)
	}

	if !reflect.DeepEqual(roundTrip, orig) {
		t.Fatalf("file round-trip mismatch:\n got: %#v\nwant: %#v", roundTrip, orig)
	}
}
