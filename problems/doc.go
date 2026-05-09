// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package problems provides reusable problem definitions for GEP experiments.
//
// It defines typed problem-facing seams over the [github.com/gmlewis/gep/v2/core]
// and [github.com/gmlewis/gep/v2/fitness] packages, separating reusable
// domain-specific fitness/problem definitions from ad-hoc experiment code.
//
// # Core types
//
//   - [Case] – a single fitness evaluation case pairing an input vector with
//     an expected output value.
//   - [BoolProblem] – a boolean logic problem defined by a set of [Case][bool]
//     fitness cases.
//   - [FloatProblem] – a floating-point regression problem defined by a set of
//     [Case][float64] fitness cases.
//
// # Constructors
//
//   - [NewBoolProblem] – creates a BoolProblem from a slice of boolean cases.
//   - [NewFloatProblem] – creates a FloatProblem from a slice of float64 cases.
//
// # Scoring functions
//
// Each problem type exposes methods that return an evolution-compatible scoring
// function (func([github.com/gmlewis/gep/v2/core.Genome][T]) float64).
// These can be passed directly as the scoringFunc argument to
// [github.com/gmlewis/gep/v2/evolution.New].
//
// Boolean scoring:
//   - [BoolProblem.NumHitsScoringFunc] – score proportional to the number of
//     correctly predicted cases.
//
// Float scoring:
//   - [FloatProblem.NumHitsAbsScoringFunc] – score based on the count of cases
//     within a given absolute error bound.
//   - [FloatProblem.MeanSquaredErrorAbsScoringFunc] – score based on
//     normalized absolute mean-squared error.
//   - [FloatProblem.RSquareScoringFunc] – score based on the R-squared
//     coefficient of determination.
//
// # Example: boolean NAND problem
//
//	cases := []problems.Case[bool]{
//	    {In: []bool{false, false}, Out: true},
//	    {In: []bool{false, true}, Out: true},
//	    {In: []bool{true, false}, Out: true},
//	    {In: []bool{true, true}, Out: false},
//	}
//	prob, _ := problems.NewBoolProblem(cases)
//	scoringFunc, _ := prob.NumHitsScoringFunc(1000.0)
//	// scoringFunc can be passed to evolution.New as its scoringFunc argument.
package problems
