package retchecking

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "retchecking checks whether return statement is used immediately after certain functions are called in the given package."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "retchecking",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {

	// list of function to be checked

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	ast.Print(pass.Fset, pass.Files[0].Decls)

	inspect.Preorder(nil, func(n ast.Node) {})
	return nil, nil
}

// check if e is in s or not
func contains(e string, s []string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

// run breadth first search on ast
func inspectBFS(pass *analysis.Pass) {
	queue := make([]*ast.Decl, 0)
	queue = append(queue)
	for len(queue) > 0 {
		//decl := queue[0]
		queue = queue[1:]
	}
}
