// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// load_grammars loads the XML representation of a GEP grammar for a target language and prints it out again.
// This program is useful to determine if the grammar is getting loaded and saved properly.
package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/gmlewis/gep/v2/grammars"
)

func main() {
	grammar, err := grammars.LoadGoMathGrammar()
	// grammar, err := grammars.LoadGoBooleanAllGatesGrammar()
	// grammar, err := grammars.LoadGoBooleanNotAndOrGatesGrammar()
	// grammar, err := grammars.LoadGoBooleanNandGatesGrammar()
	// grammar, err := grammars.LoadGoBooleanNorGatesGrammar()
	// grammar, err := grammars.LoadGoBooleanMuxSystemGrammar()
	// grammar, err := grammars.LoadGoReedMullerSystemGrammar()
	if err != nil {
		log.Fatalf("unable to load math grammar: %v", err)
	}

	// fmt.Printf("grammar:\n%#v\n", grammar)

	v, err := xml.MarshalIndent(grammar, "", "   ")
	if err != nil {
		log.Printf("unable to marshal math grammar: %v", err)
	}

	fmt.Println(string(v))
}
