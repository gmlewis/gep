// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"log"
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
	Evaluate(episodeSteps int, obs gym.Obs, action interface{}) error
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
	headSize       int
	numConstants   int
	maxIndividuals int
}

var _ Agent = &GymnasiumAgent{}

// GymnasiumAgentOption represents an option that can modify the GEP model.
type GymnasiumAgentOption func(ga *GymnasiumAgent)

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
		log.Printf("Warning: ActionSpace Discrete not handled yet: N=%v", ga.ActionSpace.N)
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
		const numTerminals = 2 // Account for added "episodeSteps"
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
		numTerminals := len(ga.ObsSpace.Subspaces) + 1 // Account for added "episodeSteps"
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

// Evaluate runs the GEP model and returns an action from an observation.
// During an episode, only Individual[0] is evaluated and rewarded.
// Once rewarded, the best-performing individuals are used to create
// a new individual used during the next episode (and therefore at position 0).
func (ga *GymnasiumAgent) Evaluate(episodeSteps int, obs gym.Obs, action interface{}) error {
	numObs := len(ga.ObsSpace.Subspaces)
	if ga.AddEpisodeSteps {
		// Add episodeSteps to the end of the observations.
		numObs++
	}
	observations := make([]int, 0, numObs)

	if err := ga.Individuals[0].Evaluate(episodeSteps, obs, action); err != nil {
		return err
	}
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
}

// Evolve evolves the GEP model based on the reward feedback
// (from negative (bad) to positive (good)) for the current
// individual (Individuals[0]).
func (ga *GymnasiumAgent) Evolve(reward float64) error {
	ga.Individuals[0].Score = reward

	if len(ga.Individuals) < 2 {
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
	gen.replication()
	gen.mutation()
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

func check(fmt string, args ...interface{}) {
	err := args[len(args)-1]
	if err != nil {
		log.Fatalf(fmt, args...)
	}
}

// For debugging...
/*
type intObsT struct {
	i int
}

func (i *intObsT) Unmarshal(dst interface{}) error {
	if v, ok := dst.(*int); ok {
		*v = i.i
	}
	return nil
}
*/
