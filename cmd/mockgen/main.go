package main

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"

	"github.com/MyNihongo/codegen"
	gen "github.com/MyNihongo/mockgen/internal/generator"
	"github.com/MyNihongo/mockgen/internal/loader"
)

type execResult struct {
	file    *codegen.File
	pkgName string
}

func main() {
	if wd, err := os.Getwd(); err != nil {
		fmt.Println(err)
	} else if result, err := execute(wd, os.Args, 1); err != nil {
		fmt.Println(err)
	} else if err = save(result.file, wd); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("mock generated: %s\n", result.pkgName)
	}
}

func execute(wd string, mocks []string, offset int) (*execResult, error) {
	if pkg, err := loader.LoadPackage(wd); err != nil {
		return nil, err
	} else {
		mockDecls := make([]*gen.MockDecl, len(mocks)-offset)

		for i := offset; i < len(mocks); i++ {
			mockName := gen.GetMockName(mocks[i])

			if typeObj := pkg.Types.Scope().Lookup(mockName.TypeName()); typeObj == nil {
				return nil, fmt.Errorf("type %s is not found in %s", mockName.TypeName(), wd)
			} else if structType, ok := typeObj.Type().Underlying().(*types.Struct); !ok {
				return nil, fmt.Errorf("type %s is not a struct", mockName.TypeName())
			} else {
				fields := make([]*gen.FieldDecl, 0)

				for j := 0; j < structType.NumFields(); j++ {
					field := structType.Field(j)

					// TODO: maybe in the future process recursively
					if field.Embedded() {
						continue
					}

					fields = append(fields, gen.NewFieldDecl(
						field.Name(),
						loader.GetTypeDeclaration(field.Type()),
					))
				}

				mockDecls[i-offset] = gen.NewMockDecl(
					mockName,
					fields,
				)
			}
		}

		if file, err := gen.GenerateMocks(wd, pkg.Name, mockDecls); err != nil {
			return nil, err
		} else {
			res := &execResult{
				file:    file,
				pkgName: pkg.Name,
			}

			return res, nil
		}
	}
}

func save(file *codegen.File, wd string) error {
	path := filepath.Join(wd, "mock_gen_test.go")
	return file.Save(path)
}
