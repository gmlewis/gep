// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package voxel

import "fmt"

// VoxelIndex identifies one occupied cell or region corner inside a design
// volume.
type VoxelIndex struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

// VoxelCell stores one occupied voxel cell and its optional material binding.
type VoxelCell struct {
	Coord    VoxelIndex `json:"coord"`
	Material string     `json:"material,omitempty"`
}

// InterfaceRegion describes one inclusive axis-aligned region within a design
// volume.
type InterfaceRegion struct {
	Name     string     `json:"name"`
	Min      VoxelIndex `json:"min"`
	Max      VoxelIndex `json:"max"`
	Material string     `json:"material,omitempty"`
	Kind     string     `json:"kind,omitempty"`
}

// DesignVolume defines the bounded voxel lattice and reserved regions for a
// candidate design.
type DesignVolume struct {
	SizeX            int               `json:"size_x"`
	SizeY            int               `json:"size_y"`
	SizeZ            int               `json:"size_z"`
	ForbiddenRegions []InterfaceRegion `json:"forbidden_regions,omitempty"`
	InterfaceRegions []InterfaceRegion `json:"interface_regions,omitempty"`
}

// Material stores reusable material metadata referenced by a voxel design.
type Material struct {
	ID         string         `json:"id"`
	Name       string         `json:"name,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}

// LoadCase stores reusable load metadata for later pilot-specific evaluation.
type LoadCase struct {
	Name   string         `json:"name"`
	Kind   string         `json:"kind"`
	Params map[string]any `json:"params,omitempty"`
}

// VoxelSpec describes high-level metadata for a voxel candidate/program.
type VoxelSpec struct {
	Name      string         `json:"name"`
	Domain    string         `json:"domain"`
	Revision  string         `json:"revision,omitempty"`
	Materials []Material     `json:"materials,omitempty"`
	LoadCases []LoadCase     `json:"load_cases,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// VoxelDesign stores one bounded voxel occupancy design.
type VoxelDesign struct {
	Volume   DesignVolume `json:"volume"`
	Occupied []VoxelCell  `json:"occupied,omitempty"`
}

// VoxelProgram is a serializable voxel-domain candidate package.
type VoxelProgram struct {
	CandidateID string      `json:"candidate_id"`
	Design      VoxelDesign `json:"design"`
	Spec        VoxelSpec   `json:"spec"`
}

// Validate returns an error when v has an empty or malformed volume definition.
func (v DesignVolume) Validate() error {
	if v.SizeX <= 0 || v.SizeY <= 0 || v.SizeZ <= 0 {
		return fmt.Errorf("design volume must be non-empty and positive; got %dx%dx%d", v.SizeX, v.SizeY, v.SizeZ)
	}

	for i, region := range v.ForbiddenRegions {
		if err := v.validateRegion("forbidden", i, region); err != nil {
			return err
		}
	}
	for i, region := range v.InterfaceRegions {
		if err := v.validateRegion("interface", i, region); err != nil {
			return err
		}
	}
	for i, forbidden := range v.ForbiddenRegions {
		for j, iface := range v.InterfaceRegions {
			if regionsOverlap(forbidden, iface) {
				return fmt.Errorf("forbidden region[%d] %q overlaps interface region[%d] %q", i, forbidden.Name, j, iface.Name)
			}
		}
	}

	return nil
}

// Validate returns an error when d has an invalid volume or occupied voxels
// outside that volume.
func (d VoxelDesign) Validate() error {
	if err := d.Volume.Validate(); err != nil {
		return err
	}

	for i, cell := range d.Occupied {
		if !d.Volume.contains(cell.Coord) {
			return fmt.Errorf("occupied[%d] coord (%d,%d,%d) is out of bounds for volume %dx%dx%d", i, cell.Coord.X, cell.Coord.Y, cell.Coord.Z, d.Volume.SizeX, d.Volume.SizeY, d.Volume.SizeZ)
		}
	}

	return nil
}

// Validate returns an error when p.Design fails voxel validation.
func (p VoxelProgram) Validate() error {
	return p.Design.Validate()
}

func (v DesignVolume) validateRegion(kind string, index int, region InterfaceRegion) error {
	if region.Max.X < region.Min.X || region.Max.Y < region.Min.Y || region.Max.Z < region.Min.Z {
		return fmt.Errorf("%s region[%d] %q has malformed bounds: max must be >= min", kind, index, region.Name)
	}
	if region.Min.X < 0 || region.Min.Y < 0 || region.Min.Z < 0 {
		return fmt.Errorf("%s region[%d] %q has malformed bounds: min coordinates must be >= 0", kind, index, region.Name)
	}
	if !v.contains(region.Min) || !v.contains(region.Max) {
		return fmt.Errorf("%s region[%d] %q is out of bounds for volume %dx%dx%d", kind, index, region.Name, v.SizeX, v.SizeY, v.SizeZ)
	}
	return nil
}

func (v DesignVolume) contains(index VoxelIndex) bool {
	return index.X >= 0 && index.X < v.SizeX &&
		index.Y >= 0 && index.Y < v.SizeY &&
		index.Z >= 0 && index.Z < v.SizeZ
}

func regionsOverlap(a, b InterfaceRegion) bool {
	return rangesOverlap(a.Min.X, a.Max.X, b.Min.X, b.Max.X) &&
		rangesOverlap(a.Min.Y, a.Max.Y, b.Min.Y, b.Max.Y) &&
		rangesOverlap(a.Min.Z, a.Max.Z, b.Min.Z, b.Max.Z)
}

func rangesOverlap(aMin, aMax, bMin, bMax int) bool {
	return aMin <= bMax && bMin <= aMax
}
