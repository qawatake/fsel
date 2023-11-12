package main

import (
	"github.com/qawatake/fsel"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(fsel.Analyzer) }
