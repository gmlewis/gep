// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package genome

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/gmlewis/gep/v2/codegen"
	"github.com/gmlewis/gep/v2/functions"
	"github.com/gmlewis/gep/v2/gene"
	"github.com/gmlewis/gep/v2/grammars"
)

func TestWriteNand(t *testing.T) {
	want := `package gepModel

func gepModel(d []bool) bool {
	var y bool

	y = (((!(d[0] && d[1])) && (d[0] || d[1])) || (!(d[1] && d[1])))

	return y
}
`

	g1 := mustGene(t, "Or.And.Not.Not.Or.And.And.d0.d1.d1.d1.d0.d1.d1.d0", functions.Bool)
	gn := New([]*gene.Gene{g1}, "Or")
	grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		t.Fatalf("unable to LoadGoBooleanAllGatesGrammar(): %v", err)
	}

	b := new(bytes.Buffer)
	if err := gn.Write(b, grammar); err != nil {
		t.Fatalf("gn.Write() error: %v", err)
	}
	if b.String() != want {
		t.Errorf("gen.Write() got:\n%v\nwant:\n%v", b.String(), want)
	}
}

func TestWriteMath(t *testing.T) {
	want := `package gepModel

import (
	"math"
)

func gepModel(d []float64) float64 {
	var y float64

	y = (d[0] * d[0])
	y += d[0]
	y += (d[0] * (d[0] * d[0]))
	y += (((d[0] * d[0]) * d[0]) * d[0])

	return y
}
`

	g1 := mustGene(t, "*.d0.d0.*.d0.*.*.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0", functions.Float64)
	g2 := mustGene(t, "d0.*.d0.*.*.d0.d0.*.d0.d0.d0.d0.d0.d0.d0.d0.d0", functions.Float64)
	g3 := mustGene(t, "*.d0.*.d0.d0.*.d0.*.d0.d0.d0.d0.d0.d0.d0.d0.d0", functions.Float64)
	g4 := mustGene(t, "*.*.d0.*.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0", functions.Float64)
	gn := New([]*gene.Gene{g1, g2, g3, g4}, "+")
	grammar, err := grammars.LoadGoMathGrammar()
	if err != nil {
		t.Fatalf("unable to LoadGoMathGrammar(): %v", err)
	}

	b := new(bytes.Buffer)
	if err := gn.Write(b, grammar); err != nil {
		t.Fatalf("gn.Write() error: %v", err)
	}
	if b.String() != want {
		t.Errorf("gen.Write() got:\n%v\nwant:\n%v", b.String(), want)
	}
}

func TestWrite6Multiplier(t *testing.T) {
	want := `package gepModel

func gepModel(d []bool) bool {
	var y bool

	y = (d[3] || ((gepNor(d[1], d[3]) || (d[2] && d[0])) || gepNand((d[1] || d[0]), d[2])))
	y = y && gepNand((gepNor((d[5] || d[5]), gepNand(d[1], d[0])) || gepNor((d[2] || d[0]), d[3])), d[0])
	y = y && gepNor((((d[1] && d[2]) && d[4]) && gepNor(d[0], d[4])), gepNor((d[1] || d[4]), gepNor(d[4], d[0])))
	y = y && (gepNor(gepNor(gepNand(d[2], d[4]), d[1]), (!(d[2]))) || gepNand(gepNand(d[1], d[3]), gepNor(d[0], d[2])))

	return y
}

func gepNand(x, y bool) bool {
	return (!(x && y))
}

func gepNor(x, y bool) bool {
	return (!(x || y))
}
`

	g1 := mustGene(t, "Or.d3.Or.Or.Nand.Nor.And.Or.d2.d1.d3.d2.d0.d1.d0.d0.d0", functions.Bool)
	g2 := mustGene(t, "Nand.Or.d0.Nor.Nor.Or.Nand.Or.d3.d5.d5.d1.d0.d2.d0.d1.d3", functions.Bool)
	g3 := mustGene(t, "Nor.And.Nor.And.Nor.Or.Nor.And.d4.d0.d4.d1.d4.d4.d0.d1.d2", functions.Bool)
	g4 := mustGene(t, "Or.Nor.Nand.Nor.Not.Nand.Nor.Nand.d1.d2.d1.d3.d0.d2.d2.d4.d4", functions.Bool)
	gn := New([]*gene.Gene{g1, g2, g3, g4}, "And")
	grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		t.Fatalf("unable to LoadGoBooleanAllGatesGrammar(): %v", err)
	}

	b := new(bytes.Buffer)
	if err := gn.Write(b, grammar); err != nil {
		t.Fatalf("gn.Write() error: %v", err)
	}
	if b.String() != want {
		t.Errorf("gen.Write() got:\n%v\nwant:\n%v", b.String(), want)
	}
}

func TestGenerateCode_MissingLinkFunctionReturnsError(t *testing.T) {
	g1 := mustGene(t, "d0", functions.Bool)
	gn := New([]*gene.Gene{g1}, "MissingLink")
	grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		t.Fatalf("unable to LoadGoBooleanAllGatesGrammar(): %v", err)
	}

	genes := make([]codegen.Expressor, len(gn.Genes))
	for i, gene := range gn.Genes {
		genes[i] = gene
	}
	_, err = codegen.Generate(codegen.Program{
		Genes:    genes,
		LinkFunc: gn.LinkFunc,
	}, grammar)
	if err == nil {
		t.Fatalf("codegen.Generate() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "MissingLink") {
		t.Fatalf("codegen.Generate() error = %q, want mention of missing link function", err)
	}
}

func TestWrite_MissingLinkFunctionReturnsError(t *testing.T) {
	g1 := mustGene(t, "d0", functions.Bool)
	gn := New([]*gene.Gene{g1}, "MissingLink")
	grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		t.Fatalf("unable to LoadGoBooleanAllGatesGrammar(): %v", err)
	}

	var b bytes.Buffer
	if err := gn.Write(&b, grammar); err == nil {
		t.Fatalf("Write() error = nil, want non-nil")
	}
}

func TestWrite_MissingLinkFunctionDoesNotWriteStdout(t *testing.T) {
	g1 := mustGene(t, "d0", functions.Bool)
	gn := New([]*gene.Gene{g1}, "MissingLink")
	grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
	if err != nil {
		t.Fatalf("unable to LoadGoBooleanAllGatesGrammar(): %v", err)
	}

	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	var b bytes.Buffer
	if err := gn.Write(&b, grammar); err == nil {
		t.Fatalf("Write() error = nil, want non-nil")
	}

	if err := w.Close(); err != nil {
		t.Fatalf("stdout close error: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("stdout read error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("Write() wrote to stdout: %q", string(out))
	}
}
