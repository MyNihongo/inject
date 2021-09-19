package pkg2

type Service2 interface {
	Foo()
}

type impl struct {
}

func GetService2() Service2 {
	return &impl{}
}

func (i *impl) Foo() {}
