// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package voxel

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func mustValidProgram() VoxelProgram {
	return VoxelProgram{
		CandidateID: "cand-voxel-01",
		Design: VoxelDesign{
			Volume: DesignVolume{
				SizeX: 4,
				SizeY: 3,
				SizeZ: 2,
				ForbiddenRegions: []InterfaceRegion{
					{
						Name: "keepout-core",
						Min:  VoxelIndex{X: 1, Y: 1, Z: 0},
						Max:  VoxelIndex{X: 1, Y: 1, Z: 1},
						Kind: "forbidden",
					},
				},
				InterfaceRegions: []InterfaceRegion{
					{
						Name:     "mount-left",
						Min:      VoxelIndex{X: 0, Y: 0, Z: 0},
						Max:      VoxelIndex{X: 0, Y: 2, Z: 1},
						Material: "aluminum",
						Kind:     "interface",
					},
				},
			},
			Occupied: []VoxelCell{
				{Coord: VoxelIndex{X: 2, Y: 0, Z: 0}, Material: "aluminum"},
				{Coord: VoxelIndex{X: 3, Y: 2, Z: 1}, Material: "polymer"},
			},
		},
		Spec: VoxelSpec{
			Name:     "cantilever-bracket",
			Domain:   "voxel",
			Revision: "v1",
			Materials: []Material{
				{ID: "aluminum", Name: "Aluminum", Properties: map[string]any{"density": 2700.0}},
				{ID: "polymer", Name: "Polymer", Properties: map[string]any{"density": 1200.0}},
			},
			LoadCases: []LoadCase{
				{Name: "pull-z", Kind: "force", Params: map[string]any{"newtons": 42.0}},
			},
			Metadata: map[string]any{"fixture": "pb-05", "train_ready": true},
		},
	}
}

func TestDesignVolumeValidateValid(t *testing.T) {
	volume := mustValidProgram().Design.Volume
	if err := volume.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
}

func TestDesignVolumeValidateEmptyOrMalformed(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*DesignVolume)
		wantErr string
	}{
		{
			name: "empty dimensions",
			mutate: func(v *DesignVolume) {
				v.SizeZ = 0
			},
			wantErr: "non-empty and positive",
		},
		{
			name: "malformed forbidden region",
			mutate: func(v *DesignVolume) {
				v.ForbiddenRegions[0].Max.X = -1
			},
			wantErr: "malformed bounds",
		},
		{
			name: "out of bounds interface region",
			mutate: func(v *DesignVolume) {
				v.InterfaceRegions[0].Max.X = v.SizeX
			},
			wantErr: "out of bounds",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			volume := mustValidProgram().Design.Volume
			tc.mutate(&volume)

			err := volume.Validate()
			if err == nil {
				t.Fatal("Validate() = nil, want error")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error %q does not mention %q", err, tc.wantErr)
			}
		})
	}
}

func TestDesignVolumeValidateOverlapForbiddenInterface(t *testing.T) {
	volume := mustValidProgram().Design.Volume
	volume.InterfaceRegions[0].Min = VoxelIndex{X: 1, Y: 1, Z: 0}
	volume.InterfaceRegions[0].Max = VoxelIndex{X: 1, Y: 1, Z: 1}

	err := volume.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want overlap error")
	}
	if !strings.Contains(err.Error(), "overlaps interface region") {
		t.Fatalf("error %q does not mention overlapping regions", err)
	}
}

func TestVoxelDesignValidateOutOfBoundsOccupiedCell(t *testing.T) {
	design := mustValidProgram().Design
	design.Occupied[1].Coord = VoxelIndex{X: design.Volume.SizeX, Y: 0, Z: 0}

	err := design.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want out-of-bounds occupied error")
	}
	if !strings.Contains(err.Error(), "occupied[1]") || !strings.Contains(err.Error(), "out of bounds") {
		t.Fatalf("error %q does not mention out-of-bounds occupied cell", err)
	}
}

func TestVoxelDesignValidateIsDeterministic(t *testing.T) {
	design := mustValidProgram().Design
	design.Occupied[0].Coord = VoxelIndex{X: -1, Y: 0, Z: 0}

	err1 := design.Validate()
	err2 := design.Validate()
	if (err1 == nil) != (err2 == nil) {
		t.Fatal("Validate() returned different nil-ness across calls")
	}
	if err1 != nil && err1.Error() != err2.Error() {
		t.Fatalf("Validate() not deterministic:\ncall1=%v\ncall2=%v", err1, err2)
	}
}

func TestVoxelProgramValidateDelegatesToDesign(t *testing.T) {
	program := mustValidProgram()
	if err := program.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}

	program.Design.Volume.SizeX = 0
	err := program.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want design validation error")
	}
}

func TestVoxelProgramJSONRoundTrip(t *testing.T) {
	orig := mustValidProgram()

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got VoxelProgram
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
