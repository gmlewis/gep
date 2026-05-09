// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package problems

import (
	"errors"

	"github.com/gmlewis/gep/v2/core"
	boolFitness "github.com/gmlewis/gep/v2/fitness/bool"
	floatFitness "github.com/gmlewis/gep/v2/fitness/float"
)

// Case is a single fitness evaluation case pairing an input vector with the
// expected output for that input. T is the domain type (e.g. bool, float64).
//
// In is the input vector fed to the genome's terminals (d0, d1, …).
// Out is the expected output value for this case.
type Case[T any] struct {
	// In is the input vector fed to the genome's terminals (d0, d1, ...).
	In []T
	// Out is the expected output value for this case.
	Out T
}

// BoolProblem is a boolean logic problem defined by a set of [Case][bool]
// fitness cases.  It provides methods that return evolution-compatible fitness
// functions so that experiment and application code can consume domain-specific
// problems without re-implementing scoring logic.
//
// Use [NewBoolProblem] to construct a BoolProblem, then call one of its
// scoring-function methods to obtain a func(core.Genome[bool]) float64 that
// can be passed directly to [github.com/gmlewis/gep/v2/evolution.New].
type BoolProblem struct {
	// Cases is the complete set of fitness evaluation cases for this problem.
	Cases []Case[bool]
}

// NewBoolProblem creates a [BoolProblem] from the given fitness cases.
// It returns an error if cases is empty.
func NewBoolProblem(cases []Case[bool]) (*BoolProblem, error) {
	if len(cases) == 0 {
		return nil, errors.New("problems.NewBoolProblem: cases must not be empty")
	}
	return &BoolProblem{Cases: cases}, nil
}

// NumHitsScoringFunc returns a fitness function that scores a genome by the
// number of cases it predicts correctly, multiplied by scaleFactor.
// The maximum possible score is scaleFactor * len(p.Cases).
//
// The returned function evaluates the genome against every case in p.Cases
// using [github.com/gmlewis/gep/v2/core.Genome.Eval], then delegates to
// [github.com/gmlewis/gep/v2/fitness/bool.NumHits] for the final score.
// A genome that fails to evaluate any case contributes a score of zero for
// that case but does not stop the overall scoring run.
//
// The returned function is compatible with the scoringFunc argument of
// [github.com/gmlewis/gep/v2/evolution.New].
func (p *BoolProblem) NumHitsScoringFunc(scaleFactor float64) (func(core.Genome[bool]) float64, error) {
	fitnessFunc, err := boolFitness.NumHits(scaleFactor)
	if err != nil {
		return nil, err
	}
	cases := p.Cases
	return func(g core.Genome[bool]) float64 {
		predicted := make([]bool, len(cases))
		target := make([]bool, len(cases))
		for i, c := range cases {
			r, err := g.Eval(c.In)
			if err != nil {
				// Evaluation failure: leave predicted[i] as zero value (false).
				// The case will be treated as incorrect.
				continue
			}
			predicted[i] = r
			target[i] = c.Out
		}
		score, err := fitnessFunc(predicted, target)
		if err != nil {
			return 0
		}
		return score
	}, nil
}

// FloatProblem is a floating-point regression problem defined by a set of
// [Case][float64] fitness cases.  It provides methods that return
// evolution-compatible fitness functions, each corresponding to a different
// error metric from [github.com/gmlewis/gep/v2/fitness/float].
//
// Use [NewFloatProblem] to construct a FloatProblem, then call one of its
// scoring-function methods to obtain a func(core.Genome[float64]) float64 that
// can be passed directly to [github.com/gmlewis/gep/v2/evolution.New].
type FloatProblem struct {
	// Cases is the complete set of fitness evaluation cases for this problem.
	Cases []Case[float64]
}

// NewFloatProblem creates a [FloatProblem] from the given fitness cases.
// It returns an error if cases is empty.
func NewFloatProblem(cases []Case[float64]) (*FloatProblem, error) {
	if len(cases) == 0 {
		return nil, errors.New("problems.NewFloatProblem: cases must not be empty")
	}
	return &FloatProblem{Cases: cases}, nil
}

// evaluateFloat evaluates g against all cases and returns predicted and target
// slices.  Cases where evaluation fails contribute a predicted value of 0.
func (p *FloatProblem) evaluateFloat(g core.Genome[float64]) (predicted, target []float64) {
	predicted = make([]float64, len(p.Cases))
	target = make([]float64, len(p.Cases))
	for i, c := range p.Cases {
		target[i] = c.Out
		r, err := g.Eval(c.In)
		if err != nil {
			continue
		}
		predicted[i] = r
	}
	return predicted, target
}

// NumHitsAbsScoringFunc returns a fitness function that scores a genome by the
// number of cases whose absolute prediction error does not exceed precision.
// scaleFactor multiplies the final hit count; the maximum possible score is
// scaleFactor * len(p.Cases).
//
// precision must be in [0, 1].  The returned function delegates to
// [github.com/gmlewis/gep/v2/fitness/float.NumHitsAbs] for the final score.
//
// The returned function is compatible with the scoringFunc argument of
// [github.com/gmlewis/gep/v2/evolution.New].
func (p *FloatProblem) NumHitsAbsScoringFunc(precision, scaleFactor float64) (func(core.Genome[float64]) float64, error) {
	fitnessFunc, err := floatFitness.NumHitsAbs(precision, scaleFactor)
	if err != nil {
		return nil, err
	}
	return func(g core.Genome[float64]) float64 {
		predicted, target := p.evaluateFloat(g)
		score, err := fitnessFunc(predicted, target)
		if err != nil {
			return 0
		}
		return score
	}, nil
}

// MeanSquaredErrorAbsScoringFunc returns a fitness function that computes the
// absolute mean-squared error between predictions and targets, then normalizes
// the result to [0, scaleFactor] using the formula
// scaleFactor / (1 + MSE/N).
//
// The returned function delegates to
// [github.com/gmlewis/gep/v2/fitness/float.MeanSquaredErrorAbs] for the final
// score.
//
// The returned function is compatible with the scoringFunc argument of
// [github.com/gmlewis/gep/v2/evolution.New].
func (p *FloatProblem) MeanSquaredErrorAbsScoringFunc(scaleFactor float64) (func(core.Genome[float64]) float64, error) {
	fitnessFunc, err := floatFitness.MeanSquaredErrorAbs(scaleFactor)
	if err != nil {
		return nil, err
	}
	return func(g core.Genome[float64]) float64 {
		predicted, target := p.evaluateFloat(g)
		score, err := fitnessFunc(predicted, target)
		if err != nil {
			return 0
		}
		return score
	}, nil
}

// RSquareScoringFunc returns a fitness function based on the square of the
// Pearson product-moment correlation coefficient (R²), normalized to
// [0, scaleFactor].
//
// The returned function delegates to
// [github.com/gmlewis/gep/v2/fitness/float.RSquare] for the final score.
//
// The returned function is compatible with the scoringFunc argument of
// [github.com/gmlewis/gep/v2/evolution.New].
func (p *FloatProblem) RSquareScoringFunc(scaleFactor float64) (func(core.Genome[float64]) float64, error) {
	fitnessFunc, err := floatFitness.RSquare(scaleFactor)
	if err != nil {
		return nil, err
	}
	return func(g core.Genome[float64]) float64 {
		predicted, target := p.evaluateFloat(g)
		score, err := fitnessFunc(predicted, target)
		if err != nil {
			return 0
		}
		return score
	}, nil
}
