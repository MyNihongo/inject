package main

import (
	"fmt"

	"github.com/MyNihongo/codegen"
)

type injectFuncData struct {
	diGraph    map[string]*pkgFuncs
	usedDecls  map[*typeDecl]string
	injectType injectType
	name       string
}

// generateServiceProvider generates the code according to the DI graph
func generateServiceProvider(pkgName string, diGraph map[string]*pkgFuncs) (*codegen.File, error) {
	file := codegen.NewFile(pkgName, "my-nihongo-di")
	imports, isSyncAdded := file.Imports(), false

	for pkgImport, pkgDecl := range diGraph {
		imports.AddImportAlias(pkgImport, pkgDecl.alias)

		for returnType, funcDecl := range pkgDecl.funcs {
			var err error
			var stmts []codegen.Stmt

			funcData := &injectFuncData{
				diGraph:    diGraph,
				usedDecls:  make(map[*typeDecl]string),
				injectType: funcDecl.injectType,
				name:       funcDecl.name,
			}

			if funcDecl.injectType == Singleton {
				if !isSyncAdded {
					imports.AddImport("sync")
					isSyncAdded = true
				}

				varName := fmt.Sprintf("impl_%s", returnType)
				file.DeclareVars(codegen.QualVar(varName, pkgDecl.alias, returnType))

				stmts, err = createInjectionStmts(funcData, pkgDecl, funcDecl, func(v codegen.Value) codegen.Stmt {
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
				stmts, err = createInjectionStmts(funcData, pkgDecl, funcDecl, func(v codegen.Value) codegen.Stmt {
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

func createInjectionStmts(funcData *injectFuncData, pkgFuncs *pkgFuncs, funcDecl *funcDecl, finalBlockFunc func(codegen.Value) codegen.Stmt) ([]codegen.Stmt, error) {
	provideFunc := codegen.QualFuncCall(pkgFuncs.alias, funcDecl.name)

	if len(funcDecl.paramDecls) == 0 {
		return []codegen.Stmt{
			finalBlockFunc(provideFunc),
		}, nil
	} else {
		stmts := make([]codegen.Stmt, 0)
		vals := make([]codegen.Value, len(funcDecl.paramDecls))

		for i, paramDecl := range funcDecl.paramDecls {
			if usedParam, ok := funcData.usedDecls[paramDecl]; ok {
				vals[i] = codegen.Identifier(usedParam)
				continue
			}

			if nestedPkgFuncs, ok := funcData.diGraph[paramDecl.pkgImport]; !ok {
				return nil, fmt.Errorf("package %s is not registered", paramDecl.pkgImport)
			} else if nestedFuncDecl, ok := nestedPkgFuncs.funcs[paramDecl.typeName]; !ok {
				return nil, fmt.Errorf("type %s is not found in the package %s", paramDecl.typeName, paramDecl.pkgImport)
			} else {
				// verify inconsistent injection types (transient into singleton)
				if nestedFuncDecl.injectType > funcData.injectType {
					return nil, fmt.Errorf("cannot inject %s (%s) into %s (%s)", getInjectionName(nestedFuncDecl.injectType), nestedFuncDecl.name, getInjectionName(funcData.injectType), funcData.name)
				}

				newParam := fmt.Sprintf("p%d", len(funcData.usedDecls))
				vals[i] = codegen.Identifier(newParam)
				funcData.usedDecls[paramDecl] = newParam

				nestedStmts, err := createInjectionStmts(funcData, nestedPkgFuncs, nestedFuncDecl, func(v codegen.Value) codegen.Stmt {
					return codegen.Assign(newParam).Values(v)
				})

				if err != nil {
					return nil, err
				} else {
					stmts = append(stmts, nestedStmts...)
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
