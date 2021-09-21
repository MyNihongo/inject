package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type injectType uint8

const (
	notSet    injectType = 0
	Singleton injectType = 1
	Transient injectType = 2
)

func main() {
	if wd, err := os.Getwd(); err != nil {
		fmt.Println(err)
	} else {
		ctx := context.Background()
		if loaded, err := loadFileContent(wd, "serviceCollection.go"); err != nil {
			fmt.Println(err)
		} else if diGraph, err := getDefinitions(ctx, wd, loaded); err != nil {
			fmt.Println(err)
		} else if file, err := generateServiceProvider(loaded.pkgName, diGraph); err != nil {
			fmt.Println(err)
		} else {
			savePath := filepath.Join(wd, "serviceProvider_gen.go")
			file.Save(savePath)
		}
	}
}
