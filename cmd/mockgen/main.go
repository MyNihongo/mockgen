package main

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/MyNihongo/mockgen/internal/loader"
)

func main() {
	if wd, err := os.Getwd(); err != nil {
		fmt.Println(err)
	} else {
		execute(wd, os.Args, 1)
	}
}

type mockDecl struct {
	*mockNameDecl
	fields []*fieldDecl
}

type mockNameDecl struct {
	typeName      string
	interfaceName string
}

type fieldDecl struct {
	name string
	*loader.TypeDecl
}

func execute(wd string, mocks []string, offset int) error {
	if pkg, err := loader.LoadPackage(wd); err != nil {
		return err
	} else {
		mockDecls := make([]*mockDecl, len(mocks)-offset)

		for i := offset; i < len(mocks); i++ {
			mockName := getMockName(mocks[i])

			if typeObj := pkg.Types.Scope().Lookup(mockName.typeName); typeObj == nil {
				return fmt.Errorf("type %s is not found in %s", mockName.typeName, wd)
			} else if structType, ok := typeObj.Type().Underlying().(*types.Struct); !ok {
				return fmt.Errorf("type %s is not a struct", mockName.typeName)
			} else {
				fields := make([]*fieldDecl, structType.NumFields())

				for j := 0; j < structType.NumFields(); j++ {
					field := structType.Field(j)

					fields[j] = &fieldDecl{
						name:     field.Name(),
						TypeDecl: loader.GetTypeDeclaration(field.Type()),
					}
				}

				mockDecls[i-offset] = &mockDecl{
					mockNameDecl: mockName,
					fields:       fields,
				}
			}
		}

		if file, err := generateMocks(wd, pkg.Name, mockDecls); err != nil {
			return err
		} else {
			path := filepath.Join(wd, "mock_gen_test.go")
			return file.Save(path)
		}
	}
}

func getMockName(mock string) *mockNameDecl {
	separatorIndex := strings.IndexByte(mock, ':')

	if separatorIndex == -1 {
		return &mockNameDecl{
			typeName: mock,
		}
	} else {
		return &mockNameDecl{
			typeName:      mock[:separatorIndex],
			interfaceName: mock[separatorIndex+1:],
		}
	}
}
