package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	const wd = ``
	mocks := []string{
		"impl1:Impl1Service",
	}

	execute(wd, mocks, 0)
}

func TestTypeDeclWithImport(t *testing.T) {
	want := &typeDecl{
		pkgImport: "github.com/MyNihongo/mockgen/examples/pkg2",
		typeName:  "Service1",
	}

	got := getTypeDeclarationFromString("github.com/MyNihongo/mockgen/examples/pkg2.Service1")

	assert.Equal(t, want, got)
}

func TestTypeDeclWithoutImport(t *testing.T) {
	want := &typeDecl{
		pkgImport: "",
		typeName:  "int64",
	}

	got := getTypeDeclarationFromString("int64")

	assert.Equal(t, want, got)
}
