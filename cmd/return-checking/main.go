package main

import (
	"github.com/BOBO1997/return-checking"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(return-checking.Analyzer) }

