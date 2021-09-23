package pkg2

import (
	"github.com/MyNihongo/inject/di/pkg3"
)

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

func GetService2(innerSrv InnerService, aa pkg3.Service3) Service2 {
	return &impl{
		innerSrv: innerSrv,
	}
}

func (i *impl) Foo() {
	i.innerSrv.Boo()
}

func (i *innerImpl) Boo() {
}
