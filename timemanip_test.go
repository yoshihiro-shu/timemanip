package timemanip_test

import (
	"testing"

	"github.com/yoshihiro-shu/timemanip"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, timemanip.Analyzer, "a")
}
