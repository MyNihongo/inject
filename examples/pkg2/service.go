package pkg2

type Service2 interface {
	Foo()
}

type InnerService interface {
	Boo()
}

type impl struct {
	innerSrv InnerService
}

type innerImpl struct {
}

func GetInnerService() InnerService {
	return &innerImpl{}
}

func GetService2(innerSrv InnerService) Service2 {
	return &impl{
		innerSrv: innerSrv,
	}
}

func (i *impl) Foo() {
	i.innerSrv.Boo()
}

func (i *innerImpl) Boo() {
}
