package main

import (
	"github.com/qawatake/faccess"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(faccess.Analyzer) }
