package custom_analyzer

import (
	"go/ast"
	"path/filepath" // Import path for handling file extensions

	"golang.org/x/tools/go/analysis"
)

// ErrCheckAnalyzer is an analysis.Analyzer that identifies calls to os.Exit
// within the main function. Such calls are generally discouraged in Go,
// especially within the main function, as they bypass normal exit handling
// and may not allow defer functions or cleanup code to run.
var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "osExitCheck",
	Doc:  "check for os.Exit() call in main function",
	Run:  run,
}

// run performs the analysis for ErrCheckAnalyzer.
// It inspects each file in the package for the main function and reports any
// occurrences of os.Exit calls within it.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Filter to analyze only .go files
		filePath := pass.Fset.Position(file.Pos()).Filename
		if filepath.Ext(filePath) != ".go" {
			continue
		}

		// Process the file for os.Exit in the main function
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok || funcDecl.Name.Name != "main" || funcDecl.Recv != nil {
				continue
			}

			// Inspect the body of the main function
			ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := selectorExpr.X.(*ast.Ident); ok && ident.Name == "os" && selectorExpr.Sel.Name == "Exit" {
						pass.Reportf(callExpr.Pos(), "usage of os.Exit in the main function is not allowed")
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
