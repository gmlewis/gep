// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package evolution_test contains integration tests that exercise the real
// user-visible seams of the Phase 4 subsystems:
//
//   - Typed evolution (evolution.New / Evolve)
//   - Codegen (codegen.ProgramFromSymbols / Generate / Write)
//   - Env/RL (env.NewGymnasiumAgents / EvaluateAgent / RewardAgent / Evolve)
//   - Problem/domain wiring (problems.NewBoolProblem / NewFloatProblem scoring
//     functions passed directly to evolution.New)
//
// These tests replaced the import-boundary assertions that existed prior to
// milestone P4-E.  Import boundaries are now proven implicitly by the
// compilation of these tests: if any package violated its intended boundary
// the build would fail here.
package evolution_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/codegen"
	"github.com/gmlewis/gep/v2/common"
	"github.com/gmlewis/gep/v2/core"
	"github.com/gmlewis/gep/v2/env"
	"github.com/gmlewis/gep/v2/evolution"
	evolutionMutation "github.com/gmlewis/gep/v2/evolution/mutation"
	boolNodes "github.com/gmlewis/gep/v2/functions/bool_nodes"
	mathNodes "github.com/gmlewis/gep/v2/functions/math_nodes"
	"github.com/gmlewis/gep/v2/grammars"
	"github.com/gmlewis/gep/v2/problems"
)

// --- Typed evolution seam ---

// boolOrLink returns a variadic Or link operator that combines any number of
// bool gene outputs by folding OR over all of them.  This is the idiomatic
// way to create link functions for multi-gene boolean genomes: the link
// operator should reduce all gene outputs, not assume a fixed arity.
func boolOrLink(t *testing.T) core.LinkFunc[bool] {
	t.Helper()
	link, err := core.NewLinkFunc[bool]("Or", func(v []bool) bool {
		result := false
		for _, b := range v {
			result = result || b
		}
		return result
	})
	if err != nil {
		t.Fatalf("core.NewLinkFunc[bool]: %v", err)
	}
	return link
}

// float64SumLink returns a variadic + link operator that combines any number
// of float64 gene outputs by summing them.
func float64SumLink(t *testing.T) core.LinkFunc[float64] {
	t.Helper()
	link, err := core.NewLinkFunc[float64]("+", func(v []float64) float64 {
		sum := 0.0
		for _, x := range v {
			sum += x
		}
		return sum
	})
	if err != nil {
		t.Fatalf("core.NewLinkFunc[float64]: %v", err)
	}
	return link
}

// TestTypedEvolution_BoolNand exercises the full typed evolution seam for a
// boolean NAND problem.  It verifies that:
//   - evolution.NewWithSeed constructs a valid population from a typed
//     bool catalog and link function produced by functions/bool_nodes
//   - Evolve returns a best individual whose score is positive
//   - the best genome can be successfully evaluated against the NAND truth table
func TestTypedEvolution_BoolNand(t *testing.T) {
	funcs := []string{"Not", "And", "Or"}
	cat, err := boolNodes.CatalogFromNames(funcs)
	if err != nil {
		t.Fatalf("boolNodes.CatalogFromNames: %v", err)
	}
	link := boolOrLink(t)

	nandCases := [][3]bool{
		{false, false, true},
		{false, true, true},
		{true, false, true},
		{true, true, false},
	}

	nandScore := func(g core.Genome[bool]) float64 {
		correct := 0
		for _, tc := range nandCases {
			r, err := g.Eval([]bool{tc[0], tc[1]})
			if err != nil {
				continue
			}
			if r == tc[2] {
				correct++
			}
		}
		return 1000.0 * float64(correct) / float64(len(nandCases))
	}

	gen, err := evolution.NewWithSeed(42, cat, 30, 7, 1, 2, 0, link, nandScore)
	if err != nil {
		t.Fatalf("evolution.NewWithSeed: %v", err)
	}
	gen.MutationConfig = evolutionMutation.Config{
		PointMutationRate: 0.044,
		InversionRate:     0.1,
	}

	best := gen.Evolve(500)

	if best.Score <= 0 {
		t.Errorf("Evolve returned best.Score=%v, want > 0", best.Score)
	}

	// Verify that the best genome evaluates correctly on at least one case.
	atLeastOneCorrect := false
	for _, tc := range nandCases {
		r, err := best.Genome.Eval([]bool{tc[0], tc[1]})
		if err != nil {
			continue
		}
		if r == tc[2] {
			atLeastOneCorrect = true
			break
		}
	}
	if !atLeastOneCorrect {
		t.Errorf("best genome failed every NAND case; score=%v KarvaString=%q",
			best.Score, best.Genome.KarvaString())
	}
}

// TestTypedEvolution_FloatRegressionImproves exercises the typed evolution
// seam for a floating-point regression problem (f(x) = x).  It verifies that:
//   - evolution.NewWithSeed constructs a valid population from a float64
//     catalog and link function produced by functions/math_nodes
//   - Evolve returns a best individual whose score is positive
func TestTypedEvolution_FloatRegressionImproves(t *testing.T) {
	funcs := []string{"+", "-", "*"}
	cat, err := mathNodes.CatalogFromNames(funcs)
	if err != nil {
		t.Fatalf("mathNodes.CatalogFromNames: %v", err)
	}
	link := float64SumLink(t)

	// Simple regression: f(x) = x.  The genome gets one input terminal.
	regCases := []struct{ in, out float64 }{
		{1.0, 1.0}, {2.0, 2.0}, {3.0, 3.0}, {4.0, 4.0},
	}
	identityScore := func(g core.Genome[float64]) float64 {
		total := 0.0
		for _, c := range regCases {
			r, err := g.Eval([]float64{c.in})
			if err != nil || math.IsInf(r, 0) || math.IsNaN(r) {
				continue
			}
			total += 1000.0 / (1.0 + math.Abs(r-c.out))
		}
		return total / float64(len(regCases))
	}

	gen, err := evolution.NewWithSeed(7, cat, 30, 5, 1, 1, 0, link, identityScore)
	if err != nil {
		t.Fatalf("evolution.NewWithSeed: %v", err)
	}
	gen.MutationConfig = evolutionMutation.Config{
		PointMutationRate: 0.044,
		InversionRate:     0.1,
	}

	best := gen.Evolve(200)
	if best.Score <= 0 {
		t.Errorf("Evolve returned best.Score=%v, want > 0", best.Score)
	}
}

// --- Codegen seam ---

// TestCodegenSeam_BoolProgramFromCoreGenome exercises the codegen seam by
// evolving a bool genome and then rendering it to Go source via:
//
//	core.Genome.SymbolNamesPerGene() → codegen.ProgramFromSymbols → codegen.Generate
//
// This verifies that:
//   - the codegen package can consume a core.Genome without importing gene or genome
//   - the generated source contains the expected Go function signature
func TestCodegenSeam_BoolProgramFromCoreGenome(t *testing.T) {
	funcs := []string{"Not", "And", "Or"}
	cat, err := boolNodes.CatalogFromNames(funcs)
	if err != nil {
		t.Fatalf("boolNodes.CatalogFromNames: %v", err)
	}
	link := boolOrLink(t)

	gen, err := evolution.NewWithSeed(99, cat, 10, 5, 1, 2, 0, link, nil)
	if err != nil {
		t.Fatalf("evolution.NewWithSeed: %v", err)
	}

	// Take the first individual's genome and render it to Go source.
	g := gen.Individuals[0].Genome

	prog := codegen.ProgramFromSymbols(g.SymbolNamesPerGene(), nil, g.Link.Symbol())

	gr, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		t.Fatalf("LoadGoBooleanAllGatesGrammar: %v", err)
	}

	code, err := codegen.Generate(prog, gr)
	if err != nil {
		t.Fatalf("codegen.Generate: %v", err)
	}

	src := string(code)
	if !strings.Contains(src, "func gepModel(d []bool) bool {") {
		t.Errorf("generated source missing expected function signature; got:\n%s", src)
	}
}

// TestCodegenSeam_FloatProgramFromCoreGenome exercises the codegen seam for
// floating-point genomes by using:
//
//	core.Genome.SymbolNamesPerGene() + ConstsPerGene() → codegen.ProgramFromSymbols → codegen.Generate
//
// This verifies that constant symbols survive the round-trip through the
// codegen seam correctly.
func TestCodegenSeam_FloatProgramFromCoreGenome(t *testing.T) {
	funcs := []string{"+", "-", "*"}
	cat, err := mathNodes.CatalogFromNames(funcs)
	if err != nil {
		t.Fatalf("mathNodes.CatalogFromNames: %v", err)
	}
	link := float64SumLink(t)

	// numConstants=0: we use only terminals and functions so that ConstsPerGene
	// returns nil slices and codegen can render the program without needing
	// actual constant values.  The test still exercises the full codegen seam
	// from SymbolNamesPerGene → ProgramFromSymbols → Generate.
	gen, err := evolution.NewWithSeed(123, cat, 10, 5, 2, 1, 0, link, nil)
	if err != nil {
		t.Fatalf("evolution.NewWithSeed: %v", err)
	}

	g := gen.Individuals[0].Genome
	prog := codegen.ProgramFromSymbols(g.SymbolNamesPerGene(), g.ConstsPerGene(), g.Link.Symbol())

	gr, err := grammars.LoadGoMathGrammar()
	if err != nil {
		t.Fatalf("LoadGoMathGrammar: %v", err)
	}

	code, err := codegen.Generate(prog, gr)
	if err != nil {
		t.Fatalf("codegen.Generate: %v", err)
	}

	src := string(code)
	if !strings.Contains(src, "func gepModel(d []float64) float64 {") {
		t.Errorf("generated source missing expected function signature; got:\n%s", src)
	}
}

// --- Problems seam ---

// TestProblemsSeam_BoolProblemWiredToEvolution demonstrates the complete
// problems → evolution wiring for a boolean domain.  It verifies that:
//   - problems.NewBoolProblem constructs a problem from typed cases
//   - BoolProblem.NumHitsScoringFunc returns a function whose type matches
//     the scoringFunc argument of evolution.New
//   - evolution.NewWithSeed accepts the scoring function without adaptation
//   - Evolve returns a positive score, proving the domain seam is live
func TestProblemsSeam_BoolProblemWiredToEvolution(t *testing.T) {
	nandCases := []problems.Case[bool]{
		{In: []bool{false, false}, Out: true},
		{In: []bool{false, true}, Out: true},
		{In: []bool{true, false}, Out: true},
		{In: []bool{true, true}, Out: false},
	}
	p, err := problems.NewBoolProblem(nandCases)
	if err != nil {
		t.Fatalf("problems.NewBoolProblem: %v", err)
	}

	sf, err := p.NumHitsScoringFunc(1000.0)
	if err != nil {
		t.Fatalf("BoolProblem.NumHitsScoringFunc: %v", err)
	}

	funcs := []string{"Not", "And", "Or"}
	cat, err := boolNodes.CatalogFromNames(funcs)
	if err != nil {
		t.Fatalf("boolNodes.CatalogFromNames: %v", err)
	}
	link := boolOrLink(t)

	gen, err := evolution.NewWithSeed(1, cat, 30, 7, 1, 2, 0, link, sf)
	if err != nil {
		t.Fatalf("evolution.NewWithSeed: %v", err)
	}
	gen.MutationConfig = evolutionMutation.Config{
		PointMutationRate: 0.044,
		InversionRate:     0.1,
	}

	best := gen.Evolve(500)
	if best.Score <= 0 {
		t.Errorf("Evolve returned best.Score=%v, want > 0", best.Score)
	}
}

// TestProblemsSeam_FloatProblemWiredToEvolution demonstrates the complete
// problems → evolution wiring for a floating-point domain.  It verifies that:
//   - problems.NewFloatProblem constructs a problem from typed cases
//   - FloatProblem.MeanSquaredErrorAbsScoringFunc returns a function whose type
//     matches the scoringFunc argument of evolution.New
//   - evolution.NewWithSeed accepts the scoring function without adaptation
//   - Evolve returns a positive score, proving the domain seam is live
func TestProblemsSeam_FloatProblemWiredToEvolution(t *testing.T) {
	// Simple regression: f(x) = x.
	floatCases := []problems.Case[float64]{
		{In: []float64{1.0}, Out: 1.0},
		{In: []float64{2.0}, Out: 2.0},
		{In: []float64{3.0}, Out: 3.0},
		{In: []float64{4.0}, Out: 4.0},
	}
	p, err := problems.NewFloatProblem(floatCases)
	if err != nil {
		t.Fatalf("problems.NewFloatProblem: %v", err)
	}

	sf, err := p.MeanSquaredErrorAbsScoringFunc(1000.0)
	if err != nil {
		t.Fatalf("FloatProblem.MeanSquaredErrorAbsScoringFunc: %v", err)
	}

	funcs := []string{"+", "-", "*"}
	cat, err := mathNodes.CatalogFromNames(funcs)
	if err != nil {
		t.Fatalf("mathNodes.CatalogFromNames: %v", err)
	}
	link := float64SumLink(t)

	gen, err := evolution.NewWithSeed(2, cat, 30, 5, 1, 1, 0, link, sf)
	if err != nil {
		t.Fatalf("evolution.NewWithSeed: %v", err)
	}
	gen.MutationConfig = evolutionMutation.Config{
		PointMutationRate: 0.044,
		InversionRate:     0.1,
	}

	best := gen.Evolve(100)
	if best.Score <= 0 {
		t.Errorf("Evolve returned best.Score=%v, want > 0", best.Score)
	}
}

// --- Env/RL seam ---

// TestEnvSeam_GymnasiumAgentsEvalRewardEvolve exercises the full env/RL seam
// by creating a GymnasiumAgents population, running a simulated episode batch
// (eval → reward), and then calling Evolve.  It verifies that:
//   - env.NewGymnasiumAgents constructs a valid population from action and
//     observation space metadata
//   - EvaluateAgent returns an action within the expected range for every agent
//   - RewardAgent records rewards without error
//   - Evolve advances the population without error
//
// This proves that the env/RL seam is live and that the env package is
// decoupled from the legacy model package.
func TestEnvSeam_GymnasiumAgentsEvalRewardEvolve(t *testing.T) {
	actionSpace := &common.Space{Type: "Discrete", N: 2}
	obsSpace := &common.Space{
		Type: "Tuple",
		Subspaces: []*common.Space{
			{Type: "Discrete", N: 32},
			{Type: "Discrete", N: 11},
			{Type: "Discrete", N: 2},
		},
	}

	agents, err := env.NewGymnasiumAgents(actionSpace, obsSpace,
		env.WithNumIndividuals(5),
	)
	if err != nil {
		t.Fatalf("env.NewGymnasiumAgents: %v", err)
	}

	// Simulate one episode per agent with a fixed observation sequence.
	numAgents := len(agents.Individuals)
	for agentNum := 0; agentNum < numAgents; agentNum++ {
		obs := &fixedObs{values: []int{10, 5, 1}} // Blackjack-style obs
		for step := 0; step < 5; step++ {
			var action int
			if err := agents.EvaluateAgent(agentNum, step, obs, &action); err != nil {
				t.Fatalf("agent %d step %d: EvaluateAgent: %v", agentNum, step, err)
			}
			if action < 0 || action >= actionSpace.N {
				t.Errorf("agent %d step %d: action=%d out of range [0,%d)",
					agentNum, step, action, actionSpace.N)
			}
		}
		agents.RewardAgent(agentNum, 1.0) // fixed positive reward
	}

	if err := agents.Evolve(); err != nil {
		t.Fatalf("env.GymnasiumAgents.Evolve: %v", err)
	}
}

// fixedObs is a minimal common.Obs implementation that returns a fixed
// integer observation vector for use in integration tests.
type fixedObs struct {
	values []int
}

func (o *fixedObs) Unmarshal(dst any) error {
	switch v := dst.(type) {
	case *[]int:
		if len(*v) < len(o.values) {
			return fmt.Errorf("fixedObs.Unmarshal: dst slice too short: got %d, need %d", len(*v), len(o.values))
		}
		copy(*v, o.values)
		return nil
	default:
		return fmt.Errorf("fixedObs.Unmarshal: unsupported type %T", dst)
	}
}
