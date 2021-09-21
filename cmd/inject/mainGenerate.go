package main

import (
	"fmt"

	"github.com/MyNihongo/codegen"
)

func generateServiceProvider(pkgName string, diGraph map[string]*pkgFuncs) (*codegen.File, error) {
	file := codegen.NewFile(pkgName, "my-nihongo-di")
	imports := file.Imports()

	for pkgImport, pkgFuncs := range diGraph {
		imports.AddImportAlias(pkgImport, pkgFuncs.alias)

		for returnType, funcDecl := range pkgFuncs.funcs {
			diFunc := file.Func(fmt.Sprintf("Provide%s", returnType)).ReturnTypes(
				codegen.QualReturnType(pkgFuncs.alias, returnType),
			)

			if len(funcDecl.paramTypes) == 0 {
				diFunc.Block(
					codegen.Return(codegen.QualFuncCall(pkgFuncs.alias, funcDecl.name)),
				)
			} else {
				// blocks := make([]codegen.Block, len(funcDecl.paramTypes))
			}
		}
	}

	return file, nil
}
