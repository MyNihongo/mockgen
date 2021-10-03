//go:generate mockgen impl1:Impl1Service impl2:Impl2Service
package mocking

import (
	"github.com/MyNihongo/mockgen/mocking/pkg1"
	"github.com/MyNihongo/mockgen/mocking/pkg2"
)

type impl1 struct {
	ser1 pkg1.Service1_1
	ser2 pkg2.Service2_1
}

type impl2 struct {
	ser11 pkg1.Service1_2
	ser3  pkg2.Service2_1
}

type Impl1Service interface {
	Foo()
}

type Impl2Service interface {
	Boo()
}

func (i *impl1) Foo() {

}

func (i *impl2) Boo() {

}
