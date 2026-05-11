// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package artifacts

import (
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/domains/circuit"
)

func testProgram() circuit.CircuitProgram {
	return circuit.CircuitProgram{
		CandidateID: "cand-circuit-01",
		Graph: circuit.CircuitGraph{
			Components: []circuit.Component{
				{
					NodeID: "vin",
					Name:   "vin",
					Type:   "source.dc",
					Outputs: []circuit.Port{
						{Node: "vin", Name: "out"},
					},
					Params: map[string]any{"volts": 5.0},
				},
				{
					NodeID: "r1",
					Name:   "r1",
					Type:   "resistor",
					Inputs: []circuit.Port{
						{Node: "vin", Name: "out"},
					},
					Outputs: []circuit.Port{
						{Node: "r1", Name: "out"},
					},
					Params: map[string]any{
						"tolerance": "5%",
						"ohms":      1000.0,
					},
				},
				{
					NodeID: "gnd",
					Name:   "gnd",
					Type:   "ground",
					Inputs: []circuit.Port{
						{Node: "r1", Name: "out"},
					},
				},
			},
		},
		Spec: circuit.CircuitSpec{
			Name:     "smoke-rc",
			Domain:   "circuit",
			Revision: "v1",
		},
	}
}

func TestJSON(t *testing.T) {
	got, err := JSON(testProgram())
	if err != nil {
		t.Fatalf("JSON() error = %v", err)
	}
	const want = "{\n" +
		"  \"candidate_id\": \"cand-circuit-01\",\n" +
		"  \"graph\": {\n" +
		"    \"components\": [\n" +
		"      {\n" +
		"        \"node_id\": \"vin\",\n" +
		"        \"name\": \"vin\",\n" +
		"        \"type\": \"source.dc\",\n" +
		"        \"outputs\": [\n" +
		"          {\n" +
		"            \"node\": \"vin\",\n" +
		"            \"name\": \"out\"\n" +
		"          }\n" +
		"        ],\n" +
		"        \"params\": {\n" +
		"          \"volts\": 5\n" +
		"        }\n" +
		"      },\n" +
		"      {\n" +
		"        \"node_id\": \"r1\",\n" +
		"        \"name\": \"r1\",\n" +
		"        \"type\": \"resistor\",\n" +
		"        \"inputs\": [\n" +
		"          {\n" +
		"            \"node\": \"vin\",\n" +
		"            \"name\": \"out\"\n" +
		"          }\n" +
		"        ],\n" +
		"        \"outputs\": [\n" +
		"          {\n" +
		"            \"node\": \"r1\",\n" +
		"            \"name\": \"out\"\n" +
		"          }\n" +
		"        ],\n" +
		"        \"params\": {\n" +
		"          \"ohms\": 1000,\n" +
		"          \"tolerance\": \"5%\"\n" +
		"        }\n" +
		"      },\n" +
		"      {\n" +
		"        \"node_id\": \"gnd\",\n" +
		"        \"name\": \"gnd\",\n" +
		"        \"type\": \"ground\",\n" +
		"        \"inputs\": [\n" +
		"          {\n" +
		"            \"node\": \"r1\",\n" +
		"            \"name\": \"out\"\n" +
		"          }\n" +
		"        ]\n" +
		"      }\n" +
		"    ]\n" +
		"  },\n" +
		"  \"spec\": {\n" +
		"    \"name\": \"smoke-rc\",\n" +
		"    \"domain\": \"circuit\",\n" +
		"    \"revision\": \"v1\"\n" +
		"  }\n" +
		"}\n"
	if string(got) != want {
		t.Fatalf("JSON() mismatch:\n--- got ---\n%s--- want ---\n%s", got, want)
	}
}

func TestSPICE(t *testing.T) {
	got, err := SPICE(testProgram())
	if err != nil {
		t.Fatalf("SPICE() error = %v", err)
	}
	const want = "" +
		"* circuit candidate: cand-circuit-01\n" +
		"* spec: smoke-rc\n" +
		".TITLE smoke-rc\n" +
		"Xvin vin.out source.dc name=\"vin\" volts=5\n" +
		"Xr1 vin.out r1.out resistor name=\"r1\" ohms=1000 tolerance=\"5%\"\n" +
		"Xgnd r1.out ground name=\"gnd\"\n" +
		".END\n"
	if got != want {
		t.Fatalf("SPICE() mismatch:\n--- got ---\n%s--- want ---\n%s", got, want)
	}
}

func TestVerilog(t *testing.T) {
	got, err := Verilog(testProgram())
	if err != nil {
		t.Fatalf("Verilog() error = %v", err)
	}
	const want = "" +
		"module smoke_rc;\n" +
		"  // candidate_id: cand-circuit-01\n" +
		"  source_dc #(.volts(5)) vin (.out(vin_out));\n" +
		"  resistor #(.ohms(1000), .tolerance(\"5%\")) r1 (.in(vin_out), .out(r1_out));\n" +
		"  ground gnd (.in(r1_out));\n" +
		"endmodule\n"
	if got != want {
		t.Fatalf("Verilog() mismatch:\n--- got ---\n%s--- want ---\n%s", got, want)
	}
}

func TestEmittersAreDeterministic(t *testing.T) {
	program := testProgram()
	json1, err := JSON(program)
	if err != nil {
		t.Fatalf("JSON #1: %v", err)
	}
	json2, err := JSON(program)
	if err != nil {
		t.Fatalf("JSON #2: %v", err)
	}
	if string(json1) != string(json2) {
		t.Fatal("JSON output is not deterministic")
	}

	spice1, err := SPICE(program)
	if err != nil {
		t.Fatalf("SPICE #1: %v", err)
	}
	spice2, err := SPICE(program)
	if err != nil {
		t.Fatalf("SPICE #2: %v", err)
	}
	if spice1 != spice2 {
		t.Fatal("SPICE output is not deterministic")
	}

	verilog1, err := Verilog(program)
	if err != nil {
		t.Fatalf("Verilog #1: %v", err)
	}
	verilog2, err := Verilog(program)
	if err != nil {
		t.Fatalf("Verilog #2: %v", err)
	}
	if verilog1 != verilog2 {
		t.Fatal("Verilog output is not deterministic")
	}
}

func TestEmittersRejectInvalidProgram(t *testing.T) {
	program := testProgram()
	program.Graph.Components[1].Inputs[0].Node = "missing"

	for _, tc := range []struct {
		name string
		fn   func(circuit.CircuitProgram) error
	}{
		{
			name: "JSON",
			fn: func(p circuit.CircuitProgram) error {
				_, err := JSON(p)
				return err
			},
		},
		{
			name: "SPICE",
			fn: func(p circuit.CircuitProgram) error {
				_, err := SPICE(p)
				return err
			},
		},
		{
			name: "Verilog",
			fn: func(p circuit.CircuitProgram) error {
				_, err := Verilog(p)
				return err
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn(program)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), "references unknown node") {
				t.Fatalf("error %q does not mention invalid graph", err)
			}
		})
	}
}
