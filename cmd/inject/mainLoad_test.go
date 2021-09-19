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

func TestLoad(t *testing.T) {
	ctx := context.Background()
	wd, _ := os.Getwd()
	filePath := getFilePath("line_split")

	want := map[string]injectType{
		"pkg1.GetService1": Singleton,
		"pkg2.GetService2": Transient,
	}

	mapping, err := loadFileContent(ctx, wd, filePath)

	assert.Nil(t, err)
	assert.Equal(t, want, mapping)
}
