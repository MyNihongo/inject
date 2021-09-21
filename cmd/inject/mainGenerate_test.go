package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileSingle(t *testing.T) {
	const want = ``

	fixture := map[string]*pkgFuncs{
		"github.com/MyNihongo/inject/examples/pkg1": {
			alias: "pkg1",
			funcs: map[string]*funcDecl{
				"Service1": {
					name:       "GetService1",
					injectType: Singleton,
					paramTypes: []*typeDecl{},
				},
			},
		},
	}

	file, err := generateServiceProvider("di", fixture)
	got := file.GoString()

	assert.Nil(t, err, want, got)
}
