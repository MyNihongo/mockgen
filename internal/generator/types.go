package generator

import (
	gen "github.com/MyNihongo/codegen"
	"github.com/MyNihongo/mockgen/internal/loader"
)

type MockDecl struct {
	*mockNameDecl
	fields []*FieldDecl
}

type mockNameDecl struct {
	typeName      string
	interfaceName string
}

type FieldDecl struct {
	name string
	*loader.TypeDecl
}

type methodValues struct {
	method       *loader.MethodDecl
	mockName     string
	params       []*gen.ParamDecl
	args         []gen.Value
	returns      []*gen.ReturnTypeDecl
	returnValues []gen.Value
}

func NewMockDecl(name *mockNameDecl, field []*FieldDecl) *MockDecl {
	return &MockDecl{
		mockNameDecl: name,
		fields:       field,
	}
}

func NewFieldDecl(name string, typeDecl *loader.TypeDecl) *FieldDecl {
	return &FieldDecl{
		name:     name,
		TypeDecl: typeDecl,
	}
}

func (m *mockNameDecl) TypeName() string {
	return m.typeName
}

func (m *mockNameDecl) InterfaceName() string {
	return m.interfaceName
}
