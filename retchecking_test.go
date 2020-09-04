package retchecking_test

import (
	"testing"

	"github.com/BOBO1997/retchecking"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, retchecking.Analyzer, "b")
}
