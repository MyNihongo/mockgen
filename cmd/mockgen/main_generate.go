package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	gen "github.com/MyNihongo/codegen"
)

const (
	assertExpectationsName = "AssertExpectations"
)

// generateMocks generates the complete code for all mocks
func generateMocks(wd, pkgName string, mocks []*mockDecl) (*gen.File, error) {
	file := gen.NewFile(pkgName, "my-nihongo-mockgen")
	file.Imports(
		gen.Import("testing"),
		gen.Import("github.com/stretchr/testify/mock"),
	)

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

		for _, field := range mock.fields {
			fieldName, mockName := field.name, fmt.Sprintf("Mock%s", field.typeName)

			fixtureStruct.AddProp(
				gen.Property(fieldName, mockName).Pointer(),
			)

			assertExpectationsFunc.AddStatement(
				gen.Identifier("f").Field(fieldName).Call(assertExpectationsName).Args(gen.Identifier("t")),
			)

			createFixtureFunc.AddStatement(
				// TODO: use the syntax
				gen.Declare(fieldName).Values(gen.Identifier(fmt.Sprintf("new(%s)", mockName))),
			)
		}

		// file.Struct(mockTypeName).Props(
		// 	gen.QualEmbeddedProperty("mock", "Mock"),
		// )
	}

	return file, nil
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

func LowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}
