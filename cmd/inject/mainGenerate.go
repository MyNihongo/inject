package main

import (
	"fmt"

	"github.com/MyNihongo/codegen"
)

// generateServiceProvider generates the code according to the DI graph
func generateServiceProvider(pkgName string, diGraph map[string]*pkgFuncs) (*codegen.File, error) {
	file := codegen.NewFile(pkgName, "my-nihongo-di")
	imports, isSyncAdded := file.Imports(), false

	for pkgImport, pkgDecl := range diGraph {
		imports.AddImportAlias(pkgImport, pkgDecl.alias)

		for returnType, funcDecl := range pkgDecl.funcs {
			var err error
			var stmts []codegen.Stmt
			usedDecls := make(map[*typeDecl]string)

			if funcDecl.injectType == Singleton {
				if isSyncAdded {
					imports.AddImport("sync")
				}

				varName := fmt.Sprintf("impl_%s", returnType)
				file.DeclareVars(codegen.QualVar(varName, pkgDecl.alias, returnType))

				stmts, err = createInjectionStmts(diGraph, pkgDecl, funcDecl, usedDecls, func(v codegen.Value) codegen.Stmt {
					return codegen.Assign(varName).Values(v)
				})

				if err != nil {
					return nil, err
				}

				stmts = []codegen.Stmt{
					codegen.QualFuncCall("sync", "DoOnce").Args(codegen.Lambda().Block(stmts...)),
					codegen.Return(codegen.Identifier(varName)),
				}
			} else {
				stmts, err = createInjectionStmts(diGraph, pkgDecl, funcDecl, usedDecls, func(v codegen.Value) codegen.Stmt {
					return codegen.Return(v)
				})

				if err != nil {
					return nil, err
				}
			}

			funcName, injectName := fmt.Sprintf("Provide%s", returnType), getInjectionName(funcDecl.injectType)
			file.CommentF("%s provides a %s instance of %s.%s", funcName, injectName, pkgImport, returnType)
			file.Func(funcName).ReturnTypes(
				codegen.QualReturnType(pkgDecl.alias, returnType),
			).Block(stmts...)
		}
	}

	return file, nil
}

func createInjectionStmts(diGraph map[string]*pkgFuncs, pkgFuncs *pkgFuncs, funcDecl *funcDecl, usedDecls map[*typeDecl]string, finalBlockFunc func(codegen.Value) codegen.Stmt) ([]codegen.Stmt, error) {
	provideFunc := codegen.QualFuncCall(pkgFuncs.alias, funcDecl.name)

	if len(funcDecl.paramDecls) == 0 {
		return []codegen.Stmt{
			finalBlockFunc(provideFunc),
		}, nil
	} else {
		var ok bool
		stmts := make([]codegen.Stmt, 0)
		vals := make([]codegen.Value, len(funcDecl.paramDecls))

		for i, paramDecl := range funcDecl.paramDecls {
			if usedParam, ok := usedDecls[paramDecl]; ok {
				vals[i] = codegen.Identifier(usedParam)
				continue
			}

			if nestedPkgFuncs, ok := diGraph[paramDecl.pkgImport]; !ok {
				return nil, fmt.Errorf("package %s is not registered", paramDecl.pkgImport)
			} else if nestedFuncDecl, ok := nestedPkgFuncs.funcs[paramDecl.typeName]; !ok {
				return nil, fmt.Errorf("type %s is not found in the package %s", paramDecl.typeName, paramDecl.pkgImport)
			} else {
				// TODO: verify against the base (singleton cannot inject transient)
				nestedFuncDecl.injectType

				newParam := fmt.Sprintf("p%d", len(usedDecls))
				vals[i] = codegen.Identifier(newParam)
				usedDecls[paramDecl] = newParam

				// TODO: assignment
				a, err := createInjectionStmts(diGraph, nestedPkgFuncs, nestedFuncDecl, usedDecls, func(v codegen.Value) codegen.Stmt {

				})

				if err != nil {
					return nil, err
				} else {
					stmts = append(stmts, a...)
				}
			}
		}

		stmts = append(stmts, provideFunc.Args(vals...))
		return stmts, nil
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
