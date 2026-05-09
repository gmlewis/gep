// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package codegen_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/codegen"
	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/grammars"
)

func mustGene(t *testing.T, karva string, funcType functions.FuncType) *gene.Gene {
	t.Helper()
	g, err := gene.New(karva, funcType)
	if err != nil {
		t.Fatalf("gene.New(%q) error: %v", karva, err)
	}
	return g
}

func TestExpression_RendersConstants(t *testing.T) {
	grammar, err := grammars.LoadGoMathGrammar()
	if err != nil {
		t.Fatalf("LoadGoMathGrammar(): %v", err)
	}

	got, err := codegen.Expression(
		[]string{"+", "d0", "c1"},
		[]float64{0.1, 0.5},
		3,
		[][]int{{1, 2}, nil, nil},
		grammar,
		make(grammars.HelperMap),
	)
	if err != nil {
		t.Fatalf("Expression(): %v", err)
	}
	if got != "(d[0]+0.5)" {
		t.Fatalf("Expression() = %q, want %q", got, "(d[0]+0.5)")
	}
}

func TestGenerate_UsesHelperFunctions(t *testing.T) {
	grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		t.Fatalf("LoadGoBooleanAllGatesGrammar(): %v", err)
	}

	code, err := codegen.Generate(codegen.Program{
		Genes: []codegen.Expressor{
			mustGene(t, "Nand.d0.d1", functions.Bool),
			mustGene(t, "Nor.d0.d1", functions.Bool),
		},
		LinkFunc: "And",
	}, grammar)
	if err != nil {
		t.Fatalf("Generate(): %v", err)
	}

	got := string(code)
	for _, want := range []string{
		"func gepModel(d []bool) bool {",
		"y = gepNand(d[0], d[1])",
		"y = y && gepNor(d[0], d[1])",
		"func gepNand(x, y bool) bool {",
		"func gepNor(x, y bool) bool {",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("Generate() missing %q in:\n%s", want, got)
		}
	}
}

func TestWrite_MissingLinkFunctionReturnsError(t *testing.T) {
	grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		t.Fatalf("LoadGoBooleanAllGatesGrammar(): %v", err)
	}

	var buf bytes.Buffer
	err = codegen.Write(&buf, codegen.Program{
		Genes:    []codegen.Expressor{mustGene(t, "d0", functions.Bool)},
		LinkFunc: "MissingLink",
	}, grammar)
	if err == nil {
		t.Fatalf("Write() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "MissingLink") {
		t.Fatalf("Write() error = %q, want mention of missing link function", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("Write() wrote %q despite error", buf.String())
	}
}
