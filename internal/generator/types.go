package generator

import (
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
