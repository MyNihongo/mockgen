package main

import (
	"fmt"
	"strings"

	gen "github.com/MyNihongo/codegen"
)

const (
	assertExpectationsName = "AssertExpectations"
	fixture                = "fixture"
	ret                    = "ret"
)

// generateMocks generates the complete code for all mocks
func generateMocks(wd, pkgName string, mocks []*mockDecl) (*gen.File, error) {
	file := gen.NewFile(pkgName, "my-nihongo-mockgen")
	file.Imports(
		gen.Import("testing"),
		gen.Import("github.com/stretchr/testify/mock"),
	)

	declProvider := NewDeclProvider()

	for _, mock := range mocks {
		fixtureName := createFixtureTypeName(mock.mockNameDecl)
		createFixtureName := fmt.Sprintf("create%s", strings.Title(fixtureName))

		// fixture
		fixtureStruct := file.Struct(fixtureName)

		// AssertExpectations
		file.CommentF("%s asserts that everything specified with On and Return was in fact called as expected. Calls may have occurred in any order.", assertExpectationsName)
		assertExpectationsFunc := file.Method(
			gen.This(fixtureName).Pointer(),
			assertExpectationsName,
		).Params(
			gen.QualParam("t", "testing", "T"),
		)

		// createFixture
		file.CommentF("%s creates a new fixture will all mocks", createFixtureName)
		createFixtureFunc := file.Func(createFixtureName).ReturnTypes(
			createFixtureReturnType(mock.mockNameDecl),
			gen.ReturnType(fixtureName).Pointer(),
		)
		initFixture := gen.InitStruct(mock.typeName).Address()
		initFixtureStmt := gen.Declare(fixture).Values(initFixture)

		initMocks := gen.InitStruct(fixtureName).Address()

		for _, field := range mock.fields {
			if methods, ok := declProvider.TryGetMock(wd, field.typeDecl); !ok {
				continue
			} else {
				fieldName, mockTypeName := field.name, fmt.Sprintf("Mock%s", field.typeName)
				generateMock(file, field, mockTypeName, methods)

				fixtureStruct.AddProp(fieldName, mockTypeName).Pointer()

				assertExpectationsFunc.AddStatement(
					gen.Identifier("f").Field(fieldName).Call(assertExpectationsName).Args(gen.Identifier("t")),
				)

				createFixtureFunc.AddStatement(
					gen.Declare(fieldName).Values(gen.New(mockTypeName)),
				)

				initFixture.AddPropValue(fieldName, gen.Identifier(fieldName))
				initMocks.AddPropValue(fieldName, gen.Identifier(fieldName))
			}
		}

		createFixtureFunc.AddStatement(initFixtureStmt)
		createFixtureFunc.AddStatement(
			gen.Return(gen.Identifier(fixture), initMocks),
		)
	}

	return file, nil
}

func generateMock(file *gen.File, field *fieldDecl, mockName string, methods []*methodDecl) {
	file.Struct(mockName).Props(
		gen.QualEmbeddedProperty("mock", "Mock"),
	)

	for _, method := range methods {
		argsVal := gen.Identifier("m").Call("Called")
		returnValues := make([]gen.Value, len(method.returns))

		params := make([]*gen.ParamDecl, len(method.params))
		returns := make([]*gen.ReturnTypeDecl, len(method.returns))

		// Params

		// Returns
		for i, returnType := range method.returns {
			var alias string
			if len(returnType.pkgImport) != 0 {
				file.AddImport(returnType.pkgImport)
				alias = getImportAlias(returnType.pkgImport)
			}

			returns[i] = gen.QualReturnType(
				alias,
				returnType.typeName,
			)
		}

		file.Method(
			gen.This(mockName),
			method.name,
		).Params(params...).ReturnTypes(returns...).Block(
			gen.Declare(ret).Values(argsVal),
			gen.Return(returnValues...),
		)
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

func getImportAlias(pkgImport string) string {
	if index := strings.LastIndexByte(pkgImport, '/'); index == -1 {
		return pkgImport
	} else {
		return pkgImport[index+1:]
	}
}
