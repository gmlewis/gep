// Package gymnasium provides pure Go implementations of Gymnasium environments.
package gymnasium

import (
	"fmt"

	"github.com/gmlewis/gep/v2/gymnasium/envs/toy_text/blackjack"
	gym "github.com/gmlewis/gym-socket-api/binding-go"
)

// Environment represents a pure Go training and execution environment.
type Environment interface {
	ActionSpace() (*gym.Space, error)
	ObservationSpace() (*gym.Space, error)
	Reset() (gym.Obs, any)
	SampleAction(action any) error
	Step(action any) (obs gym.Obs, reward float64, terminated bool, truncated bool, info any)
	Close() error
}

// Make returns a specific gymnasium environment.
func Make(environment string) (Environment, error) {
	switch environment {
	case "Blackjack-v1":
		return blackjack.New(false, false), nil
	default:
		return nil, fmt.Errorf("unknown environment %q", environment)
	}
}
