// -*- compile-command: "go run main.go"; -*-

// Copy-v0 runs a GEP algorithm on the OpenAI Gym "Copy-v0" Algorithm problem.
// https://gym.openai.com/envs/Copy-v0/
package main

import (
	"log"
	"time"

	gym "github.com/gmlewis/gym-socket-api/binding-go"
)

const (
	host        = "localhost:5001"
	environment = "Copy-v0"
)

func main() {
	env, err := gym.Make(host, environment)
	if err != nil {
		log.Fatalf("gym.Make(%q, %q): %v", host, environment, err)
	}
	defer env.Close()

	obsSpace, err := env.ObservationSpace()
	if err != nil {
		log.Fatalf("ObservationSpace(): %v", err)
	}
	log.Printf("Observation space: %+v", obsSpace)

	obs, err := env.Reset()
	if err != nil {
		log.Fatalf("Reset(): %v", err)
	}
	log.Printf("1st observation: %v", obs)

	startTime := time.Now()
	var steps int
	var done bool
	for !done {
		// TODO: Replace SampleAction with GEP algorithm.
		var action []int
		if err := env.SampleAction(&action); err != nil {
			log.Fatalf("SampleAction(): %v", err)
		}

		var reward float64
		obs, reward, done, _, err = env.Step(action)
		if err != nil {
			log.Fatalf("Step(%v): %v", action, err)
		}
		steps++
		log.Printf("Step #%v: action=%v, obs=%v, reward=%v, done=%v", steps, action, obs, reward, done)
	}

	seconds := time.Since(startTime).Seconds()
	log.Printf("Processed %v steps in %v seconds (%v sps)", steps, seconds, float64(steps)/seconds)
}
