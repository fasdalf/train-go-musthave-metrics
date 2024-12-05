package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	// определяем map подключаемых правил
	checks := map[string]bool{
		"S1001": true,
	}
	mychecks := []*analysis.Analyzer{
		//ErrCheckAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
		OSExitAnalyzer,
	}

	for _, v := range staticcheck.Analyzers {
		// добавляем в массив "нужные" проверки
		mychecks = append(mychecks, v.Analyzer)
	}

	for _, v := range simple.Analyzers {
		// добавляем в массив "нужные" проверки
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	multichecker.Main(
		mychecks...,
	)
}
