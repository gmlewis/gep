// -*- compile-command: "go run main.go"; -*-

// blackjack runs a GEP algorithm on the Gymnasium "Blackjack-v1" algorithm problem.
// https://gymnasium.farama.org/tutorials/training_agents/blackjack_tutorial/
package main

import (
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gmlewis/gep/v2/common"
	gym "github.com/gmlewis/gep/v2/gymnasium"
	"github.com/gmlewis/gep/v2/model"
)

const (
	environment  = "Blackjack-v1"
	defaultSteps = 1e4
)

var (
	debug           = flag.Bool("d", false, "Show debug information")
	episodesPerStep = flag.Int("eps", 1000, "Episodes to run per step")
	headSize        = flag.Int("hs", 100, "Head size of karva expressions")
	numConsts       = flag.Int("nc", 2, "Number of constants in karva expressions")
	numIndividuals  = flag.Int("ni", 10, "Number of individuals in population")
	numSteps        = flag.Int("s", defaultSteps, "Number of total steps to run")
	showHelp        = flag.Bool("h", false, "Show help message")
	showTime        = flag.Bool("t", false, "Display timestamps")
)

func usage() {
	log.Printf("Usage: blackjack")
	flag.PrintDefaults()
}

func showUsageAndExit(exitCode int) {
	usage()
	os.Exit(exitCode)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if !*showTime {
		log.SetFlags(0) // disable output of timestamp from log.* functions.
	}

	if *showHelp {
		showUsageAndExit(0)
	}

	log.Printf("Running %v concurrent blackjack tables for %v steps with %v episodes per step", *numIndividuals, *numSteps, *episodesPerStep)

	actionSpace, obsSpace, err := gym.GetSpaces(environment)
	log.Printf("%v Action space: %+v", environment, actionSpace)
	for i, subspace := range actionSpace.Subspaces {
		log.Printf("Action subspace[%v]: %+v", i, *subspace)
	}
	log.Printf("%v Observation space: %+v", environment, obsSpace)
	for i, subspace := range obsSpace.Subspaces {
		log.Printf("Observation subspace[%v]: %+v", i, *subspace)
	}

	opts := []model.GymnasiumAgentsOption{
		model.WithHeadSize(*headSize),
		model.WithNumConstants(*numConsts),
		model.WithNumIndividuals(*numIndividuals),
	}
	agents, err := model.NewGymnasiumAgents(actionSpace, obsSpace, opts...)
	check("NewAgents: %v", err)

	makeEvaluateAgentFunc := func(agentIdx int) func(episodeSteps int, obs common.Obs, action any) error {
		return func(episodeSteps int, obs common.Obs, action any) error {
			return agents.EvaluateAgent(agentIdx, episodeSteps, obs, action)
		}
	}
	makeRewardAgentFunc := func(agentIdx int) func(reward float64) {
		return func(reward float64) {
			agents.RewardAgent(agentIdx, reward)
		}
	}

	casino := &casinoT{
		tables: make([]*tableT, 0, *numIndividuals),
	}
	for agentIdx := 0; agentIdx < *numIndividuals; agentIdx++ {
		env, err := gym.Make(environment)
		check("gym.Make: %v", err)
		casino.tables = append(casino.tables, &tableT{
			env:           env,
			evaluateAgent: makeEvaluateAgentFunc(agentIdx),
			rewardAgent:   makeRewardAgentFunc(agentIdx),
		})
	}

	startTime := time.Now()

	for i := 1; i < *numSteps; i++ {
		casino.runEpisodes(*episodesPerStep)
		if *debug || i%(*numSteps/100) == 0 {
			agents.SortIndividuals()
			log.Printf("Step #%v: numEpisodes=%v, best: %v", i, *episodesPerStep*i, agents.Individuals[0])
		}
		agents.Evolve()
	}
	// Run one final set of episodes to generate final scores:
	casino.runEpisodes(*episodesPerStep)

	seconds := time.Since(startTime).Seconds()
	log.Printf("Processed %v steps in %v seconds (%v steps/sec)", *numSteps, seconds, float64(*numSteps)/seconds)

	// Sort the agents by score
	agents.SortIndividuals()

	log.Printf("\nFinal population:")
	for i, individual := range agents.Individuals {
		log.Printf("Individual #%v: %v", i+1, individual)
	}
	log.Printf("Done.")
}

type casinoT struct {
	tables []*tableT
}

type tableT struct {
	env           gym.Environment
	evaluateAgent func(episodeSteps int, obs common.Obs, action any) error
	rewardAgent   func(reward float64)
}

func (c *casinoT) runEpisodes(numEpisodes int) {
	// Run all tables concurrently for numEpisodes
	var wg sync.WaitGroup
	wg.Add(len(c.tables))
	for _, table := range c.tables {
		go func(table *tableT) {
			table.runEpisodes(numEpisodes)
			wg.Done()
		}(table)
	}
	wg.Wait()
}

func (t *tableT) runEpisodes(numEpisodes int) {
	lastObs, _ := t.env.Reset()
	// if *debug {
	// 	log.Printf("env.Reset: lastObj=%T=%s", lastObs, lastObs)
	// }

	// Run numEpisodes
	var totalReward float64
	var episode, episodeSteps int
	for episode < numEpisodes {
		var action int
		err := t.evaluateAgent(episodeSteps, lastObs, &action)
		check("Evaluate(%v): %v", lastObs, err)

		obs, reward, terminated, truncated, _ := t.env.Step(action)
		check("Step(%v): %v", action, err)

		totalReward += reward
		// if *debug || totalSteps%(*minSteps/100) == 0 {
		// 	log.Printf("Step #%v: episodeSteps=%v, obs=%v, action=%v, reward=%v, episodeReward=%v, terminated=%v", totalSteps, episodeSteps, lastObs, action, reward, episodeReward, terminated)
		// }

		if terminated || truncated {
			lastObs, _ = t.env.Reset()
			episode++
			episodeSteps = 0
		} else {
			lastObs = obs
			episodeSteps++
		}
	}

	t.rewardAgent(totalReward)
}

func check(fmt string, args ...any) {
	err := args[len(args)-1]
	if err != nil {
		log.Fatalf("ERROR: "+fmt, args...)
	}
}
