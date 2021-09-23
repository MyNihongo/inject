package pkg3

type Service3 interface {
	Goo()
}

type impl struct {
}

func GetService3() Service3 {
	return &impl{}
}

func (i *impl) Goo() {
}
