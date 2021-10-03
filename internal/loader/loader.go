package loader

import (
	"fmt"
	"go/types"
)

func NewDeclProvider() *declProvider {
	return &declProvider{
		mapping:      make(map[TypeDecl]bool),
		scopeMapping: make(map[string]*types.Scope),
	}
}

func (d *declProvider) TryGetMock(wd string, typeDecl *TypeDecl) ([]*MethodDecl, bool) {
	if _, ok := d.mapping[*typeDecl]; ok {
		return nil, false
	} else if scope, ok := d.getPackageScope(wd, typeDecl); !ok {
		return nil, false
	} else if typeObj := scope.Lookup(typeDecl.typeName); typeObj == nil {
		fmt.Printf("cannot find a type %s\n", typeDecl.typeName)
		return nil, false
	} else if interfaceTypeObj, ok := typeObj.Type().Underlying().(*types.Interface); !ok {
		fmt.Printf("type %s is not an interface\n", typeDecl.typeName)
		return nil, true
	} else {
		funcs := make([]*MethodDecl, interfaceTypeObj.NumMethods())

		for i := 0; i < interfaceTypeObj.NumMethods(); i++ {
			if signature, ok := interfaceTypeObj.Method(i).Type().(*types.Signature); !ok {
				fmt.Println("function type is not a signature")
				return nil, false
			} else {
				funcs[i] = &MethodDecl{
					name:    interfaceTypeObj.Method(i).Name(),
					params:  getParams(signature),
					returns: getReturns(signature),
				}
			}
		}

		d.mapping[*typeDecl] = true
		return funcs, true
	}
}

func (d *declProvider) getPackageScope(wd string, typeDecl *TypeDecl) (*types.Scope, bool) {
	var scope *types.Scope
	var ok bool

	if scope, ok = d.scopeMapping[typeDecl.pkgImport]; !ok {
		var err error
		if scope, err = LoadPackageScope(wd, typeDecl.pkgImport); err != nil {
			fmt.Println(err)
			return nil, false
		} else {
			d.scopeMapping[typeDecl.pkgImport] = scope
		}
	}

	return scope, true
}

func getParams(signature *types.Signature) []*paramDecl {
	params := signature.Params()
	decls := make([]*paramDecl, params.Len())

	for i := 0; i < params.Len(); i++ {
		param := params.At(i)

		decls[i] = &paramDecl{
			name:     param.Name(),
			TypeDecl: GetTypeDeclaration(param.Type()),
		}
	}

	return decls
}

func getReturns(signature *types.Signature) []*TypeDecl {
	results := signature.Results()
	decls := make([]*TypeDecl, results.Len())

	for i := 0; i < results.Len(); i++ {
		decls[i] = GetTypeDeclaration(results.At(i).Type())
	}

	return decls
}
