package generator

import (
	"errors"
	"fmt"
	"strings"

	gen "github.com/MyNihongo/codegen"
	"github.com/MyNihongo/mockgen/internal/loader"
)

const (
	assertExpectationsName = "AssertExpectations"
	fixture                = "fixture"
	ret                    = "ret"
	mockThis               = "m"
	call                   = "call"
)

// GenerateMocks generates the complete code for all mocks
func GenerateMocks(wd, pkgName string, mocks []*MockDecl) (*gen.File, error) {
	if len(mocks) == 0 {
		return nil, errors.New("no mocks found")
	}

	file := gen.NewFile(pkgName, "my-nihongo-mockgen")
	file.Imports(
		gen.Import("testing"),
		gen.Import("github.com/stretchr/testify/mock"),
	)

	declProvider := loader.NewDeclProvider()

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
			gen.QualParam("t", "testing", "T").Pointer(),
		)

		// createFixture
		file.CommentF("%s creates a new fixture with all mocks", createFixtureName)
		createFixtureFunc := file.Func(createFixtureName).ReturnTypes(
			createFixtureReturnType(mock.mockNameDecl),
			gen.ReturnType(fixtureName).Pointer(),
		)

		initFixture := gen.InitStruct(mock.typeName).Address()
		initFixtureStmt := gen.Declare(fixture).Values(initFixture)
		initMocks := gen.InitStruct(fixtureName).Address()

		for _, field := range mock.fields {
			fieldName, mockTypeName := field.name, fmt.Sprintf("Mock%s", field.TypeName())

			// struct declaration
			fixtureStruct.AddProp(fieldName, mockTypeName).Pointer()

			// init a fixture
			createFixtureFunc.AddStatement(
				gen.Declare(fieldName).Values(gen.New(mockTypeName)),
			)

			initFixture.AddPropValue(fieldName, gen.Identifier(fieldName))
			initMocks.AddPropValue(fieldName, gen.Identifier(fieldName))

			// assert expectations
			assertExpectationsFunc.AddStatement(
				gen.Identifier("f").Field(fieldName).Call(assertExpectationsName).Args(gen.Identifier("t")),
			)

			if methods, ok := declProvider.TryGetMock(wd, field.TypeDecl); ok {
				generateMock(file, field, mockTypeName, methods)
			}
		}

		createFixtureFunc.AddStatement(initFixtureStmt)
		createFixtureFunc.AddStatement(
			gen.Return(gen.Identifier(fixture), initMocks),
		)
	}

	return file, nil
}

func generateMock(file *gen.File, field *FieldDecl, mockName string, methods []*loader.MethodDecl) {
	file.Struct(mockName).Props(
		gen.QualEmbeddedProperty("mock", "Mock"),
	)

	for _, method := range methods {
		params := make([]*gen.ParamDecl, method.LenParams())
		returns := make([]*gen.ReturnTypeDecl, method.LenReturns())

		args := make([]gen.Value, method.LenParams())
		returnValues := make([]gen.Value, method.LenReturns())

		// Params
		for i, param := range method.Params() {
			params[i] = gen.QualParam(
				param.Name(),
				addImportAlias(file, param.PkgImport()),
				param.TypeName(),
			).SetIsPointer(param.IsPointer())

			args[i] = gen.Identifier(param.Name())
		}

		// Returns
		for i, returnType := range method.Returns() {
			alias := addImportAlias(file, returnType.PkgImport())

			returns[i] = gen.QualReturnType(
				alias,
				returnType.TypeName(),
			).SetIsPointer(returnType.IsPointer())

			returnValues[i] = createReturnValue(returnType, alias, i)
		}

		vals := &methodValues{
			method:       method,
			mockName:     mockName,
			params:       params,
			args:         args,
			returns:      returns,
			returnValues: returnValues,
		}

		generateAssertNotCalled(file, vals)
		generateMethodImpl(file, vals)
		generateMethodSetup(file, vals)
	}
}

func generateMethodImpl(file *gen.File, vals *methodValues) {
	var callArgsStmt gen.Stmt
	if callArgsValue := gen.Identifier(mockThis).Call("Called").Args(vals.args...); vals.method.LenReturns() == 0 {
		callArgsStmt = callArgsValue
	} else {
		callArgsStmt = gen.Declare(ret).Values(callArgsValue)
	}

	file.Method(
		gen.This(vals.mockName).Pointer(),
		vals.method.Name(),
	).Params(vals.params...).ReturnTypes(vals.returns...).Block(
		callArgsStmt,
		gen.Return(vals.returnValues...),
	)
}

func generateAssertNotCalled(file *gen.File, vals *methodValues) {
	const (
		t            = "t"
		paramsOffset = 1
		argsOffset   = 2
	)

	params := make([]*gen.ParamDecl, len(vals.params)+paramsOffset)
	args := make([]gen.Value, len(vals.args)+argsOffset)

	params[0] = gen.QualParam(t, "testing", "T").Pointer()
	args[0] = gen.Identifier(t)
	args[1] = gen.String(vals.method.Name())

	for i := 0; i < len(vals.params); i++ {
		params[i+paramsOffset] = vals.params[i]
		args[i+argsOffset] = vals.args[i]
	}

	file.Method(
		gen.This(vals.mockName).Pointer(),
		fmt.Sprintf("Assert%sNotCalled", vals.method.Name()),
	).Params(params...).Block(
		gen.Identifier("m").Call("AssertNotCalled").Args(args...),
	)
}

func generateMethodSetup(file *gen.File, vals *methodValues) {
	args := make([]gen.Value, len(vals.args)+1)
	args[0] = gen.String(vals.method.Name())

	for i, arg := range vals.args {
		args[i+1] = arg
	}

	methodSetup := file.Method(
		gen.This(vals.mockName).Pointer(),
		fmt.Sprintf("On%s", vals.method.Name()),
	).Params(vals.params...)

	var callSetupStmt gen.Stmt
	var returnValues []gen.Value

	if callSetupValue := gen.Identifier(mockThis).Call("On").Args(args...); len(vals.returns) != 0 {
		setupReturnsName := fmt.Sprintf("setup_%s_%s", vals.mockName, vals.method.Name())

		methodSetup.ReturnTypes(
			gen.ReturnType(setupReturnsName).Pointer(),
		)

		callSetupStmt = gen.Declare(call).Values(callSetupValue)
		returnValues = []gen.Value{
			gen.InitStruct(setupReturnsName).Props(
				gen.PropValue(call, gen.Identifier(call)),
			).Address(),
		}

		generateMethodReturnSetup(file, setupReturnsName, vals.returns)
	} else {
		callSetupStmt = callSetupValue
		returnValues = make([]gen.Value, 0)
	}

	methodSetup.Block(
		callSetupStmt,
		gen.Return(returnValues...),
	)
}

func generateMethodReturnSetup(file *gen.File, setupReturnsName string, returns []*gen.ReturnTypeDecl) {
	params := make([]*gen.ParamDecl, len(returns))
	args := make([]gen.Value, len(returns))

	for i, ret := range returns {
		argName := fmt.Sprintf("param%d", i+1)

		args[i] = gen.Identifier(argName)
		params[i] = gen.QualParam(argName, ret.GetTypeAlias(), ret.GetTypeName()).SetIsPointer(ret.GetIsPointer())
	}

	file.Struct(setupReturnsName).Props(
		gen.QualProperty(call, "mock", "Call").Pointer(),
	)

	file.Method(
		gen.This(setupReturnsName).Pointer(),
		"Return",
	).Params(params...).Block(
		gen.Identifier("s").Field(call).Call("Return").Args(args...),
	)
}
