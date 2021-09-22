package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileTransientSingle(t *testing.T) {
	const want = `// Code generated by my-nihongo-di. DO NOT EDIT.
package di
import pkg1 "github.com/MyNihongo/inject/examples/pkg1"
// ProvideService1 provides a transient instance of github.com/MyNihongo/inject/examples/pkg1.Service1
func ProvideService1()pkg1.Service1{
return pkg1.GetService1()
}
`
	fixture := map[string]*pkgFuncs{
		"github.com/MyNihongo/inject/examples/pkg1": {
			alias: "pkg1",
			funcs: map[string]*funcDecl{
				"Service1": {
					name:       "GetService1",
					injectType: Transient,
					paramTypes: []*typeDecl{},
				},
			},
		},
	}

	file, err := generateServiceProvider("di", fixture)
	got := file.GoString()

	assert.Nil(t, err, want, got)
	assert.Equal(t, want, got)
}

func TestFileSingletonSingle(t *testing.T) {
	const want = `// Code generated by my-nihongo-di. DO NOT EDIT.
package di
import pkg1 "github.com/MyNihongo/inject/examples/pkg1"
var impl_Service1 pkg1.Service1
// ProvideService1 provides a singleton instance of github.com/MyNihongo/inject/examples/pkg1.Service1
func ProvideService1()pkg1.Service1{
sync.DoOnce(func (){
impl_Service1=pkg1.GetService1()
})
return impl_Service1
}
`
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
	assert.Equal(t, want, got)
}

func TestFileErrorIfParamNotFound(t *testing.T) {
	fixture := map[string]*pkgFuncs{
		"github.com/MyNihongo/inject/examples/pkg1": {
			alias: "pkg1",
			funcs: map[string]*funcDecl{
				"Service1": {
					name:       "GetService1",
					injectType: Singleton,
					paramTypes: []*typeDecl{
						{
							pkgImport: "github.com/MyNihongo/inject/examples/pkg2",
							typeName:  "Service2",
						},
					},
				},
			},
		},
	}

	file, err := generateServiceProvider("di", fixture)

	assert.Nil(t, file)
	assert.Error(t, err, "test")
}
