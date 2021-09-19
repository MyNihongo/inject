//go:generate inject
package examples

import (
	"github.com/MyNihongo/inject"
	"github.com/MyNihongo/inject/examples/pkg1"
	"github.com/MyNihongo/inject/examples/pkg2"
)

func BuildServiceProvider() {
	inject.NewServiceCollection().
		AddSingleton(pkg1.GetService1).
		AddTransient(pkg2.GetService2)
}