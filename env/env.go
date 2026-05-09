// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package env

import (
	"cmp"
	"errors"
	"fmt"
	"math/rand"
	"sort"

	"github.com/gmlewis/gep/v2/common"
	"github.com/gmlewis/gep/v2/functions"
	bn "github.com/gmlewis/gep/v2/functions/bool_nodes"
	in "github.com/gmlewis/gep/v2/functions/int_nodes"
	mn "github.com/gmlewis/gep/v2/functions/math_nodes"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/genome"
)

const (
	defaultHeadSize       = 10
	defaultNumConstants   = 1
	defaultNumIndividuals = 10
)

// GymnasiumAgents implements the Agent interface for Gymnasium.
// See: https://github.com/Farama-Foundation/Gymnasium
//
// A GymnasiumAgents population runs one individual per concurrent episode.
// After each episode batch, agents are scored via [RewardAgent] and then
// evolved by [Evolve] so that the population improves over time.
//
// Use [NewGymnasiumAgents] to create a population from action- and
// observation-space metadata.
type GymnasiumAgents struct {
	// ActionSpace describes the action space of the target Gymnasium environment.
	ActionSpace *common.Space
	// ObsSpace describes the observation space of the target Gymnasium environment.
	ObsSpace *common.Space
	// Individuals holds the current population of GEP genomes, one per agent slot.
	Individuals []*genome.Genome

	// options
	appendEpisodeSteps bool
	debug              bool
	headSize           int
	numConstants       int
	numIndividuals     int
}

// GymnasiumAgentsOption is a functional option that configures a [GymnasiumAgents]
// population.  Pass one or more options to [NewGymnasiumAgents].
type GymnasiumAgentsOption func(ga *GymnasiumAgents)

// WithAppendEpisodeSteps returns an option that appends the current episode-step
// count to each observation vector before it is passed to the agent.  This is
// useful when the environment's optimal policy depends on elapsed time.
func WithAppendEpisodeSteps() GymnasiumAgentsOption {
	return func(ga *GymnasiumAgents) {
		ga.appendEpisodeSteps = true
	}
}

// WithDebug returns an option that enables verbose debug logging during agent
// evaluation and evolution.
func WithDebug() GymnasiumAgentsOption {
	return func(ga *GymnasiumAgents) {
		ga.debug = true
	}
}

// WithHeadSize returns an option that overrides the head size of each GEP
// gene.  Larger head sizes allow more complex decision programs at the cost of
// a larger search space.
func WithHeadSize(headSize int) GymnasiumAgentsOption {
	return func(ga *GymnasiumAgents) {
		ga.headSize = headSize
	}
}

// WithNumConstants returns an option that sets the number of numerical
// constants available inside each gene.
func WithNumConstants(numConstants int) GymnasiumAgentsOption {
	return func(ga *GymnasiumAgents) {
		ga.numConstants = numConstants
	}
}

// WithNumIndividuals returns an option that sets the maximum number of genomes
// (agents) in the population.
func WithNumIndividuals(numIndividuals int) GymnasiumAgentsOption {
	return func(ga *GymnasiumAgents) {
		ga.numIndividuals = numIndividuals
	}
}

// NewGymnasiumAgents creates a randomly initialised population of GEP agents
// from the given action and observation spaces.  Use [GymnasiumAgentsOption]
// values to customise head size, population size, and other knobs.
func NewGymnasiumAgents(actionSpace, obsSpace *common.Space, opts ...GymnasiumAgentsOption) (*GymnasiumAgents, error) {
	ga := &GymnasiumAgents{
		ActionSpace:    actionSpace,
		ObsSpace:       obsSpace,
		headSize:       defaultHeadSize,
		numConstants:   defaultNumConstants,
		numIndividuals: defaultNumIndividuals,
	}

	for _, f := range opts {
		f(ga)
	}

	var err error
	ga.Individuals, err = ga.newIndividuals()
	if err != nil {
		return nil, err
	}

	return ga, nil
}

// newIndividuals builds the initial random population based on the configured
// action and observation spaces.
func (ga *GymnasiumAgents) newIndividuals() ([]*genome.Genome, error) {
	numGenes := 1
	switch ga.ActionSpace.Type {
	case "Discrete":
		// N=2 means the output can take two values: 0 or 1.
	case "Tuple":
		numGenes = len(ga.ActionSpace.Subspaces)
	default:
		return nil, fmt.Errorf("ActionSpace type %v not yet implemented", ga.ActionSpace.Type)
	}

	switch ga.ObsSpace.Type {
	case "Discrete":
		funcWeights, err := gene.AllSymbolsEqualWeights(functions.Int)
		if err != nil {
			return nil, err
		}
		numTerminals := 1
		if ga.appendEpisodeSteps {
			numTerminals++
		}
		return newIndividuals(
			funcWeights,
			functions.Int,
			ga.numIndividuals,
			ga.headSize,
			numGenes,
			numTerminals,
			ga.numConstants,
			"tuple",
		)

	case "Tuple":
		funcWeights, err := gene.AllSymbolsEqualWeights(functions.Int)
		if err != nil {
			return nil, err
		}
		numTerminals := len(ga.ObsSpace.Subspaces)
		if ga.appendEpisodeSteps {
			numTerminals++
		}
		return newIndividuals(
			funcWeights,
			functions.Int,
			ga.numIndividuals,
			ga.headSize,
			numGenes,
			numTerminals,
			ga.numConstants,
			"tuple",
		)

	default:
		return nil, fmt.Errorf("ObservationSpace type %v not yet implemented", ga.ObsSpace.Type)
	}
}

// newIndividuals creates a slice of randomly initialised genomes.
func newIndividuals(
	fs []gene.FuncWeight,
	funcType functions.FuncType,
	numIndividuals,
	headSize,
	numGenesPerGenome,
	numTerminals,
	numConstants int,
	linkFunc string,
) ([]*genome.Genome, error) {
	n, err := maxArity(fs, funcType)
	if err != nil {
		return nil, err
	}
	tailSize := headSize*(n-1) + 1
	individuals := make([]*genome.Genome, numIndividuals)
	for i := range individuals {
		genes := make([]*gene.Gene, numGenesPerGenome)
		for j := range genes {
			genes[j] = gene.RandomNew(headSize, tailSize, numTerminals, numConstants, fs, funcType, nil)
		}
		individuals[i] = genome.New(genes, linkFunc, nil)
	}
	return individuals, nil
}

// maxArity determines the maximum arity (number of input terminals) among the
// symbols in fs for the given funcType.
func maxArity(fs []gene.FuncWeight, funcType functions.FuncType) (int, error) {
	var lookup functions.FuncMap
	switch funcType {
	case functions.Bool:
		lookup = bn.BoolAllGates
	case functions.Int:
		lookup = in.Int
	case functions.Float64:
		lookup = mn.Math
	default:
		return 0, fmt.Errorf("unknown funcType: %v", funcType)
	}

	r := 0
	var errs []error
	for _, f := range fs {
		if fn, ok := lookup[f.Symbol]; ok {
			if fn.Terminals() > r {
				r = fn.Terminals()
			}
		} else {
			errs = append(errs, fmt.Errorf("unable to find symbol %v in function map", f.Symbol))
		}
	}
	return r, errors.Join(errs...)
}

// EvaluateAgent runs the GEP genome for agent at agentIdx against the given
// observation and populates action with the result.
//
// During an episode only Individual #agentIdx is evaluated.  The action is
// clamped to the valid range of the corresponding action space.
func (ga *GymnasiumAgents) EvaluateAgent(agentIdx, episodeSteps int, obs common.Obs, action any) error {
	observations, err := ga.processObservations(episodeSteps, obs)
	if err != nil {
		return err
	}

	if err := ga.Individuals[agentIdx].Evaluate(observations, action); err != nil {
		return err
	}

	switch v := action.(type) {
	case *[]int:
		for i, val := range *v {
			(*v)[i] = clamp(val, 0, ga.ActionSpace.Subspaces[i].N-1)
		}
	case *int:
		*v = clamp(*v, 0, ga.ActionSpace.N-1)
	default:
		return fmt.Errorf("agent.Evaluate: action type '%T' not yet supported", v)
	}
	return nil
}

// RewardAgent records a reward score for the agent at agentIdx after an
// episode.  reward can be any float64; the range −1000 ≤ reward ≤ 1000 works
// well in practice.
func (ga *GymnasiumAgents) RewardAgent(agentIdx int, reward float64) {
	ga.Individuals[agentIdx].Score = reward
}

// SortIndividuals sorts the population by descending Score (best performer
// first).
func (ga *GymnasiumAgents) SortIndividuals() {
	sort.Slice(ga.Individuals, func(i, j int) bool {
		return ga.Individuals[i].Score > ga.Individuals[j].Score
	})
}

// Evolve advances the population by one generation.  It sorts individuals by
// score, preserves an elitist copy of the best performer, applies mutation and
// crossover, and then re-inserts the elite copy in place of the weakest
// individual.  All scores are reset to zero after evolution so that the next
// round of episodes starts with a clean slate.
func (ga *GymnasiumAgents) Evolve() error {
	ga.SortIndividuals()

	// Preserve a copy of the best-performing individual so it survives
	// mutation and crossover (elitism).
	bestInd, err := ga.Individuals[0].Dup()
	if err != nil {
		return err
	}

	if err := mutation(ga.Individuals); err != nil {
		return err
	}
	if err := crossover(ga.Individuals); err != nil {
		return err
	}

	if len(ga.Individuals) == 0 {
		return errors.New("programming error: no individuals available after evolution operators")
	}
	// Overwrite the weakest individual with the elite copy.
	ga.Individuals[len(ga.Individuals)-1] = bestInd

	if len(ga.Individuals) != ga.numIndividuals {
		return fmt.Errorf("programming error: got %v individuals, want %v", len(ga.Individuals), ga.numIndividuals)
	}
	// Reset scores so the next round starts fresh.
	for _, individual := range ga.Individuals {
		individual.Score = 0
	}

	return nil
}

// processObservations converts an [common.Obs] into a flat []int that the GEP
// genome can evaluate.  When [WithAppendEpisodeSteps] is active, the current
// step count is appended as an extra input.
func (ga *GymnasiumAgents) processObservations(episodeSteps int, obs common.Obs) ([]int, error) {
	resultLen := len(ga.ObsSpace.Subspaces)
	if ga.appendEpisodeSteps {
		resultLen++
	}

	result := make([]int, resultLen)
	if err := obs.Unmarshal(&result); err != nil {
		return nil, err
	}

	if ga.appendEpisodeSteps {
		result[resultLen-1] = episodeSteps
	}

	return result, nil
}

// mutation applies point mutation to a random subset of the population.
func mutation(individuals []*genome.Genome) error {
	if len(individuals) == 0 {
		return nil
	}
	numMutations := 1 + rand.Intn(len(individuals))
	for range numMutations {
		genomeNum := rand.Intn(len(individuals))
		numPointMutations := 1 + rand.Intn(2)
		if err := individuals[genomeNum].Mutate(numPointMutations); err != nil {
			return err
		}
	}
	return nil
}

// crossover applies one-point crossover to random pairs of genomes in the
// population.
func crossover(individuals []*genome.Genome) error {
	if len(individuals) < 2 {
		return nil
	}
	numCrossovers := 1 + rand.Intn(len(individuals)-1)
	for range numCrossovers {
		genomeNum1 := rand.Intn(len(individuals))
		var genomeNum2 int
		for {
			genomeNum2 = rand.Intn(len(individuals))
			if genomeNum2 != genomeNum1 {
				break
			}
		}
		if err := singleCrossover(individuals, genomeNum1, genomeNum2); err != nil {
			return err
		}
	}
	return nil
}

// singleCrossover swaps a randomly-chosen head segment between two genomes.
func singleCrossover(individuals []*genome.Genome, idx1, idx2 int) error {
	genome1 := individuals[idx1]
	genome2 := individuals[idx2]

	geneIdx1 := rand.Intn(len(genome1.Genes))
	gene1 := genome1.Genes[geneIdx1]
	geneIdx2 := rand.Intn(len(genome2.Genes))
	gene2 := genome2.Genes[geneIdx2]

	if len(gene1.Symbols) != len(gene2.Symbols) || gene1.HeadSize != gene2.HeadSize {
		return fmt.Errorf("programming error: gene1: %v symbols (headSize=%v), gene2: %v symbols (headSize=%v)",
			len(gene1.Symbols), gene1.HeadSize, len(gene2.Symbols), gene2.HeadSize)
	}

	symbolIdx := rand.Intn(gene1.HeadSize)
	head1, tail1 := gene1.Symbols[:gene1.HeadSize], gene1.Symbols[gene1.HeadSize:]
	head2, tail2 := gene2.Symbols[:gene2.HeadSize], gene2.Symbols[gene2.HeadSize:]
	newSyms1 := append([]string{}, head2[symbolIdx:]...)
	newSyms1 = append(newSyms1, head1[:symbolIdx]...)
	newSyms1 = append(newSyms1, tail1...)
	newSyms2 := append([]string{}, head1[symbolIdx:]...)
	newSyms2 = append(newSyms2, head2[:symbolIdx]...)
	newSyms2 = append(newSyms2, tail2...)

	if len(newSyms1) != len(newSyms2) || len(newSyms1) != len(gene1.Symbols) {
		return fmt.Errorf("programming error: newSyms1: %v symbols, newSyms2: %v symbols, gene1: %v symbols",
			len(newSyms1), len(newSyms2), len(gene1.Symbols))
	}

	gene1.Symbols = newSyms1
	gene2.Symbols = newSyms2
	gene1.InvalidateCache()
	gene2.InvalidateCache()
	genome1.SymbolMap = nil
	genome2.SymbolMap = nil
	return nil
}

// clamp returns v clamped to [min, max].
func clamp[T cmp.Ordered](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
