package main

import (
	"context"
	"fmt"
	"os"
)

type injectType uint8

const (
	notSet    injectType = 0
	Singleton injectType = 1
	Transient injectType = 2
)

func main() {
	if wd, err := os.Getwd(); err != nil {
		fmt.Print(err)
	} else {
		ctx := context.Background()
		if loaded, err := loadFileContent(ctx, wd, "serviceCollection.go"); err != nil {
			fmt.Print(err)
		} else {
			getDefinitions(loaded)
		}
	}
}
