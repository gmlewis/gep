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
	const want = `{
  "candidate_id": "cand-circuit-01",
  "graph": {
    "components": [
      {
        "node_id": "vin",
        "name": "vin",
        "type": "source.dc",
        "outputs": [
          {
            "node": "vin",
            "name": "out"
          }
        ],
        "params": {
          "volts": 5
        }
      },
      {
        "node_id": "r1",
        "name": "r1",
        "type": "resistor",
        "inputs": [
          {
            "node": "vin",
            "name": "out"
          }
        ],
        "outputs": [
          {
            "node": "r1",
            "name": "out"
          }
        ],
        "params": {
          "ohms": 1000,
          "tolerance": "5%"
        }
      },
      {
        "node_id": "gnd",
        "name": "gnd",
        "type": "ground",
        "inputs": [
          {
            "node": "r1",
            "name": "out"
          }
        ]
      }
    ]
  },
  "spec": {
    "name": "smoke-rc",
    "domain": "circuit",
    "revision": "v1"
  }
}
`
	if string(got) != want {
		t.Fatalf("JSON() mismatch:\n--- got ---\n%s--- want ---\n%s", got, want)
	}
}

func TestSPICE(t *testing.T) {
	got, err := SPICE(testProgram())
	if err != nil {
		t.Fatalf("SPICE() error = %v", err)
	}
	const want = `* circuit candidate: cand-circuit-01
* spec: smoke-rc
.TITLE smoke-rc
Xvin vin.out source.dc name="vin" volts=5
Xr1 vin.out r1.out resistor name="r1" ohms=1000 tolerance="5%"
Xgnd r1.out ground name="gnd"
.END
`
	if got != want {
		t.Fatalf("SPICE() mismatch:\n--- got ---\n%s--- want ---\n%s", got, want)
	}
}

func TestVerilog(t *testing.T) {
	got, err := Verilog(testProgram())
	if err != nil {
		t.Fatalf("Verilog() error = %v", err)
	}
	const want = `module smoke_rc;
  // candidate_id: cand-circuit-01
  source_dc #(.volts(5)) vin (.out(vin_out));
  resistor #(.ohms(1000), .tolerance("5%")) r1 (.in(vin_out), .out(r1_out));
  ground gnd (.in(r1_out));
endmodule
`
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
