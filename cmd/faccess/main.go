package main

import (
	"github.com/qawatake/fderef"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(fderef.Analyzer) }
