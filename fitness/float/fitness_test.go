package fitness

import (
	"math"
	"testing"
)

var predicted = []float64{0, 0.25, 0.75, 1, 0, 0.25, 0.75, 1, 0, 0.25, 0.75, 1}
var target = []float64{0, 0, 0, 0, 0.5, 0.5, 0.5, 0.5, 1, 1, 1, 1}

func TestNumHitsAbs(t *testing.T) {
	tests := []struct{ precision, want float64 }{
		{precision: 0.0, want: 2},
		{precision: 0.25, want: 6},
		{precision: 0.5, want: 8},
		{precision: 0.75, want: 10},
		{precision: 1.0, want: 12},
	}
	for i, test := range tests {
		f, err := NumHitsAbs(test.precision, 1)
		if err != nil {
			t.Errorf("NumHitsAbs test %v: got NewHits error %v, want nil", i, err)
		}
		got, err := f(predicted, target)
		if err != nil {
			t.Errorf("NumHitsAbs test %v: got error %v, want nil", i, err)
		}
		if got != test.want {
			t.Errorf("NumHitsAbs test %v: got result %v, want %v", i, got, test.want)
		}
	}
}

func TestNumHitsRel(t *testing.T) {
	tests := []struct{ precision, want float64 }{
		{precision: 0.0, want: 2},
		{precision: 0.25, want: 3},
		{precision: 0.5, want: 5},
		{precision: 0.75, want: 6},
		{precision: 1.0, want: 9},
	}
	for i, test := range tests {
		f, err := NumHitsRel(test.precision, 1)
		if err != nil {
			t.Errorf("NumHitsRel test %v: got NewHits error %v, want nil", i, err)
		}
		got, err := f(predicted, target)
		if err != nil {
			t.Errorf("NumHitsRel test %v: got error %v, want nil", i, err)
		}
		if got != test.want {
			t.Errorf("NumHitsRel test %v: got result %v, want %v", i, got, test.want)
		}
	}
}

func TestSelectionRangeAbs(t *testing.T) {
	tests := []struct{ selectionRange, want float64 }{
		{selectionRange: 0.0, want: -5.5},
		{selectionRange: 0.25, want: -2.5},
		{selectionRange: 0.5, want: 0.5},
		{selectionRange: 0.75, want: 3.5},
		{selectionRange: 1.0, want: 6.5},
		{selectionRange: 2.0, want: 18.5},
	}
	for i, test := range tests {
		f, err := SelectionRangeAbs(test.selectionRange, 1)
		if err != nil {
			t.Errorf("SelectionRangeAbs %v: got NewHits error %v, want nil", i, err)
		}
		got, err := f(predicted, target)
		if err != nil {
			t.Errorf("SelectionRangeAbs %v: got error %v, want nil", i, err)
		}
		if got != test.want {
			t.Errorf("SelectionRangeAbs %v: got result %v, want %v", i, got, test.want)
		}
	}
}

func TestSelectionRangeRel(t *testing.T) {
	tests := []struct{ selectionRange, want float64 }{
		{selectionRange: 0.0, want: -5},
		{selectionRange: 0.25, want: -2.75},
		{selectionRange: 0.5, want: -0.5},
		{selectionRange: 0.75, want: 1.75},
		{selectionRange: 1.0, want: 4},
		{selectionRange: 2.0, want: 13},
	}
	for i, test := range tests {
		f, err := SelectionRangeRel(test.selectionRange, 1)
		if err != nil {
			t.Errorf("SelectionRangeRel %v: got NewHits error %v, want nil", i, err)
		}
		got, err := f(predicted, target)
		if err != nil {
			t.Errorf("SelectionRangeRel %v: got error %v, want nil", i, err)
		}
		if got != test.want {
			t.Errorf("SelectionRangeRel %v: got result %v, want %v", i, got, test.want)
		}
	}
}

func TestMeanSquaredErrorAbs(t *testing.T) {
	f, err := MeanSquaredErrorAbs(1000)
	if err != nil {
		t.Errorf("MeanSquaredErrorAbs: got NewHits error %v, want nil", err)
	}
	got, err := f(predicted, target)
	if err != nil {
		t.Errorf("MeanSquaredErrorAbs: got error %v, want nil", err)
	}
	if want := 755.9055; math.Abs(got-want) > 1e-4 {
		t.Errorf("MeanSquaredErrorAbs: got result %v, want %v", got, want)
	}
}

func TestMeanSquaredErrorAbsRoot(t *testing.T) {
	f, err := MeanSquaredErrorAbsRoot(1000)
	if err != nil {
		t.Errorf("MeanSquaredErrorAbs: got NewHits error %v, want nil", err)
	}
	got, err := f(predicted, target)
	if err != nil {
		t.Errorf("MeanSquaredErrorAbs: got error %v, want nil", err)
	}
	if want := 859.0757; math.Abs(got-want) > 1e-4 {
		t.Errorf("MeanSquaredErrorAbs: got result %v, want %v", got, want)
	}
}

func TestMeanSquaredErrorRel(t *testing.T) {
	f, err := MeanSquaredErrorRel(1000)
	if err != nil {
		t.Errorf("MeanSquaredErrorRel: got NewHits error %v, want nil", err)
	}
	got, err := f(predicted, target)
	if err != nil {
		t.Errorf("MeanSquaredErrorRel: got error %v, want nil", err)
	}
	if want := 676.0563; math.Abs(got-want) > 1e-4 {
		t.Errorf("MeanSquaredErrorRel: got result %v, want %v", got, want)
	}
}

func TestMeanSquaredErrorRelRoot(t *testing.T) {
	f, err := MeanSquaredErrorRelRoot(1000)
	if err != nil {
		t.Errorf("MeanSquaredErrorRel: got NewHits error %v, want nil", err)
	}
	got, err := f(predicted, target)
	if err != nil {
		t.Errorf("MeanSquaredErrorRel: got error %v, want nil", err)
	}
	if want := 833.4540; math.Abs(got-want) > 1e-4 {
		t.Errorf("MeanSquaredErrorRel: got result %v, want %v", got, want)
	}
}

func TestRSquare(t *testing.T) {
	target2 := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0, 1.1, 1.2}
	f, err := RSquare(1000)
	if err != nil {
		t.Errorf("RSquare: got NewHits error %v, want nil", err)
	}
	got, err := f(predicted, target2)
	if err != nil {
		t.Errorf("RSquare: got error %v, want nil", err)
	}
	if want := 102.7972; math.Abs(got-want) > 1e-4 {
		t.Errorf("RSquare: got result %v, want %v", got, want)
	}
}
