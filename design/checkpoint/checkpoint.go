// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package checkpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gmlewis/gep/v2/design"
	"github.com/gmlewis/gep/v2/design/novelty"
	"github.com/gmlewis/gep/v2/design/objectives"
)

// currentSchemaVersion is the version written into every new snapshot.
// Increment this when the on-disk format changes in a breaking way.
const currentSchemaVersion = 1

// EliteRecord holds the metadata saved for one elite candidate at the time a
// checkpoint is written.
type EliteRecord struct {
	// CandidateID is the stable identifier of the candidate.
	CandidateID string `json:"candidate_id"`
	// Generation is the generation number at which this candidate became elite.
	Generation int `json:"generation"`
	// AggregateScore is the final scalar score used for ranking.
	AggregateScore float64 `json:"aggregate_score"`
	// Breakdown holds the per-objective contributions and hard-fail flag.
	Breakdown objectives.ScoreBreakdown `json:"breakdown"`
	// Behavior is the behavior vector archived for novelty tracking.
	Behavior novelty.BehaviorVector `json:"behavior,omitempty"`
}

// AggregateSnapshot records aggregate scoring state at checkpoint time.
type AggregateSnapshot struct {
	// BestScore is the highest aggregate score seen so far.
	BestScore float64 `json:"best_score"`
	// MeanScore is the mean aggregate score over the current population.
	MeanScore float64 `json:"mean_score,omitempty"`
	// TotalCandidatesEvaluated is the cumulative count of candidates evaluated.
	TotalCandidatesEvaluated int `json:"total_candidates_evaluated,omitempty"`
}

// NoveltySnapshot is the serializable state of a novelty archive. It stores
// the archive entries and the configuration used to create the archive so that
// an identical [novelty.Archive] can be reconstructed on replay.
type NoveltySnapshot struct {
	// Config is the archive configuration (K, MaxSize). Distance is not
	// serialized; callers must supply an appropriate DistanceFunc on replay.
	Config novelty.ArchiveConfig `json:"config"`
	// Entries contains the archive entries in insertion order.
	Entries []novelty.ArchiveEntry `json:"entries,omitempty"`
}

// Snapshot is the top-level checkpoint value. It combines all durable run
// state into a single JSON-serializable value.
type Snapshot struct {
	// SchemaVersion identifies the on-disk format. The only currently defined
	// value is 1.
	SchemaVersion int `json:"schema_version"`
	// Manifest is the run-level manifest at the time of the checkpoint.
	Manifest design.RunManifest `json:"manifest"`
	// Elites holds the elite candidate records retained at checkpoint time.
	Elites []EliteRecord `json:"elites,omitempty"`
	// Aggregate records aggregate scoring state.
	Aggregate AggregateSnapshot `json:"aggregate"`
	// Novelty records the serializable state of the novelty archive.
	Novelty NoveltySnapshot `json:"novelty"`
	// ArtifactRefs lists artifact references associated with this checkpoint.
	ArtifactRefs []design.ArtifactRef `json:"artifact_refs,omitempty"`
}

var (
	errNilWriter      = errors.New("nil checkpoint writer")
	errNilReader      = errors.New("nil checkpoint reader")
	errNilSnapshot    = errors.New("nil snapshot")
	errUnknownVersion = errors.New("unknown checkpoint schema version")
)

// Save writes snap to w as indented JSON. w must not be nil and snap must not
// be nil. Save always sets snap.SchemaVersion to the current version before
// encoding.
func Save(w io.Writer, snap *Snapshot) error {
	if w == nil {
		return errNilWriter
	}
	if snap == nil {
		return errNilSnapshot
	}
	snap.SchemaVersion = currentSchemaVersion
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("encode checkpoint JSON: %w", err)
	}
	return nil
}

// SaveFile writes snap to the named file as indented JSON, creating or
// truncating the file as needed.
func SaveFile(filename string, snap *Snapshot) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create checkpoint file %q: %w", filename, err)
	}
	defer f.Close()
	if err := Save(f, snap); err != nil {
		return fmt.Errorf("write checkpoint file %q: %w", filename, err)
	}
	return nil
}

// Load decodes one snapshot from r. It returns an error if r is nil, if the
// JSON is malformed, or if the schema version is not recognized.
func Load(r io.Reader) (*Snapshot, error) {
	if r == nil {
		return nil, errNilReader
	}
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var snap Snapshot
	if err := dec.Decode(&snap); err != nil {
		return nil, fmt.Errorf("decode checkpoint JSON: %w", err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return nil, errors.New("checkpoint JSON must contain exactly one object")
	}
	if snap.SchemaVersion != currentSchemaVersion {
		return nil, fmt.Errorf("%w: %d", errUnknownVersion, snap.SchemaVersion)
	}
	return &snap, nil
}

// LoadFile decodes one snapshot from the named JSON file.
func LoadFile(filename string) (*Snapshot, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open checkpoint file %q: %w", filename, err)
	}
	defer f.Close()
	return Load(f)
}

// ReplayManifest returns the RunManifest stored in snap. It is a convenience
// helper for callers that want to seed a new run from an existing checkpoint
// without having access to the original evaluator.
func ReplayManifest(snap *Snapshot) (design.RunManifest, error) {
	if snap == nil {
		return design.RunManifest{}, errNilSnapshot
	}
	return snap.Manifest, nil
}

// ReplayElites returns the elite records stored in snap. Callers can use
// these to bootstrap a resumed run's elite pool without re-evaluating old
// candidates.
func ReplayElites(snap *Snapshot) ([]EliteRecord, error) {
	if snap == nil {
		return nil, errNilSnapshot
	}
	result := make([]EliteRecord, len(snap.Elites))
	copy(result, snap.Elites)
	return result, nil
}

// ReplayNoveltyEntries returns the archive entries stored in snap together
// with the ArchiveConfig used when the checkpoint was written. Callers must
// supply an appropriate DistanceFunc when constructing a new novelty.Archive
// because function values are not serializable.
//
// The returned entries are in the original insertion order so that
// novelty.Archive.Add can be called in sequence to reproduce the same archive.
func ReplayNoveltyEntries(snap *Snapshot) (novelty.ArchiveConfig, []novelty.ArchiveEntry, error) {
	if snap == nil {
		return novelty.ArchiveConfig{}, nil, errNilSnapshot
	}
	entries := make([]novelty.ArchiveEntry, len(snap.Novelty.Entries))
	copy(entries, snap.Novelty.Entries)
	return snap.Novelty.Config, entries, nil
}
