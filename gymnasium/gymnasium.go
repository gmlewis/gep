// Package gymnasium provides pure Go implementations of Gymnasium environments.
package gymnasium

import (
	"fmt"

	"github.com/gmlewis/gep/v2/common"
	"github.com/gmlewis/gep/v2/gymnasium/envs/toy_text/blackjack"
)

// Environment represents a pure Go training and execution environment.
type Environment interface {
	ActionSpace() (*common.Space, error)
	ObservationSpace() (*common.Space, error)
	Reset() (common.Obs, any)
	SampleAction(action any) error
	Step(action any) (obs common.Obs, reward float64, terminated bool, truncated bool, info any)
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

// GetSpaces returns the ActionSpace and ObsSpace for an environment.
func GetSpaces(environment string) (actionSpace, obsSpace *common.Space, err error) {
	env, err := Make(environment)
	if err != nil {
		return nil, nil, err
	}
	defer env.Close()

	actionSpace, err = env.ActionSpace()
	if err != nil {
		return nil, nil, err
	}

	obsSpace, err = env.ObservationSpace()
	if err != nil {
		return nil, nil, err
	}

	return actionSpace, obsSpace, nil
}
