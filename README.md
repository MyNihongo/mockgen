## Mockgen: Mock generation in Go
Mockgen is an utility for [testify](https://github.com/stretchr/testify) which generates the mundane mock initialisation code such as:

- mock classes of all services of a struct;
- type-safe setup methods (`On("FunctionName")` and `Return()`);
- `AssertExpectations` for all services of a struct.

### Installing
Install the CLI for code generation
```go
go get github.com/MyNihongo/mockgen/cmd/mockgen
```
### Add `go generate` for required mocks
```go
//go:generate mockgen impl1:Impl1Service
package mocking

import (
	"github.com/MyNihongo/mockgen/mocking/pkg1"
	"github.com/MyNihongo/mockgen/mocking/pkg2"
)

type impl1 struct {
	ser1 pkg1.Service1_1
	ser2 pkg2.Service2_1
}

type Impl1Service interface {
	Foo()
}

func (i *impl1) Foo() {
	i.ser1.Foo("string", 12)
	i.ser2.Foo("string1", "string2")
}
```
This will generate quite a lot of code, so please go to [mock_gen_test.go](examples/mock_gen_test.go) to see what is actually generated.

#### Command samples
Unlimited number of services can be passed for generation.  
A single entry may be `structType` or `structType:interfaceType`.
```sh
mockgen service1:Service1 service2:Service2 service3:Service3
```
Generate mocks for *service1*, *service2* and *service3* with interfaces.
```sh
mockgen service1:Service1 service2
```
Generate a mock for *service1* with its interface; Generate a mock for *service2* without an interface.

## Keep the package reference
Command `go mod tidy` will remove the reference of Mockgen by default. In order to keep the reference create a dummy file somewhere and import a dummy type.
```go
package pkg

import "github.com/MyNihongo/mockgen"

var _ mockgen.Reference
```