package cascade_test

import (
	"testing"

	"github.com/nikolaydubina/go-lint-cascade/analysis/cascade"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, cascade.Analyzer, "a")
}
