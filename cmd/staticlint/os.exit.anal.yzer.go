package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OSExitAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for os.Exit in main()",
	Run:  runOSExitAnalyzer,
}

func runOSExitAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.String() != "main" {
			continue
		}

		isMain := false
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				if s, ok := x.Fun.(*ast.SelectorExpr); ok {
					if s.Sel.Name == "Exit" {
						if i, ok := s.X.(*ast.Ident); ok {
							if i.Name == "os" {
								pass.Reportf(s.Pos(), "call to os.Exit in main")
							}
						}
					}
				}
			case *ast.FuncDecl:
				isMain = x.Name.Name == "main"
				return isMain
			}
			return true
		})
	}
	return nil, nil
}
