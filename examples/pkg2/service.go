package pkg2

import (
	"context"
	"time"
)

type MyOutput struct {
	val bool
}

type MyInput struct {
	val1 int
	val2 string
}

type Service2_1 interface {
	Foo(arg1, arg2 string) (string, int, error)
}

type Service2_2 interface {
	Now(ctx context.Context) time.Time
}

type Service2_3 interface {
	Foo(param *MyInput) *MyOutput
}
