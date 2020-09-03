package retchecking

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/gostaticanalysis/analysisutil"
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
	functions := make([]types.Object, 0)
	functions = append(functions, analysisutil.LookupFromImports(pass.Pkg.Imports(), "net/http", "Error"))
	functions = append(functions, analysisutil.ObjectOf(pass, "a", "fError"))
	//functions = append(functions, analysisutil.ObjectOf(pass, "fError"))
	fmt.Println(functions)
	ast.Print(pass.Fset, pass.Files[0].Decls)

	for key, item := range pass.TypesInfo.Implicits {
		fmt.Println("key: ", key, " item: ", item)
	}

	inspect.Preorder(nil, func(node ast.Node) {
		switch block := node.(type) {
		case *ast.BlockStmt:
			for i, stmt := range block.List {
				switch expr := stmt.(type) {
				case *ast.ExprStmt:
					switch x := expr.X.(type) {
					case *ast.CallExpr:
						var xobj types.Object
						switch fun := x.Fun.(type) {
						case *ast.SelectorExpr:
							xobj = pass.TypesInfo.Selections[fun].Obj()
						default:
							xobj = pass.TypesInfo.Implicits[x]
						}
						fmt.Println("xobj: ", xobj)
						if i == len(block.List)-1 {
							fmt.Println("NG")
							pass.Reportf(x.Pos(), "NG")
						} else if i < len(block.List)-1 && isObjInList(xobj, functions) {
							switch block.List[i+1].(type) {
							case *ast.ReturnStmt:
								fmt.Println("OK")
								pass.Reportf(x.Pos(), "OK")
							default:
								fmt.Println("NG")
								pass.Reportf(x.Pos(), "NG")
							}
						} else {
							continue
						}
					}
				}
			}
		}
	})
	return nil, nil
}

func isObjInList(xobj types.Object, functions []types.Object) bool {
	for _, obj := range functions {
		if xobj == obj {
			return true
		}
	}
	return false
}

func isStringInList(xobj string, functions []string) bool {
	for _, obj := range functions {
		if xobj == obj {
			return true
		}
	}
	return false
}
