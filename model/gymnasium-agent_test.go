package model

import (
	"math/rand"
	"testing"

	gym "github.com/gmlewis/gym-socket-api/binding-go"
	jsoniter "github.com/json-iterator/go"
)

var jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary

type stepT []int

func (s stepT) Unmarshal(dst interface{}) error {
	buf, err := jsoncomp.Marshal(s)
	if err != nil {
		return err
	}
	return jsoncomp.Unmarshal(buf, dst)
}

func TestNewGymnasiumAgent(t *testing.T) {
	tests := []struct {
		name        string
		actionSpace *gym.Space
		obsSpace    *gym.Space
	}{
		{
			name:        "Blackjack",
			actionSpace: &gym.Space{Type: "Discrete", N: 2},
			obsSpace: &gym.Space{
				Type: "Tuple",
				Subspaces: []*gym.Space{
					{Type: "Discrete", N: 32},
					{Type: "Discrete", N: 11},
					{Type: "Discrete", N: 2},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewGymnasiumAgent(tt.actionSpace, tt.obsSpace)
			if err != nil {
				t.Fatal(err)
			}

			step := func() gym.Obs {
				n := len(tt.obsSpace.Subspaces)
				lastObs := make(stepT, n)
				for i := 0; i < n; i++ {
					lastObs[i] = rand.Intn(tt.obsSpace.Subspaces[i].N)
				}
				return lastObs
			}

			for episode := 1; episode <= 1000; episode++ {
				var action int
				for i := 0; i < 10; i++ {
					lastObs := step()
					if err := agent.Evaluate(i, lastObs, &action); err != nil {
						t.Fatal(err)
					}

					if action < 0 || action >= tt.actionSpace.N {
						t.Fatalf("episode=%v: agent.Evaluate returned bad action: %v, want 0<#<%v", episode, action, tt.actionSpace.N)
					}
				}
			}

			randomReward := rand.Float64()*2000 - 1000
			if err := agent.Evolve(randomReward); err != nil {
				t.Fatal(err)
			}
		})
	}
}
