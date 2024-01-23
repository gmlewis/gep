// -*- compile-command: "go run main.go"; -*-

// blackjack runs a GEP algorithm on the Gymnasium "Blackjack-v1" algorithm problem.
// https://gymnasium.farama.org/tutorials/training_agents/blackjack_tutorial/
package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/gmlewis/gep/v2/model"
	gym "github.com/gmlewis/gym-socket-api/binding-go"
)

const (
	host            = "localhost:5001"
	environment     = "Blackjack-v1"
	defaultMinSteps = 1000 // 1e6
	defaultMaxSteps = 2000 // 1e8
)

var (
	debug    = flag.Bool("d", false, "Show debug information")
	showHelp = flag.Bool("h", false, "Show help message")
	showTime = flag.Bool("t", false, "Display timestamps")
	minSteps = flag.Int("min", defaultMinSteps, "Minimum number of steps to run")
	maxSteps = flag.Int("max", defaultMaxSteps, "Maximum number of steps to run")
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
	log.SetFlags(0) // disable output of timestamp from log.* functions.
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	env, err := gym.Make(host, environment)
	check("gym.Make(%q, %q): %v", host, environment, err)
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

	gep, err := model.NewGymnasium(actionSpace, obsSpace)
	check("NewGymnasium: %v", err)

	lastObs, err := env.Reset()
	check("Reset: %v", err)
	log.Printf("env.Reset: lastObj=%T=%s", lastObs, lastObs)

	startTime := time.Now()
	var stepsSinceReset int
	var steps int
	var done bool
	var reward float64
	for (!done || reward < 1.0 || steps < *minSteps) && steps < *maxSteps {
		if done {
			lastObs, err = env.Reset()
			check("Reset: %v", err)
			stepsSinceReset = 0
		}

		var action int
		// TODO: integrate Evaluate into loop.
		// err := gep.Evaluate(stepsSinceReset, lastObs, &action)
		// check("Evaluate(%v): %v", lastObs, err)

		var obs gym.Obs
		obs, reward, done, _, err = env.Step(action)
		check("Step(%v): %v", action, err)
		steps++
		stepsSinceReset++
		if *debug || steps%(*minSteps/100) == 0 {
			log.Printf("Step #%v: ssr=%v, obs=%v, action=%v, reward=%v, done=%v", steps, stepsSinceReset, lastObs, action, reward, done)
		}
		lastObs = obs

		err = gep.Evolve(reward)
		check("Evolve(%v): %v", reward, err)
	}

	seconds := time.Since(startTime).Seconds()
	log.Printf("Processed %v steps in %v seconds (%v steps/sec)", steps, seconds, float64(steps)/seconds)
}

func check(fmt string, args ...interface{}) {
	err := args[len(args)-1]
	if err != nil {
		log.Fatalf("ERROR: "+fmt, args...)
	}
}
