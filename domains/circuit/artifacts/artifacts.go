// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package artifacts

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/gmlewis/gep/v2/domains/circuit"
)

// JSON emits canonical indented JSON for one [circuit.CircuitProgram].
func JSON(program circuit.CircuitProgram) ([]byte, error) {
	if err := program.Validate(); err != nil {
		return nil, fmt.Errorf("validate circuit program: %w", err)
	}
	data, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal circuit program JSON: %w", err)
	}
	return append(data, '\n'), nil
}

// SPICE emits a deterministic SPICE-style netlist for one
// [circuit.CircuitProgram].
func SPICE(program circuit.CircuitProgram) (string, error) {
	if err := program.Validate(); err != nil {
		return "", fmt.Errorf("validate circuit program: %w", err)
	}
	var b strings.Builder
	title := firstNonEmpty(program.Spec.Name, program.CandidateID, "circuit")
	fmt.Fprintf(&b, "* circuit candidate: %s\n", firstNonEmpty(program.CandidateID, "unnamed"))
	if program.Spec.Name != "" {
		fmt.Fprintf(&b, "* spec: %s\n", program.Spec.Name)
	}
	fmt.Fprintf(&b, ".TITLE %s\n", title)
	for _, c := range program.Graph.Components {
		fmt.Fprintf(&b, "X%s", c.NodeID)
		for _, p := range c.Inputs {
			fmt.Fprintf(&b, " %s.%s", p.Node, p.Name)
		}
		for _, p := range c.Outputs {
			fmt.Fprintf(&b, " %s.%s", p.Node, p.Name)
		}
		fmt.Fprintf(&b, " %s name=%s", c.Type, marshalScalar(c.Name))
		for _, key := range sortedKeys(c.Params) {
			fmt.Fprintf(&b, " %s=%s", key, marshalScalar(c.Params[key]))
		}
		b.WriteByte('\n')
	}
	b.WriteString(".END\n")
	return b.String(), nil
}

// Verilog emits deterministic structural-Verilog-style text for one
// [circuit.CircuitProgram].
func Verilog(program circuit.CircuitProgram) (string, error) {
	if err := program.Validate(); err != nil {
		return "", fmt.Errorf("validate circuit program: %w", err)
	}
	var b strings.Builder
	moduleName := sanitizeIdentifier(firstNonEmpty(program.Spec.Name, program.CandidateID, "circuit_design"))
	fmt.Fprintf(&b, "module %s;\n", moduleName)
	if program.CandidateID != "" {
		fmt.Fprintf(&b, "  // candidate_id: %s\n", program.CandidateID)
	}
	for _, c := range program.Graph.Components {
		fmt.Fprintf(&b, "  %s", sanitizeIdentifier(c.Type))
		if paramText := verilogParamList(c.Params); paramText != "" {
			fmt.Fprintf(&b, " #(%s)", paramText)
		}
		fmt.Fprintf(&b, " %s (%s);\n", sanitizeIdentifier(string(c.NodeID)), verilogPortList(c))
	}
	b.WriteString("endmodule\n")
	return b.String(), nil
}

func verilogParamList(params map[string]any) string {
	if len(params) == 0 {
		return ""
	}
	parts := make([]string, 0, len(params))
	for _, key := range sortedKeys(params) {
		parts = append(parts, fmt.Sprintf(".%s(%s)", sanitizeIdentifier(key), marshalScalar(params[key])))
	}
	return strings.Join(parts, ", ")
}

func verilogPortList(c circuit.Component) string {
	ports := make([]string, 0, len(c.Inputs)+len(c.Outputs))
	for i, p := range c.Inputs {
		ports = append(ports, fmt.Sprintf(".%s(%s)", instancePortName("in", i, len(c.Inputs)), signalName(p)))
	}
	for i, p := range c.Outputs {
		ports = append(ports, fmt.Sprintf(".%s(%s)", instancePortName("out", i, len(c.Outputs)), signalName(p)))
	}
	return strings.Join(ports, ", ")
}

func instancePortName(prefix string, index int, count int) string {
	if count == 1 {
		return prefix
	}
	return fmt.Sprintf("%s%d", prefix, index)
}

func signalName(p circuit.Port) string {
	return sanitizeIdentifier(string(p.Node) + "_" + p.Name)
}

func sortedKeys(m map[string]any) []string {
	if len(m) == 0 {
		return nil
	}
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func marshalScalar(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%q", fmt.Sprint(v))
	}
	return string(data)
}

func sanitizeIdentifier(s string) string {
	if s == "" {
		return "unnamed"
	}
	var b strings.Builder
	for i, r := range s {
		if unicode.IsLetter(r) || r == '_' || (i > 0 && unicode.IsDigit(r)) {
			b.WriteRune(r)
			continue
		}
		if i == 0 && unicode.IsDigit(r) {
			b.WriteString("n_")
			b.WriteRune(r)
			continue
		}
		b.WriteByte('_')
	}
	return b.String()
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
