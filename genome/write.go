package genome

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"go/format"

	"github.com/gmlewis/gep/functions"
	"github.com/gmlewis/gep/grammars"
)

type dump struct {
	w      io.Writer
	gr     *grammars.Grammar
	fm     functions.FuncMap
	genome *Genome
	subs   map[string]string
}

func (g *Genome) Write(w io.Writer, grammar *grammars.Grammar) {
	d := &dump{
		gr:     grammar,
		genome: g,
		subs: map[string]string{
			"CHARX": "X",
		},
	}
	code, err := d.generateCode()
	if err != nil {
		fmt.Printf("error generating code: %v", err)
	}
	fmt.Fprintf(w, "%s", code)
}

func (d *dump) generateCode() ([]byte, error) {
	var buf bytes.Buffer
	d.w = &buf
	// d.write("// GML: d.gr.Open\n")
	d.write(d.gr.Open)
	for _, h := range d.gr.Headers {
		if h.Type != "default" {
			continue
		}
		// d.write(fmt.Sprintf("// GML: d.gr.Headers: h=%#v\n", h))
		d.write(h.Chardata)
		d.write(d.gr.Endline)
	}
	for _, t := range d.gr.Tempvars {
		if t.Type != "default" {
			continue
		}
		// d.write(fmt.Sprintf("// GML: d.gr.Tempvars: t=%#v\n", t))
		d.write(t.Chardata)
		d.subs["tempvarname"] = t.Varname
		d.write(d.gr.Endline)
	}
	// Generate the expression, keeping track of any helper functions that are needed.
	helpers := make(grammars.HelperMap)
	s, ok := d.gr.Functions.FuncMap[d.genome.LinkFunc]
	if !ok {
		return nil, fmt.Errorf("unable to find grammar linking function: %v\n", s.Symbol())
	}
	glf, ok := s.(*grammars.Function)
	exps := []string{""}
	for i, e := range d.genome.Genes {
		// d.write(fmt.Sprintf("// GML: d.genome.Genes: e=%#v\n", e))
		exp, err := e.Expression(d.gr, helpers)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			// d.write(fmt.Sprintf("// GML: len(d.genome.Genes)=%v\n", len(d.genome.Genes)))
			merge := strings.Replace(glf.Uniontype, "{tempvarname}", d.subs["tempvarname"], -1)
			merge = strings.Replace(merge, "{member}", exp, -1)
			merge = strings.Replace(merge, "{symbol}", glf.SymbolName, -1)
			exps = append(exps, merge)
		} else {
			// d.write(fmt.Sprintf("// GML: len(d.genome.Genes)=%v\n", len(d.genome.Genes)))
			exps = append(exps, d.subs["tempvarname"]+" = "+exp)
		}
	}
	exps = append(exps, "") // blank line
	fmt.Fprintln(d.w, strings.Join(exps, "\n"))
	for _, f := range d.gr.Footers {
		if f.Type != "default" {
			continue
		}
		// d.write(fmt.Sprintf("// GML: d.gr.Footers=%#v\n", f))
		d.write(f.Chardata)
		d.write(d.gr.Endline)
	}
	if len(helpers) > 0 { // Write out the helpers
		keys := make([]string, 0, len(helpers))
		for k, _ := range helpers {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			d.write(d.gr.Endline)
			d.write(helpers[k])
		}
	}

	clean, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.Bytes(), err
	}
	return clean, nil
}

func (d *dump) write(s string) {
	s = strings.Replace(s, "{CRLF}", "\n", -1)
	s = strings.Replace(s, "{TAB}", "\t", -1)
	for k, v := range d.subs {
		s = strings.Replace(s, fmt.Sprintf("{%v}", k), v, -1)
	}
	fmt.Fprint(d.w, s)
}
