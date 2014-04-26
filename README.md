# Gene Expression Programming (GEP) in Go

This is an independent implementation of the Gene Expression Programming (GEP)
machine learning algorithm created by Dr. Cândida Ferreira.
It was written in the Go programming language by Glenn Lewis (gmlewis@google.com).

For more information, please visit http://www.gepsoft.com/

Here is a concise summary of GEP:

- http://www.gene-expression-programming.com/webpapers/GEP.pdf

For the full book written on the topic, it is available here:

- http://www.springer.com/3-540-32796-7

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
$ time $GOPATH/bin/nand
Stopping after generation #0
nand solution: Not.And.And.Or.d0.And.And.d0.d0.d1.d0.d0.d1.d1.d0, score=1000

real 0m0.022s
user 0m0.000s
sys  0m0.020s
```

# Symbolic Regression Experiment

To run the Symbolic Regression experiment:

```
$ go get github.com/gmlewis/gep/experiments/symbolic_regression
$ time $GOPATH/bin/symbolic_regression
Stopping after generation #2462
(a^4 + a^3 + a^2 + a) solution: *.+.+./.*.*.d0.d0.d0.d0.d0.d0.d0, score=9965.59143500096

real 0m0.186s
user 0m0.180s
sys  0m0.000s
```

----------------------------------------------------------------------

# Unit Tests

To run unit tests, you can do this:

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

Please note that the grammars here are not current... I need to
download the latest from Gepsoft and replace these.

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
