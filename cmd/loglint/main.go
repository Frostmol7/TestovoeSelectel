package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"TestovoeSelectel/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
