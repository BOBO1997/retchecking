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
	functions := []string{"fError", "S.Error", "http.Error", "f"}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	/*
		fmt.Println()
		ast.Print(pass.Fset, pass.Files[0].Decls)
		fmt.Println()
	*/

	inspect.Preorder(nil, func(n ast.Node) {
		//fmt.Println(n.End())
		//ast.Print(pass.Fset, n)
		switch n := n.(type) {
		case *ast.BlockStmt: // inspect a block statement
			for i, tobechecked := range n.List {
				switch tobechecked := tobechecked.(type) {
				case *ast.ExprStmt: // expression statement
					//ast.Print(nil, tobechecked)
					switch funccall := tobechecked.X.(type) { // check function calling
					case *ast.CallExpr:
						if contains(funccall.Fun.(*ast.Ident).Name, functions) { //
							//fmt.Println("put sth")
							if i+1 == len(n.List) {
								pass.Reportf(funccall.Pos(), "NG")
							} else {
								switch n.List[i+1].(type) {
								case *ast.ReturnStmt:
									continue
								default:
									pass.Reportf(funccall.Pos(), "NG")
								}
							}
						}
					}
				}
			}
		}
	})

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
