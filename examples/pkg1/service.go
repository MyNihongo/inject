package pkg1

type Service1 interface {
	Foo()
}

type impl struct {
}

func GetService1() Service1 {
	return &impl{}
}

func (i *impl) Foo() {}
