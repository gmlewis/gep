// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package scenarios

import (
	"reflect"
	"strings"
	"testing"
)

// --- helpers ---

// makeRegistry builds a ScenarioRegistry from a slice of ScenarioSets.
func makeRegistry(sets ...ScenarioSet) *ScenarioRegistry {
	return &ScenarioRegistry{Sets: sets}
}

// makeSet builds a ScenarioSet with the given name and scenarios.
func makeSet(name string, scenarios ...Scenario) ScenarioSet {
	return ScenarioSet{Name: name, Scenarios: scenarios}
}

// sc is a shorthand for building a Scenario with only ID and split.
func sc(id ScenarioID, split ScenarioSplit) Scenario {
	return Scenario{ID: id, Split: split}
}

// --- LoadScenarioSet tests ---

func TestLoadScenarioSetNilReader(t *testing.T) {
	_, err := LoadScenarioSet(nil)
	if err == nil {
		t.Fatal("expected error for nil reader, got nil")
	}
}

func TestLoadScenarioSetMalformedJSON(t *testing.T) {
	_, err := LoadScenarioSet(strings.NewReader("{not valid json"))
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

func TestLoadScenarioSetTwoObjects(t *testing.T) {
	_, err := LoadScenarioSet(strings.NewReader(`{"name":"a","scenarios":[]}{}` + "\n"))
	if err == nil {
		t.Fatal("expected error for two-object stream, got nil")
	}
}

func TestLoadScenarioSetUnknownField(t *testing.T) {
	_, err := LoadScenarioSet(strings.NewReader(`{"name":"a","unknown_field":true,"scenarios":[]}`))
	if err == nil {
		t.Fatal("expected error for unknown field, got nil")
	}
}

func TestLoadScenarioSetEmpty(t *testing.T) {
	set, err := LoadScenarioSet(strings.NewReader(`{"name":"empty","scenarios":[]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if set.Name != "empty" {
		t.Fatalf("Name = %q, want %q", set.Name, "empty")
	}
	if len(set.Scenarios) != 0 {
		t.Fatalf("Scenarios count = %d, want 0", len(set.Scenarios))
	}
}

// --- LoadScenarioSetFile tests ---

func TestLoadScenarioSetFileNotFound(t *testing.T) {
	_, err := LoadScenarioSetFile("testdata/does_not_exist.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadScenarioSetFileSmokeFixture(t *testing.T) {
	set, err := LoadScenarioSetFile("testdata/set_smoke.json")
	if err != nil {
		t.Fatalf("LoadScenarioSetFile smoke: %v", err)
	}
	if set.Name != "smoke" {
		t.Fatalf("Name = %q, want %q", set.Name, "smoke")
	}
	if got, want := len(set.Scenarios), 5; got != want {
		t.Fatalf("Scenarios count = %d, want %d", got, want)
	}
}

func TestLoadScenarioSetFileParamsFixture(t *testing.T) {
	set, err := LoadScenarioSetFile("testdata/set_params.json")
	if err != nil {
		t.Fatalf("LoadScenarioSetFile params: %v", err)
	}
	if set.Name != "params" {
		t.Fatalf("Name = %q, want %q", set.Name, "params")
	}
	if got, want := len(set.Scenarios), 3; got != want {
		t.Fatalf("Scenarios count = %d, want %d", got, want)
	}
	// Verify params round-tripped for first scenario.
	p := set.Scenarios[0].Params
	if p == nil {
		t.Fatal("Params is nil for first scenario, want non-nil")
	}
	if _, ok := p["load"]; !ok {
		t.Fatal("Params missing 'load' key")
	}
}

// --- Fixture split and ordering ---

func TestSmokeFixtureSplitAssignments(t *testing.T) {
	set, err := LoadScenarioSetFile("testdata/set_smoke.json")
	if err != nil {
		t.Fatalf("load smoke: %v", err)
	}
	reg := makeRegistry(*set)

	trains := reg.BySplit(Train)
	vals := reg.BySplit(Validation)
	tests := reg.BySplit(Test)

	if got, want := len(trains), 2; got != want {
		t.Errorf("Train count = %d, want %d", got, want)
	}
	if got, want := len(vals), 1; got != want {
		t.Errorf("Validation count = %d, want %d", got, want)
	}
	if got, want := len(tests), 2; got != want {
		t.Errorf("Test count = %d, want %d", got, want)
	}
}

func TestSmokeFixtureOrdering(t *testing.T) {
	set, err := LoadScenarioSetFile("testdata/set_smoke.json")
	if err != nil {
		t.Fatalf("load smoke: %v", err)
	}
	reg := makeRegistry(*set)

	trains := reg.BySplit(Train)
	if len(trains) < 2 {
		t.Fatalf("expected at least 2 train scenarios, got %d", len(trains))
	}
	// Insertion order must be preserved.
	if trains[0].ID != "s-train-01" {
		t.Errorf("trains[0].ID = %q, want %q", trains[0].ID, "s-train-01")
	}
	if trains[1].ID != "s-train-02" {
		t.Errorf("trains[1].ID = %q, want %q", trains[1].ID, "s-train-02")
	}
}

func TestSmokeFixtureValidatesClean(t *testing.T) {
	set, err := LoadScenarioSetFile("testdata/set_smoke.json")
	if err != nil {
		t.Fatalf("load smoke: %v", err)
	}
	reg := makeRegistry(*set)
	if err := reg.Validate(); err != nil {
		t.Fatalf("Validate smoke fixture = %v, want nil", err)
	}
}

// --- ScenarioRegistry.Validate tests ---

func TestValidateEmptyRegistry(t *testing.T) {
	reg := makeRegistry()
	if err := reg.Validate(); err != nil {
		t.Fatalf("Validate empty registry = %v, want nil", err)
	}
}

func TestValidateValidRegistry(t *testing.T) {
	reg := makeRegistry(
		makeSet("s1",
			sc("a", Train),
			sc("b", Validation),
		),
		makeSet("s2",
			sc("c", Test),
		),
	)
	if err := reg.Validate(); err != nil {
		t.Fatalf("Validate valid registry = %v, want nil", err)
	}
}

func TestValidateCrossSpitDuplicate(t *testing.T) {
	reg := makeRegistry(
		makeSet("s1",
			sc("dup", Train),
		),
		makeSet("s2",
			sc("dup", Validation),
		),
	)
	err := reg.Validate()
	if err == nil {
		t.Fatal("expected error for cross-split duplicate, got nil")
	}
	if !strings.Contains(err.Error(), "dup") {
		t.Errorf("error %q does not mention scenario ID", err.Error())
	}
}

func TestValidateSameSplitDuplicate(t *testing.T) {
	reg := makeRegistry(
		makeSet("s1",
			sc("same", Train),
			sc("same", Train),
		),
	)
	err := reg.Validate()
	if err == nil {
		t.Fatal("expected error for same-split duplicate, got nil")
	}
	if !strings.Contains(err.Error(), "same") {
		t.Errorf("error %q does not mention scenario ID", err.Error())
	}
}

func TestValidateDuplicateAcrossSets(t *testing.T) {
	reg := makeRegistry(
		makeSet("s1", sc("x", Train)),
		makeSet("s2", sc("x", Train)),
	)
	err := reg.Validate()
	if err == nil {
		t.Fatal("expected error for duplicate ID across sets in same split, got nil")
	}
}

func TestValidateIsDeterministic(t *testing.T) {
	reg := makeRegistry(
		makeSet("s1",
			sc("a", Train),
			sc("b", Validation),
			sc("c", Test),
		),
	)
	err1 := reg.Validate()
	err2 := reg.Validate()
	if (err1 == nil) != (err2 == nil) {
		t.Fatal("Validate returned different nil-ness on repeated calls")
	}
	if err1 != nil && err1.Error() != err2.Error() {
		t.Fatalf("Validate not deterministic:\n call1=%v\n call2=%v", err1, err2)
	}
}

// --- ScenarioRegistry.ByID tests ---

func TestByIDFound(t *testing.T) {
	reg := makeRegistry(
		makeSet("s1",
			sc("first", Train),
			sc("second", Validation),
		),
	)
	got, ok := reg.ByID("second")
	if !ok {
		t.Fatal("ByID(\"second\") = false, want true")
	}
	if got.Split != Validation {
		t.Errorf("got Split = %q, want %q", got.Split, Validation)
	}
}

func TestByIDMissing(t *testing.T) {
	reg := makeRegistry(makeSet("s1", sc("a", Train)))
	_, ok := reg.ByID("missing")
	if ok {
		t.Fatal("ByID(\"missing\") = true, want false")
	}
}

func TestByIDFirstSetWins(t *testing.T) {
	// When the same ID appears in two different sets (before validation),
	// ByID should return the first one encountered.
	reg := makeRegistry(
		makeSet("s1", Scenario{ID: "x", Split: Train, Tags: []string{"first"}}),
		makeSet("s2", Scenario{ID: "x", Split: Test, Tags: []string{"second"}}),
	)
	got, ok := reg.ByID("x")
	if !ok {
		t.Fatal("ByID = false, want true")
	}
	if len(got.Tags) == 0 || got.Tags[0] != "first" {
		t.Errorf("ByID returned wrong entry: %+v", got)
	}
}

// --- ScenarioRegistry.BySplit tests ---

func TestBySplitEmpty(t *testing.T) {
	reg := makeRegistry(makeSet("s1", sc("a", Train)))
	result := reg.BySplit(Validation)
	if len(result) != 0 {
		t.Fatalf("BySplit(Validation) = %d items, want 0", len(result))
	}
}

func TestBySplitReturnsCopy(t *testing.T) {
	reg := makeRegistry(makeSet("s1", sc("a", Train)))
	r1 := reg.BySplit(Train)
	if len(r1) == 0 {
		t.Fatal("expected at least one train scenario")
	}
	// Mutating the returned slice must not affect future calls.
	r1[0].ID = "mutated"
	r2 := reg.BySplit(Train)
	if r2[0].ID == "mutated" {
		t.Fatal("BySplit returned a reference instead of a copy")
	}
}

func TestBySplitPreservesDeclarationOrder(t *testing.T) {
	reg := makeRegistry(
		makeSet("s1",
			sc("t1", Train),
			sc("v1", Validation),
			sc("t2", Train),
		),
		makeSet("s2",
			sc("t3", Train),
		),
	)
	trains := reg.BySplit(Train)
	wantIDs := []ScenarioID{"t1", "t2", "t3"}
	got := make([]ScenarioID, len(trains))
	for i, s := range trains {
		got[i] = s.ID
	}
	if !reflect.DeepEqual(got, wantIDs) {
		t.Errorf("BySplit order = %v, want %v", got, wantIDs)
	}
}

// --- Multi-set registry with both fixtures ---

func TestMultiSetRegistry(t *testing.T) {
	smoke, err := LoadScenarioSetFile("testdata/set_smoke.json")
	if err != nil {
		t.Fatalf("load smoke: %v", err)
	}
	params, err := LoadScenarioSetFile("testdata/set_params.json")
	if err != nil {
		t.Fatalf("load params: %v", err)
	}

	reg := makeRegistry(*smoke, *params)

	// All IDs are unique across fixtures, so Validate must pass.
	if err := reg.Validate(); err != nil {
		t.Fatalf("Validate multi-set registry = %v, want nil", err)
	}

	// Train scenarios come from both sets.
	trains := reg.BySplit(Train)
	// smoke: 2 train, params: 2 train
	if got, want := len(trains), 4; got != want {
		t.Fatalf("Train count = %d, want %d", got, want)
	}
}

// --- Deterministic loading of fixtures ---

func TestFixtureLoadingIsDeterministic(t *testing.T) {
	load := func() *ScenarioSet {
		t.Helper()
		set, err := LoadScenarioSetFile("testdata/set_smoke.json")
		if err != nil {
			t.Fatalf("load: %v", err)
		}
		return set
	}

	s1 := load()
	s2 := load()

	if !reflect.DeepEqual(s1, s2) {
		t.Fatal("fixture loads are not deterministic")
	}
}

func TestBySplitIsDeterministic(t *testing.T) {
	set, err := LoadScenarioSetFile("testdata/set_smoke.json")
	if err != nil {
		t.Fatalf("load smoke: %v", err)
	}
	reg := makeRegistry(*set)

	r1 := reg.BySplit(Train)
	r2 := reg.BySplit(Train)
	if !reflect.DeepEqual(r1, r2) {
		t.Fatal("BySplit not deterministic")
	}
}
