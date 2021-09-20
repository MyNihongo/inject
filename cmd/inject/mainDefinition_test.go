package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
