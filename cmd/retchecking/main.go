package main

import (
	"github.com/BOBO1997/retchecking"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(retchecking.Analyzer) }
