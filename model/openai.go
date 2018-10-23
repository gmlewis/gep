// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model

import (
	"log"

	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/genome"
	gym "github.com/gmlewis/gym-socket-api/binding-go"
)

const (
	headSize     = 10
	tailSize     = 10
	numConstants = 6
)

// Model represents a GEP model.
type Model interface {
	Evaluate(obs gym.Obs, action interface{}) error
	Evolve(reward float64) error
}

// openai implements the Model interface.
type openai struct {
	actionSpace *gym.Space
	obsSpace    *gym.Space
	genome      *genome.Genome
}

// ForOpenAI returns a Model based upon the action and observation spaces.
func ForOpenAI(actionSpace, obsSpace *gym.Space) (Model, error) {
	m := &openai{actionSpace: actionSpace, obsSpace: obsSpace}

	numGenes := 1
	switch actionSpace.Type {
	case "Discrete":
	case "Tuple":
		numGenes = actionSpace.N
	case "MultiBinary":
	case "MultiDiscrete":
	case "Box":
	}

	var genes []*gene.Gene
	switch obsSpace.Type {
	case "Discrete":
		funcs := []gene.FuncWeight{
			{Symbol: "+", Weight: 1},
			{Symbol: "-", Weight: 5},
			{Symbol: "*", Weight: 5},
		}
		for i := 0; i < numGenes; i++ {
			genes = append(genes, gene.RandomNew(headSize, tailSize, 1, numConstants, funcs, functions.Int))
		}
		m.genome = genome.New(genes, "tuple")
	case "Tuple":
	case "MultiBinary":
	case "MultiDiscrete":
	case "Box":
	}
	log.Printf("numGenes=%v, genes=%v", numGenes, genes)

	return m, nil
}

// Evaluate runs the model and returns an action from an observation.
func (o *openai) Evaluate(obs gym.Obs, action interface{}) error {
	return nil
}

// Evolve evolves the model based on the reward feedback (from -1 (bad) to 1 (good)).
func (o *openai) Evolve(reward float64) error {
	return nil
}
