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
	headSize     = 10
	numConstants = 1
	maxGenomes   = 10
)

// Model represents a GEP model.
type Model interface {
	Evaluate(obs gym.Obs, action interface{}) error
	Evolve(reward float64) error
}

// Gymnasium implements the Model interface for Gymnasium.
// See: https://github.com/Farama-Foundation/Gymnasium
type Gymnasium struct {
	ActionSpace *gym.Space
	ObsSpace    *gym.Space
	Genomes     []*genome.Genome
}

// NewGymnasium returns a Model based upon the action and observation spaces.
func NewGymnasium(actionSpace, obsSpace *gym.Space) (*Gymnasium, error) {
	o := &Gymnasium{ActionSpace: actionSpace, ObsSpace: obsSpace}

	numGenes := 1
	switch actionSpace.Type {
	case "Discrete":
		log.Printf("Warning: ActionSpace Discrete not handled yet: N=%v", actionSpace.N)
		// N=2 means that the output can have two values: 0, 1
	case "Tuple":
		numGenes = len(actionSpace.Subspaces)
	// case "MultiBinary":
	// case "MultiDiscrete":
	// case "Box":
	default:
		return nil, fmt.Errorf("ActionSpace type %v not yet implemented", actionSpace.Type)
	}

	// var genes []*gene.Gene
	switch obsSpace.Type {
	case "Discrete":
		funcs := []gene.FuncWeight{
			{Symbol: "+", Weight: 1},
			{Symbol: "-", Weight: 5},
			{Symbol: "*", Weight: 5},
		}
		// for i := 0; i < numGenes; i++ {
		// 	genes = append(genes, gene.RandomNew(headSize, tailSize, 1, numConstants, funcs, functions.Int))
		// }
		// o.Genomes = append(o.Genomes, genome.New(genes, "tuple"))
		const numTerminals = 2 // Account for "stepsSinceReset"
		gen := New(funcs, functions.Int, 1, headSize, numGenes, numTerminals, numConstants, "tuple", nil)
		o.Genomes = gen.Genomes

	case "Tuple":
		funcType := functions.Int
		var numInputs int
		for i, subspace := range obsSpace.Subspaces {
			numInputs += subspace.N
			if subspace.Type != "Discrete" {
				log.Printf("WARNING: Subspace[%v] unhandled Type=%q", i, subspace.Type)
			}
		}

		funcs := []gene.FuncWeight{
			{Symbol: "+", Weight: 1},
			{Symbol: "-", Weight: 5},
			{Symbol: "*", Weight: 5},
		}
		const numTerminals = 2 // Account for "stepsSinceReset"
		numGenomes := len(obsSpace.Subspaces)
		gen := New(funcs, funcType, numGenomes, headSize, numGenes, numInputs, numConstants, "tuple", nil)
		o.Genomes = gen.Genomes

	// case "MultiBinary":
	// case "MultiDiscrete":
	// case "Box":
	default:
		return nil, fmt.Errorf("ObservationSpace type %v not yet implemented", obsSpace.Type)
	}
	// log.Printf("numGenes=%v, genes=%#v", numGenes, genes)

	return o, nil
}

// Evaluate runs the model and returns an action from an observation.
func (o *Gymnasium) Evaluate(stepsSinceReset int, obs gym.Obs, action interface{}) error {
	if err := o.Genomes[0].Evaluate(stepsSinceReset, obs, action); err != nil {
		return err
	}
	// log.Printf("Evaluate: obs=%v, raw action=%v", obs, action)

	switch o.ActionSpace.Type {
	// case "Discrete":
	case "Tuple":
		if vals, ok := action.(*[]int); ok {
			for i, v := range *vals {
				if o.ActionSpace.Subspaces[i].Type == "Discrete" {
					if v < 0 {
						v = -v
					}
					(*vals)[i] = v % o.ActionSpace.Subspaces[i].N
				} else {
					return fmt.Errorf("subspace %v not yet supported", o.ActionSpace.Subspaces[i])
				}
			}
		} else {
			return fmt.Errorf("action type %T not yet supported", action)
		}
	// case "MultiBinary":
	// case "MultiDiscrete":
	// case "Box":
	default:
		return fmt.Errorf("ActionSpace type %v not yet implemented", o.ActionSpace.Type)
	}

	// log.Printf("Evaluate: obs=%v, final action=%v", obs, action)

	return nil
}

// Evolve evolves the model based on the reward feedback (from -1 (bad) to 1 (good)).
func (o *Gymnasium) Evolve(reward float64) error {
	o.Genomes[0].Score = (reward + 1.0) * 500.0 // -1..1 => 0..1000

	// Find the best genome so far.
	sort.Slice(o.Genomes, func(a, b int) bool { return o.Genomes[a].Score > o.Genomes[b].Score })

	if len(o.Genomes) < maxGenomes || o.Genomes[0].Score < 500 {
		oai, err := NewGymnasium(o.ActionSpace, o.ObsSpace)
		check("NewGymnasium: %v", err)
		o.Genomes = append([]*genome.Genome{oai.Genomes[0]}, o.Genomes...)
	} else {
		ng := o.Genomes[0].Dup()
		gen := &Generation{Genomes: o.Genomes}
		gen.replication()
		gen.mutation()
		// o.Genomes = append([]*genome.Genome{ng}, gen.Genomes...)
		o.Genomes = append(gen.Genomes[0:maxGenomes-1], ng)

		// for i, g := range o.Genomes {
		// 	log.Printf("genome[%v]=%v", i, g)
		// 	for in := 0; in < o.ObsSpace.N; in++ {
		// 		action := []int{}
		// 		intObs := &intObsT{i: in}
		// 		if err := g.Evaluate(intObs, &action); err != nil {
		// 			log.Fatalf("g.Evaluate: %v", err)
		// 		}
		// 		log.Printf("genome[%v](%v) = %v", i, in, action)
		// 	}
		// }
		// log.Fatalf("Debug stop.")
	}

	if len(o.Genomes) > maxGenomes {
		o.Genomes = o.Genomes[0:maxGenomes]
	}

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
