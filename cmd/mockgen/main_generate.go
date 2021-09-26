package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	gen "github.com/MyNihongo/codegen"
)

const (
	mockContainerName      = "mockContainer"
	assertExpectationsName = "AssertExpectations"
)

func generateMocks(wd, pkgName string, mocks []*mockDecl) (*gen.File, error) {
	file := gen.NewFile(pkgName, "my-nihongo-mockgen")
	file.Imports(
		gen.Import("testing"),
		gen.Import("github.com/stretchr/testify/mock"),
	)

	mockContainer := file.Struct(mockContainerName)

	file.CommentF("%s asserts that everything specified with On and Return was in fact called as expected. Calls may have occurred in any order.", assertExpectationsName)
	assertExpectationsFunc := file.Method(
		gen.This(mockContainerName).Pointer(),
		assertExpectationsName,
	).Params(
		gen.QualParam("t", "testing", "T"),
	)

	for _, mock := range mocks {
		mockTypeName := createMockTypeName(mock.mockNameDecl)
		fieldName := LowerFirst(mockTypeName)

		mockContainer.AddProp(gen.Property(fieldName, mockTypeName).Pointer())
		assertExpectationsFunc.AddStatement(gen.Identifier("m").Field(fieldName).Call(assertExpectationsName).Args(gen.Identifier("t")))

		file.Struct(mockTypeName).Props(
			gen.QualEmbeddedProperty("mock", "Mock"),
		)
	}

	return file, nil
}

func createMockTypeName(mockName *mockNameDecl) string {
	var name string
	if len(mockName.interfaceName) != 0 {
		name = mockName.interfaceName
	} else {
		name = mockName.typeName
	}

	return fmt.Sprintf("Mock%s", strings.Title(name))
}

func LowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}
