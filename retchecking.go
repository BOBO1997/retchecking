package retchecking

import (
	"fmt"
	"go/ast"
	"go/token"
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
	ast.Print(pass.Fset, pass.Files[0].Decls)

	// TODO: use map[types.Object]struct{} instead of []types.Object
	functions := make([]types.Object, 0)
	functions = append(functions, analysisutil.LookupFromImports(pass.Pkg.Imports(), "net/http", "Error"))
	functions = append(functions, analysisutil.ObjectOf(pass, "a", "fError"))
	//functions = append(functions, pass.Pkg.Scope().Lookup("Error"))
	functions = append(functions, analysisutil.MethodOf(analysisutil.TypeOf(pass, "a", "S"), "Error"))
	for _, function := range functions {
		fmt.Println(function)
	}

	/*
		fmt.Println("\nDefs")
		fmt.Println(pass.TypesInfo.Defs)
		for key, item := range pass.TypesInfo.Defs {
			fmt.Println("key: ", key, " item: ", item)
		}
		fmt.Println("\nScopes")
		for key, item := range pass.TypesInfo.Scopes {
			fmt.Println("key: ", key, " item: ", item)
		}
		fmt.Println("\nImplicits")
		for key, item := range pass.TypesInfo.Implicits {
			fmt.Println("key: ", key, " item: ", item)
		}
		fmt.Println("\nSelections")
		for key, item := range pass.TypesInfo.Selections {
			fmt.Println("key: ", key, " item: ", item)
		}
		fmt.Println("\nInitOrder")
		for key, item := range pass.TypesInfo.InitOrder {
			fmt.Println("key: ", key, " item: ", item)
		}
		fmt.Println("\nTypes")
		for key, item := range pass.TypesInfo.Types {
			fmt.Println("key: ", key, " item: ", item)
		}
		fmt.Println("\nUses")
		for key, item := range pass.TypesInfo.Uses {
			fmt.Println("key: ", key, " item: ", item)
		}
	*/

	// TODO: add a handler for *ast.AssignStmt and closure

	inspect.Preorder(nil, func(node ast.Node) {
		switch block := node.(type) {
		case *ast.BlockStmt:
			for i, stmt := range block.List {
				fmt.Println()
				switch expr := stmt.(type) {
				case *ast.ExprStmt: // if a function is called
					switch x := expr.X.(type) {
					case *ast.CallExpr:
						var xobj types.Object
						var pos token.Pos
						switch fun := x.Fun.(type) {
						case *ast.SelectorExpr: // if the function is a method of a structure
							//xobj = pass.TypesInfo.Selections[fun].Obj()
							fmt.Println("fun.Sel.Obj: ", fun.Sel)
							xobj = pass.TypesInfo.Uses[fun.Sel]
							pos = fun.Sel.Pos()
						default: // if the function is a normal function
							//xobj = pass.TypesInfo.Implicits[x]
							fmt.Println("fun.(*ast.Ident): ", fun.(*ast.Ident))
							xobj = pass.TypesInfo.Uses[fun.(*ast.Ident)]
							pos = fun.(*ast.Ident).Pos()
						}
						fmt.Println("xobj: ", xobj)
						if i == len(block.List)-1 {
							fmt.Println("The last function without return", pass.Fset.Position(pos), "\nNG")
							pass.Reportf(pos, "NG")
						} else if i < len(block.List)-1 && isObjInList(xobj, functions) {
							switch block.List[i+1].(type) {
							case *ast.ReturnStmt:
								fmt.Println(pass.Fset.Position(pos), "\nOK")
								pass.Reportf(pos, "OK")
							default:
								fmt.Println("function without return", pass.Fset.Position(pos), "\nNG")
								fmt.Println("NG")
								pass.Reportf(pos, "NG")
							}
						} else {
							fmt.Println(pass.Fset.Position(pos), "\npass")
							pass.Reportf(pos, "OK")
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
