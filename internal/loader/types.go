package loader

import "go/types"

type paramDecl struct {
	name string
	*TypeDecl
}

type MethodDecl struct {
	name    string
	params  []*paramDecl
	returns []*TypeDecl
}

type declProvider struct {
	mapping      map[TypeDecl]bool
	scopeMapping map[string]*types.Scope
}

type TypeDecl struct {
	pkgImport string
	typeName  string
	isPointer bool
}

func NewTypeDecl(pkgImport, typeName string, isPointer bool) *TypeDecl {
	return &TypeDecl{
		pkgImport: pkgImport,
		typeName:  typeName,
		isPointer: isPointer,
	}
}

func (p *paramDecl) Name() string {
	return p.name
}

func (m *MethodDecl) Name() string {
	return m.name
}

func (m *MethodDecl) Params() []*paramDecl {
	return m.params
}

func (m *MethodDecl) LenParams() int {
	return len(m.params)
}

func (m *MethodDecl) Returns() []*TypeDecl {
	return m.returns
}

func (m *MethodDecl) LenReturns() int {
	return len(m.returns)
}

func (t *TypeDecl) PkgImport() string {
	return t.pkgImport
}

func (t *TypeDecl) TypeName() string {
	return t.typeName
}
