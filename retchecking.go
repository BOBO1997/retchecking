package retchecking

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	mapset "github.com/deckarep/golang-set"
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
	fileName := pass.Files[0].Name.Name
	//ast.Print(pass.Fset, pass.Files[0].Decls)

	// DONE: use map[types.Object]struct{} instead of []types.Object
	funcSetObj := mapset.NewSet()
	if obj := analysisutil.LookupFromImports(pass.Pkg.Imports(), "net/http", "Error"); obj != nil {
		funcSetObj.Add(obj)
	}
	if obj := analysisutil.ObjectOf(pass, fileName, "fError"); obj != nil {
		funcSetObj.Add(obj)
	}
	if obj := analysisutil.ObjectOf(pass, fileName, "f"); obj != nil {
		funcSetObj.Add(obj)
	}
	if obj := analysisutil.MethodOf(analysisutil.TypeOf(pass, fileName, "S"), "Error"); obj != nil {
		funcSetObj.Add(obj)
	}

	funcSetIdent := mapset.NewSet()

	// for debugging
	/*
		fmt.Println(funcSetObj)
		fmt.Println(funcSetObj.Cardinality())
		fmt.Println(funcSetObj.Contains(analysisutil.LookupFromImports(pass.Pkg.Imports(), "net/http", "Error")))
		fmt.Println(funcSetObj.Contains(analysisutil.LookupFromImports(pass.Pkg.Imports(), pass.Files[0].Name.Name, "a")))
	*/

	// TODO: add a handler for *ast.AssignStmt, *ast.FuncLit, *ast.IndexExpr and *ast.SelectorExpr
	// TODO: *ast.AssignStmt, add this to funcSetObj
	// DONE: *ast.FuncLit, add this to the switch handler

	// depth first search on ast nodes, with referring each element of the list of the node.
	// time complexity of this part is three times larger than the simple depth first search
	inspect.Preorder(nil, func(node ast.Node) {
		switch block := node.(type) {
		case *ast.BlockStmt:
			for i, stmt := range block.List {
				switch expr := stmt.(type) {
				case *ast.AssignStmt: // handling with assign statement, := and =
					//fmt.Println("AssignStmt")
					// need to refer to RHS, since it is unnecessary to verify v if RHS is not the function we are paying attention
					if len(expr.Lhs) == len(expr.Rhs) { // ! restrict the condition
						for i, v := range expr.Lhs {
							switch v := v.(type) {
							case *ast.Ident:
								//fmt.Println("this is ident", v, pass.TypesInfo.Defs[v])
								if obj := pass.TypesInfo.Defs[v]; obj != nil {
									switch rhs := expr.Rhs[i].(type) {
									case *ast.Ident:
										if funcSetObj.Contains(rhs) {
											funcSetIdent.Add(v)
										}
									case *ast.SelectorExpr:
										if funcSetObj.Contains(pass.TypesInfo.Uses[rhs.Sel]) {
											funcSetIdent.Add(v)
										}
									}
								}
							case *ast.IndexExpr:
								pass.Reportf(v.Pos(), "Undefined Process")
								// TODO
								// ! This cannot be resolved!
							case *ast.SelectorExpr:
								if obj := pass.TypesInfo.Defs[v.Sel]; obj != nil {
									switch rhs := expr.Rhs[i].(type) {
									case *ast.Ident:
										if funcSetObj.Contains(rhs) {
											funcSetIdent.Add(v) // correct?
										}
									case *ast.SelectorExpr:
										if funcSetObj.Contains(pass.TypesInfo.Uses[rhs.Sel]) {
											funcSetIdent.Add(v) // correct?
										}
									}
									funcSetObj.Add(obj)
								}
								// TODO
								// ! This cannot be resolved!
							default:
								pass.Reportf(v.Pos(), "Undefined Process")
							}
							//fmt.Println(funcSetIdent)
							//fmt.Println(funcSetObj)
						}
					}
				case *ast.ExprStmt: // if a function is called
					switch x := expr.X.(type) {
					case *ast.CallExpr:
						var xobj types.Object
						var pos token.Pos
						var varFlag bool
						switch fun := x.Fun.(type) {
						case *ast.IndexExpr:
							// TODO
							// ! This cannot be resolved!
						case *ast.SelectorExpr: // if the function is a method of a structure
							//fmt.Println("fun.Sel.Obj: ", fun.Sel)
							xobj = pass.TypesInfo.Uses[fun.Sel]
							pos = fun.Sel.Pos()
						case *ast.FuncLit:
							//fmt.Println("fun is an *ast.FuncLit: ", fun)
							pos = fun.Pos()
							pass.Reportf(pos, "OK")
							continue
						case *ast.Ident: // if the function is a normal function or variable
							//fmt.Println("fun is an *ast.Ident: ", fun)
							xobj = pass.TypesInfo.Uses[fun]
							pos = fun.Pos()
							varFlag = funcSetIdent.Contains(fun) // why false?
							//fmt.Println(fun, funcSetIdent, varFlag)
						default:
							//fmt.Println("fun: ", fun)
							pos = fun.Pos()
							//pass.Reportf(pos, "Undefined Process")
						}
						//fmt.Println("xobj: ", xobj)
						if i == len(block.List)-1 {
							//fmt.Println("The last function without return", pass.Fset.Position(pos), "\nNG")
							pass.Reportf(pos, "NG")
						} else if i < len(block.List)-1 && (funcSetObj.Contains(xobj) || varFlag) {
							switch block.List[i+1].(type) {
							case *ast.ReturnStmt:
								//fmt.Println(pass.Fset.Position(pos), "\nOK")
								//pass.Reportf(pos, "OK")
							default:
								//fmt.Println("function without return", pass.Fset.Position(pos), "\nNG")
								//fmt.Println("NG")
								pass.Reportf(pos, "NG")
							}
						} else {
							//fmt.Println(pass.Fset.Position(pos), "\npass")
							//pass.Reportf(pos, "OK")
							continue
						}
					}
				}
			}
		}
	})

	//dumpTypesInfo(pass)
	return nil, nil
}

// old lib, unused
func isObjInList(xobj types.Object, functions []types.Object) bool {
	for _, obj := range functions {
		if xobj == obj {
			return true
		}
	}
	return false
}

// old lib, unused
func isStringInList(xobj string, functions []string) bool {
	for _, obj := range functions {
		if xobj == obj {
			return true
		}
	}
	return false
}

// for debugging, unused
func dumpTypesInfo(pass *analysis.Pass) {
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
}
