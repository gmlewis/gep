package fitness

import "testing"

var predicted = []bool{false, false, false, true, true, false, true, true}
var target = []bool{false, false, false, false, true, true, true, true}

func TestNumHits(t *testing.T) {
	f, err := NumHits(1)
	if err != nil {
		t.Errorf("NumHits: got NewHits error %v, want nil", err)
	}
	got, err := f(predicted, target)
	if err != nil {
		t.Errorf("NumHits: got error %v, want nil", err)
	}
	if want := 6.0; got != want {
		t.Errorf("NumHits: got result %v, want %v", got, want)
	}
}
