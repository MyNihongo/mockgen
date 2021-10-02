package pkg1

type Service1 interface {
	Foo(arg1 string, arg2 int16) string
	Boo(arg1 string) (uint64, error)
}

type Service11 interface {
	Foo(arg1 string, arg2 int16)
	Boo(arg1 string)
}
