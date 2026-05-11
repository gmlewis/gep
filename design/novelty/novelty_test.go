// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package novelty

import (
	"math"
	"reflect"
	"testing"
)

// --- helpers ---

func vec(vals ...float64) BehaviorVector { return BehaviorVector(vals) }

func entry(label string, vals ...float64) ArchiveEntry {
	return ArchiveEntry{Behavior: vec(vals...), Label: label}
}

// --- SquaredEuclidean tests ---

func TestSquaredEuclideanSameVector(t *testing.T) {
	if got := SquaredEuclidean(vec(1, 2, 3), vec(1, 2, 3)); got != 0 {
		t.Fatalf("SquaredEuclidean same vector = %v, want 0", got)
	}
}

func TestSquaredEuclideanKnownValue(t *testing.T) {
	// (3-0)^2 + (4-0)^2 = 9 + 16 = 25
	if got := SquaredEuclidean(vec(3, 4), vec(0, 0)); got != 25 {
		t.Fatalf("SquaredEuclidean = %v, want 25", got)
	}
}

func TestSquaredEuclideanSymmetric(t *testing.T) {
	a, b := vec(1, 2, 3), vec(4, 5, 6)
	if SquaredEuclidean(a, b) != SquaredEuclidean(b, a) {
		t.Fatal("SquaredEuclidean must be symmetric")
	}
}

func TestSquaredEuclideanDifferentLengths(t *testing.T) {
	// a=[1,2], b=[1,2,3]: missing dim treated as 0; extra dim contributes 3^2=9
	got := SquaredEuclidean(vec(1, 2), vec(1, 2, 3))
	if got != 9 {
		t.Fatalf("SquaredEuclidean mismatched lengths = %v, want 9", got)
	}
}

// --- EuclideanDistance tests ---

func TestEuclideanDistance(t *testing.T) {
	// 3-4-5 right triangle
	got := EuclideanDistance(vec(3, 4), vec(0, 0))
	if math.Abs(got-5) > 1e-12 {
		t.Fatalf("EuclideanDistance = %v, want 5", got)
	}
}

// --- NewArchive / Len tests ---

func TestNewArchiveDefaults(t *testing.T) {
	a := NewArchive(ArchiveConfig{})
	if a.Len() != 0 {
		t.Fatalf("new archive Len = %d, want 0", a.Len())
	}
}

func TestNewArchiveDefaultsK(t *testing.T) {
	// K=0 should default to 1.
	a := NewArchive(ArchiveConfig{K: 0})
	if a.cfg.K != 1 {
		t.Fatalf("default K = %d, want 1", a.cfg.K)
	}
}

// --- Empty archive behavior ---

func TestScoreEmptyArchive(t *testing.T) {
	a := NewArchive(ArchiveConfig{K: 3})
	result := a.Score(vec(1, 2, 3))

	if result.Score != 0 {
		t.Fatalf("empty archive Score = %v, want 0", result.Score)
	}
	if len(result.NeighborDistances) != 0 {
		t.Fatalf("empty archive NeighborDistances = %v, want []", result.NeighborDistances)
	}
}

// --- Archive insertion ---

func TestAddIncreasesLen(t *testing.T) {
	a := NewArchive(ArchiveConfig{K: 1})
	a.Add(entry("c1", 1, 0))
	if a.Len() != 1 {
		t.Fatalf("Len after 1 Add = %d, want 1", a.Len())
	}
	a.Add(entry("c2", 2, 0))
	if a.Len() != 2 {
		t.Fatalf("Len after 2 Adds = %d, want 2", a.Len())
	}
}

func TestAddRespectMaxSize(t *testing.T) {
	a := NewArchive(ArchiveConfig{K: 1, MaxSize: 2})
	a.Add(entry("c1", 1, 0))
	a.Add(entry("c2", 2, 0))
	a.Add(entry("c3", 3, 0)) // should be silently dropped
	if a.Len() != 2 {
		t.Fatalf("Len after 3 Adds with MaxSize=2 = %d, want 2", a.Len())
	}
}

func TestAddMaxSizeZeroMeansUnbounded(t *testing.T) {
	a := NewArchive(ArchiveConfig{K: 1, MaxSize: 0})
	for i := range 100 {
		a.Add(entry("", float64(i)))
	}
	if a.Len() != 100 {
		t.Fatalf("Len = %d, want 100 (unbounded)", a.Len())
	}
}

func TestEntriesReturnsCopy(t *testing.T) {
	a := NewArchive(ArchiveConfig{K: 1})
	a.Add(entry("c1", 1))
	snapshot := a.Entries()
	// mutating the snapshot must not affect the archive
	snapshot[0].Label = "mutated"
	if a.Entries()[0].Label != "c1" {
		t.Fatal("Entries() should return a copy, not a reference")
	}
}

// --- k-NN novelty score calculation ---

func TestScoreOneNeighbor(t *testing.T) {
	a := NewArchive(ArchiveConfig{K: 1})
	a.Add(entry("origin", 0, 0))

	// query at (3,4): squared Euclidean distance = 9+16 = 25
	result := a.Score(vec(3, 4))
	if result.Score != 25 {
		t.Fatalf("Score = %v, want 25", result.Score)
	}
	if len(result.NeighborDistances) != 1 || result.NeighborDistances[0] != 25 {
		t.Fatalf("NeighborDistances = %v, want [25]", result.NeighborDistances)
	}
}

func TestScoreKNearestMean(t *testing.T) {
	// archive: points at distance 1, 2, 3, 4 from query
	a := NewArchive(ArchiveConfig{K: 3})
	a.Add(entry("d1", 1)) // distance to origin = 1
	a.Add(entry("d2", 2)) // distance = 4
	a.Add(entry("d3", 3)) // distance = 9
	a.Add(entry("d4", 4)) // distance = 16

	// k=3 nearest: distances 1, 4, 9 -> mean = (1+4+9)/3 = 14/3
	result := a.Score(vec(0))
	want := (1.0 + 4.0 + 9.0) / 3.0
	if math.Abs(result.Score-want) > 1e-12 {
		t.Fatalf("Score = %v, want %v", result.Score, want)
	}
	if len(result.NeighborDistances) != 3 {
		t.Fatalf("NeighborDistances len = %d, want 3", len(result.NeighborDistances))
	}
}

func TestScoreFewerThanKEntries(t *testing.T) {
	// K=5 but only 2 entries: mean over both
	a := NewArchive(ArchiveConfig{K: 5})
	a.Add(entry("a", 1))
	a.Add(entry("b", 3))

	// distances from origin: 1, 9 -> mean = 5
	result := a.Score(vec(0))
	if result.Score != 5 {
		t.Fatalf("Score = %v, want 5", result.Score)
	}
	if len(result.NeighborDistances) != 2 {
		t.Fatalf("NeighborDistances len = %d, want 2", len(result.NeighborDistances))
	}
}

func TestScoreNeighborDistancesAscending(t *testing.T) {
	a := NewArchive(ArchiveConfig{K: 4})
	a.Add(entry("far", 10))
	a.Add(entry("near", 1))
	a.Add(entry("mid", 5))
	a.Add(entry("zero", 0))

	result := a.Score(vec(0))
	for i := 1; i < len(result.NeighborDistances); i++ {
		if result.NeighborDistances[i] < result.NeighborDistances[i-1] {
			t.Fatalf("NeighborDistances not ascending: %v", result.NeighborDistances)
		}
	}
}

// --- Stable / deterministic novelty scores ---

func TestScoreDeterministicForFixedArchive(t *testing.T) {
	a := NewArchive(ArchiveConfig{K: 2})
	a.Add(entry("p1", 1, 0))
	a.Add(entry("p2", 0, 1))
	a.Add(entry("p3", 2, 2))

	query := vec(1, 1)
	got1 := a.Score(query)
	got2 := a.Score(query)

	if !reflect.DeepEqual(got1, got2) {
		t.Fatalf("Score not deterministic:\n got1=%#v\n got2=%#v", got1, got2)
	}
}

func TestScoreDoesNotMutateArchive(t *testing.T) {
	a := NewArchive(ArchiveConfig{K: 2})
	a.Add(entry("p1", 1))
	a.Add(entry("p2", 3))
	lenBefore := a.Len()

	_ = a.Score(vec(0))

	if a.Len() != lenBefore {
		t.Fatalf("Score mutated archive: Len was %d, now %d", lenBefore, a.Len())
	}
}

func TestScoreSameResultAfterRepeatedScoreCalls(t *testing.T) {
	a := NewArchive(ArchiveConfig{K: 3})
	a.Add(entry("a", 1, 0, 0))
	a.Add(entry("b", 0, 1, 0))
	a.Add(entry("c", 0, 0, 1))

	query := vec(0.5, 0.5, 0.5)
	first := a.Score(query)
	for range 10 {
		if !reflect.DeepEqual(a.Score(query), first) {
			t.Fatal("Score returned different result on repeated call")
		}
	}
}

// --- Custom distance function ---

func TestCustomDistanceFunc(t *testing.T) {
	// Manhattan distance
	manhattan := func(a, b BehaviorVector) float64 {
		var sum float64
		n := len(a)
		if len(b) > n {
			n = len(b)
		}
		for i := range n {
			var ai, bi float64
			if i < len(a) {
				ai = a[i]
			}
			if i < len(b) {
				bi = b[i]
			}
			d := ai - bi
			if d < 0 {
				d = -d
			}
			sum += d
		}
		return sum
	}

	a := NewArchive(ArchiveConfig{K: 1, Distance: manhattan})
	a.Add(entry("p", 3, 4))

	// Manhattan distance from (0,0) to (3,4) = 7
	result := a.Score(vec(0, 0))
	if result.Score != 7 {
		t.Fatalf("custom distance Score = %v, want 7", result.Score)
	}
}

// --- ArchiveGrowth correctness ---

func TestArchiveGrowthIsCorrect(t *testing.T) {
	a := NewArchive(ArchiveConfig{K: 1})

	// Empty archive: score = 0.
	if r := a.Score(vec(0)); r.Score != 0 {
		t.Fatalf("empty archive score = %v, want 0", r.Score)
	}

	// After one insertion the score should reflect that distance.
	a.Add(entry("c1", 5))
	r1 := a.Score(vec(0))
	if r1.Score != 25 { // squared euclidean: 5^2
		t.Fatalf("score after 1 insertion = %v, want 25", r1.Score)
	}

	// Adding a closer entry reduces the k=1 nearest-neighbor score.
	a.Add(entry("c2", 2))
	r2 := a.Score(vec(0))
	if r2.Score != 4 { // 2^2=4
		t.Fatalf("score after 2 insertions = %v, want 4", r2.Score)
	}
}
