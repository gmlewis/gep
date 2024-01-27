// -*- compile-command: "go run main.go"; -*-

// blackjack runs a GEP algorithm on the Gymnasium "Blackjack-v1" algorithm problem.
// https://gymnasium.farama.org/tutorials/training_agents/blackjack_tutorial/
package main

import (
	"flag"
	"log"
	"os"
	"time"

	gym "github.com/gmlewis/gep/v2/gymnasium"
	"github.com/gmlewis/gep/v2/model"
)

const (
	environment     = "Blackjack-v1"
	defaultMinSteps = 1e8
	defaultMaxSteps = 2e8
)

var (
	debug          = flag.Bool("d", false, "Show debug information")
	episodesPerInd = flag.Int("epi", 1000, "Episodes per individual") // TODO
	headSize       = flag.Int("hs", 100, "Head size of karva expressions")
	maxIndividuals = flag.Int("mi", 100, "Maximum individuals in population")
	numConsts      = flag.Int("nc", 2, "Number of constants in karva expressions")
	numEnvs        = flag.Int("ne", 1000, "Number of concurrent Blackjack environments")
	showHelp       = flag.Bool("h", false, "Show help message")
	showTime       = flag.Bool("t", false, "Display timestamps")
	minSteps       = flag.Int("min", defaultMinSteps, "Minimum number of steps to run")
	maxSteps       = flag.Int("max", defaultMaxSteps, "Maximum number of steps to run")
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

	if *debug {
		log.Printf("Running %v concurrent blackjack tables with %v episodes per individual", *numEnvs, *episodesPerInd)
	}

	env, err := gym.Make(environment)
	check("gym.Make(%q, %q): %v", environment, err)
	defer env.Close()

	actionSpace, err := env.ActionSpace()
	check("ActionSpace: %v", err)
	log.Printf("Action space: %+v", actionSpace)
	for i, subSpace := range actionSpace.Subspaces {
		log.Printf("SubSpace[%v]: %+v", i, subSpace)
	}

	obsSpace, err := env.ObservationSpace()
	check("ObservationSpace: %v", err)
	log.Printf("Observation space: %+v", obsSpace)
	for i, sub := range obsSpace.Subspaces {
		log.Printf("Observation subspace[%v]: %+v", i, *sub)
	}

	opts := []model.GymnasiumAgentOption{
		model.WithHeadSize(*headSize),
		model.WithNumConstants(*numConsts),
		model.WithMaxIndividuals(*maxIndividuals),
	}
	agent, err := model.NewGymnasiumAgent(actionSpace, obsSpace, opts...)
	check("NewGymnasium: %v", err)

	lastObs, _ := env.Reset()
	log.Printf("env.Reset: lastObj=%T=%s", lastObs, lastObs)

	startTime := time.Now()
	var episodeReward float64
	var episodeSteps int
	var totalSteps int
	for (episodeReward < 1.0 || totalSteps < *minSteps) && totalSteps < *maxSteps {
		var action int
		// TODO: integrate Evaluate into loop.
		// err := agent.Evaluate(episodeSteps, lastObs, &action)
		// check("Evaluate(%v): %v", lastObs, err)
		err := env.SampleAction(&action)
		check("env.SampleAction: %v", err)

		obs, reward, terminated, truncated, _ := env.Step(action)
		check("Step(%v): %v", action, err)
		episodeReward += reward
		totalSteps++
		episodeSteps++
		if *debug || totalSteps%(*minSteps/100) == 0 {
			log.Printf("Step #%v: episodeSteps=%v, obs=%v, action=%v, reward=%v, episodeReward=%v, terminated=%v", totalSteps, episodeSteps, lastObs, action, reward, episodeReward, terminated)
		}
		lastObs = obs

		if terminated || truncated {
			err := agent.Evolve(episodeReward)
			check("Evolve(%v): %v", episodeReward, err)
			lastObs, _ = env.Reset()
			episodeSteps = 0
			episodeReward = 0
		}
	}

	seconds := time.Since(startTime).Seconds()
	log.Printf("Processed %v steps in %v seconds (%v steps/sec)", totalSteps, seconds, float64(totalSteps)/seconds)

	log.Printf("\nFinal population:")
	for i, individual := range agent.Individuals {
		log.Printf("Individual #%v: %v", i+1, individual)
	}
	log.Printf("Done.")
}

func check(fmt string, args ...any) {
	err := args[len(args)-1]
	if err != nil {
		log.Fatalf("ERROR: "+fmt, args...)
	}
}
