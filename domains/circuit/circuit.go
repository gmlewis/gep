// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package circuit

import "fmt"

// NodeID is the stable identifier for one graph node.
type NodeID string

// Port identifies one named port on a graph node.
type Port struct {
	// Node is the node identifier that owns this port.
	Node NodeID `json:"node"`
	// Name is the component-local port name.
	Name string `json:"name"`
}

// Component is one typed graph node in a circuit candidate.
type Component struct {
	// NodeID is this component's unique node identifier within a graph.
	NodeID NodeID `json:"node_id"`
	// Name is the human-readable component name.
	Name string `json:"name"`
	// Type is the stable component type identifier.
	Type string `json:"type"`
	// Inputs lists component input port references.
	Inputs []Port `json:"inputs,omitempty"`
	// Outputs lists component output port references.
	Outputs []Port `json:"outputs,omitempty"`
	// Params stores optional component-specific parameters.
	Params map[string]any `json:"params,omitempty"`
}

// CircuitGraph is one candidate circuit topology.
type CircuitGraph struct {
	// Components is the ordered list of graph components.
	Components []Component `json:"components"`
}

// CircuitSpec describes high-level metadata for a circuit candidate/program.
type CircuitSpec struct {
	// Name is the stable spec identifier.
	Name string `json:"name"`
	// Domain names the target domain (for PB-03 this is expected to be "circuit").
	Domain string `json:"domain"`
	// Revision is an optional schema/version marker.
	Revision string `json:"revision,omitempty"`
	// Metadata stores optional JSON-serializable metadata fields.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// CircuitConstraint is one declarative circuit-structure constraint.
type CircuitConstraint struct {
	// Kind is the stable constraint type identifier.
	Kind string `json:"kind"`
	// Enabled controls whether the constraint is active.
	Enabled bool `json:"enabled"`
	// Params stores optional JSON-serializable constraint parameters.
	Params map[string]any `json:"params,omitempty"`
}

// CircuitProgram is a serializable circuit-domain candidate package.
type CircuitProgram struct {
	// CandidateID is the stable identifier of the candidate.
	CandidateID string `json:"candidate_id"`
	// Graph is the candidate topology.
	Graph CircuitGraph `json:"graph"`
	// Spec is optional descriptive metadata for the candidate topology.
	Spec CircuitSpec `json:"spec"`
	// Constraints lists optional declarative constraints associated with the
	// candidate.
	Constraints []CircuitConstraint `json:"constraints,omitempty"`
}

// Validate returns an error when the graph is structurally invalid.
//
// Validation checks are deterministic and return the first encountered error in
// component declaration order:
//   - duplicate component node IDs,
//   - missing component names/types,
//   - illegal port references to unknown nodes.
func (g CircuitGraph) Validate() error {
	seen := make(map[NodeID]struct{}, len(g.Components))
	for i, c := range g.Components {
		if c.NodeID == "" {
			return fmt.Errorf("component[%d]: missing node_id", i)
		}
		if _, ok := seen[c.NodeID]; ok {
			return fmt.Errorf("component[%d]: duplicate node_id %q", i, c.NodeID)
		}
		seen[c.NodeID] = struct{}{}

		if c.Name == "" {
			return fmt.Errorf("component[%d] node %q: missing component name", i, c.NodeID)
		}
		if c.Type == "" {
			return fmt.Errorf("component[%d] node %q: missing component type", i, c.NodeID)
		}
	}

	for i, c := range g.Components {
		if err := validatePorts("input", i, c.NodeID, c.Inputs, seen); err != nil {
			return err
		}
		if err := validatePorts("output", i, c.NodeID, c.Outputs, seen); err != nil {
			return err
		}
	}

	return nil
}

func validatePorts(kind string, componentIndex int, componentNode NodeID, ports []Port, nodes map[NodeID]struct{}) error {
	for portIndex, p := range ports {
		if p.Node == "" {
			return fmt.Errorf("component[%d] node %q: %s port[%d] missing node reference", componentIndex, componentNode, kind, portIndex)
		}
		if _, ok := nodes[p.Node]; !ok {
			return fmt.Errorf("component[%d] node %q: %s port[%d] references unknown node %q", componentIndex, componentNode, kind, portIndex, p.Node)
		}
	}
	return nil
}

// Validate returns an error when p.Graph fails structural validation.
func (p CircuitProgram) Validate() error {
	return p.Graph.Validate()
}
