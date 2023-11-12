package main

import (
	"github.com/qawatake/fderef"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(fderef.Analyzer) }
