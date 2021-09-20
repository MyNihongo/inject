package main

import (
	"fmt"
	"strings"
)

type injectDecl struct {
	function   string
	injectType injectType
}

func getDefinitions(loaded *loadResult) {
	// mapping := make(map[string]string)
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
