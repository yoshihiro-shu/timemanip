// Package timemanip provides an analyzer that detects usage of time.Time manipulation methods.
package timemanip

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is the timemanip analyzer.
var Analyzer = &analysis.Analyzer{
	Name:     "timemanip",
	Doc:      "detects usage of time.Time manipulation methods (Add, AddDate, Sub, Truncate, Round)",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// forbiddenMethods contains the set of time.Time methods that should not be used.
var forbiddenMethods = map[string]bool{
	"Add":      true,
	"AddDate":  true,
	"Sub":      true,
	"Truncate": true,
	"Round":    true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		methodName := sel.Sel.Name
		if !forbiddenMethods[methodName] {
			return
		}

		// Check if this is a method call (not a function call)
		selection, ok := pass.TypesInfo.Selections[sel]
		if !ok {
			return
		}

		// Get the receiver type
		recvType := selection.Recv()
		if !isTimeType(recvType) {
			return
		}

		// Check for nolint directive
		if hasNolintDirective(pass, call.Pos()) {
			return
		}

		pass.Reportf(call.Pos(), "use of time.Time.%s is not allowed", methodName)
	})

	return nil, nil
}

// hasNolintDirective checks if the given position has a //nolint:timemanip comment.
func hasNolintDirective(pass *analysis.Pass, pos token.Pos) bool {
	file := pass.Fset.File(pos)
	if file == nil {
		return false
	}

	line := file.Line(pos)

	for _, f := range pass.Files {
		for _, cg := range f.Comments {
			for _, c := range cg.List {
				commentLine := file.Line(c.Pos())
				if commentLine != line {
					continue
				}

				text := c.Text
				// Check for //nolint:timemanip or //nolint (all linters)
				if strings.Contains(text, "nolint:timemanip") || strings.Contains(text, "nolint:all") {
					return true
				}
				// Check for //nolint without specific linter (disables all)
				if strings.TrimSpace(text) == "//nolint" || strings.HasPrefix(strings.TrimSpace(text), "//nolint ") {
					return true
				}
			}
		}
	}

	return false
}

// isTimeType checks if the given type is time.Time or *time.Time.
func isTimeType(t types.Type) bool {
	// Unwrap alias types (Go 1.22+)
	t = types.Unalias(t)

	// Handle pointer types
	if ptr, ok := t.(*types.Pointer); ok {
		t = types.Unalias(ptr.Elem())
	}

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}

	return pkg.Path() == "time" && obj.Name() == "Time"
}
