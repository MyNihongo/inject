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
		pkgName: "examples",
		imports: []*importStmt{
			{path: "github.com/MyNihongo/inject/examples/pkg1"},
			{path: "github.com/MyNihongo/inject/examples/pkg2"},
		},
		injects: map[string]injectType{
			"pkg1.GetService1": Singleton,
			"pkg2.GetService2": Transient,
		},
	}

	ctx := context.Background()
	wd, _ := os.Getwd()

	fileNames := []string{
		"func_line_split",
		"func_inline",
		"func_arg_split",
	}

	for _, fileName := range fileNames {
		filePath := getFilePath(fileName)
		mapping, err := loadFileContent(ctx, wd, filePath)

		assert.Nil(t, err)
		assert.Equal(t, want, mapping)
	}
}

func TestLoadImportSingle(t *testing.T) {
	want := &loadResult{
		pkgName: "examples",
		imports: []*importStmt{
			{path: "github.com/MyNihongo/inject/examples/pkg1"},
		},
		injects: map[string]injectType{},
	}

	wd, _ := os.Getwd()
	ctx, filePath := context.Background(), getFilePath("import_single")

	mapping, err := loadFileContent(ctx, wd, filePath)

	assert.Nil(t, err)
	assert.Equal(t, want, mapping)
}

func TestLoadImportSingleAlias(t *testing.T) {
	want := &loadResult{
		pkgName: "examples",
		imports: []*importStmt{
			{
				alias: "my_alias",
				path:  "github.com/MyNihongo/inject/examples/pkg1",
			},
		},
		injects: map[string]injectType{},
	}

	wd, _ := os.Getwd()
	ctx, filePath := context.Background(), getFilePath("import_single_alias")

	mapping, err := loadFileContent(ctx, wd, filePath)

	assert.Nil(t, err)
	assert.Equal(t, want, mapping)
}
