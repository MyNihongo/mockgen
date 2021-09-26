package main

import (
	"fmt"
	"go/types"
)

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

func (d *declProvider) TryGetMock(wd string, typeDecl *typeDecl) (string, bool) {
	if _, ok := d.mapping[*typeDecl]; ok {
		return "", false
	} else if scope, ok := d.getPackageScope(wd, typeDecl); !ok {
		return "", false
	} else if typeObj := scope.Lookup(typeDecl.typeName); typeObj == nil {
		fmt.Printf("cannot find a type %s\n", typeDecl.typeName)
		return "", false
	} else if interfaceTypeObj, ok := typeObj.Type().Underlying().(*types.Interface); !ok {
		fmt.Printf("type %s is not an interface\n", typeDecl.typeName)
		return "", true
	} else {
		for i := 0; i < interfaceTypeObj.NumMethods(); i++ {
			a := interfaceTypeObj.Method(i)
			fmt.Println(a)
		}

		d.mapping[*typeDecl] = true
		return "", true
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
