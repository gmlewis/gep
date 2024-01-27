package model

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gmlewis/gep/v2/common"
	"github.com/google/go-cmp/cmp"
)

type obsT []int

func (o obsT) Unmarshal(dst any) error {
	switch v := dst.(type) {
	case *[]int:
		(*v)[0] = o[0]
		(*v)[1] = o[1]
		(*v)[2] = o[2]
	default:
		return fmt.Errorf("unsupported obs type %T", dst)
	}
	return nil
}

func TestNewGymnasiumAgents(t *testing.T) {
	tests := []struct {
		name        string
		actionSpace *common.Space
		obsSpace    *common.Space
	}{
		{
			name:        "Blackjack",
			actionSpace: &common.Space{Type: "Discrete", N: 2},
			obsSpace: &common.Space{
				Type: "Tuple",
				Subspaces: []*common.Space{
					{Type: "Discrete", N: 32},
					{Type: "Discrete", N: 11},
					{Type: "Discrete", N: 2},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agents, err := NewGymnasiumAgents(tt.actionSpace, tt.obsSpace)
			if err != nil {
				t.Fatal(err)
			}

			step := func() common.Obs {
				n := len(tt.obsSpace.Subspaces)
				lastObs := make(obsT, n)
				for i := 0; i < n; i++ {
					lastObs[i] = rand.Intn(tt.obsSpace.Subspaces[i].N)
				}
				return lastObs
			}

			for agentNum := 0; agentNum < defaultNumIndividuals; agentNum++ {
				for episode := 1; episode <= 10; episode++ {
					var action int
					for episodeStep := 0; episodeStep < 10; episodeStep++ {
						lastObs := step()
						if err := agents.EvaluateAgent(agentNum, episodeStep, lastObs, &action); err != nil {
							t.Fatal(err)
						}

						if action < 0 || action >= tt.actionSpace.N {
							t.Fatalf("episode=%v: agent.Evaluate returned bad action: %v, want 0<#<%v", episode, action, tt.actionSpace.N)
						}
					}
					randomReward := rand.Float64()*2000 - 1000
					agents.RewardAgent(agentNum, randomReward)
				}
			}

			if err := agents.Evolve(); err != nil {
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
		actionSpace        *common.Space
		obsSpace           *common.Space
		want               []int
	}{
		{
			name:        "Blackjack",
			actionSpace: &common.Space{Type: "Discrete", N: 2},
			obsSpace: &common.Space{
				Type: "Tuple",
				Subspaces: []*common.Space{
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
			actionSpace:        &common.Space{Type: "Discrete", N: 2},
			obsSpace: &common.Space{
				Type: "Tuple",
				Subspaces: []*common.Space{
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
			var opts []GymnasiumAgentsOption
			if tt.appendEpisodeSteps {
				opts = append(opts, WithAppendEpisodeSteps())
			}
			agents, err := NewGymnasiumAgents(tt.actionSpace, tt.obsSpace, opts...)
			if err != nil {
				t.Fatal(err)
			}

			step := func() common.Obs {
				n := len(tt.obsSpace.Subspaces)
				lastObs := make(obsT, n)
				for i := 0; i < n; i++ {
					lastObs[i] = i
				}
				return lastObs
			}

			lastObs := step()
			got, err := agents.processObservations(episodeStep, lastObs)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("processObservations mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
