package main

import "fmt"

type declProvider struct {
	mapping map[typeDecl]string
}

func (d *declProvider) TryGetMock(wd string, typeDecl *typeDecl) (string, bool) {
	if _, ok := d.mapping[*typeDecl]; ok {
		return "", false
	} else if scope, err := loadPackageScope(wd, typeDecl.pkgImport); err != nil {
		fmt.Println(err)
		return "", false
	} else if typeObj := scope.Lookup(typeDecl.typeName); typeObj == nil {
		fmt.Printf("cannot find a type %s\n", typeDecl.typeName)
		return "", false
	} else {
		return "", true
	}
}
