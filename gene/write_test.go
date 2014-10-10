package gene

import (
	"testing"

	"github.com/gmlewis/gep/grammars"
)

func TestExpression(t *testing.T) {
	want := "(((!((d[0] && d[1]))) && (d[0] || d[1])) || (!((d[1] && d[1]))))"
	g := New("Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0")
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
