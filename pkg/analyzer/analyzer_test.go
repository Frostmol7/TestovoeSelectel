package analyzer_test

import (
    "testing"

    "golang.org/x/tools/go/analysis/analysistest"
    "TestovoeSelectel/pkg/analyzer"
)

func TestAnalyzer(t *testing.T) {
    analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "testdata.go")
}
