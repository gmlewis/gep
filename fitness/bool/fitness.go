// Package fitness provides common fitness functions.
package fitness

import "errors"

// BoolFunc ...
type BoolFunc func(predicted, target []bool) (float64, error)

// LengthErr is returned when the predicted and target slices are empty or not the same length.
var LengthErr = errors.New("length error")

// NumHits returns a fitness function that counts the number of correct values.
// The maximum fitness possible is N*scaleFactor, where N=len(target).
func NumHits(scaleFactor float64) (BoolFunc, error) {
	return func(predicted, target []bool) (float64, error) {
		if len(predicted) == 0 || len(target) == 0 || len(predicted) != len(target) {
			return 0, LengthErr
		}
		result := 0.0
		for i, t := range target {
			if predicted[i] == t {
				result += 1
			}
		}
		return result * scaleFactor, nil
	}, nil
}
