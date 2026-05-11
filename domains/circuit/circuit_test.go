// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package circuit

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func mustValidGraph() CircuitGraph {
	return CircuitGraph{
		Components: []Component{
			{
				NodeID: "vin",
				Name:   "vin",
				Type:   "source.dc",
				Outputs: []Port{
					{Node: "vin", Name: "out"},
				},
			},
			{
				NodeID: "r1",
				Name:   "r1",
				Type:   "resistor",
				Inputs: []Port{
					{Node: "vin", Name: "out"},
				},
				Outputs: []Port{
					{Node: "r1", Name: "out"},
				},
				Params: map[string]any{"ohms": 1000.0},
			},
			{
				NodeID: "gnd",
				Name:   "gnd",
				Type:   "ground",
				Inputs: []Port{
					{Node: "r1", Name: "out"},
				},
			},
		},
	}
}

func TestCircuitGraphValidateValid(t *testing.T) {
	g := mustValidGraph()
	if err := g.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
}

func TestCircuitGraphValidateDuplicateNodeID(t *testing.T) {
	g := mustValidGraph()
	g.Components[2].NodeID = "r1"

	err := g.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want duplicate node error")
	}
	if !strings.Contains(err.Error(), "duplicate node_id") {
		t.Fatalf("error %q does not mention duplicate node_id", err)
	}
}

func TestCircuitGraphValidateMissingComponentName(t *testing.T) {
	g := mustValidGraph()
	g.Components[1].Name = ""

	err := g.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want missing component name error")
	}
	if !strings.Contains(err.Error(), "missing component name") {
		t.Fatalf("error %q does not mention missing component name", err)
	}
}

func TestCircuitGraphValidateMissingComponentType(t *testing.T) {
	g := mustValidGraph()
	g.Components[1].Type = ""

	err := g.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want missing component type error")
	}
	if !strings.Contains(err.Error(), "missing component type") {
		t.Fatalf("error %q does not mention missing component type", err)
	}
}

func TestCircuitGraphValidateIllegalPortReference(t *testing.T) {
	g := mustValidGraph()
	g.Components[1].Inputs[0].Node = "missing-node"

	err := g.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want illegal port reference error")
	}
	if !strings.Contains(err.Error(), "references unknown node") {
		t.Fatalf("error %q does not mention unknown node", err)
	}
}

func TestCircuitGraphValidateIsDeterministic(t *testing.T) {
	g := mustValidGraph()
	g.Components[1].Inputs[0].Node = "missing-node"

	err1 := g.Validate()
	err2 := g.Validate()
	if (err1 == nil) != (err2 == nil) {
		t.Fatal("Validate() returned different nil-ness across calls")
	}
	if err1 != nil && err1.Error() != err2.Error() {
		t.Fatalf("Validate() not deterministic:\ncall1=%v\ncall2=%v", err1, err2)
	}
}

func TestCircuitProgramValidateDelegatesToGraph(t *testing.T) {
	p := CircuitProgram{Graph: mustValidGraph()}
	if err := p.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}

	p.Graph.Components[0].Type = ""
	err := p.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want graph validation error")
	}
}

func TestCircuitProgramJSONRoundTrip(t *testing.T) {
	orig := CircuitProgram{
		CandidateID: "cand-circuit-01",
		Graph:       mustValidGraph(),
		Spec: CircuitSpec{
			Name:     "smoke-rc",
			Domain:   "circuit",
			Revision: "v1",
			Metadata: map[string]any{"source": "fixture", "train": true},
		},
		Constraints: []CircuitConstraint{
			{Kind: "max_component_count", Enabled: true, Params: map[string]any{"max": 16.0}},
			{Kind: "allowed_types", Enabled: true, Params: map[string]any{"types": []any{"source.dc", "resistor", "ground"}}},
		},
	}

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got CircuitProgram
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if !reflect.DeepEqual(orig, got) {
		t.Fatalf("round-trip mismatch:\norig=%+v\n got=%+v", orig, got)
	}

	if err := got.Validate(); err != nil {
		t.Fatalf("Validate() after JSON round-trip = %v, want nil", err)
	}
}
