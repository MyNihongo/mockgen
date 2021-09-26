//go:generate mockgen impl1:Impl1Service impl2:Impl2Service
package examples

import (
	"github.com/MyNihongo/mockgen/examples/pkg1"
	"github.com/MyNihongo/mockgen/examples/pkg2"
)

type impl1 struct {
	ser1 pkg1.Service1
	ser2 pkg2.Service2
}

type impl2 struct {
	ser11 pkg1.Service11
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
