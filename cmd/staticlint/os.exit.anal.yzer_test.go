package main

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestOSExitAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OSExitAnalyzer, "./os.exit/main")
	analysistest.Run(t, analysistest.TestData(), OSExitAnalyzer, "./os.exit/nonmain")
}
