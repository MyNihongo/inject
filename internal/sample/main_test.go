//go:generate inject
package sample

import (
	"github.com/MyNihongo/inject"
	"github.com/MyNihongo/inject/internal/sample/pkg1"
	"github.com/MyNihongo/inject/internal/sample/pkg2"
)

func BuildServiceProvider() {
	inject.NewServiceCollection().
		AddSingleton(pkg1.GetService1).
		AddTransient(pkg2.GetService2)
}
