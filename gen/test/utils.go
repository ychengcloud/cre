package test

import (
	"go/ast"
	"log"
)

func contain(s []*ast.Ident, e string) bool {
	for _, ident := range s {
		if ident.Name == e {
			return true
		}
	}
	return false
}

func Exist(f *ast.File, structName, fieldType, fieldName string) bool {
	var found bool
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == structName {
				switch y := x.Type.(type) {
				case *ast.StructType:
					for _, field := range y.Fields.List {
						var ident *ast.Ident
						switch z := field.Type.(type) {
						case *ast.Ident:
							ident = z
						case *ast.StarExpr:
							ident = z.X.(*ast.Ident)
						case *ast.ArrayType:
							switch se := z.Elt.(type) {
							case *ast.Ident:
								ident = se
							case *ast.StarExpr:
								ident = se.X.(*ast.Ident)
							}
						}
						// log.Printf("field %s, %#v, %s\n", field.Names, field.Type, ident)
						if ident.Name == fieldType && contain(field.Names, fieldName) {
							found = true
						}
					}
					return false
				}
			}
		}
		return true
	})
	return found
}

//解析函数调用代码段
func extractFuncCallInFunc(stmts []ast.Stmt) (name string, body any) {
	for _, stmt := range stmts {
		log.Printf("stmt: %#v\n", stmt)
		if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
			if call, ok := exprStmt.X.(*ast.CallExpr); ok {
				log.Printf("call: %#v\n", call)
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
					name = fun.Sel.Name
					body = fun.Sel.Obj.Data
				}
			}
		}
	}
	return
}

func ExistFunc(f *ast.File, structName, funcName string) any {
	funcs := []*ast.FuncDecl{}
	for _, d := range f.Decls {
		if fn, isFn := d.(*ast.FuncDecl); isFn {
			funcs = append(funcs, fn)
		}
	}

	for _, function := range funcs {
		extractFuncCallInFunc(function.Body.List)
		log.Printf("ExistFuncs: %s, %s, %#v\n", structName, funcName, function.Body)
		if function.Name.Name == funcName {
			return function.Body
		}
	}

	return nil
}
