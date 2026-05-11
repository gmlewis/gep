// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package novelty

import (
	"math"
	"sort"
)

// BehaviorVector encodes the observable behavior of one candidate design as a
// slice of float64 values. The length of all vectors in a single archive must
// be consistent, but the package does not enforce this at the type level.
type BehaviorVector []float64

// ArchiveEntry pairs a BehaviorVector with an optional label used for
// book-keeping (e.g. a candidate ID or generation number).
type ArchiveEntry struct {
	// Behavior is the behavior vector for this archive entry.
	Behavior BehaviorVector `json:"behavior"`
	// Label is an optional human-readable identifier for the candidate.
	Label string `json:"label,omitempty"`
}

// DistanceFunc measures the dissimilarity between two BehaviorVectors.
// It must be symmetric and return a non-negative value.
// The default distance function used by NewArchive is SquaredEuclidean.
type DistanceFunc func(a, b BehaviorVector) float64

// SquaredEuclidean returns the squared Euclidean distance between two
// BehaviorVectors. If the vectors differ in length the shorter one is
// zero-padded conceptually: missing dimensions contribute their full squared
// value from the longer vector.
func SquaredEuclidean(a, b BehaviorVector) float64 {
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	var sum float64
	for i := range n {
		var ai, bi float64
		if i < len(a) {
			ai = a[i]
		}
		if i < len(b) {
			bi = b[i]
		}
		d := ai - bi
		sum += d * d
	}
	return sum
}

// EuclideanDistance returns the Euclidean distance between two BehaviorVectors.
// If the vectors differ in length the shorter one is zero-padded conceptually.
func EuclideanDistance(a, b BehaviorVector) float64 {
	return math.Sqrt(SquaredEuclidean(a, b))
}

// ArchiveConfig controls how the novelty archive grows and how k-nearest-
// neighbor novelty is computed.
type ArchiveConfig struct {
	// K is the number of nearest neighbors used when computing the novelty
	// score. If K is zero it defaults to 1.
	K int `json:"k"`
	// MaxSize caps the total number of entries stored in the archive. When the
	// archive is full, Add is a no-op. A value of zero means unbounded.
	MaxSize int `json:"max_size,omitempty"`
	// Distance is the distance function used for k-NN search. If nil,
	// SquaredEuclidean is used.
	Distance DistanceFunc `json:"-"`
}

// NoveltyResult is the outcome of one novelty-score query.
type NoveltyResult struct {
	// Score is the mean distance to the K nearest neighbors in the archive.
	// It is zero when the archive is empty.
	Score float64 `json:"score"`
	// NeighborDistances contains the distances to each of the K nearest
	// neighbors, in ascending order. It is empty when the archive is empty.
	NeighborDistances []float64 `json:"neighbor_distances,omitempty"`
}

// Archive is the live novelty store. It is not safe for concurrent use without
// external synchronization.
type Archive struct {
	cfg      ArchiveConfig
	entries  []ArchiveEntry
	distFunc DistanceFunc
}

// NewArchive creates a new Archive with the given configuration.
// If cfg.Distance is nil, SquaredEuclidean is used.
// If cfg.K is zero, it defaults to 1.
func NewArchive(cfg ArchiveConfig) *Archive {
	distFunc := cfg.Distance
	if distFunc == nil {
		distFunc = SquaredEuclidean
	}
	k := cfg.K
	if k <= 0 {
		k = 1
	}
	return &Archive{
		cfg: ArchiveConfig{
			K:        k,
			MaxSize:  cfg.MaxSize,
			Distance: distFunc,
		},
		distFunc: distFunc,
	}
}

// Len returns the number of entries currently in the archive.
func (a *Archive) Len() int {
	return len(a.entries)
}

// Entries returns a copy of the current archive entries in insertion order.
func (a *Archive) Entries() []ArchiveEntry {
	result := make([]ArchiveEntry, len(a.entries))
	copy(result, a.entries)
	return result
}

// Add inserts a new entry into the archive. If the archive has reached its
// MaxSize the entry is silently dropped. A MaxSize of zero means unbounded.
func (a *Archive) Add(entry ArchiveEntry) {
	if a.cfg.MaxSize > 0 && len(a.entries) >= a.cfg.MaxSize {
		return
	}
	a.entries = append(a.entries, entry)
}

// Score computes the novelty score for query against the current archive
// contents using the configured distance function and K value.
//
// The novelty score is the mean distance to the K nearest neighbors. When the
// archive contains fewer than K entries the mean is computed over all available
// entries. When the archive is empty Score returns a zero NoveltyResult.
func (a *Archive) Score(query BehaviorVector) NoveltyResult {
	if len(a.entries) == 0 {
		return NoveltyResult{}
	}

	// Compute distances to all archive entries.
	dists := make([]float64, len(a.entries))
	for i, e := range a.entries {
		dists[i] = a.distFunc(query, e.Behavior)
	}

	// Sort distances to find the K nearest.
	sort.Float64s(dists)

	k := a.cfg.K
	if k > len(dists) {
		k = len(dists)
	}

	neighbors := dists[:k]
	var total float64
	for _, d := range neighbors {
		total += d
	}

	neighborCopy := make([]float64, k)
	copy(neighborCopy, neighbors)

	return NoveltyResult{
		Score:             total / float64(k),
		NeighborDistances: neighborCopy,
	}
}
