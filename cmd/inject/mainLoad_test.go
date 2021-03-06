package main

import (
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
		imports: map[string]string{
			"pkg1": "github.com/MyNihongo/inject/examples/pkg1",
			"pkg2": "github.com/MyNihongo/inject/examples/pkg2",
		},
		injects: map[string]injectType{
			"pkg1.GetService1": Singleton,
			"pkg2.GetService2": Transient,
		},
	}

	wd, _ := os.Getwd()

	fileNames := []string{
		"func_line_split",
		"func_inline",
		"func_arg_split",
	}

	for _, fileName := range fileNames {
		filePath := getFilePath(fileName)
		mapping, err := loadFileContent(wd, filePath)

		assert.Nil(t, err)
		assert.Equal(t, want, mapping)
	}
}

func TestLoadImportSingle(t *testing.T) {
	want := &loadResult{
		pkgName: "examples",
		imports: map[string]string{
			"pkg1": "github.com/MyNihongo/inject/examples/pkg1",
		},
		injects: map[string]injectType{},
	}

	wd, _ := os.Getwd()
	filePath := getFilePath("import_single")

	mapping, err := loadFileContent(wd, filePath)

	assert.Nil(t, err)
	assert.Equal(t, want, mapping)
}

func TestLoadImportSingleAlias(t *testing.T) {
	want := &loadResult{
		pkgName: "examples",
		imports: map[string]string{
			"my_alias": "github.com/MyNihongo/inject/examples/pkg1",
		},
		injects: map[string]injectType{},
	}

	wd, _ := os.Getwd()
	filePath := getFilePath("import_single_alias")

	mapping, err := loadFileContent(wd, filePath)

	assert.Nil(t, err)
	assert.Equal(t, want, mapping)
}

func TestLoadImportMultiline(t *testing.T) {
	want := &loadResult{
		pkgName: "examples",
		imports: map[string]string{
			"pkg1":     "github.com/MyNihongo/inject/examples/pkg1",
			"my_alias": "github.com/MyNihongo/inject/examples/pkg2",
		},
		injects: map[string]injectType{},
	}

	wd, _ := os.Getwd()
	filePath := getFilePath("import_multiline")

	mapping, err := loadFileContent(wd, filePath)

	assert.Nil(t, err)
	assert.Equal(t, want, mapping)
}
