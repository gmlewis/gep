// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package design

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// RunManifest stores durable metadata for one applied-design run.
type RunManifest struct {
	RunID          string                 `json:"run_id,omitempty"`
	RunConfig      RunConfig              `json:"run_config"`
	Seeds          []SeedRecord           `json:"seeds,omitempty"`
	Artifacts      []ArtifactRef          `json:"artifacts,omitempty"`
	ScenarioSplits []ScenarioSplitSummary `json:"scenario_splits,omitempty"`
}

// RunConfig records run-level configuration used to reproduce a run.
type RunConfig struct {
	Domain         string `json:"domain,omitempty"`
	Experiment     string `json:"experiment,omitempty"`
	PopulationSize int    `json:"population_size,omitempty"`
	NumGenerations int    `json:"num_generations,omitempty"`
}

// ArtifactRef references an emitted artifact from a run.
type ArtifactRef struct {
	Name   string `json:"name,omitempty"`
	Path   string `json:"path,omitempty"`
	Format string `json:"format,omitempty"`
	Kind   string `json:"kind,omitempty"`
}

// ScenarioSplitSummary records aggregate scenario metadata for one split.
type ScenarioSplitSummary struct {
	SplitName     string `json:"split_name,omitempty"`
	ScenarioCount int    `json:"scenario_count,omitempty"`
	Source        string `json:"source,omitempty"`
}

// SeedRecord records one deterministic seed used during a run.
type SeedRecord struct {
	Name    string `json:"name,omitempty"`
	Value   int64  `json:"value"`
	Purpose string `json:"purpose,omitempty"`
}

var (
	errNilReader   = errors.New("nil run manifest reader")
	errNilWriter   = errors.New("nil run manifest writer")
	errNilManifest = errors.New("nil run manifest")
)

// LoadRunManifest decodes one run manifest from JSON.
func LoadRunManifest(r io.Reader) (*RunManifest, error) {
	if r == nil {
		return nil, errNilReader
	}

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var m RunManifest
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("decode run manifest JSON: %w", err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return nil, errors.New("run manifest JSON must contain exactly one object")
	}
	return &m, nil
}

// LoadRunManifestFile decodes one run manifest from a JSON file.
func LoadRunManifestFile(filename string) (*RunManifest, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open run manifest file %q: %w", filename, err)
	}
	defer f.Close()

	return LoadRunManifest(f)
}

// WriteRunManifest writes one run manifest as indented JSON.
func WriteRunManifest(w io.Writer, m *RunManifest) error {
	if w == nil {
		return errNilWriter
	}
	if m == nil {
		return errNilManifest
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(m); err != nil {
		return fmt.Errorf("encode run manifest JSON: %w", err)
	}
	return nil
}

// WriteRunManifestFile writes one run manifest as indented JSON to a file.
func WriteRunManifestFile(filename string, m *RunManifest) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create run manifest file %q: %w", filename, err)
	}
	defer f.Close()

	if err := WriteRunManifest(f, m); err != nil {
		return fmt.Errorf("write run manifest file %q: %w", filename, err)
	}
	return nil
}
