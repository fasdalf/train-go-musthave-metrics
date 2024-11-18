package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOSExitAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OSExitAnalyzer, "./os.exit/main")
	analysistest.Run(t, analysistest.TestData(), OSExitAnalyzer, "./os.exit/nonmain")
}
