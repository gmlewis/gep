// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model

import (
	"cmp"
	"fmt"
	"sort"

	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/genome"
	gym "github.com/gmlewis/gym-socket-api/binding-go"
)

const (
	defaultHeadSize       = 10
	defaultNumConstants   = 1
	defaultMaxIndividuals = 100
)

// Agent represents the methods provided by a GEP model.
type Agent interface {
	Evaluate(episodeSteps int, obs gym.Obs, action any) error
	Evolve(reward float64) error
}

// GymnasiumAgent implements the Agent interface for Gymnasium.
// See: https://github.com/Farama-Foundation/Gymnasium
// A gymnasium trains only a single individual at a time during each episode.
// As more individuals are run and rewarded (at the end of their episode),
// new individuals are created by mixing the genetic makeup of the other
// top-performing individuals.
type GymnasiumAgent struct {
	ActionSpace *gym.Space
	ObsSpace    *gym.Space
	Individuals []*genome.Genome

	// options
	appendEpisodeSteps bool
	headSize           int
	numConstants       int
	maxIndividuals     int
}

var _ Agent = &GymnasiumAgent{}

// GymnasiumAgentOption represents an option that can modify the GEP model.
type GymnasiumAgentOption func(ga *GymnasiumAgent)

// WithAppendEpisodeSteps adds an option to append the current number of
// episode steps to the observations sent to the agent.
func WithAppendEpisodeSteps() GymnasiumAgentOption {
	return func(ga *GymnasiumAgent) {
		ga.appendEpisodeSteps = true
	}
}

// WithHeadSize adds an option to change the head size of the GEP model.
// This determines the maximum complexity of the generated algorithm.
func WithHeadSize(headSize int) GymnasiumAgentOption {
	return func(ga *GymnasiumAgent) {
		ga.headSize = headSize
	}
}

// WithNumConstants adds an option to change the number of constants in the GEP model.
func WithNumConstants(numConstants int) GymnasiumAgentOption {
	return func(ga *GymnasiumAgent) {
		ga.numConstants = numConstants
	}
}

// WithMaxIndividuals adds an option to change the maximum number of genomes in the GEP model.
func WithMaxIndividuals(maxIndividuals int) GymnasiumAgentOption {
	return func(ga *GymnasiumAgent) {
		ga.maxIndividuals = maxIndividuals
	}
}

// NewGymnasiumAgent returns an Agent based upon the action and observation spaces.
func NewGymnasiumAgent(actionSpace, obsSpace *gym.Space, opts ...GymnasiumAgentOption) (*GymnasiumAgent, error) {
	ga := &GymnasiumAgent{
		ActionSpace:    actionSpace,
		ObsSpace:       obsSpace,
		headSize:       defaultHeadSize,
		numConstants:   defaultNumConstants,
		maxIndividuals: defaultMaxIndividuals,
	}

	for _, f := range opts {
		f(ga)
	}

	individual, err := ga.newIndividual()
	if err != nil {
		return nil, err
	}
	ga.Individuals = append(ga.Individuals, individual)
	return ga, nil
}

func (ga *GymnasiumAgent) newIndividual() (*genome.Genome, error) {
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
		funcs := []gene.FuncWeight{
			{Symbol: "+", Weight: 1},
			{Symbol: "-", Weight: 5},
			{Symbol: "*", Weight: 5},
		}
		// for i := 0; i < numGenes; i++ {
		// 	genes = append(genes, gene.RandomNew(headSize, tailSize, 1, numConstants, funcs, functions.Int))
		// }
		// o.Individuals = append(o.Individuals, genome.New(genes, "tuple"))
		numTerminals := 1
		if ga.appendEpisodeSteps {
			numTerminals++
		}
		gen := New(
			funcs,
			functions.Int,
			1, // always start with one individual in the gymnasium for the first episode
			ga.headSize,
			numGenes,
			numTerminals,
			ga.numConstants,
			"tuple",
			nil)
		return gen.Individuals[0], nil

	case "Tuple":
		funcType := functions.Int
		funcs := []gene.FuncWeight{
			{Symbol: "+", Weight: 1},
			{Symbol: "-", Weight: 5},
			{Symbol: "*", Weight: 5},
		}
		numTerminals := len(ga.ObsSpace.Subspaces)
		if ga.appendEpisodeSteps {
			numTerminals++
		}
		gen := New(
			funcs,
			funcType,
			1, // always start with one individual in the gymnasium for the first episode
			ga.headSize,
			numGenes,
			numTerminals,
			ga.numConstants,
			"tuple",
			nil)
		return gen.Individuals[0], nil

	// case "MultiBinary":
	// case "MultiDiscrete":
	// case "Box":
	default:
		return nil, fmt.Errorf("ObservationSpace type %v not yet implemented", ga.ObsSpace.Type)
	}
}

// Evaluate runs the GEP model and returns an action from an observation
// by populating the passed-in reference.
//
// During an episode, only Individual #0 is evaluated and rewarded.
// Once rewarded, in [Evolve], the best-performing individuals are used to create
// a new individual used during the next episode (and therefore at position 0).
func (ga *GymnasiumAgent) Evaluate(episodeSteps int, obs gym.Obs, action any) error {
	observations, err := ga.processObservations(episodeSteps, obs)
	if err != nil {
		return err
	}

	if err := ga.Individuals[0].Evaluate(observations, action); err != nil {
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

	/*
		// log.Printf("Evaluate: obs=%v, raw action=%v", obs, action)

		switch ga.ActionSpace.Type {
		// case "Discrete":
		case "Tuple":
			if vals, ok := action.(*[]int); ok {
				for i, v := range *vals {
					if ga.ActionSpace.Subspaces[i].Type == "Discrete" {
						if v < 0 {
							v = -v
						}
						(*vals)[i] = v % ga.ActionSpace.Subspaces[i].N
					} else {
						return fmt.Errorf("subspace %v not yet supported", ga.ActionSpace.Subspaces[i])
					}
				}
			} else {
				return fmt.Errorf("action type %T not yet supported", action)
			}
		// case "MultiBinary":
		// case "MultiDiscrete":
		// case "Box":
		default:
			return fmt.Errorf("ActionSpace type %v not yet implemented", ga.ActionSpace.Type)
		}

		// log.Printf("Evaluate: obs=%v, final action=%v", obs, action)

		return nil
	*/
}

// Evolve evolves the GEP model based on the reward feedback
// (ranging from negative (bad) to positive (good)) for the current (#0)
// individual.
func (ga *GymnasiumAgent) Evolve(reward float64) error {
	ga.Individuals[0].Score = reward // should this be += ?

	if len(ga.Individuals) < ga.maxIndividuals {
		// Create a completely new random individual
		individual, err := ga.newIndividual()
		if err != nil {
			return err
		}
		ga.Individuals = append([]*genome.Genome{individual}, ga.Individuals...)
		return nil
	}

	// Sort individuals by best performance
	sort.Slice(ga.Individuals, func(i, j int) bool {
		return ga.Individuals[i].Score > ga.Individuals[j].Score
	})

	// Preserve a copy of the best performing individual which
	// will be added back into the population after replication,
	// mutation, and crossover.
	ng := ga.Individuals[0].Dup()
	gen := &Generation{Individuals: ga.Individuals}
	// gen.replication()
	// Only perform mutation and crossover on the new individual #0
	// which will be evaluted during the next upcoming episode.
	gen.singleMutation(0)
	// gen.crossover() // TODO
	ga.Individuals = append(gen.Individuals[0:ga.maxIndividuals-1], ng)

	// for i, g := range ga.Individuals {
	// 	log.Printf("genome[%v]=%v", i, g)
	// 	for in := 0; in < ga.ObsSpace.N; in++ {
	// 		action := []int{}
	// 		intObs := &intObsT{i: in}
	// 		if err := g.Evaluate(intObs, &action); err != nil {
	// 			log.Fatalf("g.Evaluate: %v", err)
	// 		}
	// 		log.Printf("genome[%v](%v) = %v", i, in, action)
	// 	}
	// }
	// log.Fatalf("Debug stop.")

	// if len(ga.Individuals) > ga.maxIndividuals {
	// 	ga.Individuals = ga.Individuals[0:ga.maxIndividuals]
	// }

	return nil
}

func (ga *GymnasiumAgent) processObservations(episodeSteps int, obs gym.Obs) ([]int, error) {
	var result []int
	if err := obs.Unmarshal(&result); err != nil {
		return nil, err
	}

	if ga.appendEpisodeSteps {
		result = append(result, episodeSteps)
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
