package inject

type ServiceCollection interface {
	AddSingleton(impl interface{}) ServiceCollection
	AddTransient(impl interface{}) ServiceCollection
}

func NewServiceCollection() ServiceCollection {
	panic("not implemented, run code generation")
}
