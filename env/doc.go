// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package env provides the repository's dedicated environment / reinforcement
// learning (RL) subsystem.
//
// It owns the orchestration of Gymnasium-compatible agent populations:
// creating initial populations, evaluating individual agents against
// environment observations, rewarding agents with episode scores, and
// evolving the population across episodes.  The gymnasium package remains
// focused only on pure environment definitions and metadata.
//
// # Core types
//
//   - [GymnasiumAgents] – a population of GEP genomes that act as RL agents
//     inside a Gymnasium-compatible environment.
//   - [GymnasiumAgentsOption] – functional option type for [NewGymnasiumAgents].
//
// # Constructors
//
//   - [NewGymnasiumAgents] – builds a random initial population from action-
//     and observation-space metadata.
//
// # Functional options
//
//   - [WithAppendEpisodeSteps] – appends the current episode-step count to the
//     observation vector (useful for environments where step count matters).
//   - [WithDebug] – enables verbose debug logging.
//   - [WithHeadSize] – controls the maximum complexity of the generated
//     decision programs.
//   - [WithNumConstants] – sets the number of numerical constants available
//     inside each gene.
//   - [WithNumIndividuals] – sets the population size.
//
// # Agent lifecycle
//
//	agents, _ := env.NewGymnasiumAgents(actionSpace, obsSpace,
//	    env.WithNumIndividuals(100),
//	    env.WithHeadSize(20),
//	)
//	// Per step: evaluate one agent, accumulate reward, then evolve.
//	agents.EvaluateAgent(agentIdx, episodeSteps, obs, &action)
//	agents.RewardAgent(agentIdx, reward)
//	agents.Evolve()
package env
