//go:generate inject
package di

import (
	"github.com/MyNihongo/inject/di/pkg1"

	"github.com/MyNihongo/inject/di/pkg2"
	"github.com/MyNihongo/inject/di/pkg3"

	"github.com/MyNihongo/inject"
)

func BuildServiceProvider() {
	inject.NewServiceCollection().
		AddSingleton(pkg1.GetService1).
		AddTransient(pkg2.GetService2).
		AddTransient(pkg2.GetInnerService).
		AddTransient(pkg3.GetService3)
}
