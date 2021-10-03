package generator

import (
	"fmt"
	"strings"

	gen "github.com/MyNihongo/codegen"
	"github.com/MyNihongo/mockgen/internal/loader"
)

func GetMockName(mock string) *mockNameDecl {
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

func createFixtureTypeName(mockName *mockNameDecl) string {
	var name string
	if len(mockName.interfaceName) != 0 {
		name = mockName.interfaceName
	} else {
		name = mockName.typeName
	}

	return fmt.Sprintf("fixture%s", strings.Title(name))
}

func createFixtureReturnType(mockName *mockNameDecl) *gen.ReturnTypeDecl {
	if len(mockName.interfaceName) != 0 {
		return gen.ReturnType(mockName.interfaceName)
	} else {
		return gen.ReturnType(mockName.typeName).Pointer()
	}
}

func addImportAlias(file *gen.File, pkgImport string) string {
	if len(pkgImport) == 0 {
		return ""
	} else {
		file.AddImport(pkgImport)

		if index := strings.LastIndexByte(pkgImport, '/'); index == -1 {
			return pkgImport
		} else {
			return pkgImport[index+1:]
		}
	}
}

func createReturnValue(returnType *loader.TypeDecl, alias string, index int) gen.Value {
	variable, arg := gen.Identifier(ret), gen.Int(index)

	var funcName string
	switch returnType.TypeName() {
	case "error":
		funcName = "Error"
	case "bool":
		funcName = "Bool"
	case "int":
		funcName = "Int"
	case "string":
		funcName = "String"
	}

	if len(funcName) != 0 {
		return variable.Call(funcName).Args(arg)
	} else {
		return variable.Call("Get").Args(arg).
			CastQual(alias, returnType.TypeName())
	}
}
