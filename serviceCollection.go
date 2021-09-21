package inject

type serviceCollection interface {
	AddSingleton(impl interface{}) serviceCollection
	AddTransient(impl interface{}) serviceCollection
}

func NewServiceCollection() serviceCollection {
	panic("not implemented, run code generation")
}
