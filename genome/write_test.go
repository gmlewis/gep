package genome

import (
	"bytes"
	"testing"

	"github.com/gmlewis/gep/gene"
	"github.com/gmlewis/gep/grammars"
)

func TestWriteNand(t *testing.T) {
	want := `package gepModel

func gepModel(d []bool) bool {
	blTemp := false

	blTemp = (((!(d[0] && d[1])) && (d[0] || d[1])) || (!(d[1] && d[1])))

	return blTemp
}
`

	g1 := gene.New("Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0")
	gn := New([]gene.Gene{g1}, "Or")
	grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		t.Fatalf("unable to LoadGoBooleanAllGatesGrammar(): %v", err)
	}

	b := new(bytes.Buffer)
	gn.Write(b, grammar)
	if b.String() != want {
		t.Errorf("gen.Write() got %v, want %v", b.String(), want)
	}
}
