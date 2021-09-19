package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getFilePath(fileName string) string {
	return path.Join("test", fmt.Sprintf("%s.txt", fileName))
}

func TestLoadLineSplit(t *testing.T) {
	want := &loadResult{
		injects: map[string]injectType{
			"pkg1.GetService1": Singleton,
			"pkg2.GetService2": Transient,
		},
	}

	ctx := context.Background()
	wd, _ := os.Getwd()

	fileNames := []string{
		"line_split",
		"inline",
		"arg_split",
	}

	for _, fileName := range fileNames {
		filePath := getFilePath(fileName)
		mapping, err := loadFileContent(ctx, wd, filePath)

		assert.Nil(t, err)
		assert.Equal(t, want, mapping)
	}
}
