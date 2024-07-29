package fsel_test

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/qawatake/fsel"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	t.Parallel()
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, fsel.Analyzer, "a")
}
