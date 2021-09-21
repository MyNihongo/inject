//go:generate inject
package examples

import (
	"github.com/MyNihongo/inject"
	"github.com/MyNihongo/inject/examples/pkg1"
	"github.com/MyNihongo/inject/examples/pkg2"
	"github.com/MyNihongo/inject/examples/pkg3"
)

func BuildServiceProvider() {
	inject.NewServiceCollection().
		AddSingleton(pkg1.GetService1).
		AddTransient(pkg2.GetService2).
		AddTransient(pkg2.GetInnerService).
		AddTransient(pkg3.GetService3)
}
