// В пакете osexit определяется кастомный анализатор, проверяющий
// наличие функции os.Exit() в функции main пакета main
package osexit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer - анализатор, который проверяет наличие os.Exit() в функции main пакета main
var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "checks for the presence of os.Exit() in the main function of the main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if strings.HasSuffix(filename, "_test.go") || !strings.HasSuffix(filename, ".go") || file.Name.Name != "main" {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := selector.X.(*ast.Ident); ok {
						if ident.Name == "os" && selector.Sel.Name == "Exit" {
							pass.Reportf(ident.NamePos, "direct call to os.Exit in main package main function")
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
