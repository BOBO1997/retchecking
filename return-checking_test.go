package return-checking_test

import (
	"testing"

	"github.com/BOBO1997/return-checking"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, return-checking.Analyzer, "a")
}

