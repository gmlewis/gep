// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gene

import (
	"testing"

	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/grammars"
)

func TestExpression_Bool(t *testing.T) {
	want := "(((!((d[0] && d[1]))) && (d[0] || d[1])) || (!((d[1] && d[1]))))"
	g := New("Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0", functions.Bool)
	grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		t.Fatalf("unable to LoadGoBooleanAllGatesGrammar(): %v", err)
	}

	helpers := make(grammars.HelperMap)
	got, err := g.Expression(grammar, helpers)
	if err != nil {
		t.Fatalf("g.Expression error: %v", err)
	}

	if got != want {
		t.Errorf("g.Expression got %q, want %q", got, want)
	}

	if len(helpers) != 0 {
		t.Errorf("helpers got length %v, want 0", len(helpers))
	}
}

func TestConstants_Float64(t *testing.T) {
	want := "(d[0]+0.50)"
	g := New("+.d0.c1.+.+.+.+.d0.d1.d1.d1.d0.d1.d1.d0", functions.Float64)
	g.Constants = []float64{0.1, 0.5}
	grammar, err := grammars.LoadGoMathGrammar()
	if err != nil {
		t.Fatalf("unable to LoadGoMathGrammar(): %v", err)
	}

	helpers := make(grammars.HelperMap)
	got, err := g.Expression(grammar, helpers)
	if err != nil {
		t.Fatalf("g.Expression error: %v", err)
	}

	if got != want {
		t.Errorf("g.Expression got %q, want %q", got, want)
	}

	if len(helpers) != 0 {
		t.Errorf("helpers got length %v, want 0", len(helpers))
	}
}
