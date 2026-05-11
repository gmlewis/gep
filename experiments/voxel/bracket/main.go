// -*- compile-command: "go test ./experiments/voxel/bracket/..."; -*-
// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/design"
	"github.com/gmlewis/gep/v2/design/checkpoint"
	"github.com/gmlewis/gep/v2/design/objectives"
	"github.com/gmlewis/gep/v2/design/promotion"
	designscenarios "github.com/gmlewis/gep/v2/design/scenarios"
	"github.com/gmlewis/gep/v2/domains/voxel"
	voxelartifacts "github.com/gmlewis/gep/v2/domains/voxel/artifacts"
	voxelscenarios "github.com/gmlewis/gep/v2/domains/voxel/scenarios"
	"github.com/gmlewis/gep/v2/evolution"
	boolNodes "github.com/gmlewis/gep/v2/functions/bool_nodes"
)

type runConfig struct {
	Seed           int64
	PopulationSize int
	Generations    int
	OutputDir      string
}

type runResult struct {
	CandidateID   string
	Score         float64
	Karva         string
	OccupiedCells int
	Promoted      bool
}

const (
	candidateJSONArtifactName    = "candidate.json"
	candidateOBJArtifactName     = "candidate.obj"
	candidateSummaryArtifactName = "candidate.txt"
	promotionReportArtifactName  = "promotion_report.json"
	runManifestArtifactName      = "run_manifest.json"
	checkpointArtifactName       = "checkpoint.json"
)

var voxelObjectiveDefs = []objectives.ObjectiveDef{{
	Name:   "voxel_structural_score",
	Weight: 1,
	Kind:   objectives.Soft,
}}

func main() {
	cfg := runConfig{}
	flag.Int64Var(&cfg.Seed, "seed", 20260511, "deterministic evolution seed")
	flag.IntVar(&cfg.PopulationSize, "population", 64, "evolution population size")
	flag.IntVar(&cfg.Generations, "generations", 100, "maximum evolution generations")
	flag.StringVar(&cfg.OutputDir, "out", filepath.Join("artifacts", "voxel", "bracket"), "output directory for emitted artifacts")
	flag.Parse()

	result, err := runPilot(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bracket: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("bracket complete: candidate=%s score=%.2f occupied=%d promoted=%t karva=%s\n", result.CandidateID, result.Score, result.OccupiedCells, result.Promoted, result.Karva)
	fmt.Printf("artifacts written to %s\n", cfg.OutputDir)
}

func runPilot(cfg runConfig) (runResult, error) {
	if cfg.PopulationSize <= 0 {
		return runResult{}, errors.New("population must be > 0")
	}
	if cfg.Generations <= 0 {
		return runResult{}, errors.New("generations must be > 0")
	}

	set, registry, err := loadScenarioRegistry()
	if err != nil {
		return runResult{}, err
	}
	trainScenarios := registry.BySplit(designscenarios.Train)
	validationScenarios := registry.BySplit(designscenarios.Validation)
	if len(trainScenarios) == 0 {
		return runResult{}, errors.New("voxel fixture scenarios contain no train split")
	}
	if len(validationScenarios) == 0 {
		return runResult{}, errors.New("voxel fixture scenarios contain no validation split")
	}

	catalog, err := boolNodes.CatalogFromNames([]string{"Not", "And", "Or", "Xor"})
	if err != nil {
		return runResult{}, fmt.Errorf("build boolean function catalog: %w", err)
	}
	link, err := boolNodes.LinkFuncFrom("Or")
	if err != nil {
		return runResult{}, fmt.Errorf("build link function: %w", err)
	}

	population, err := evolution.NewWithSeed(cfg.Seed, catalog, cfg.PopulationSize, 6, 1, 4, 0, link, func(g core.Genome[bool]) float64 {
		return scoreCandidate(g, trainScenarios)
	})
	if err != nil {
		return runResult{}, fmt.Errorf("create seeded population: %w", err)
	}
	population.StopFunc = func(best evolution.Individual[bool]) bool {
		return best.Score >= 950
	}

	candidateID := fmt.Sprintf("voxel-bracket-seed-%d", cfg.Seed)
	manifest := design.RunManifest{
		RunID: candidateID,
		RunConfig: design.RunConfig{
			Domain:         "voxel",
			Experiment:     "bracket",
			PopulationSize: cfg.PopulationSize,
			NumGenerations: cfg.Generations,
		},
		Seeds: []design.SeedRecord{{
			Name:    "evolution_seed",
			Value:   cfg.Seed,
			Purpose: "deterministic evolution seed",
		}},
		ScenarioSplits: []design.ScenarioSplitSummary{
			{SplitName: string(designscenarios.Train), ScenarioCount: len(trainScenarios), Source: set.Source},
			{SplitName: string(designscenarios.Validation), ScenarioCount: len(validationScenarios), Source: set.Source},
		},
	}

	best := population.Evolve(cfg.Generations)
	program, err := decodeVoxelProgram(candidateID, best.Genome, trainScenarios[0])
	if err != nil {
		return runResult{}, err
	}
	artifactRefs, err := exportArtifacts(program, cfg.OutputDir)
	if err != nil {
		return runResult{}, err
	}
	manifest.Artifacts = append(manifest.Artifacts, artifactRefs...)

	trainResults, err := evaluateScenarios(best.Genome, trainScenarios)
	if err != nil {
		return runResult{}, err
	}
	validationResults, err := evaluateScenarios(best.Genome, validationScenarios)
	if err != nil {
		return runResult{}, err
	}
	report := promotion.Evaluate(candidateID, []promotion.AcceptanceCriterion{
		{Split: designscenarios.Train, MinAggregateScore: 700},
		{Split: designscenarios.Validation, MinAggregateScore: 700},
	}, []promotion.SplitEvalSummary{
		promotion.SummarizeSplit(designscenarios.Train, trainResults),
		promotion.SummarizeSplit(designscenarios.Validation, validationResults),
	})

	promotionReportPath := filepath.Join(cfg.OutputDir, promotionReportArtifactName)
	if err := writeJSONFile(promotionReportPath, report); err != nil {
		return runResult{}, fmt.Errorf("write promotion report: %w", err)
	}
	manifest.Artifacts = append(manifest.Artifacts, design.ArtifactRef{
		Name:   "promotion_report",
		Path:   promotionReportPath,
		Format: "json",
		Kind:   "promotion.report",
	})

	manifestPath := filepath.Join(cfg.OutputDir, runManifestArtifactName)
	if err := design.WriteRunManifestFile(manifestPath, &manifest); err != nil {
		return runResult{}, fmt.Errorf("write run manifest: %w", err)
	}
	manifest.Artifacts = append(manifest.Artifacts, design.ArtifactRef{
		Name:   "run_manifest",
		Path:   manifestPath,
		Format: "json",
		Kind:   "run.manifest",
	})
	if err := design.WriteRunManifestFile(manifestPath, &manifest); err != nil {
		return runResult{}, fmt.Errorf("rewrite run manifest with full artifact refs: %w", err)
	}

	validationSummary := promotion.SummarizeSplit(designscenarios.Validation, validationResults)
	checkpointPath := filepath.Join(cfg.OutputDir, checkpointArtifactName)
	checkpointRefs := append([]design.ArtifactRef(nil), manifest.Artifacts...)
	checkpointRefs = append(checkpointRefs, design.ArtifactRef{
		Name:   "checkpoint",
		Path:   checkpointPath,
		Format: "json",
		Kind:   "run.checkpoint",
	})
	snapshot := &checkpoint.Snapshot{
		Manifest: manifest,
		Elites: []checkpoint.EliteRecord{{
			CandidateID:    candidateID,
			Generation:     cfg.Generations,
			AggregateScore: validationSummary.MeanAggregateScore,
			Breakdown:      validationResults[0].Breakdown,
		}},
		Aggregate: checkpoint.AggregateSnapshot{
			BestScore:                best.Score,
			MeanScore:                validationSummary.MeanAggregateScore,
			TotalCandidatesEvaluated: cfg.PopulationSize * cfg.Generations,
		},
		ArtifactRefs: checkpointRefs,
	}
	if err := checkpoint.SaveFile(checkpointPath, snapshot); err != nil {
		return runResult{}, fmt.Errorf("write checkpoint snapshot: %w", err)
	}
	if _, err := checkpoint.LoadFile(checkpointPath); err != nil {
		return runResult{}, fmt.Errorf("reload checkpoint snapshot: %w", err)
	}

	return runResult{
		CandidateID:   candidateID,
		Score:         best.Score,
		Karva:         best.Genome.KarvaString(),
		OccupiedCells: len(program.Design.Occupied),
		Promoted:      report.Promoted,
	}, nil
}

func loadScenarioRegistry() (*designscenarios.ScenarioSet, *designscenarios.ScenarioRegistry, error) {
	set, err := voxelscenarios.LoadFixtureSet()
	if err != nil {
		return nil, nil, fmt.Errorf("load voxel fixture scenarios: %w", err)
	}
	registry := &designscenarios.ScenarioRegistry{Sets: []designscenarios.ScenarioSet{*set}}
	if err := registry.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validate voxel fixture scenarios: %w", err)
	}
	return set, registry, nil
}

func loadTrainScenarios() ([]designscenarios.Scenario, error) {
	_, registry, err := loadScenarioRegistry()
	if err != nil {
		return nil, err
	}
	train := registry.BySplit(designscenarios.Train)
	if len(train) == 0 {
		return nil, errors.New("voxel fixture scenarios contain no train split")
	}
	return train, nil
}

func scoreCandidate(genome core.Genome[bool], scenarios []designscenarios.Scenario) float64 {
	if len(scenarios) == 0 {
		return 0
	}
	var total float64
	for _, sc := range scenarios {
		total += scoreScenario(genome, sc)
	}
	return total / float64(len(scenarios))
}

func scoreScenario(genome core.Genome[bool], scenario designscenarios.Scenario) float64 {
	score, _ := scoreScenarioAggregate(genome, scenario)
	return score
}

func evaluateScenarios(genome core.Genome[bool], scenarios []designscenarios.Scenario) ([]objectives.AggregateResult, error) {
	results := make([]objectives.AggregateResult, 0, len(scenarios))
	for _, sc := range scenarios {
		score, rejected := scoreScenarioAggregate(genome, sc)
		results = append(results, objectives.Score(voxelObjectiveDefs, map[string]float64{
			"voxel_structural_score": score,
		}, rejected, 0))
	}
	return results, nil
}

func scoreScenarioAggregate(genome core.Genome[bool], scenario designscenarios.Scenario) (score float64, hardFailed bool) {
	program, err := decodeVoxelProgram("score", genome, scenario)
	if err != nil {
		return 0, true
	}
	occupied := len(program.Design.Occupied)
	maxCells := scenarioMaxCells(scenario)
	if maxCells > 0 && occupied > maxCells {
		return 0, true
	}

	leftCoverage := interfaceCoverage(program.Design, "anchor_left")
	rightCoverage := interfaceCoverage(program.Design, "load_right")
	connectivity := 0.0
	if hasPathLeftToRight(program.Design) {
		connectivity = 1
	}

	targetCells := maxCells
	if targetCells <= 0 {
		v := program.Design.Volume
		targetCells = (v.SizeX * v.SizeY * v.SizeZ) / 2
	}
	densityFit := 1 - math.Abs(float64(occupied-targetCells))/float64(maxInt(targetCells, 1))
	if densityFit < 0 {
		densityFit = 0
	}

	return 500*connectivity + 150*leftCoverage + 150*rightCoverage + 200*densityFit, false
}

func decodeVoxelProgram(candidateID string, genome core.Genome[bool], scenario designscenarios.Scenario) (voxel.VoxelProgram, error) {
	volume, err := scenarioVolume(scenario)
	if err != nil {
		return voxel.VoxelProgram{}, fmt.Errorf("decode voxel program: %w", err)
	}
	design, err := decodeDesign(genome, volume)
	if err != nil {
		return voxel.VoxelProgram{}, fmt.Errorf("decode voxel program: %w", err)
	}
	program := voxel.VoxelProgram{
		CandidateID: candidateID,
		Design:      design,
		Spec: voxel.VoxelSpec{
			Name:   "bracket",
			Domain: "voxel",
			Metadata: map[string]any{
				"scenario_id": string(scenario.ID),
			},
		},
	}
	if err := program.Validate(); err != nil {
		return voxel.VoxelProgram{}, fmt.Errorf("decode voxel program: %w", err)
	}
	return program, nil
}

func decodeDesign(genome core.Genome[bool], volume voxel.DesignVolume) (voxel.VoxelDesign, error) {
	occupied := make([]voxel.VoxelCell, 0, volume.SizeX*volume.SizeY*volume.SizeZ)
	for z := 0; z < volume.SizeZ; z++ {
		for y := 0; y < volume.SizeY; y++ {
			for x := 0; x < volume.SizeX; x++ {
				coord := voxel.VoxelIndex{X: x, Y: y, Z: z}
				if isForbiddenCell(volume.ForbiddenRegions, coord) {
					continue
				}
				fill := isSpineCell(coord)
				if !fill {
					features := voxelFeatures(volume, coord)
					value, err := genome.Eval(features)
					if err != nil {
						return voxel.VoxelDesign{}, fmt.Errorf("evaluate voxel occupancy at (%d,%d,%d): %w", x, y, z, err)
					}
					fill = value
				}
				if !fill {
					continue
				}
				occupied = append(occupied, voxel.VoxelCell{Coord: coord, Material: "steel"})
			}
		}
	}

	design := voxel.VoxelDesign{
		Volume:   volume,
		Occupied: occupied,
	}
	if err := design.Validate(); err != nil {
		return voxel.VoxelDesign{}, err
	}
	return design, nil
}

func voxelFeatures(volume voxel.DesignVolume, coord voxel.VoxelIndex) []bool {
	return []bool{
		coord.X <= volume.SizeX/3,
		coord.X >= (2*volume.SizeX)/3,
		coord.Y <= volume.SizeY/2,
		coord.Z == 0,
	}
}

func isSpineCell(coord voxel.VoxelIndex) bool {
	return coord.Y == 0 && coord.Z == 0
}

func scenarioVolume(scenario designscenarios.Scenario) (voxel.DesignVolume, error) {
	rawVolume, ok := scenario.Params["volume"]
	if !ok {
		return voxel.DesignVolume{}, errors.New("scenario missing volume params")
	}
	volumeMap, ok := rawVolume.(map[string]any)
	if !ok {
		return voxel.DesignVolume{}, fmt.Errorf("scenario volume params have type %T, want object", rawVolume)
	}
	sizeX, ok := numericParam(volumeMap["size_x"])
	if !ok || sizeX <= 0 {
		return voxel.DesignVolume{}, errors.New("scenario volume size_x must be > 0")
	}
	sizeY, ok := numericParam(volumeMap["size_y"])
	if !ok || sizeY <= 0 {
		return voxel.DesignVolume{}, errors.New("scenario volume size_y must be > 0")
	}
	sizeZ, ok := numericParam(volumeMap["size_z"])
	if !ok || sizeZ <= 0 {
		return voxel.DesignVolume{}, errors.New("scenario volume size_z must be > 0")
	}

	volume := voxel.DesignVolume{
		SizeX: sizeX,
		SizeY: sizeY,
		SizeZ: sizeZ,
		InterfaceRegions: []voxel.InterfaceRegion{
			{
				Name: "anchor_left",
				Min:  voxel.VoxelIndex{X: 0, Y: 0, Z: 0},
				Max:  voxel.VoxelIndex{X: 0, Y: sizeY - 1, Z: sizeZ - 1},
				Kind: "interface",
			},
			{
				Name: "load_right",
				Min:  voxel.VoxelIndex{X: sizeX - 1, Y: 0, Z: 0},
				Max:  voxel.VoxelIndex{X: sizeX - 1, Y: sizeY - 1, Z: sizeZ - 1},
				Kind: "interface",
			},
		},
	}
	if sizeX > 3 && sizeY > 2 {
		volume.ForbiddenRegions = []voxel.InterfaceRegion{{
			Name: "keepout_mid",
			Min:  voxel.VoxelIndex{X: sizeX / 2, Y: 1, Z: 0},
			Max:  voxel.VoxelIndex{X: sizeX / 2, Y: sizeY - 2, Z: sizeZ - 1},
			Kind: "forbidden",
		}}
	}
	if err := volume.Validate(); err != nil {
		return voxel.DesignVolume{}, err
	}
	return volume, nil
}

func scenarioMaxCells(scenario designscenarios.Scenario) int {
	raw, ok := scenario.Params["max_cells"]
	if !ok {
		return 0
	}
	value, ok := numericParam(raw)
	if !ok || value < 0 {
		return 0
	}
	return value
}

func numericParam(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case float64:
		return int(n), true
	default:
		return 0, false
	}
}

func interfaceCoverage(design voxel.VoxelDesign, interfaceName string) float64 {
	var region voxel.InterfaceRegion
	found := false
	for _, r := range design.Volume.InterfaceRegions {
		if r.Name == interfaceName {
			region = r
			found = true
			break
		}
	}
	if !found {
		return 0
	}

	total := 0
	covered := 0
	occupied := occupiedSet(design.Occupied)
	for z := region.Min.Z; z <= region.Max.Z; z++ {
		for y := region.Min.Y; y <= region.Max.Y; y++ {
			for x := region.Min.X; x <= region.Max.X; x++ {
				total++
				if _, ok := occupied[voxel.VoxelIndex{X: x, Y: y, Z: z}]; ok {
					covered++
				}
			}
		}
	}
	if total == 0 {
		return 0
	}
	return float64(covered) / float64(total)
}

func hasPathLeftToRight(design voxel.VoxelDesign) bool {
	occupied := occupiedSet(design.Occupied)
	if len(occupied) == 0 {
		return false
	}

	queue := make([]voxel.VoxelIndex, 0, len(occupied))
	seen := make(map[voxel.VoxelIndex]struct{}, len(occupied))
	for coord := range occupied {
		if coord.X == 0 {
			queue = append(queue, coord)
			seen[coord] = struct{}{}
		}
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.X == design.Volume.SizeX-1 {
			return true
		}
		for _, next := range neighbors(cur) {
			if !contains(design.Volume, next) {
				continue
			}
			if _, ok := occupied[next]; !ok {
				continue
			}
			if _, ok := seen[next]; ok {
				continue
			}
			seen[next] = struct{}{}
			queue = append(queue, next)
		}
	}
	return false
}

func occupiedSet(cells []voxel.VoxelCell) map[voxel.VoxelIndex]struct{} {
	set := make(map[voxel.VoxelIndex]struct{}, len(cells))
	for _, cell := range cells {
		set[cell.Coord] = struct{}{}
	}
	return set
}

func neighbors(coord voxel.VoxelIndex) []voxel.VoxelIndex {
	return []voxel.VoxelIndex{
		{X: coord.X + 1, Y: coord.Y, Z: coord.Z},
		{X: coord.X - 1, Y: coord.Y, Z: coord.Z},
		{X: coord.X, Y: coord.Y + 1, Z: coord.Z},
		{X: coord.X, Y: coord.Y - 1, Z: coord.Z},
		{X: coord.X, Y: coord.Y, Z: coord.Z + 1},
		{X: coord.X, Y: coord.Y, Z: coord.Z - 1},
	}
}

func contains(volume voxel.DesignVolume, coord voxel.VoxelIndex) bool {
	return coord.X >= 0 && coord.X < volume.SizeX &&
		coord.Y >= 0 && coord.Y < volume.SizeY &&
		coord.Z >= 0 && coord.Z < volume.SizeZ
}

func isForbiddenCell(forbidden []voxel.InterfaceRegion, coord voxel.VoxelIndex) bool {
	for _, region := range forbidden {
		if inRegion(coord, region) {
			return true
		}
	}
	return false
}

func inRegion(coord voxel.VoxelIndex, region voxel.InterfaceRegion) bool {
	return coord.X >= region.Min.X && coord.X <= region.Max.X &&
		coord.Y >= region.Min.Y && coord.Y <= region.Max.Y &&
		coord.Z >= region.Min.Z && coord.Z <= region.Max.Z
}

func exportArtifacts(program voxel.VoxelProgram, outputDir string) ([]design.ArtifactRef, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("create output directory %q: %w", outputDir, err)
	}

	jsonBytes, err := voxelartifacts.JSON(program)
	if err != nil {
		return nil, fmt.Errorf("emit JSON artifact: %w", err)
	}
	jsonPath := filepath.Join(outputDir, candidateJSONArtifactName)
	if err := os.WriteFile(jsonPath, jsonBytes, 0o644); err != nil {
		return nil, fmt.Errorf("write JSON artifact: %w", err)
	}

	objText, err := voxelartifacts.OBJ(program)
	if err != nil {
		return nil, fmt.Errorf("emit OBJ artifact: %w", err)
	}
	objPath := filepath.Join(outputDir, candidateOBJArtifactName)
	if err := os.WriteFile(objPath, []byte(objText), 0o644); err != nil {
		return nil, fmt.Errorf("write OBJ artifact: %w", err)
	}

	summaryText, err := voxelartifacts.Summary(program)
	if err != nil {
		return nil, fmt.Errorf("emit summary artifact: %w", err)
	}
	summaryPath := filepath.Join(outputDir, candidateSummaryArtifactName)
	if err := os.WriteFile(summaryPath, []byte(summaryText), 0o644); err != nil {
		return nil, fmt.Errorf("write summary artifact: %w", err)
	}

	return []design.ArtifactRef{
		{Name: "promoted_voxel_json", Path: jsonPath, Format: "json", Kind: "voxel.program"},
		{Name: "promoted_voxel_obj", Path: objPath, Format: "obj", Kind: "voxel.mesh"},
		{Name: "promoted_voxel_summary", Path: summaryPath, Format: "txt", Kind: "voxel.summary"},
	}, nil
}

func writeJSONFile(filename string, value any) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create %q: %w", filename, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		return fmt.Errorf("encode %q: %w", filename, err)
	}
	return nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
