// Package blackjack implements the Blackjack environement in Go.
// It is inspired by: https://gymnasium.farama.org/tutorials/training_agents/blackjack_tutorial/
// and: https://github.com/Farama-Foundation/Gymnasium/blob/main/gymnasium/envs/toy_text/blackjack.py
//
// Blackjack is a card game where the goal is to beat the dealer by obtaining cards
// that sum to closer to 21 (without going over 21) than the dealers cards.
//
// # Description
// The game starts with the dealer having one face up and one face down card,
// while the player has two face up cards. All cards are drawn from an infinite deck
// (i.e. with replacement).
//
// The card values are:
// - Face cards (Jack, Queen, King) have a point value of 10.
// - Aces can either count as 11 (called a 'usable ace') or 1.
// - Numerical cards (2-9) have a value equal to their number.
//
// The player has the sum of cards held. The player can request
// additional cards (hit) until they decide to stop (stick) or exceed 21 (bust,
// immediate loss).
//
// After the player sticks, the dealer reveals their facedown card, and draws cards
// until their sum is 17 or greater. If the dealer goes bust, the player wins.
//
// If neither the player nor the dealer busts, the outcome (win, lose, draw) is
// decided by whose sum is closer to 21.
//
// This environment corresponds to the version of the blackjack problem
// described in Example 5.1 in Reinforcement Learning: An Introduction
// by Sutton and Barto [<a href="#blackjack_ref">1</a>].
//
// # Action Space
// The action shape is `(1,)` in the range `{0, 1}` indicating
// whether to stick or hit.
//
// - 0: Stick
// - 1: Hit
//
// # Observation Space
// The observation consists of a 3-tuple containing: the player's current sum,
// the value of the dealer's one showing card (1-10 where 1 is ace),
// and whether the player holds a usable ace (0 or 1).
//
// The observation is returned as `(int(), int(), int())`.
//
// # Starting State
// The starting state is initialised in the following range.
//
//	| Observation               | Min  | Max  |
//	|---------------------------|------|------|
//	| Player current sum        |  4   |  12  |
//	| Dealer showing card value |  2   |  11  |
//	| Usable Ace                |  0   |  1   |
//
// # Rewards
// - win game: +1
// - lose game: -1
// - draw game: 0
// - win game with natural blackjack:
// +1.5 (if <a href="#nat">natural</a> is True)
// +1 (if <a href="#nat">natural</a> is False)
//
// # Episode End
// The episode ends if the following happens:
//
// - Termination:
// 1. The player hits and the sum of hand exceeds 21.
// 2. The player sticks.
//
// An ace will always be counted as usable (11) unless it busts the player.
//
// # Information
//
// No additional information is returned.
//
// # Arguments
//
// ```python
// import gymnasium as gym
// gym.make('Blackjack-v1', natural=False, sab=False)
// ```
//
// <a id="nat"></a>`natural=False`: Whether to give an additional reward for
// starting with a natural blackjack, i.e. starting with an ace and ten (sum is 21).
//
// <a id="sab"></a>`sab=False`: Whether to follow the exact rules outlined in the book by
// Sutton and Barto. If `sab` is `True`, the keyword argument `natural` will be ignored.
// If the player achieves a natural blackjack and the dealer does not, the player
// will win (i.e. get a reward of +1). The reverse rule does not apply.
// If both the player and the dealer get a natural, it will be a draw (i.e. reward 0).
//
// # References
// <a id="blackjack_ref"></a>[1] R. Sutton and A. Barto, “Reinforcement Learning:
// An Introduction” 2020. [Online]. Available: [http://www.incompleteideas.net/book/RLbook2020.pdf](http://www.incompleteideas.net/book/RLbook2020.pdf)
package blackjack

import (
	"fmt"
	"log"
	"math/rand"
	"slices"

	gym "github.com/gmlewis/gym-socket-api/binding-go"
	jsoniter "github.com/json-iterator/go"
)

var jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary

// Environment represents a Blackjack environment.
type Environment struct {
	natural bool
	sab     bool
	dealer  []int
	player  []int
}

// New returns a new Blackjack environment.
func New(natural, sab bool) *Environment {
	return &Environment{
		natural: natural,
		sab:     sab,
	}
}

// Step performs the provided action and returns the next observation,
// the reward for this action, and whether this episode is terminated.
// If the action is invalid, truncated will be true. Info is always empty.
func (e *Environment) Step(action int) (obs gym.Obs, reward float64, terminated bool, truncated bool, info any) {
	if action < 0 || action > 1 {
		log.Printf("ERROR: Blackjack: invalid action=%v; must be 0 (stay) or 1 (hit).", action)
		return obs, reward, terminated, true, nil
	}

	if action == 1 { // hit: add a card to players hand and return.
		e.player = append(e.player, drawCard())
		if isBust(e.player) {
			terminated = true
			reward = -1
		}
		return e.getObs(), reward, terminated, false, nil
	}

	// stay: play out the dealer's hand, and score.
	terminated = true
	for sumHand(e.dealer) < 17 {
		e.dealer = append(e.dealer, drawCard())
	}
	reward = cmp(score(e.player), score(e.dealer))
	if e.sab && isNatural(e.player) && !isNatural(e.dealer) {
		// Player automatically wins. Rules consistent with S&B.
		reward = 1
	} else if !e.sab && e.natural && isNatural(e.player) && reward == 1 {
		// Natural gives extra points, but doesn't autowin. Legacy implementation.
		reward = 1.5
	}

	return e.getObs(), reward, terminated, false, nil
}

// Reset resets to a brand new episode.
func (e *Environment) Reset() (obs gym.Obs, info any) {
	e.dealer = drawHand()
	e.player = drawHand()
	return e.getObs(), nil
}

type obsT [3]int

var _ gym.Obs = obsT{}

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

func (e *Environment) getObs() (obs gym.Obs) {
	result := obsT{sumHand(e.player), e.dealer[0], 0}
	if usableAce(e.player) {
		result[2] = 1
	}
	return result
}

// 1 = Ace, 2-10 = Number cards, Jack/Queen/King = 10
var deck = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10}

func sum(arr []int) (result int) {
	for _, v := range arr {
		result += v
	}
	return result
}

func cmp(a, b int) float64 {
	var v1, v2 float64
	if a > b {
		v1 = 1
	}
	if a < b {
		v2 = 1
	}
	return v1 - v2
}

func drawCard() int {
	return deck[rand.Intn(len(deck))]
}

func drawHand() []int {
	return []int{drawCard(), drawCard()}
}

func usableAce(hand []int) bool {
	oneInHand := slices.Contains(hand, 1)
	okSum := sum(hand)+10 <= 21
	return oneInHand && okSum
}

func sumHand(hand []int) int {
	if usableAce(hand) {
		return sum(hand) + 10
	}
	return sum(hand)
}

func isBust(hand []int) bool {
	return sumHand(hand) > 21
}

func score(hand []int) int {
	total := sumHand(hand)
	if total > 21 {
		return 0
	}
	return total
}

func isNatural(hand []int) bool {
	if len(hand) != 2 {
		return false
	}
	return (hand[0] == 1 && hand[1] == 10) || (hand[1] == 1 && hand[0] == 10)
}
