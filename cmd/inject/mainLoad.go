package main

import (
	"context"
	"os"
	"path"
)

// loadFileContent loads the container definition
func loadFileContent(ctx context.Context, wd string) error {
	filePath := path.Join(wd, "serviceCollection.go")
	if file, err := os.Open(filePath); err != nil {
		return err
	} else {

	}
}
