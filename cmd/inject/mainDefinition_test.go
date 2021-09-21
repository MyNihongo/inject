package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getExamplesWd() string {
	wd, _ := os.Getwd()
	cmdDir := filepath.Join("cmd", "inject")

	dirIndex := strings.LastIndex(wd, cmdDir)
	return filepath.Join(wd[:dirIndex], "examples")
}

func TestGroupingSamePackageOne(t *testing.T) {
	want := map[string][]*injectDecl{
		"di": {
			{
				function:   "createFoo",
				injectType: Singleton,
			},
		},
	}

	fixture := &loadResult{
		pkgName: "di",
		imports: map[string]string{},
		injects: map[string]injectType{
			"createFoo": Singleton,
		},
	}

	got, err := getInjectionGrouping(fixture)

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestGroupingSamePackageMultiple(t *testing.T) {
	want := map[string][]*injectDecl{
		"di": {
			{
				function:   "createFoo",
				injectType: Singleton,
			},
			{
				function:   "createBoo",
				injectType: Transient,
			},
		},
	}

	fixture := &loadResult{
		pkgName: "di",
		imports: map[string]string{},
		injects: map[string]injectType{
			"createFoo": Singleton,
			"createBoo": Transient,
		},
	}

	got, err := getInjectionGrouping(fixture)

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestGroupingAnotherPackageOne(t *testing.T) {
	want := map[string][]*injectDecl{
		"github.com/MyNihongo/inject/examples/pkg1": {
			{
				function:   "CreateFoo",
				injectType: Singleton,
			},
		},
	}

	fixture := &loadResult{
		pkgName: "di",
		imports: map[string]string{
			"pkg1": "github.com/MyNihongo/inject/examples/pkg1",
		},
		injects: map[string]injectType{
			"pkg1.CreateFoo": Singleton,
		},
	}

	got, err := getInjectionGrouping(fixture)

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestGroupingAnotherPackageMultiple(t *testing.T) {
	want := map[string][]*injectDecl{
		"github.com/MyNihongo/inject/examples/pkg1": {
			{
				function:   "CreateFoo",
				injectType: Singleton,
			},
			{
				function:   "CreateBoo",
				injectType: Transient,
			},
		},
		"github.com/MyNihongo/inject/examples/pkg2": {
			{
				function:   "CreateFoo",
				injectType: Transient,
			},
			{
				function:   "CreateBoo",
				injectType: Singleton,
			},
		},
	}

	fixture := &loadResult{
		pkgName: "di",
		imports: map[string]string{
			"pkg1": "github.com/MyNihongo/inject/examples/pkg1",
			"pkg2": "github.com/MyNihongo/inject/examples/pkg2",
		},
		injects: map[string]injectType{
			"pkg1.CreateFoo": Singleton,
			"pkg1.CreateBoo": Transient,
			"pkg2.CreateFoo": Transient,
			"pkg2.CreateBoo": Singleton,
		},
	}

	got, err := getInjectionGrouping(fixture)

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestGroupingErrorIfNoImport(t *testing.T) {
	fixture := &loadResult{
		pkgName: "di",
		imports: map[string]string{},
		injects: map[string]injectType{
			"pkg.CreateFoo": Singleton,
		},
	}

	got, err := getInjectionGrouping(fixture)

	assert.Nil(t, got)
	assert.Error(t, err, "package import for pkg not found")
}

func TestGetTypeDeclaration(t *testing.T) {
	want := &typeDecl{
		pkgImport: "github.com/MyNihongo/inject/examples/pkg1",
		typeName:  "Service1",
	}

	got := getTypeDeclarationString("github.com/MyNihongo/inject/examples/pkg1.Service1")

	assert.Equal(t, want, got)
}

func TestDefinitions(t *testing.T) {
	want := map[string]map[string]*funcDecl{
		"github.com/MyNihongo/inject/examples/pkg1": {
			"Service1": {
				name:       "GetService1",
				paramTypes: []*typeDecl{},
				injectType: Singleton,
			},
		},
		"github.com/MyNihongo/inject/examples/pkg2": {
			"Service2": {
				name: "GetService2",
				paramTypes: []*typeDecl{
					{
						pkgImport: "github.com/MyNihongo/inject/examples/pkg2",
						typeName:  "InnerService",
					},
					{
						pkgImport: "github.com/MyNihongo/inject/examples/pkg3",
						typeName:  "Service3",
					},
				},
				injectType: Transient,
			},
			"InnerService": {
				name:       "GetInnerService",
				paramTypes: []*typeDecl{},
				injectType: Transient,
			},
		},
		"github.com/MyNihongo/inject/examples/pkg3": {
			"Service3": {
				name:       "GetService3",
				paramTypes: []*typeDecl{},
				injectType: Transient,
			},
		},
	}

	ctx, wd := context.Background(), getExamplesWd()
	fixture, _ := loadFileContent(wd, "serviceCollection.go")

	got, err := getDefinitions(ctx, wd, fixture)

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}
