// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package grammars parses XML GEP grammar representations for a particular output language.
// It can then be used to generate code from Karva expressions for that language.
package grammars

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gmlewis/gep/v2/functions"
)

// Functions is a collection of Functions available in the language grammar.
type Functions struct {
	Count     int        `xml:"count,attr"`
	Functions []Function `xml:"function"`

	// Lookup table of function symbol name to function definition
	FuncMap functions.FuncMap `xml:"-"`
}

// Function represents a single function.
type Function struct {
	// Idx is the index of the function in the XML grammar.
	Idx int `xml:"idx,attr"`
	// SymbolName is the Karva symbol used to represent the function.
	SymbolName string `xml:"symbol,attr"`
	// TerminalCount specifies the number of input terminals to the function.
	TerminalCount int `xml:"terminals,attr"`
	// Uniontype (optional) determines how the diadic function is rendered.
	Uniontype string `xml:"uniontype,attr"`
	// Chardata is the action function rendering for the given language.
	Chardata string `xml:",chardata"`
}

// Symbol returns the symbol name of the function.
func (f *Function) Symbol() string {
	return f.SymbolName
}

// Terminals returns the terminal count of the function.
func (f *Function) Terminals() int {
	return f.TerminalCount
}

// BoolFunction allows FuncMap to implement interace functions.FuncMap.
func (f *Function) BoolFunction([]bool) bool {
	return false
}

// IntFunction allows FuncMap to implement interace functions.FuncMap.
func (f *Function) IntFunction([]int) int {
	return 0
}

// Float64Function allows FuncMap to implement interace functions.FuncMap.
func (f *Function) Float64Function([]float64) float64 {
	return 0.0
}

// VectorIntFunction allows FuncMap to implement interace functions.FuncMap.
func (f *Function) VectorIntFunction([]functions.VectorInt) functions.VectorInt {
	return functions.VectorInt{}
}

// Replacement determines how a function can be replaced.
type Replacement struct {
	Type      string `xml:"type,attr"`
	Replace   string `xml:"replace,attr"`
	Indexzero bool   `xml:"indexzero,attr"`
	Indexone  bool   `xml:"indexone,attr"`
	Chardata  string `xml:",chardata"`
}

// Categories lists the types of grammar available.
type Categories struct {
	Functioncall      Functioncall      `xml:"functioncall"`
	TransformFunction TransformFunction `xml:"transformfunction"`
	Switch            Switch            `xml:"switch"`
	Case              Format            `xml:"case"`
	Equality          Format            `xml:"equality"`
}

// Functioncall represents a function call.
type Functioncall struct {
	Call string `xml:"call,attr"`
}

// TransformFunction represents the transformation of a function.
type TransformFunction struct {
	Header    string `xml:"header,attr"`
	Footer    string `xml:"footer,attr"`
	Prototype string `xml:"prototype,attr"`
}

// Switch represents the grammar for a switch call.
type Switch struct {
	Special            string `xml:"special,attr"`
	Top                string `xml:"top,attr"`
	Bottom             string `xml:"bottom,attr"`
	Categoricaldefault string `xml:"categoricaldefault,attr"`
	Numericaldefault   string `xml:"numericaldefault,attr"`
}

// Format represents how an item is formatted.
type Format struct {
	Format string `xml:"format,attr"`
}

// Transformation represents a function transformation.
type Transformation struct {
	Name         string `xml:"name,attr"`
	Call         string `xml:"call,attr"`
	Itemformat   string `xml:"itemformat,attr"`
	Prototype    string `xml:"prototype,attr"`
	Declarations string `xml:"declarations,attr"`
	Chardata     string `xml:",chardata"`
}

// Constant represents a constant.
type Constant struct {
	Type       string `xml:"type,attr"`
	Replace    string `xml:"replace,attr"`
	Labelindex int    `xml:"labelindex,attr"`
	Chardata   string `xml:",chardata"`
}

// Tempvar represents a temporary variable.
type Tempvar struct {
	Type     string `xml:"type,attr"`
	Typename string `xml:"typename,attr"`
	Varname  string `xml:"varname,attr"`
	Chardata string `xml:",chardata"`
}

// Helper represents a helper function.
type Helper struct {
	Replaces  string `xml:"replaces,attr"`
	Prototype string `xml:"prototype,attr"`
	Chardata  string `xml:",chardata"`
}

// Testing represents a testing function.
type Testing struct {
	Prototype Prototype `xml:"prototype"`
	Method    Method    `xml:"method"`
}

// Prototype is a function signature.
type Prototype struct {
	Paramsformat string `xml:"paramsformat,attr"`
	Chardata     string `xml:",chardata"`
}

// Method represents a method.
type Method struct {
	Callformat string `xml:"callformat,attr"`
	Listformat string `xml:"listformat,attr"`
	Chardata   string `xml:",chardata"`
}

// OrderItem is a single element of the code structure.
type OrderItem struct {
	Name string `xml:"name,attr"`
}

// HelperMap maps helper replacement functions to their definitions.
type HelperMap map[string]string

// Helpers represents helper functions for the target language.
type Helpers struct {
	Count       int      `xml:"count,attr"`
	Declaration string   `xml:"declaration,attr"`
	Assignment  string   `xml:"assignment,attr"`
	Helpers     []Helper `xml:"helper"`

	// Lookup table of helper symbol name to helper definition
	HelperMap HelperMap `xml:"-"`
}

// Keyword is a reserved keyword in the target language.
type Keyword struct {
	Chardata string `xml:",chardata"`
}

// LinkingFunctions are special functions used by GEP for the target language.
type LinkingFunctions struct {
	Count            int      `xml:"count,attr"`
	LinkingFunctions []Helper `xml:"linkingFunction"`
}

// BasicFunctions are simple functions used by GEP for the target language.
type BasicFunctions struct {
	Count          int      `xml:"count,attr"`
	BasicFunctions []Helper `xml:"basicFunction"`
}

// Grammar represents the complete XML specification for rendering a Karva string for the given target language.
type Grammar struct {
	XMLName  xml.Name `xml:"grammar"`
	Comments string   `xml:",comment"`
	Name     string   `xml:"name,attr"`
	Version  string   `xml:"version,attr"`
	Ext      string   `xml:"ext,attr"`
	Type     string   `xml:"type,attr"`
	// Functions lists all possible functions available in the language grammar.
	Functions Functions `xml:"functions"`
	// Order specifies the order structure of the program.
	Order []OrderItem `xml:"order>item"`
	// Open contains the text for the start of the generated source code file.
	Open string `xml:"open"`
	// Close contains the text for the end of the generated source code file.
	Close string `xml:"close"`
	// Headers/Subheaders contains any required headers such as imports needed near the top of the generated source code file.
	Headers    []Replacement `xml:"headers>header"`
	Subheaders []Replacement `xml:"subheaders>subheader"`
	// RandomConstants defines the syntax for a random constant.
	RandomConstants []Replacement `xml:"randomconstants>randomconst"`
	// Categories defines the special features of the target language.
	Categories Categories `xml:"categories"`
	// Transformations/ReverseTransformations lists any special transformations for the target language.
	Transformations        []Transformation `xml:"transformations>transformation"`
	ReverseTransformations []Transformation `xml:"reversetransformations>transformation"`
	// Constants defines the syntax for constants for the target language.
	Constants []Constant `xml:"constants>constant"`
	// Tempvars defines the syntax for temporary variables for the target language.
	Tempvars []Tempvar `xml:"tempvars>tempvar"`
	// Endline specifies how the end of a line is rendered for the target language.
	Endline string `xml:"endline"`
	// Indent is the number of tab characters to add to each line in the code block.
	Indent int `xml:"indent"`
	// Parenstype specifies the type of parentheses to use for arrays in the target language.
	// 0=(), 1=[]
	Parenstype int `xml:"parenstype"`
	// Footers defines special code structures for the target language.
	Footers []Replacement `xml:"footers>footer"`
	// Helpers defines the complete set of helper functions used in the Functions above.
	Helpers Helpers `xml:"helpers"`
	// Keywords lists the keywords for the target language.
	Keywords []Keyword `xml:"keywords>keyword"`
	// Commentmark defines how a comment is rendered in the target language.
	Commentmark string `xml:"commentmark"`
	// LinkingFunctions lists special linking functions used by GEP for the target language.
	LinkingFunctions LinkingFunctions `xml:"linkingFunctions"`
	// BasicFunctions lists basic functions used by GEP for the target language.
	BasicFunctions BasicFunctions `xml:"basicFunctions"`
	Ddfcomment     string         `xml:"ddfcomment"`
	Udfcomment     string         `xml:"udfcomment"`
	// Testing defines a function used for testing a method for the target language.
	Testing Testing `xml:"testing"`
}

func loadGrammar(path string) (*Grammar, error) {
	v := &Grammar{}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("unable to read file %q: %q", path, err)
		return nil, err
	}

	err = xml.Unmarshal(data, &v)
	if err != nil {
		log.Printf("error unmarshaling %q: %q", path, err)
		return nil, err
	}

	// Build the function map lookups for fast access
	v.Functions.FuncMap = make(functions.FuncMap, len(v.Functions.Functions))
	for i, f := range v.Functions.Functions {
		v.Functions.FuncMap[f.SymbolName] = &v.Functions.Functions[i]
	}

	// Build the helpers map lookups for fast access
	v.Helpers.HelperMap = make(HelperMap, len(v.Helpers.Helpers))
	for _, h := range v.Helpers.Helpers {
		v.Helpers.HelperMap[h.Replaces] = h.Chardata
	}

	return v, nil
}

func getPath(filename string) string {
	_, packageFile, _, _ := runtime.Caller(0)
	grammarPath := filepath.Dir(packageFile)

	// Support Travis CI automated builds by searching for files
	dirs := strings.Split(os.Getenv("GOPATH"), ":")
	for _, dir := range dirs {
		name := filepath.Join(dir, "src", grammarPath, filename)
		if _, err := os.Stat(name); err == nil {
			return name
		}
	}
	name := filepath.Join(grammarPath, filename)
	if _, err := os.Stat(name); err == nil {
		return name
	}
	return filename
}

// LoadGoMathGrammar loads the floating-point math grammer for Go as the target language.
func LoadGoMathGrammar() (*Grammar, error) {
	path := getPath("go.Math.00.default.grm.xml")
	return loadGrammar(path)
}

// LoadGoBooleanAllGatesGrammar loads the general boolean grammer for Go as the target language.
func LoadGoBooleanAllGatesGrammar() (*Grammar, error) {
	path := getPath("go.Boolean.01.AllGates.grm.xml")
	return loadGrammar(path)
}

// LoadGoBooleanNotAndOrGatesGrammar loads the specialized boolean grammer for Go as the target language.
func LoadGoBooleanNotAndOrGatesGrammar() (*Grammar, error) {
	path := getPath("go.Boolean.02.NotAndOrGates.grm.xml")
	return loadGrammar(path)
}

// LoadGoBooleanNandGatesGrammar loads the specialized boolean grammer for Go as the target language.
func LoadGoBooleanNandGatesGrammar() (*Grammar, error) {
	path := getPath("go.Boolean.03.NandGates.grm.xml")
	return loadGrammar(path)
}

// LoadGoBooleanNorGatesGrammar loads the specialized boolean grammer for Go as the target language.
func LoadGoBooleanNorGatesGrammar() (*Grammar, error) {
	path := getPath("go.Boolean.04.NorGates.grm.xml")
	return loadGrammar(path)
}

// LoadGoBooleanMuxSystemGrammar loads the specialized boolean grammer for Go as the target language.
func LoadGoBooleanMuxSystemGrammar() (*Grammar, error) {
	path := getPath("go.Boolean.05.MuxSystem.grm.xml")
	return loadGrammar(path)
}

// LoadGoReedMullerSystemGrammar loads the specialized boolean grammer for Go as the target language.
func LoadGoReedMullerSystemGrammar() (*Grammar, error) {
	path := getPath("go.Boolean.06.ReedMullerSystem.grm.xml")
	return loadGrammar(path)
}
