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

type funcDecl struct {
}

type typeDecl struct {
	pkgImport string
	typeName  string
}

func getDefinitions(ctx context.Context, wd string, loaded *loadResult) error {
	diDecl := make(map[string]map[string][]*funcDecl)

	if grouping, err := getInjectionGrouping(loaded); err != nil {
		return err
	} else {
		for pkgImport, injections := range grouping {
			if scope, err := loadPackage(ctx, wd, pkgImport); err != nil {
				return err
			} else {
				for _, injection := range injections {
					if typeObj := scope.Lookup(injection.function); typeObj == nil {
						return fmt.Errorf("cannot find a func %s in the package %s", injection.function, pkgImport)
					} else if funcDecl, ok := typeObj.(*types.Func); !ok {
						return fmt.Errorf("%s is not a function", injection.function)
					} else if signature, ok := funcDecl.Type().(*types.Signature); !ok {
						return fmt.Errorf("cannot retrieve a signature of %s", injection.function)
					} else {
						// Return types
						var returnType string
						if returnTypes := signature.Results(); returnTypes == nil || returnTypes.Len() != 1 {
							return fmt.Errorf("func %s does not return a single value", injection.function)
						} else {
							fmt.Println(returnTypes.At(0).Type())
						}

						// Params
						if params := signature.Params(); params != nil {
							for i := 0; i < params.Len(); i++ {
								param := params.At(i)

								fmt.Println(param.Type())
							}
						}
					}
				}
			}
		}

		return nil
	}
}

// getInjectionGrouping creates a grouping of injection calles by their packages
func getInjectionGrouping(loaded *loadResult) (map[string][]*injectDecl, error) {
	grouping := make(map[string][]*injectDecl)

	for inject, injectType := range loaded.injects {
		var alias, function, pkgImport string
		if dotIndex := strings.IndexByte(inject, '.'); dotIndex == -1 {
			alias, function = loaded.pkgName, inject
			pkgImport = loaded.pkgName
		} else {
			alias, function = inject[:dotIndex], inject[dotIndex+1:]
		}

		if len(pkgImport) == 0 {
			var ok bool
			if pkgImport, ok = loaded.imports[alias]; !ok {
				return nil, fmt.Errorf("package import for %s not found", alias)
			}
		}

		decl := &injectDecl{
			function:   function,
			injectType: injectType,
		}

		if slice, ok := grouping[pkgImport]; ok {
			grouping[pkgImport] = append(slice, decl)
		} else {
			grouping[pkgImport] = []*injectDecl{decl}
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

func getTypeDeclaration(t *types.Type) *typeDecl {
	strVal := fmt.Sprint(t)
	return getTypeDeclarationString(strVal)
}

func getTypeDeclarationString(strVal string) *typeDecl {
	typeSeparator := strings.LastIndexByte(strVal, '.')

	return &typeDecl{
		pkgImport: strVal[:typeSeparator],
		typeName:  strVal[typeSeparator+1:],
	}
}
