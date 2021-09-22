package main

import (
	"fmt"

	"github.com/MyNihongo/codegen"
)

// generateServiceProvider generates the code according to the DI graph
func generateServiceProvider(pkgName string, diGraph map[string]*pkgFuncs) (*codegen.File, error) {
	file := codegen.NewFile(pkgName, "my-nihongo-di")
	imports, isSyncAdded := file.Imports(), false

	for pkgImport, pkgFuncs := range diGraph {
		imports.AddImportAlias(pkgImport, pkgFuncs.alias)

		for returnType, funcDecl := range pkgFuncs.funcs {
			var err error
			var stmts []codegen.Stmt
			usedDecls := make(map[*typeDecl]string)

			if funcDecl.injectType == Singleton {
				if isSyncAdded {
					imports.AddImport("sync")
				}

				varName := fmt.Sprintf("impl_%s", returnType)
				file.DeclareVars(codegen.QualVar(varName, pkgFuncs.alias, returnType))

				stmts, err = createInjectionStmts(diGraph, pkgFuncs, funcDecl, func(v codegen.Value) codegen.Stmt {
					return codegen.Assign(varName).Values(v)
				})

				stmts = []codegen.Stmt{
					codegen.QualFuncCall("sync", "DoOnce").Args(codegen.Lambda().Block(stmts...)),
					codegen.Return(codegen.Identifier(varName)),
				}
			} else {
				stmts, err = createInjectionStmts(diGraph, pkgFuncs, funcDecl, func(v codegen.Value) codegen.Stmt {
					return codegen.Return(v)
				})
			}

			funcName, injectName := fmt.Sprintf("Provide%s", returnType), getInjectionName(funcDecl.injectType)
			file.CommentF("%s provides a %s instance of %s.%s", funcName, injectName, pkgImport, returnType)
			file.Func(funcName).ReturnTypes(
				codegen.QualReturnType(pkgFuncs.alias, returnType),
			).Block(stmts...)
		}
	}

	return file, nil
}

func createInjectionStmts(diGraph map[string]*pkgFuncs, pkgFuncs *pkgFuncs, funcDecl *funcDecl, finalBlockFunc func(codegen.Value) codegen.Stmt) ([]codegen.Stmt, error) {
	provideFunc := codegen.QualFuncCall(pkgFuncs.alias, funcDecl.name)

	if len(funcDecl.paramDecls) == 0 {
		return []codegen.Stmt{
			finalBlockFunc(provideFunc),
		}, nil
	} else {
		var ok bool
		for _, paramDecl := range funcDecl.paramDecls {

			if pkgFuncs, ok = diGraph[paramDecl.pkgImport]; !ok {
				return nil, fmt.Errorf("package %s is not registered", paramDecl.pkgImport)
			} else if funcDecl, ok = pkgFuncs.funcs[paramDecl.typeName]; !ok {
				return nil, fmt.Errorf("type %s is not found in the package %s", paramDecl.typeName, paramDecl.pkgImport)
			} else {

			}
		}
	}
}

func getInjectionName(injectType injectType) string {
	switch injectType {
	case Singleton:
		return "singleton"
	case Transient:
		return "transient"
	default:
		panic(fmt.Sprintf("unknown injection type: %d", injectType))
	}
}

// if len(funcDecl.paramTypes) == 0 {
// 	diFunc.Block(
// 		codegen.Return(codegen.QualFuncCall(pkgFuncs.alias, funcDecl.name)),
// 	)
// } else {
// 	// blocks := make([]codegen.Block, len(funcDecl.paramTypes))
// }
