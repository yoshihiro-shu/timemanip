// Package timemanip provides an analyzer that detects usage of time.Time manipulation methods.
package timemanip

import (
	"go/ast"
	"go/types"

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

		pass.Reportf(call.Pos(), "use of time.Time.%s is not allowed", methodName)
	})

	return nil, nil
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
