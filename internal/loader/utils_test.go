package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeDeclWithImport(t *testing.T) {
	want := &TypeDecl{
		pkgImport: "github.com/MyNihongo/mockgen/examples/pkg2",
		typeName:  "Service1",
		isPointer: false,
	}

	got := getTypeDeclarationFromString("github.com/MyNihongo/mockgen/examples/pkg2.Service1")

	assert.Equal(t, want, got)
}

func TestTypeDeclWithImportPointer(t *testing.T) {
	want := &TypeDecl{
		pkgImport: "github.com/MyNihongo/mockgen/examples/pkg2",
		typeName:  "Service1",
		isPointer: true,
	}

	got := getTypeDeclarationFromString("*github.com/MyNihongo/mockgen/examples/pkg2.Service1")

	assert.Equal(t, want, got)
}

func TestTypeDeclWithoutImport(t *testing.T) {
	want := &TypeDecl{
		pkgImport: "",
		typeName:  "int64",
		isPointer: false,
	}

	got := getTypeDeclarationFromString("int64")

	assert.Equal(t, want, got)
}

func TestTypeDeclWithoutImportPointer(t *testing.T) {
	want := &TypeDecl{
		pkgImport: "",
		typeName:  "int64",
		isPointer: true,
	}

	got := getTypeDeclarationFromString("*int64")

	assert.Equal(t, want, got)
}
