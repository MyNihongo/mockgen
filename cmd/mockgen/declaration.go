package main

import (
	"fmt"
	"go/types"
)

type paramDecl struct {
	name string
	*typeDecl
}

type methodDecl struct {
	params  []*paramDecl
	returns []*typeDecl
}

type declProvider struct {
	mapping      map[typeDecl]bool
	scopeMapping map[string]*types.Scope
}

func NewDeclProvider() *declProvider {
	return &declProvider{
		mapping:      make(map[typeDecl]bool),
		scopeMapping: make(map[string]*types.Scope),
	}
}

func (d *declProvider) TryGetMock(wd string, typeDecl *typeDecl) ([]*methodDecl, bool) {
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
		funcs := make([]*methodDecl, interfaceTypeObj.NumMethods())

		for i := 0; i < interfaceTypeObj.NumMethods(); i++ {
			if signature, ok := interfaceTypeObj.Method(i).Type().(*types.Signature); !ok {
				fmt.Println("function type is not a signature")
				return nil, false
			} else {
				funcs[i] = &methodDecl{
					params:  getParams(signature),
					returns: getReturns(signature),
				}
			}
		}

		d.mapping[*typeDecl] = true
		return funcs, true
	}
}

func (d *declProvider) getPackageScope(wd string, typeDecl *typeDecl) (*types.Scope, bool) {
	var scope *types.Scope
	var ok bool

	if scope, ok = d.scopeMapping[typeDecl.pkgImport]; !ok {
		var err error
		if scope, err = loadPackageScope(wd, typeDecl.pkgImport); err != nil {
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
			typeDecl: getTypeDeclaration(param.Type()),
		}
	}

	return decls
}

func getReturns(signature *types.Signature) []*typeDecl {
	results := signature.Results()
	decls := make([]*typeDecl, results.Len())

	for i := 0; i < results.Len(); i++ {
		decls[i] = getTypeDeclaration(results.At(i).Type())
	}

	return decls
}
