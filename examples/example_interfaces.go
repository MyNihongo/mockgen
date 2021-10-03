package mocking

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
