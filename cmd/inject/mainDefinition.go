package main

import (
	"context"
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type injectDecl struct {
	function   string
	injectType injectType
}

type pkgInjections struct {
	alias      string
	injections []*injectDecl
}

type pkgFuncs struct {
	alias string
	funcs map[typeNameDecl]*funcDecl
}

type funcDecl struct {
	name       string
	paramDecls []*typeDecl
	injectType injectType
}

type typeDecl struct {
	pkgImport string
	typeName  typeNameDecl
}

type typeNameDecl struct {
	typeName  string
	isPointer bool
}

func getDefinitions(ctx context.Context, wd string, loaded *loadResult) (map[string]*pkgFuncs, error) {
	// First key - package import; Second key - return type; Value - function parameters and inject type
	diDecl := make(map[string]*pkgFuncs)

	if grouping, err := getInjectionGrouping(loaded); err != nil {
		return nil, err
	} else {
		for pkgImport, pkgInjections := range grouping {
			if scope, err := loadPackage(ctx, wd, pkgImport); err != nil {
				return nil, err
			} else {
				for _, injection := range pkgInjections.injections {
					if typeObj := scope.Lookup(injection.function); typeObj == nil {
						return nil, fmt.Errorf("cannot find a func %s in the package %s", injection.function, pkgImport)
					} else if funcType, ok := typeObj.(*types.Func); !ok {
						return nil, fmt.Errorf("%s is not a function", injection.function)
					} else if signature, ok := funcType.Type().(*types.Signature); !ok {
						return nil, fmt.Errorf("cannot retrieve a signature of %s", injection.function)
					} else {
						// Return types
						var returnType *typeDecl
						if returnTypes := signature.Results(); returnTypes == nil || returnTypes.Len() != 1 {
							return nil, fmt.Errorf("func %s does not return a single value", injection.function)
						} else {
							returnType = getTypeDeclaration(returnTypes.At(0).Type())
						}

						var pkgGrouping *pkgFuncs
						if pkgGrouping, ok = diDecl[returnType.pkgImport]; !ok {
							pkgGrouping = &pkgFuncs{
								alias: pkgInjections.alias,
								funcs: make(map[typeNameDecl]*funcDecl, 1),
							}
							diDecl[returnType.pkgImport] = pkgGrouping
						}

						// Params
						var paramTypes []*typeDecl
						if params := signature.Params(); params != nil {
							paramTypes = make([]*typeDecl, params.Len())
							for i := 0; i < params.Len(); i++ {
								t := params.At(i).Type()
								paramTypes[i] = getTypeDeclaration(t)
							}
						} else {
							paramTypes = make([]*typeDecl, 0)
						}

						pkgGrouping.funcs[returnType.typeName] = &funcDecl{
							name:       injection.function,
							paramDecls: paramTypes,
							injectType: injection.injectType,
						}
					}
				}
			}
		}

		return diDecl, nil
	}
}

// getInjectionGrouping creates a grouping of injection calles by their packages
func getInjectionGrouping(loaded *loadResult) (map[string]*pkgInjections, error) {
	grouping := make(map[string]*pkgInjections)

	for inject, injectType := range loaded.injects {
		var alias, function, pkgImport string
		if dotIndex := strings.IndexByte(inject, '.'); dotIndex == -1 {
			function = inject
		} else {
			alias, function = inject[:dotIndex], inject[dotIndex+1:]
		}

		if len(alias) != 0 {
			var ok bool
			if pkgImport, ok = loaded.imports[alias]; !ok {
				return nil, fmt.Errorf("package import for %s not found", alias)
			}
		}

		decl := &injectDecl{
			function:   function,
			injectType: injectType,
		}

		if pkgInject, ok := grouping[pkgImport]; ok {
			pkgInject.injections = append(pkgInject.injections, decl)
		} else {
			grouping[pkgImport] = &pkgInjections{
				alias:      alias,
				injections: []*injectDecl{decl},
			}
		}
	}

	return grouping, nil
}

func loadPackage(ctx context.Context, wd, pkgImport string) (*types.Scope, error) {
	cfg := &packages.Config{
		Context: ctx,
		Dir:     wd,
		Mode:    packages.NeedTypes | packages.NeedImports,
	}

	if pkgs, err := packages.Load(cfg, pkgImport); err != nil {
		return nil, err
	} else if len(pkgs) != 1 {
		return nil, fmt.Errorf("cannot resolve a single package for %s. Count: %d", pkgImport, len(pkgs))
	} else {
		return pkgs[0].Types.Scope(), nil
	}
}

func getTypeDeclaration(t types.Type) *typeDecl {
	strVal := fmt.Sprint(t)
	return getTypeDeclarationString(strVal)
}

func getTypeDeclarationString(strVal string) *typeDecl {
	typeSeparator := strings.LastIndexByte(strVal, '.')
	pkgImport := strVal[:typeSeparator]

	var isPointer bool
	if strings.HasPrefix(pkgImport, "*") {
		pkgImport = pkgImport[1:]
		isPointer = true
	}

	return &typeDecl{
		pkgImport: pkgImport,
		typeName: typeNameDecl{
			typeName:  strVal[typeSeparator+1:],
			isPointer: isPointer,
		},
	}
}
