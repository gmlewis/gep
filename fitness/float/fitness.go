// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package fitness provides common fitness functions.
package fitness

import (
	"errors"
	"math"
)

// FloatFunc ...
type FloatFunc func(predicted, target []float64) (float64, error)

// ErrLength is returned when the predicted and target slices are empty or not the same length.
var ErrLength = errors.New("length error")

// NumHitsAbs returns a fitness function that favors models that perform well for all
// fitness cases within a certain absolute error (that is, the precision that is chosen
// for the evolved models - a number between 0 and 1, inclusive) of the correct value.
// The maximum fitness possible is N*scaleFactor, where N=len(target).
func NumHitsAbs(precision, scaleFactor float64) (FloatFunc, error) {
	if precision < 0 || precision > 1 {
		return nil, errors.New("invalid precision, should be 0-1")
	}
	return func(predicted, target []float64) (float64, error) {
		if len(predicted) == 0 || len(target) == 0 || len(predicted) != len(target) {
			return 0, ErrLength
		}
		result := 0.0
		for i, t := range target {
			e := math.Abs(predicted[i] - t)
			if e <= precision {
				result++
			}
		}
		return result * scaleFactor, nil
	}, nil
}

// NumHitsRel returns a fitness function that favors models that perform well for all
// fitness cases within a certain relative error (that is, the precision that is chosen
// for the evolved models - a number between 0 and 1, inclusive) of the correct value.
// The maximum fitness possible is N*scaleFactor, where N=len(target).
func NumHitsRel(precision, scaleFactor float64) (FloatFunc, error) {
	if precision < 0 || precision > 1 {
		return nil, errors.New("invalid precision, should be 0-1")
	}
	return func(predicted, target []float64) (float64, error) {
		if len(predicted) == 0 || len(target) == 0 || len(predicted) != len(target) {
			return 0, ErrLength
		}
		result := 0.0
		for i, t := range target {
			if t == 0 {
				if predicted[i] == 0 {
					result++
				}
				continue
			}
			e := math.Abs((predicted[i] - t) / t)
			if e <= precision {
				result++
			}
		}
		return result * scaleFactor, nil
	}, nil
}

// SelectionRangeAbs returns a fitness function that is used as a limit for
// selection to operate, above which the performance of a program on a particular
// fitness case contributes nothing to its fitness. The precision is the
// limit for improvement as it allows the fine-tuning of the evolved solutions as
// accurately as possible.
// The maximum fitness possible is N*selectionRange*scaleFactor, where N=len(target).
func SelectionRangeAbs(selectionRange, scaleFactor float64) (FloatFunc, error) {
	return func(predicted, target []float64) (float64, error) {
		if len(predicted) == 0 || len(target) == 0 || len(predicted) != len(target) {
			return 0, ErrLength
		}
		result := 0.0
		for i, t := range target {
			e := math.Abs(predicted[i] - t)
			result += (selectionRange - e)
		}
		return result * scaleFactor, nil
	}, nil
}

// SelectionRangeRel returns a fitness function that is used as a limit for
// selection to operate, above which the performance of a program on a particular
// fitness case contributes nothing to its fitness. The precision is the
// limit for improvement as it allows the fine-tuning of the evolved solutions as
// accurately as possible. The error is calculated relative to the target value.
// The maximum fitness possible is N*selectionRange*scaleFactor, where N=len(target).
func SelectionRangeRel(selectionRange, scaleFactor float64) (FloatFunc, error) {
	return func(predicted, target []float64) (float64, error) {
		if len(predicted) == 0 || len(target) == 0 || len(predicted) != len(target) {
			return 0, ErrLength
		}
		result := 0.0
		for i, t := range target {
			if t == 0 {
				if predicted[i] == 0 {
					result += selectionRange
				}
				continue
			}
			e := math.Abs((predicted[i] - t) / t)
			result += (selectionRange - e)
		}
		return result * scaleFactor, nil
	}, nil
}

// MeanSquaredErrorAbs returns a fitness function that calculates the mean square
// error and is normalized from 0 to scaleFactor.
func MeanSquaredErrorAbs(scaleFactor float64) (FloatFunc, error) {
	return func(predicted, target []float64) (float64, error) {
		if len(predicted) == 0 || len(target) == 0 || len(predicted) != len(target) {
			return 0, ErrLength
		}
		result := 0.0
		for i, t := range target {
			e := predicted[i] - t
			result += (e * e)
		}
		result = scaleFactor / (1 + result/float64(len(predicted)))
		return result, nil
	}, nil
}

// MeanSquaredErrorAbsRoot returns a fitness function that calculates the square root
// of the mean square error and is normalized from 0 to scaleFactor.
func MeanSquaredErrorAbsRoot(scaleFactor float64) (FloatFunc, error) {
	return func(predicted, target []float64) (float64, error) {
		if len(predicted) == 0 || len(target) == 0 || len(predicted) != len(target) {
			return 0, ErrLength
		}
		result := 0.0
		for i, t := range target {
			e := predicted[i] - t
			result += (e * e)
		}
		result = scaleFactor / (1 + math.Sqrt(result)/float64(len(predicted)))
		return result, nil
	}, nil
}

// MeanSquaredErrorRelRoot returns a fitness function that calculates the square root
// of the mean square relative error and is normalized from 0 to scaleFactor.
func MeanSquaredErrorRelRoot(scaleFactor float64) (FloatFunc, error) {
	return func(predicted, target []float64) (float64, error) {
		if len(predicted) == 0 || len(target) == 0 || len(predicted) != len(target) {
			return 0, ErrLength
		}
		result := 0.0
		for i, t := range target {
			if t == 0 {
				result += (predicted[i] * predicted[i])
				continue
			}
			e := (predicted[i] - t) / t
			result += (e * e)
		}
		result = scaleFactor / (1 + math.Sqrt(result)/float64(len(predicted)))
		return result, nil
	}, nil
}

// MeanSquaredErrorRel returns a fitness function that calculates the mean square
// relative error and is normalized from 0 to scaleFactor.
func MeanSquaredErrorRel(scaleFactor float64) (FloatFunc, error) {
	return func(predicted, target []float64) (float64, error) {
		if len(predicted) == 0 || len(target) == 0 || len(predicted) != len(target) {
			return 0, ErrLength
		}
		result := 0.0
		for i, t := range target {
			if t == 0 {
				result += (predicted[i] * predicted[i])
				continue
			}
			e := (predicted[i] - t) / t
			result += (e * e)
		}
		result = scaleFactor / (1 + result/float64(len(predicted)))
		return result, nil
	}, nil
}

// RSquare returns a fitness function that is based on the standard R-square, which returns
// the square of the Pearson product moment correlation coefficient. The return value
// is normalized from 0 to scaleFactor.
func RSquare(scaleFactor float64) (FloatFunc, error) {
	return func(predicted, target []float64) (float64, error) {
		if len(predicted) == 0 || len(target) == 0 || len(predicted) != len(target) {
			return 0, ErrLength
		}
		var sumP, sumT, sumPP, sumTP, sumTT float64
		for i, t := range target {
			sumP += predicted[i]
			sumT += t
			sumPP += (predicted[i] * predicted[i])
			sumTP += (predicted[i] * t)
			sumTT += (t * t)
		}
		n := float64(len(target))
		d1 := n*sumTT - sumT*sumT
		d2 := n*sumPP - sumP*sumP
		d := math.Sqrt(d1 * d2)
		if d == 0 {
			return 0, nil
		}
		n1 := n * sumTP
		n2 := sumT * sumP
		r := (n1 - n2) / d
		result := scaleFactor * r * r
		return result, nil
	}, nil
}
