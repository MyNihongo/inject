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
		loadFileContent(ctx, wd, "serviceCollection.go")
	}
}
