package main

import (
	"fmt"

	"github.com/MyNihongo/codegen"
)

func generateServiceProvider(pkgName string, diGraph map[string]*pkgFuncs) (*codegen.File, error) {
	file := codegen.NewFile(pkgName, "my-nihongo-di")
	imports, isSyncAdded := file.Imports(), false

	for pkgImport, pkgFuncs := range diGraph {
		imports.AddImportAlias(pkgImport, pkgFuncs.alias)

		for returnType, funcDecl := range pkgFuncs.funcs {
			var stmts []codegen.Stmt
			if funcDecl.injectType == Singleton {
				if isSyncAdded {
					imports.AddImport("sync")
				}

				varName := fmt.Sprintf("impl_%s", returnType)
				file.DeclareVars(codegen.QualVar(varName, pkgFuncs.alias, returnType))

				stmts = generateInjectionStats(pkgFuncs, funcDecl, func(v codegen.Value) codegen.Stmt {
					return codegen.Assign(varName).Values(v)
				})

				stmts = []codegen.Stmt{
					codegen.QualFuncCall("sync", "DoOnce").Args(codegen.Lambda().Block(stmts...)),
					codegen.Return(codegen.Identifier(varName)),
				}
			} else {
				stmts = generateInjectionStats(pkgFuncs, funcDecl, func(v codegen.Value) codegen.Stmt {
					return codegen.Return(v)
				})
			}

			file.Func(fmt.Sprintf("Provide%s", returnType)).ReturnTypes(
				codegen.QualReturnType(pkgFuncs.alias, returnType),
			).Block(stmts...)
		}
	}

	return file, nil
}

func generateInjectionStats(pkgFuncs *pkgFuncs, funcDecl *funcDecl, finalBlockFunc func(codegen.Value) codegen.Stmt) []codegen.Stmt {
	provideFunc := codegen.QualFuncCall(pkgFuncs.alias, funcDecl.name)

	if len(funcDecl.paramTypes) == 0 {
		return []codegen.Stmt{
			finalBlockFunc(provideFunc),
		}
	} else {
		panic("aaa")
	}
}

// if len(funcDecl.paramTypes) == 0 {
// 	diFunc.Block(
// 		codegen.Return(codegen.QualFuncCall(pkgFuncs.alias, funcDecl.name)),
// 	)
// } else {
// 	// blocks := make([]codegen.Block, len(funcDecl.paramTypes))
// }
