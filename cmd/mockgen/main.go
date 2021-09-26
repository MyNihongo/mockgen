package main

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
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
	*typeDecl
}

type typeDecl struct {
	pkgImport string
	typeName  string
}

func execute(wd string, mocks []string, offset int) error {
	if scope, err := loadPackageScope(wd); err != nil {
		return err
	} else {
		mockDecls := make([]*mockDecl, len(mocks)-offset)

		for i := offset; i < len(mocks); i++ {
			mockName := getMockName(mocks[i])

			if typeObj := scope.Lookup(mockName.typeName); typeObj == nil {
				return fmt.Errorf("type %s is not found in %s", mockName.typeName, wd)
			} else if structType, ok := typeObj.Type().Underlying().(*types.Struct); !ok {
				return fmt.Errorf("type %s is not a struct", mockName.typeName)
			} else {
				fields := make([]*fieldDecl, structType.NumFields())

				for j := 0; j < structType.NumFields(); j++ {
					field := structType.Field(j)

					fields[j] = &fieldDecl{
						name:     field.Name(),
						typeDecl: getTypeDeclaration(field.Type()),
					}
				}

				mockDecls[i-offset] = &mockDecl{
					mockNameDecl: mockName,
					fields:       fields,
				}
			}
		}

		if file, err := generateMocks(mockDecls); err != nil {
			return err
		} else {
			path := filepath.Join(wd, "mock_gen.go")
			return file.Save(path)
		}
	}
}

func loadPackageScope(wd string) (*types.Scope, error) {
	cfg := &packages.Config{
		Dir:  wd,
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps,
	}

	if packages, err := packages.Load(cfg); err != nil {
		return nil, err
	} else if len(packages) != 1 {
		return nil, fmt.Errorf("cannot identify a unique package in %s", wd)
	} else {
		return packages[0].Types.Scope(), nil
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

func getTypeDeclaration(typeName types.Type) *typeDecl {
	strVal := typeName.String()
	typeSeparator := strings.LastIndexByte(strVal, '.')

	return &typeDecl{
		pkgImport: strVal[:typeSeparator],
		typeName:  strVal[typeSeparator+1:],
	}
}
