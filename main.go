package main

import (
	"github.com/nikolaydubina/go-lint-cascade/analysis/cascade"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(cascade.Analyzer) }
