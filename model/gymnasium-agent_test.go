package model

import (
	"math/rand"
	"testing"

	gym "github.com/gmlewis/gym-socket-api/binding-go"
	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
)

var jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary

type stepT []int

func (s stepT) Unmarshal(dst any) error {
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

func TestProcessObservations(t *testing.T) {
	const episodeStep = 4

	tests := []struct {
		name               string
		appendEpisodeSteps bool
		actionSpace        *gym.Space
		obsSpace           *gym.Space
		want               []int
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
			want: []int{0, 1, 2},
		},
		{
			name:               "Blackjack with episode step",
			appendEpisodeSteps: true,
			actionSpace:        &gym.Space{Type: "Discrete", N: 2},
			obsSpace: &gym.Space{
				Type: "Tuple",
				Subspaces: []*gym.Space{
					{Type: "Discrete", N: 32},
					{Type: "Discrete", N: 11},
					{Type: "Discrete", N: 2},
				},
			},
			want: []int{0, 1, 2, episodeStep},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts []GymnasiumAgentOption
			if tt.appendEpisodeSteps {
				opts = append(opts, WithAppendEpisodeSteps())
			}
			agent, err := NewGymnasiumAgent(tt.actionSpace, tt.obsSpace, opts...)
			if err != nil {
				t.Fatal(err)
			}

			step := func() gym.Obs {
				n := len(tt.obsSpace.Subspaces)
				lastObs := make(stepT, n)
				for i := 0; i < n; i++ {
					lastObs[i] = i
				}
				return lastObs
			}

			lastObs := step()
			got, err := agent.processObservations(episodeStep, lastObs)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("processObservations mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
