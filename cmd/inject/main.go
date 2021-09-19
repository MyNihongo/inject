package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	if wd, err := os.Getwd(); err != nil {
		fmt.Print(err)
	} else {
		ctx := context.Background()
		loadFileContent(ctx, wd)
	}
}
