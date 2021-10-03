package main

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"

	gen "github.com/MyNihongo/mockgen/internal/generator"
	"github.com/MyNihongo/mockgen/internal/loader"
)

func main() {
	if wd, err := os.Getwd(); err != nil {
		fmt.Println(err)
	} else {
		execute(wd, os.Args, 1)
	}
}

func execute(wd string, mocks []string, offset int) error {
	if pkg, err := loader.LoadPackage(wd); err != nil {
		return err
	} else {
		mockDecls := make([]*gen.MockDecl, len(mocks)-offset)

		for i := offset; i < len(mocks); i++ {
			mockName := gen.GetMockName(mocks[i])

			if typeObj := pkg.Types.Scope().Lookup(mockName.TypeName()); typeObj == nil {
				return fmt.Errorf("type %s is not found in %s", mockName.TypeName(), wd)
			} else if structType, ok := typeObj.Type().Underlying().(*types.Struct); !ok {
				return fmt.Errorf("type %s is not a struct", mockName.TypeName())
			} else {
				fields := make([]*gen.FieldDecl, structType.NumFields())

				for j := 0; j < structType.NumFields(); j++ {
					field := structType.Field(j)

					fields[j] = gen.NewFieldDecl(
						field.Name(),
						loader.GetTypeDeclaration(field.Type()),
					)
				}

				mockDecls[i-offset] = gen.NewMockDecl(
					mockName,
					fields,
				)
			}
		}

		if file, err := gen.GenerateMocks(wd, pkg.Name, mockDecls); err != nil {
			return err
		} else {
			path := filepath.Join(wd, "mock_gen_test.go")
			return file.Save(path)
		}
	}
}
