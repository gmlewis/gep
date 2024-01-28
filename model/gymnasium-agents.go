// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model

import (
	"cmp"
	"fmt"
	"log"
	"sort"

	"github.com/gmlewis/gep/v2/common"
	"github.com/gmlewis/gep/v2/functions"
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
// A gymnasium trains only a single individual at a time during each episode.
// As more individuals are run and rewarded (at the end of their episode),
// new individuals are created by mixing the genetic makeup of the other
// top-performing individuals.
type GymnasiumAgents struct {
	ActionSpace *common.Space
	ObsSpace    *common.Space
	Individuals []*genome.Genome

	// options
	appendEpisodeSteps bool
	debug              bool
	headSize           int
	numConstants       int
	numIndividuals     int
}

// GymnasiumAgentsOption represents an option that can modify the GEP model.
type GymnasiumAgentsOption func(ga *GymnasiumAgents)

// WithAppendEpisodeSteps adds an option to append the current number of
// episode steps to the observations sent to the agent.
func WithAppendEpisodeSteps() GymnasiumAgentsOption {
	return func(ga *GymnasiumAgents) {
		ga.appendEpisodeSteps = true
	}
}

// WithDebug prints debug information during run.
func WithDebug() GymnasiumAgentsOption {
	return func(ga *GymnasiumAgents) {
		ga.debug = true
	}
}

// WithHeadSize adds an option to change the head size of the GEP model.
// This determines the maximum complexity of the generated algorithm.
func WithHeadSize(headSize int) GymnasiumAgentsOption {
	return func(ga *GymnasiumAgents) {
		ga.headSize = headSize
	}
}

// WithNumConstants adds an option to change the number of constants in the GEP model.
func WithNumConstants(numConstants int) GymnasiumAgentsOption {
	return func(ga *GymnasiumAgents) {
		ga.numConstants = numConstants
	}
}

// WithNumIndividuals adds an option to change the maximum number of genomes in the GEP model.
func WithNumIndividuals(numIndividuals int) GymnasiumAgentsOption {
	return func(ga *GymnasiumAgents) {
		ga.numIndividuals = numIndividuals
	}
}

// NewGymnasiumAgents returns an Agent based upon the action and observation spaces.
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

func (ga *GymnasiumAgents) newIndividuals() ([]*genome.Genome, error) {
	numGenes := 1
	switch ga.ActionSpace.Type {
	case "Discrete":
		// N=2 means that the output can have two values: 0, 1
	case "Tuple":
		numGenes = len(ga.ActionSpace.Subspaces)
	// case "MultiBinary":
	// case "MultiDiscrete":
	// case "Box":
	default:
		return nil, fmt.Errorf("ActionSpace type %v not yet implemented", ga.ActionSpace.Type)
	}

	// var genes []*gene.Gene
	switch ga.ObsSpace.Type {
	case "Discrete":
		funcWeights := gene.AllSymbolsEqualWeights(functions.Int)
		// for i := 0; i < numGenes; i++ {
		// 	genes = append(genes, gene.RandomNew(headSize, tailSize, 1, numConstants, funcs, functions.Int))
		// }
		// o.Individuals = append(o.Individuals, genome.New(genes, "tuple"))
		numTerminals := 1
		if ga.appendEpisodeSteps {
			numTerminals++
		}
		gen := New(
			funcWeights,
			functions.Int,
			ga.numIndividuals,
			ga.headSize,
			numGenes,
			numTerminals,
			ga.numConstants,
			"tuple",
			nil,
			ga.debug)
		return gen.Individuals, nil

	case "Tuple":
		funcType := functions.Int
		funcWeights := gene.AllSymbolsEqualWeights(functions.Int)
		numTerminals := len(ga.ObsSpace.Subspaces)
		if ga.appendEpisodeSteps {
			numTerminals++
		}
		gen := New(
			funcWeights,
			funcType,
			ga.numIndividuals,
			ga.headSize,
			numGenes,
			numTerminals,
			ga.numConstants,
			"tuple",
			nil,
			ga.debug)
		return gen.Individuals, nil

	// case "MultiBinary":
	// case "MultiDiscrete":
	// case "Box":
	default:
		return nil, fmt.Errorf("ObservationSpace type %v not yet implemented", ga.ObsSpace.Type)
	}
}

// EvaluateAgent runs the GEP model for a single individual
// and returns an action from an observation
// by populating the passed-in reference.
//
// During an episode, only Individual #agentIdx is evaluated.
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
		before := *v
		*v = clamp(*v, 0, ga.ActionSpace.N-1)
		if ga.debug {
			log.Printf("EvaluateAgent(agentIdx=%v, obs=%+v)=%v => clamp(0,%v) => %v", agentIdx, observations, before, ga.ActionSpace.N-1, *v)
		}
	default:
		return fmt.Errorf("agent.Evaluate: action type '%T' not yet supported", v)
	}
	return nil
}

// RewardAgent rewards a single agent after evaluation.
// reward can be any float64 but the range -1000 <= reward <= 1000 works nicely.
func (ga *GymnasiumAgents) RewardAgent(agentIdx int, reward float64) {
	ga.Individuals[agentIdx].Score = reward
}

// SortIndividuals sorts all agents based upon their Score.
func (ga *GymnasiumAgents) SortIndividuals() {
	// Sort individuals by best performance (best first)
	sort.Slice(ga.Individuals, func(i, j int) bool {
		return ga.Individuals[i].Score > ga.Individuals[j].Score
	})
}

// Evolve evolves the GEP model based on all individual scores
// (ranging from negative (bad) to positive (good)).
func (ga *GymnasiumAgents) Evolve() error {
	ga.SortIndividuals()

	// Preserve a copy of the best performing individual which
	// will be added back into the population after replication,
	// mutation, and crossover.
	bestInd := ga.Individuals[0].Dup()
	gen := &Generation{Individuals: ga.Individuals, debug: ga.debug}
	// gen.replication()  // This seems to eliminate all diversity - investigate
	gen.mutation()
	gen.crossover()
	gen.Individuals[ga.numIndividuals-1] = bestInd // Overwrite lowest performer
	ga.Individuals = gen.Individuals

	if len(ga.Individuals) != ga.numIndividuals {
		log.Fatalf("programming error: got %v individuals, want %v", len(ga.Individuals), ga.numIndividuals)
	}
	// Reset all scores to zero
	for _, individual := range ga.Individuals {
		individual.Score = 0
	}

	return nil
}

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

func clamp[T cmp.Ordered](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
