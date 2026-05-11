// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package scenarios

import (
	"bytes"
	"reflect"
	"testing"

	designscenarios "github.com/gmlewis/gep/v2/design/scenarios"
)

func TestLoadFixtureSet(t *testing.T) {
	set, err := LoadFixtureSet()
	if err != nil {
		t.Fatalf("LoadFixtureSet() error = %v", err)
	}
	if set.Name != "voxel-smoke" {
		t.Fatalf("Name = %q, want %q", set.Name, "voxel-smoke")
	}
	if got, want := len(set.Scenarios), 4; got != want {
		t.Fatalf("scenario count = %d, want %d", got, want)
	}
}

func TestFixtureSetLoadsViaDesignScenarios(t *testing.T) {
	set, err := designscenarios.LoadScenarioSet(bytes.NewReader(fixtureSetJSON))
	if err != nil {
		t.Fatalf("LoadScenarioSet() error = %v", err)
	}

	reg := &designscenarios.ScenarioRegistry{Sets: []designscenarios.ScenarioSet{*set}}
	if err := reg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if got, want := len(reg.BySplit(designscenarios.Train)), 2; got != want {
		t.Fatalf("train count = %d, want %d", got, want)
	}
	if got, want := len(reg.BySplit(designscenarios.Validation)), 1; got != want {
		t.Fatalf("validation count = %d, want %d", got, want)
	}
	if got, want := len(reg.BySplit(designscenarios.Test)), 1; got != want {
		t.Fatalf("test count = %d, want %d", got, want)
	}
}

func TestLoadFixtureSetIsDeterministic(t *testing.T) {
	set1, err := LoadFixtureSet()
	if err != nil {
		t.Fatalf("LoadFixtureSet #1: %v", err)
	}
	set2, err := LoadFixtureSet()
	if err != nil {
		t.Fatalf("LoadFixtureSet #2: %v", err)
	}
	if !reflect.DeepEqual(set1, set2) {
		t.Fatalf("LoadFixtureSet() not deterministic:\nset1=%+v\nset2=%+v", set1, set2)
	}
}
