package main

import (
	"github.com/MyNihongo/codegen"
)

func generateServiceProvider(pkgName string, diGraph map[string]*pkgFuncs) (*codegen.File, error) {
	file := codegen.NewFile(pkgName, "my-nihongo-di")
	// imports := file.Imports()

	// for pkgImport, funcs := range diGraph {

	// 	for _, v := range funcs {

	// 	}
	// }

	return file, nil
}
