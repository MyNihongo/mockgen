package loader

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

func LoadPackage(wd string, patterns ...string) (*packages.Package, error) {
	cfg := &packages.Config{
		Dir:  wd,
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedName,
	}

	if packages, err := packages.Load(cfg, patterns...); err != nil {
		return nil, err
	} else if len(packages) != 1 {
		return nil, fmt.Errorf("cannot identify a unique package in %s", wd)
	} else {
		return packages[0], nil
	}
}

func LoadPackageScope(wd string, patterns ...string) (*types.Scope, error) {
	if pkg, err := LoadPackage(wd, patterns...); err != nil {
		return nil, err
	} else {
		return pkg.Types.Scope(), nil
	}
}

func GetTypeDeclaration(typeName types.Type) *TypeDecl {
	return getTypeDeclarationFromString(typeName.String())
}

func getTypeDeclarationFromString(strVal string) *TypeDecl {
	if typeSeparator := strings.LastIndexByte(strVal, '.'); typeSeparator == -1 {
		return &TypeDecl{
			typeName: strVal,
		}
	} else {
		return &TypeDecl{
			pkgImport: strVal[:typeSeparator],
			typeName:  strVal[typeSeparator+1:],
		}
	}
}
