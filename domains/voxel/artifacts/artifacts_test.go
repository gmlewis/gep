// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package artifacts

import (
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/domains/voxel"
)

func testProgram() voxel.VoxelProgram {
	return voxel.VoxelProgram{
		CandidateID: "cand-voxel-01",
		Design: voxel.VoxelDesign{
			Volume: voxel.DesignVolume{
				SizeX: 4,
				SizeY: 3,
				SizeZ: 2,
				ForbiddenRegions: []voxel.InterfaceRegion{
					{
						Name: "keepout-core",
						Min:  voxel.VoxelIndex{X: 1, Y: 1, Z: 0},
						Max:  voxel.VoxelIndex{X: 1, Y: 1, Z: 1},
						Kind: "forbidden",
					},
				},
				InterfaceRegions: []voxel.InterfaceRegion{
					{
						Name:     "mount-left",
						Min:      voxel.VoxelIndex{X: 0, Y: 0, Z: 0},
						Max:      voxel.VoxelIndex{X: 0, Y: 2, Z: 1},
						Material: "aluminum",
						Kind:     "interface",
					},
				},
			},
			Occupied: []voxel.VoxelCell{
				{Coord: voxel.VoxelIndex{X: 2, Y: 0, Z: 0}, Material: "aluminum"},
				{Coord: voxel.VoxelIndex{X: 3, Y: 2, Z: 1}, Material: "polymer"},
			},
		},
		Spec: voxel.VoxelSpec{
			Name:     "cantilever-bracket",
			Domain:   "voxel",
			Revision: "v1",
			Materials: []voxel.Material{
				{ID: "aluminum", Name: "Aluminum", Properties: map[string]any{"density": 2700.0}},
				{ID: "polymer", Name: "Polymer", Properties: map[string]any{"density": 1200.0}},
			},
			LoadCases: []voxel.LoadCase{
				{Name: "pull-z", Kind: "force", Params: map[string]any{"newtons": 42.0}},
			},
			Metadata: map[string]any{"fixture": "pb-06", "train_ready": true},
		},
	}
}

func TestJSON(t *testing.T) {
	got, err := JSON(testProgram())
	if err != nil {
		t.Fatalf("JSON() error = %v", err)
	}
	const want = "{\n" +
		"  \"candidate_id\": \"cand-voxel-01\",\n" +
		"  \"design\": {\n" +
		"    \"volume\": {\n" +
		"      \"size_x\": 4,\n" +
		"      \"size_y\": 3,\n" +
		"      \"size_z\": 2,\n" +
		"      \"forbidden_regions\": [\n" +
		"        {\n" +
		"          \"name\": \"keepout-core\",\n" +
		"          \"min\": {\n" +
		"            \"x\": 1,\n" +
		"            \"y\": 1,\n" +
		"            \"z\": 0\n" +
		"          },\n" +
		"          \"max\": {\n" +
		"            \"x\": 1,\n" +
		"            \"y\": 1,\n" +
		"            \"z\": 1\n" +
		"          },\n" +
		"          \"kind\": \"forbidden\"\n" +
		"        }\n" +
		"      ],\n" +
		"      \"interface_regions\": [\n" +
		"        {\n" +
		"          \"name\": \"mount-left\",\n" +
		"          \"min\": {\n" +
		"            \"x\": 0,\n" +
		"            \"y\": 0,\n" +
		"            \"z\": 0\n" +
		"          },\n" +
		"          \"max\": {\n" +
		"            \"x\": 0,\n" +
		"            \"y\": 2,\n" +
		"            \"z\": 1\n" +
		"          },\n" +
		"          \"material\": \"aluminum\",\n" +
		"          \"kind\": \"interface\"\n" +
		"        }\n" +
		"      ]\n" +
		"    },\n" +
		"    \"occupied\": [\n" +
		"      {\n" +
		"        \"coord\": {\n" +
		"          \"x\": 2,\n" +
		"          \"y\": 0,\n" +
		"          \"z\": 0\n" +
		"        },\n" +
		"        \"material\": \"aluminum\"\n" +
		"      },\n" +
		"      {\n" +
		"        \"coord\": {\n" +
		"          \"x\": 3,\n" +
		"          \"y\": 2,\n" +
		"          \"z\": 1\n" +
		"        },\n" +
		"        \"material\": \"polymer\"\n" +
		"      }\n" +
		"    ]\n" +
		"  },\n" +
		"  \"spec\": {\n" +
		"    \"name\": \"cantilever-bracket\",\n" +
		"    \"domain\": \"voxel\",\n" +
		"    \"revision\": \"v1\",\n" +
		"    \"materials\": [\n" +
		"      {\n" +
		"        \"id\": \"aluminum\",\n" +
		"        \"name\": \"Aluminum\",\n" +
		"        \"properties\": {\n" +
		"          \"density\": 2700\n" +
		"        }\n" +
		"      },\n" +
		"      {\n" +
		"        \"id\": \"polymer\",\n" +
		"        \"name\": \"Polymer\",\n" +
		"        \"properties\": {\n" +
		"          \"density\": 1200\n" +
		"        }\n" +
		"      }\n" +
		"    ],\n" +
		"    \"load_cases\": [\n" +
		"      {\n" +
		"        \"name\": \"pull-z\",\n" +
		"        \"kind\": \"force\",\n" +
		"        \"params\": {\n" +
		"          \"newtons\": 42\n" +
		"        }\n" +
		"      }\n" +
		"    ],\n" +
		"    \"metadata\": {\n" +
		"      \"fixture\": \"pb-06\",\n" +
		"      \"train_ready\": true\n" +
		"    }\n" +
		"  }\n" +
		"}\n"
	if string(got) != want {
		t.Fatalf("JSON() mismatch:\n--- got ---\n%s--- want ---\n%s", got, want)
	}
}

func TestOBJ(t *testing.T) {
	got, err := OBJ(testProgram())
	if err != nil {
		t.Fatalf("OBJ() error = %v", err)
	}
	const want = "" +
		"# voxel candidate: cand-voxel-01\n" +
		"# spec: cantilever-bracket\n" +
		"o cantilever_bracket\n" +
		"v 2 0 0\n" +
		"v 3 0 0\n" +
		"v 3 1 0\n" +
		"v 2 1 0\n" +
		"v 2 0 1\n" +
		"v 3 0 1\n" +
		"v 3 1 1\n" +
		"v 2 1 1\n" +
		"# material: aluminum\n" +
		"f 1 4 3 2\n" +
		"f 5 6 7 8\n" +
		"f 1 2 6 5\n" +
		"f 3 4 8 7\n" +
		"f 1 5 8 4\n" +
		"f 2 3 7 6\n" +
		"v 3 2 1\n" +
		"v 4 2 1\n" +
		"v 4 3 1\n" +
		"v 3 3 1\n" +
		"v 3 2 2\n" +
		"v 4 2 2\n" +
		"v 4 3 2\n" +
		"v 3 3 2\n" +
		"# material: polymer\n" +
		"f 9 12 11 10\n" +
		"f 13 14 15 16\n" +
		"f 9 10 14 13\n" +
		"f 11 12 16 15\n" +
		"f 9 13 16 12\n" +
		"f 10 11 15 14\n"
	if got != want {
		t.Fatalf("OBJ() mismatch:\n--- got ---\n%s--- want ---\n%s", got, want)
	}
}

func TestSummary(t *testing.T) {
	got, err := Summary(testProgram())
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	const want = "" +
		"candidate:  cand-voxel-01\n" +
		"spec:       cantilever-bracket (v1)\n" +
		"volume:     4x3x2  (24 cells total)\n" +
		"occupied:   2 cell(s)\n" +
		"forbidden:  1 region(s)\n" +
		"interface:  1 region(s)\n" +
		"materials:  2\n" +
		"load cases: 1\n"
	if got != want {
		t.Fatalf("Summary() mismatch:\n--- got ---\n%s--- want ---\n%s", got, want)
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

	obj1, err := OBJ(program)
	if err != nil {
		t.Fatalf("OBJ #1: %v", err)
	}
	obj2, err := OBJ(program)
	if err != nil {
		t.Fatalf("OBJ #2: %v", err)
	}
	if obj1 != obj2 {
		t.Fatal("OBJ output is not deterministic")
	}

	summary1, err := Summary(program)
	if err != nil {
		t.Fatalf("Summary #1: %v", err)
	}
	summary2, err := Summary(program)
	if err != nil {
		t.Fatalf("Summary #2: %v", err)
	}
	if summary1 != summary2 {
		t.Fatal("Summary output is not deterministic")
	}
}

func TestEmittersRejectInvalidProgram(t *testing.T) {
	program := testProgram()
	program.Design.Volume.SizeX = 0

	for _, tc := range []struct {
		name string
		fn   func(voxel.VoxelProgram) error
	}{
		{
			name: "JSON",
			fn: func(p voxel.VoxelProgram) error {
				_, err := JSON(p)
				return err
			},
		},
		{
			name: "OBJ",
			fn: func(p voxel.VoxelProgram) error {
				_, err := OBJ(p)
				return err
			},
		},
		{
			name: "Summary",
			fn: func(p voxel.VoxelProgram) error {
				_, err := Summary(p)
				return err
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn(program)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), "non-empty and positive") {
				t.Fatalf("error %q does not mention invalid design", err)
			}
		})
	}
}
