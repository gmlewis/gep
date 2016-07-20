// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package fitness provides common fitness functions.
package fitness

import "errors"

// BoolFunc ...
type BoolFunc func(predicted, target []bool) (float64, error)

// ErrLength is returned when the predicted and target slices are empty or not the same length.
var ErrLength = errors.New("length error")

// NumHits returns a fitness function that counts the number of correct values.
// The maximum fitness possible is N*scaleFactor, where N=len(target).
func NumHits(scaleFactor float64) (BoolFunc, error) {
	return func(predicted, target []bool) (float64, error) {
		if len(predicted) == 0 || len(target) == 0 || len(predicted) != len(target) {
			return 0, ErrLength
		}
		result := 0.0
		for i, t := range target {
			if predicted[i] == t {
				result++
			}
		}
		return result * scaleFactor, nil
	}, nil
}
