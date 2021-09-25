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

func TestDefinitionSamePackageOne(t *testing.T) {
	want := map[string]*pkgInjections{
		"": {
			alias: "",
			injections: []*injectDecl{
				{
					function:   "createFoo",
					injectType: Singleton,
				},
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
	assert.EqualValues(t, want, got)
}

func TestDefinitionSamePackageMultiple(t *testing.T) {
	want := map[string]*pkgInjections{
		"": {
			alias: "",
			injections: []*injectDecl{
				{
					function:   "createFoo",
					injectType: Singleton,
				},
				{
					function:   "createBoo",
					injectType: Transient,
				},
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
	assert.EqualValues(t, want, got)
}

func TestDefinitionAnotherPackageOne(t *testing.T) {
	want := map[string]*pkgInjections{
		"github.com/MyNihongo/inject/examples/pkg1": {
			alias: "pkg1",
			injections: []*injectDecl{
				{
					function:   "CreateFoo",
					injectType: Singleton,
				},
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
	assert.EqualValues(t, want, got)
}

func TestDefinitionAnotherPackageMultiple(t *testing.T) {
	want := map[string]*pkgInjections{
		"github.com/MyNihongo/inject/examples/pkg1": {
			alias: "my_pkg1",
			injections: []*injectDecl{
				{
					function:   "CreateFoo",
					injectType: Singleton,
				},
				{
					function:   "CreateBoo",
					injectType: Transient,
				},
			},
		},
		"github.com/MyNihongo/inject/examples/pkg2": {
			alias: "not_my_pkg2",
			injections: []*injectDecl{
				{
					function:   "CreateFoo",
					injectType: Transient,
				},
				{
					function:   "CreateBoo",
					injectType: Singleton,
				},
			},
		},
	}

	fixture := &loadResult{
		pkgName: "di",
		imports: map[string]string{
			"my_pkg1":     "github.com/MyNihongo/inject/examples/pkg1",
			"not_my_pkg2": "github.com/MyNihongo/inject/examples/pkg2",
		},
		injects: map[string]injectType{
			"my_pkg1.CreateFoo":     Singleton,
			"my_pkg1.CreateBoo":     Transient,
			"not_my_pkg2.CreateFoo": Transient,
			"not_my_pkg2.CreateBoo": Singleton,
		},
	}

	got, err := getInjectionGrouping(fixture)

	assert.Nil(t, err)
	assert.EqualValues(t, want, got)
}

func TestDefinitionErrorIfNoImport(t *testing.T) {
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
		typeName: typeNameDecl{
			typeName:  "Service1",
			isPointer: false,
		},
	}

	got := getTypeDeclarationString("github.com/MyNihongo/inject/examples/pkg1.Service1")

	assert.Equal(t, want, got)
}

func TestGetTypeDeclarationPointer(t *testing.T) {
	want := &typeDecl{
		pkgImport: "github.com/MyNihongo/inject/examples/pkg1",
		typeName: typeNameDecl{
			typeName:  "Service1",
			isPointer: true,
		},
	}

	got := getTypeDeclarationString("*github.com/MyNihongo/inject/examples/pkg1.Service1")

	assert.Equal(t, want, got)
}

func TestDefinitions(t *testing.T) {
	want := map[string]*pkgFuncs{
		"github.com/MyNihongo/inject/di/pkg1": {
			alias: "pkg1",
			funcs: map[typeNameDecl]*funcDecl{
				{typeName: "Service1", isPointer: false}: {
					name:       "GetService1",
					paramDecls: []*typeDecl{},
					injectType: Singleton,
				},
			},
		},
		"github.com/MyNihongo/inject/di/pkg2": {
			alias: "pkg2",
			funcs: map[typeNameDecl]*funcDecl{
				{typeName: "Service2", isPointer: false}: {
					name: "GetService2",
					paramDecls: []*typeDecl{
						{
							pkgImport: "github.com/MyNihongo/inject/di/pkg2",
							typeName: typeNameDecl{
								typeName:  "InnerService",
								isPointer: false,
							},
						},
						{
							pkgImport: "github.com/MyNihongo/inject/di/pkg3",
							typeName: typeNameDecl{
								typeName:  "Service3",
								isPointer: false,
							},
						},
					},
					injectType: Transient,
				},
				{typeName: "InnerService", isPointer: false}: {
					name:       "GetInnerService",
					paramDecls: []*typeDecl{},
					injectType: Transient,
				},
			},
		},
		"github.com/MyNihongo/inject/di/pkg3": {
			alias: "pkg3",
			funcs: map[typeNameDecl]*funcDecl{
				{typeName: "Service3", isPointer: false}: {
					name:       "GetService3",
					paramDecls: []*typeDecl{},
					injectType: Transient,
				},
			},
		},
	}

	ctx, wd := context.Background(), getExamplesWd()
	fixture, _ := loadFileContent(wd, "serviceCollection.go")

	got, err := getDefinitions(ctx, wd, fixture)

	assert.Nil(t, err)
	assert.EqualValues(t, want, got)
}
