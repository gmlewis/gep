// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package artifacts

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gmlewis/gep/v2/domains/voxel"
)

// JSON emits canonical indented JSON for one [voxel.VoxelProgram].
func JSON(program voxel.VoxelProgram) ([]byte, error) {
	if err := program.Validate(); err != nil {
		return nil, fmt.Errorf("validate voxel program: %w", err)
	}
	data, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal voxel program JSON: %w", err)
	}
	return append(data, '\n'), nil
}

// OBJ emits a Wavefront OBJ mesh for one [voxel.VoxelProgram]. Each occupied
// voxel cell is rendered as an axis-aligned unit cube. Vertex coordinates are
// integers (one unit equals one voxel). The output is deterministic and
// requires no external geometry kernel.
func OBJ(program voxel.VoxelProgram) (string, error) {
	if err := program.Validate(); err != nil {
		return "", fmt.Errorf("validate voxel program: %w", err)
	}

	var b strings.Builder
	name := firstNonEmpty(program.Spec.Name, program.CandidateID, "voxel_design")
	fmt.Fprintf(&b, "# voxel candidate: %s\n", firstNonEmpty(program.CandidateID, "unnamed"))
	if program.Spec.Name != "" {
		fmt.Fprintf(&b, "# spec: %s\n", program.Spec.Name)
	}
	fmt.Fprintf(&b, "o %s\n", sanitizeIdentifier(name))

	// Emit one unit cube per occupied cell. Vertices for cell (x,y,z) run from
	// (x,y,z) to (x+1,y+1,z+1). Faces use 1-based OBJ vertex indices.
	vertexBase := 1
	for _, cell := range program.Design.Occupied {
		x, y, z := cell.Coord.X, cell.Coord.Y, cell.Coord.Z
		// 8 corners of the unit cube at (x,y,z).
		//   v0 = (x,   y,   z  )   v1 = (x+1, y,   z  )
		//   v2 = (x+1, y+1, z  )   v3 = (x,   y+1, z  )
		//   v4 = (x,   y,   z+1)   v5 = (x+1, y,   z+1)
		//   v6 = (x+1, y+1, z+1)   v7 = (x,   y+1, z+1)
		corners := [8][3]int{
			{x, y, z},
			{x + 1, y, z},
			{x + 1, y + 1, z},
			{x, y + 1, z},
			{x, y, z + 1},
			{x + 1, y, z + 1},
			{x + 1, y + 1, z + 1},
			{x, y + 1, z + 1},
		}
		for _, c := range corners {
			fmt.Fprintf(&b, "v %d %d %d\n", c[0], c[1], c[2])
		}
		// Emit a material comment when the cell has a material binding.
		if cell.Material != "" {
			fmt.Fprintf(&b, "# material: %s\n", cell.Material)
		}
		// 6 quad faces using 1-based vertex indices relative to vertexBase.
		b0 := vertexBase
		faces := [6][4]int{
			{b0 + 0, b0 + 3, b0 + 2, b0 + 1}, // -Z face
			{b0 + 4, b0 + 5, b0 + 6, b0 + 7}, // +Z face
			{b0 + 0, b0 + 1, b0 + 5, b0 + 4}, // -Y face
			{b0 + 2, b0 + 3, b0 + 7, b0 + 6}, // +Y face
			{b0 + 0, b0 + 4, b0 + 7, b0 + 3}, // -X face
			{b0 + 1, b0 + 2, b0 + 6, b0 + 5}, // +X face
		}
		for _, f := range faces {
			fmt.Fprintf(&b, "f %d %d %d %d\n", f[0], f[1], f[2], f[3])
		}
		vertexBase += 8
	}
	return b.String(), nil
}

// Summary emits a concise human-readable plain-text overview for one
// [voxel.VoxelProgram].
func Summary(program voxel.VoxelProgram) (string, error) {
	if err := program.Validate(); err != nil {
		return "", fmt.Errorf("validate voxel program: %w", err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "candidate:  %s\n", firstNonEmpty(program.CandidateID, "(unnamed)"))
	if program.Spec.Name != "" {
		fmt.Fprintf(&b, "spec:       %s", program.Spec.Name)
		if program.Spec.Revision != "" {
			fmt.Fprintf(&b, " (%s)", program.Spec.Revision)
		}
		b.WriteByte('\n')
	}
	v := program.Design.Volume
	fmt.Fprintf(&b, "volume:     %dx%dx%d  (%d cells total)\n",
		v.SizeX, v.SizeY, v.SizeZ, v.SizeX*v.SizeY*v.SizeZ)
	fmt.Fprintf(&b, "occupied:   %d cell(s)\n", len(program.Design.Occupied))
	if len(v.ForbiddenRegions) > 0 {
		fmt.Fprintf(&b, "forbidden:  %d region(s)\n", len(v.ForbiddenRegions))
	}
	if len(v.InterfaceRegions) > 0 {
		fmt.Fprintf(&b, "interface:  %d region(s)\n", len(v.InterfaceRegions))
	}
	if len(program.Spec.Materials) > 0 {
		fmt.Fprintf(&b, "materials:  %d\n", len(program.Spec.Materials))
	}
	if len(program.Spec.LoadCases) > 0 {
		fmt.Fprintf(&b, "load cases: %d\n", len(program.Spec.LoadCases))
	}
	return b.String(), nil
}

func sanitizeIdentifier(s string) string {
	if s == "" {
		return "unnamed"
	}
	var b strings.Builder
	for i, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r == '_':
			b.WriteRune(r)
		case i > 0 && r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
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
