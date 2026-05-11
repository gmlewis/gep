// -*- compile-command: "go test ./experiments/control/mass_spring_damper/..."; -*-
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
	CandidateID      string
	Score            float64
	Karva            string
	Promoted         bool
	TrainMeanScore   float64
	ValidMeanScore   float64
	TestMeanScore    float64
	FinalAbsPosition float64
}

type plantScenario struct {
	Mass           float64
	SpringK        float64
	DampingC       float64
	TimeStep       float64
	Steps          int
	TargetPosition float64
	InitialPos     float64
	InitialVel     float64
	Disturbance    float64
	MaxForce       float64
	MaxAbsPosition float64
	SettleBand     float64
}

type simulationResult struct {
	MeanAbsError   float64
	MeanAbsControl float64
	FinalError     float64
	FinalPosition  float64
	HardFailed     bool
}

type controllerPolicyArtifact struct {
	CandidateID string             `json:"candidate_id"`
	Domain      string             `json:"domain"`
	Experiment  string             `json:"experiment"`
	Karva       string             `json:"karva"`
	Terminals   []string           `json:"terminals"`
	Scores      map[string]float64 `json:"scores,omitempty"`
}

const (
	controllerJSONArtifactName  = "controller_policy.json"
	controllerSummaryArtifact   = "controller_summary.txt"
	promotionReportArtifactName = "promotion_report.json"
	runManifestArtifactName     = "run_manifest.json"
	checkpointArtifactName      = "checkpoint.json"
	scenarioSetSource           = "embedded://experiments/control/mass_spring_damper/scenarios"
)

var controlObjectiveDefs = []objectives.ObjectiveDef{{
	Name:   "tracking_score",
	Weight: 1,
	Kind:   objectives.Soft,
}}

func main() {
	cfg := runConfig{}
	flag.Int64Var(&cfg.Seed, "seed", 20260511, "deterministic evolution seed")
	flag.IntVar(&cfg.PopulationSize, "population", 72, "evolution population size")
	flag.IntVar(&cfg.Generations, "generations", 140, "maximum evolution generations")
	flag.StringVar(&cfg.OutputDir, "out", filepath.Join("artifacts", "control", "mass_spring_damper"), "output directory for emitted artifacts")
	flag.Parse()

	result, err := runPilot(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mass_spring_damper: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("mass_spring_damper complete: candidate=%s score=%.2f promoted=%t karva=%s\n", result.CandidateID, result.Score, result.Promoted, result.Karva)
	fmt.Printf("split means: train=%.2f validation=%.2f test=%.2f final_abs_position=%.4f\n", result.TrainMeanScore, result.ValidMeanScore, result.TestMeanScore, result.FinalAbsPosition)
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
	testScenarios := registry.BySplit(designscenarios.Test)
	if len(trainScenarios) == 0 {
		return runResult{}, errors.New("control fixture scenarios contain no train split")
	}
	if len(validationScenarios) == 0 {
		return runResult{}, errors.New("control fixture scenarios contain no validation split")
	}
	if len(testScenarios) == 0 {
		return runResult{}, errors.New("control fixture scenarios contain no test split")
	}

	candidateID := fmt.Sprintf("mass-spring-damper-seed-%d", cfg.Seed)
	manifest := design.RunManifest{
		RunID: candidateID,
		RunConfig: design.RunConfig{
			Domain:         "control",
			Experiment:     "mass_spring_damper",
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
			{SplitName: string(designscenarios.Test), ScenarioCount: len(testScenarios), Source: set.Source},
		},
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
		return best.Score >= 930
	}

	best := population.Evolve(cfg.Generations)
	trainResults, err := evaluateScenarios(best.Genome, trainScenarios)
	if err != nil {
		return runResult{}, err
	}
	validationResults, err := evaluateScenarios(best.Genome, validationScenarios)
	if err != nil {
		return runResult{}, err
	}
	testResults, err := evaluateScenarios(best.Genome, testScenarios)
	if err != nil {
		return runResult{}, err
	}

	trainSummary := promotion.SummarizeSplit(designscenarios.Train, trainResults)
	validationSummary := promotion.SummarizeSplit(designscenarios.Validation, validationResults)
	testSummary := promotion.SummarizeSplit(designscenarios.Test, testResults)

	report := promotion.Evaluate(candidateID, []promotion.AcceptanceCriterion{
		{Split: designscenarios.Train, MinAggregateScore: 250},
		{Split: designscenarios.Validation, MinAggregateScore: 250},
		{Split: designscenarios.Test, MinAggregateScore: 250},
	}, []promotion.SplitEvalSummary{trainSummary, validationSummary, testSummary})

	policy := decodeControllerPolicy(candidateID, best.Genome)
	policy.Scores = map[string]float64{
		"train_mean":      trainSummary.MeanAggregateScore,
		"validation_mean": validationSummary.MeanAggregateScore,
		"test_mean":       testSummary.MeanAggregateScore,
	}
	artifactRefs, err := exportArtifacts(policy, cfg.OutputDir)
	if err != nil {
		return runResult{}, err
	}
	manifest.Artifacts = append(manifest.Artifacts, artifactRefs...)

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

	checkpointPath := filepath.Join(cfg.OutputDir, checkpointArtifactName)
	checkpointRefs := append([]design.ArtifactRef(nil), manifest.Artifacts...)
	checkpointRefs = append(checkpointRefs, design.ArtifactRef{
		Name:   "checkpoint",
		Path:   checkpointPath,
		Format: "json",
		Kind:   "run.checkpoint",
	})

	validationPlant, err := scenarioPlant(validationScenarios[0])
	if err != nil {
		return runResult{}, err
	}
	validationSim, err := simulateController(best.Genome, validationPlant)
	if err != nil {
		return runResult{}, err
	}

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
		CandidateID:      candidateID,
		Score:            best.Score,
		Karva:            best.Genome.KarvaString(),
		Promoted:         report.Promoted,
		TrainMeanScore:   trainSummary.MeanAggregateScore,
		ValidMeanScore:   validationSummary.MeanAggregateScore,
		TestMeanScore:    testSummary.MeanAggregateScore,
		FinalAbsPosition: math.Abs(validationSim.FinalPosition),
	}, nil
}

func loadScenarioRegistry() (*designscenarios.ScenarioSet, *designscenarios.ScenarioRegistry, error) {
	set := &designscenarios.ScenarioSet{
		Name:   "mass_spring_damper_smoke",
		Source: scenarioSetSource,
		Scenarios: []designscenarios.Scenario{
			{
				ID:    "msd_train_disturbance_pos",
				Split: designscenarios.Train,
				Params: map[string]any{
					"mass":             1.0,
					"spring_k":         2.8,
					"damping_c":        0.7,
					"dt":               0.03,
					"steps":            140,
					"target_position":  1.0,
					"initial_position": 0.2,
					"initial_velocity": 0.0,
					"disturbance":      0.15,
					"max_force":        2.0,
					"max_abs_position": 4.0,
					"settle_band":      0.08,
				},
			},
			{
				ID:    "msd_train_disturbance_neg",
				Split: designscenarios.Train,
				Params: map[string]any{
					"mass":             1.0,
					"spring_k":         3.1,
					"damping_c":        0.8,
					"dt":               0.03,
					"steps":            140,
					"target_position":  1.0,
					"initial_position": -0.1,
					"initial_velocity": 0.1,
					"disturbance":      -0.20,
					"max_force":        2.2,
					"max_abs_position": 4.0,
					"settle_band":      0.08,
				},
			},
			{
				ID:    "msd_validation_corner",
				Split: designscenarios.Validation,
				Params: map[string]any{
					"mass":             1.1,
					"spring_k":         3.2,
					"damping_c":        0.9,
					"dt":               0.03,
					"steps":            150,
					"target_position":  1.0,
					"initial_position": 0.35,
					"initial_velocity": -0.15,
					"disturbance":      0.10,
					"max_force":        2.2,
					"max_abs_position": 4.0,
					"settle_band":      0.08,
				},
			},
			{
				ID:    "msd_test_holdout",
				Split: designscenarios.Test,
				Params: map[string]any{
					"mass":             0.9,
					"spring_k":         2.6,
					"damping_c":        0.6,
					"dt":               0.03,
					"steps":            170,
					"target_position":  1.0,
					"initial_position": -0.4,
					"initial_velocity": 0.2,
					"disturbance":      0.25,
					"max_force":        2.3,
					"max_abs_position": 4.5,
					"settle_band":      0.08,
				},
			},
		},
	}
	registry := &designscenarios.ScenarioRegistry{Sets: []designscenarios.ScenarioSet{*set}}
	if err := registry.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validate control fixture scenarios: %w", err)
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
		return nil, errors.New("control fixture scenarios contain no train split")
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
		results = append(results, objectives.Score(controlObjectiveDefs, map[string]float64{
			"tracking_score": score,
		}, rejected, 0))
	}
	return results, nil
}

func scoreScenarioAggregate(genome core.Genome[bool], scenario designscenarios.Scenario) (score float64, hardFailed bool) {
	plant, err := scenarioPlant(scenario)
	if err != nil {
		return 0, true
	}
	sim, err := simulateController(genome, plant)
	if err != nil {
		return 0, true
	}
	if sim.HardFailed {
		return 0, true
	}

	tracking := 1.0 / (1.0 + sim.MeanAbsError)
	terminal := 1.0 / (1.0 + math.Abs(sim.FinalError))
	effort := 1.0 / (1.0 + sim.MeanAbsControl)
	return 500*tracking + 350*terminal + 150*effort, false
}

func scenarioPlant(scenario designscenarios.Scenario) (plantScenario, error) {
	params := scenario.Params
	mass, ok := floatParam(params, "mass")
	if !ok || mass <= 0 {
		return plantScenario{}, errors.New("scenario mass must be > 0")
	}
	springK, ok := floatParam(params, "spring_k")
	if !ok || springK < 0 {
		return plantScenario{}, errors.New("scenario spring_k must be >= 0")
	}
	dampingC, ok := floatParam(params, "damping_c")
	if !ok || dampingC < 0 {
		return plantScenario{}, errors.New("scenario damping_c must be >= 0")
	}
	timeStep, ok := floatParam(params, "dt")
	if !ok || timeStep <= 0 {
		return plantScenario{}, errors.New("scenario dt must be > 0")
	}
	steps, ok := intParam(params, "steps")
	if !ok || steps <= 0 {
		return plantScenario{}, errors.New("scenario steps must be > 0")
	}
	targetPosition, ok := floatParam(params, "target_position")
	if !ok {
		return plantScenario{}, errors.New("scenario target_position is required")
	}
	initialPos, ok := floatParam(params, "initial_position")
	if !ok {
		return plantScenario{}, errors.New("scenario initial_position is required")
	}
	initialVel, ok := floatParam(params, "initial_velocity")
	if !ok {
		return plantScenario{}, errors.New("scenario initial_velocity is required")
	}
	disturbance, ok := floatParam(params, "disturbance")
	if !ok {
		return plantScenario{}, errors.New("scenario disturbance is required")
	}
	maxForce, ok := floatParam(params, "max_force")
	if !ok || maxForce <= 0 {
		return plantScenario{}, errors.New("scenario max_force must be > 0")
	}
	maxAbsPosition, ok := floatParam(params, "max_abs_position")
	if !ok || maxAbsPosition <= 0 {
		return plantScenario{}, errors.New("scenario max_abs_position must be > 0")
	}
	settleBand, ok := floatParam(params, "settle_band")
	if !ok || settleBand <= 0 {
		settleBand = 0.05
	}

	return plantScenario{
		Mass:           mass,
		SpringK:        springK,
		DampingC:       dampingC,
		TimeStep:       timeStep,
		Steps:          steps,
		TargetPosition: targetPosition,
		InitialPos:     initialPos,
		InitialVel:     initialVel,
		Disturbance:    disturbance,
		MaxForce:       maxForce,
		MaxAbsPosition: maxAbsPosition,
		SettleBand:     settleBand,
	}, nil
}

func simulateController(genome core.Genome[bool], plant plantScenario) (simulationResult, error) {
	if err := genome.Validate(); err != nil {
		return simulationResult{}, fmt.Errorf("simulate controller: %w", err)
	}

	position := plant.InitialPos
	velocity := plant.InitialVel
	var totalAbsError float64
	var totalAbsControl float64

	for step := 0; step < plant.Steps; step++ {
		errorSignal := plant.TargetPosition - position
		features := []bool{
			errorSignal >= 0,
			velocity >= 0,
			plant.Disturbance >= 0,
			math.Abs(errorSignal) <= plant.SettleBand,
		}
		controllerOn, err := genome.Eval(features)
		if err != nil {
			return simulationResult{}, fmt.Errorf("simulate controller step %d: %w", step, err)
		}
		control := -plant.MaxForce
		if controllerOn {
			control = plant.MaxForce
		}
		totalAbsError += math.Abs(errorSignal)
		totalAbsControl += math.Abs(control)

		acceleration := (control + plant.Disturbance - plant.SpringK*position - plant.DampingC*velocity) / plant.Mass
		velocity += acceleration * plant.TimeStep
		position += velocity * plant.TimeStep
		if math.Abs(position) > plant.MaxAbsPosition {
			return simulationResult{HardFailed: true, FinalPosition: position}, nil
		}
	}

	finalError := plant.TargetPosition - position
	return simulationResult{
		MeanAbsError:   totalAbsError / float64(plant.Steps),
		MeanAbsControl: totalAbsControl / float64(plant.Steps),
		FinalError:     finalError,
		FinalPosition:  position,
	}, nil
}

func decodeControllerPolicy(candidateID string, genome core.Genome[bool]) controllerPolicyArtifact {
	return controllerPolicyArtifact{
		CandidateID: candidateID,
		Domain:      "control",
		Experiment:  "mass_spring_damper",
		Karva:       genome.KarvaString(),
		Terminals: []string{
			"d0=position_error_nonnegative",
			"d1=velocity_nonnegative",
			"d2=disturbance_nonnegative",
			"d3=within_settle_band",
		},
	}
}

func exportArtifacts(policy controllerPolicyArtifact, outputDir string) ([]design.ArtifactRef, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("create output directory %q: %w", outputDir, err)
	}

	jsonPath := filepath.Join(outputDir, controllerJSONArtifactName)
	if err := writeJSONFile(jsonPath, policy); err != nil {
		return nil, fmt.Errorf("write controller policy artifact: %w", err)
	}

	summaryText := fmt.Sprintf("candidate: %s\nkarva: %s\ntrain_mean: %.3f\nvalidation_mean: %.3f\ntest_mean: %.3f\n", policy.CandidateID, policy.Karva, policy.Scores["train_mean"], policy.Scores["validation_mean"], policy.Scores["test_mean"])
	summaryPath := filepath.Join(outputDir, controllerSummaryArtifact)
	if err := os.WriteFile(summaryPath, []byte(summaryText), 0o644); err != nil {
		return nil, fmt.Errorf("write controller summary artifact: %w", err)
	}

	return []design.ArtifactRef{
		{Name: "promoted_controller_policy", Path: jsonPath, Format: "json", Kind: "control.policy"},
		{Name: "promoted_controller_summary", Path: summaryPath, Format: "txt", Kind: "control.summary"},
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

func floatParam(params map[string]any, key string) (float64, bool) {
	raw, ok := params[key]
	if !ok {
		return 0, false
	}
	switch v := raw.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}

func intParam(params map[string]any, key string) (int, bool) {
	raw, ok := params[key]
	if !ok {
		return 0, false
	}
	switch v := raw.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}
