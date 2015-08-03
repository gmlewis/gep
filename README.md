# Gene Expression Programming (GEP) in Go

This is an independent implementation of the Gene Expression Programming (GEP)
machine learning algorithm created by Dr. Cândida Ferreira.
It was written in the Go programming language by Glenn Lewis (gmlewis@google.com).

For more information, please visit http://www.gepsoft.com/

Here is a concise summary of GEP:

- http://www.gene-expression-programming.com/webpapers/GEP.pdf

The definitive book about GEP is available here:

- http://www.springer.com/3-540-32796-7

The lightning talk I gave about GEP at GopherCon 2014 is available here:

- http://go-talks.appspot.com/github.com/gmlewis/gep/gophercon2014/gep-in-go/gep-in-go.slide#1

## Status
[![GoDoc](https://godoc.org/github.com/gmlewis/gep?status.svg)](https://godoc.org/github.com/gmlewis/gep)
[![Build Status](https://travis-ci.org/gmlewis/gep.png)](https://travis-ci.org/gmlewis/gep)

----------------------------------------------------------------------

# NAND Function Experiment

To build and run this code, it may help to understand this presentation,
specifically about Go workspaces: http://talks.golang.org/2012/tutorial.slide#9

For example, to set up your GOPATH (once only):

```
$ export GOPATH=$HOME/go
$ mkdir -p $GOPATH/src
```

To run the NAND gate GEP experiment:

```
$ go get github.com/gmlewis/gep/experiments/nand
$ $GOPATH/bin/nand
Stopping after generation #0

// gepModel is auto-generated Go source code for the
// nand solution karva expression:
// "Not.And.Or.And.Or.And.d1.d0.d1.d1.d0.d0.d1.d1.d0", score=1000
package gepModel

func gepModel(d []bool) bool {
	y := false

	y = (!(((d[1] || d[1]) || (d[0] && d[0])) && (d[1] && d[0])))

	return y
}
```

# Symbolic Regression Experiment

To run the Symbolic Regression experiment:

```
$ go get github.com/gmlewis/gep/experiments/symbolic_regression
$ $GOPATH/bin/symbolic_regression
Stopping after generation #86

// gepModel is auto-generated Go source code for the
// (a^4 + a^3 + a^2 + a) solution karva expression:
// "*.d0.d0.d0.d0.d0.d0.*.d0.d0.d0.d0.d0.d0.d0.d0.d0|+|*.*.d0.d0.d0.d0.d0.*.d0.d0.d0.d0.d0.d0.d0.d0.d0|+|d0.d0.d0.*.*.d0.*.*.d0.d0.d0.d0.d0.d0.d0.d0.d0|+|*.*.*.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0.d0", score=11965.591435001414
package gepModel

import (
	"math"
)

func gepModel(d []float64) float64 {
	y := 0.0

	y = (d[0] * d[0])
	y += ((d[0] * d[0]) * d[0])
	y += d[0]
	y += ((d[0] * d[0]) * (d[0] * d[0]))

	return y
}
```

----------------------------------------------------------------------

# Unit Tests

To run unit tests, type:

```
$ go test github.com/gmlewis/gep/...
```

----------------------------------------------------------------------

# Grammars

While converting the C++ grammars to Go grammars, it was useful to load
in the XML files, parse them, and then dump them out to compare input
versus output.  This helped to weed out errors.

For example:

```
$ go get github.com/gmlewis/gep/experiments/load_grammars
$ $GOPATH/bin/load_grammars
```

----------------------------------------------------------------------

Enjoy!

----------------------------------------------------------------------

# License

Copyright 2014 Glenn Lewis gmlewis@google.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
