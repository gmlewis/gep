// -*- compile-command: "go run main.go"; -*-

// Copy-v0 runs a GEP algorithm on the OpenAI Gym "Copy-v0" Algorithm problem.
// https://gym.openai.com/envs/Copy-v0/
package main

import (
	"log"

	gym "github.com/openai/gym-http-api/binding-go"
)

const (
	baseURL = "http://localhost:5000"
	env     = "Copy-v0"
)

func main() {
	c, err := gym.NewClient(baseURL)
	if err != nil {
		log.Fatalf("gym.NewClient(%q): %v", baseURL, err)
	}

	id, err := c.Create(env)
	if err != nil {
		log.Fatalf("Unable to create environment %q: %v", env, err)
	}
	defer c.Close(id)
	log.Printf("id=%q", id)

	actionSpace, err := c.ActionSpace(id)
	if err != nil {
		log.Fatalf("ActionSpace(%q): %v", id, err)
	}
	log.Printf("Action space: %+v", actionSpace)
	observationSpace, err := c.ObservationSpace(id)
	if err != nil {
		log.Fatalf("ObservationSpace(%q): %v", id, err)
	}
	log.Printf("Observation space: %+v", observationSpace)

	if err := c.StartMonitor(id, "/tmp/"+env, true, false, false); err != nil {
		log.Fatalf(`StartMonitor(%q, "/tmp/%v"): %v`, id, env, err)
	}

	obs, err := c.Reset(id)
	if err != nil {
		log.Fatalf("Reset(%q): %v", id, err)
	}
	log.Printf("1st observation: %v", obs)
	var done bool
	for !done {
		// TODO: Replace SampleAction with GEP algorithm.
		act, err := c.SampleAction(id)
		if err != nil {
			log.Fatalf("SampleAction(%q): %v", id, err)
		}

		var reward float64
		obs, reward, done, _, err = c.Step(id, act, false)
		if err != nil {
			log.Fatalf("Step(%q, %v, false): %v", id, act, err)
		}
		log.Printf("reward=%v, observatioan=%v", reward, obs)
	}

	if err := c.CloseMonitor(id); err != nil {
		log.Fatalf("CloseMonitor(%q): %v", id, err)
	}
}
