// -*- compile-command: "go run main.go"; -*-

// Copy-v0 runs a GEP algorithm on the OpenAI Gym "Copy-v0" Algorithm problem.
// https://gym.openai.com/envs/Copy-v0/
package main

import (
	"flag"
	"log"
	"time"

	gym "github.com/gmlewis/gym-socket-api/binding-go"
)

const (
	host            = "localhost:5001"
	environment     = "Copy-v0"
	defaultMinSteps = 1e4
	defaultMaxSteps = 1e6
)

var (
	minSteps = flag.Int("min", defaultMinSteps, "Minimum number of steps to run")
	maxSteps = flag.Int("max", defaultMaxSteps, "Maximum number of steps to run")
)

func main() {
	flag.Parse()

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

	lastObs, err := env.Reset()
	check("Reset: %v", err)

	startTime := time.Now()
	var (
		steps  int
		done   bool
		reward float64
	)
	for (!done || reward < 1.0 || steps < *minSteps) && steps < *maxSteps {
		if done {
			lastObs, err = env.Reset()
			check("Reset: %v", err)
		}

		// TODO: Replace SampleAction with GEP algorithm.
		var action []int
		err := env.SampleAction(&action)
		check("SampleAction: %v", err)

		var obs gym.Obs
		obs, reward, done, _, err = env.Step(action)
		check("Step(%v): %v", action, err)
		steps++
		log.Printf("Step #%v: obs=%v, action=%v, reward=%v, done=%v", steps, lastObs, action, reward, done)
		lastObs = obs

		// TODO: Send reward to GEP.
	}

	seconds := time.Since(startTime).Seconds()
	log.Printf("Processed %v steps in %v seconds (%v steps/sec)", steps, seconds, float64(steps)/seconds)
}

func check(fmt string, args ...interface{}) {
	err := args[len(args)-1]
	if err != nil {
		log.Fatalf(fmt, args...)
	}
}
