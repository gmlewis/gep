package gene

import (
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/grammars"
)

func TestBuildMathTreeRejectsEqualLengthSymbolIndex(t *testing.T) {
	g := &Gene{
		Symbols:   []string{"d0"},
		SymbolMap: map[string]int{},
	}

	f := g.buildMathTree(len(g.Symbols), nil)
	if got := f([]float64{42}); got != 0 {
		t.Fatalf("buildMathTree() returned %v, want 0", got)
	}
}

func TestBuildVectorIntTreeRejectsEqualLengthSymbolIndex(t *testing.T) {
	g := &Gene{
		Symbols:   []string{"d0"},
		SymbolMap: map[string]int{},
	}

	f := g.buildVectorIntTree(len(g.Symbols), nil)
	if got := f([]VectorInt{{1, 2, 3}}); len(got) != 0 {
		t.Fatalf("buildVectorIntTree() returned %v, want empty vector", got)
	}
}

func TestBuildExpRejectsEqualLengthSymbolIndex(t *testing.T) {
	g := &Gene{Symbols: []string{"d0"}}

	_, err := g.buildExp(len(g.Symbols), nil, &grammars.Grammar{}, nil)
	if err == nil || !strings.Contains(err.Error(), "bad symbolIndex") {
		t.Fatalf("buildExp() error = %v, want bad symbolIndex error", err)
	}
}

func TestBuildExpRejectsOutOfRangeTerminalSymbol(t *testing.T) {
	g := &Gene{
		Symbols:      []string{"d1"},
		numTerminals: 1,
	}

	_, err := g.buildExp(0, nil, &grammars.Grammar{}, nil)
	if err == nil || !strings.Contains(err.Error(), "exceeds number of terminals") {
		t.Fatalf("buildExp() error = %v, want terminal bounds error", err)
	}
}

func TestBuildExpRejectsOutOfRangeConstantSymbol(t *testing.T) {
	g := &Gene{
		Symbols:      []string{"c1"},
		Constants:    []float64{3.14},
		numTerminals: 1,
	}

	_, err := g.buildExp(0, nil, &grammars.Grammar{}, nil)
	if err == nil || !strings.Contains(err.Error(), "exceeds length of constant slice") {
		t.Fatalf("buildExp() error = %v, want constant bounds error", err)
	}
}
