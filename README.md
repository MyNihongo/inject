## Inject: Compile-time DI initialization in Go
Inject is a compile-time dependency-injection library. Supported injection types are:
- *Singleton* - single instance for the entire application;
- *Transient* - each time a new instance is created.

The library verifies injection consistency (transient instances cannot be injected into singletons because transient will become singleton).
### Installing
1. Install the CLI for code generation
```go
go get github.com/MyNihongo/inject/cmd/inject
```
2. Add the library itself
```go
go get github.com/MyNihongo/inject
```
Ensure that `$GOPATH/bin` is added to your `$PATH`.
### Define a service collection
1. Create a new file with the name `serviceCollection.go`. Inject will look for a file with this name;
2. In the file define a function `BuildServiceProvider`. Inject will look for a function with this name. In this function all services should be declared.
3. Add `//go:generate inject` to the header of the file.
4. Run `go generate .\...`
```go
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
```
This will generate the following code in the file `serviceProvider_gen.go`.
```go
// Code generated by my-nihongo-di. DO NOT EDIT.
package examples

import (
	pkg1 "github.com/MyNihongo/inject/examples/pkg1"
	pkg2 "github.com/MyNihongo/inject/examples/pkg2"
	pkg3 "github.com/MyNihongo/inject/examples/pkg3"
	"sync"
)

var impl_Service1 pkg1.Service1

// ProvideService1 provides a singleton instance of github.com/MyNihongo/inject/examples/pkg1.Service1
func ProvideService1() pkg1.Service1 {
	sync.DoOnce(func() {
		impl_Service1 = pkg1.GetService1()
	})
	return impl_Service1
}

// ProvideInnerService provides a transient instance of github.com/MyNihongo/inject/examples/pkg2.InnerService
func ProvideInnerService() pkg2.InnerService {
	return pkg2.GetInnerService()
}

// ProvideService2 provides a transient instance of github.com/MyNihongo/inject/examples/pkg2.Service2
func ProvideService2() pkg2.Service2 {
	p0 := ProvideInnerService()
	p1 := ProvideService3()
	return pkg2.GetService2(p0, p1)
}

// ProvideService3 provides a transient instance of github.com/MyNihongo/inject/examples/pkg3.Service3
func ProvideService3() pkg3.Service3 {
	return pkg3.GetService3()
}

```