package pkg2

import (
	"context"
	"time"
)

type Service2_1 interface {
	Foo(arg1, arg2 string) (string, int, error)
}

type Service2_2 interface {
	Now(ctx context.Context) time.Time
}
