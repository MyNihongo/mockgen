package pkg1

type Service1_1 interface {
	Foo(param1 string, param2 int16) string
	Boo(param string) (uint64, error)
}

type Service1_2 interface {
	Foo(param1 string, param2 int16) (int, bool)
	Boo(param string)
}
