package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeDeclWithImport(t *testing.T) {
	want := &TypeDecl{
		pkgImport: "github.com/MyNihongo/mockgen/examples/pkg2",
		typeName:  "Service1",
	}

	got := getTypeDeclarationFromString("github.com/MyNihongo/mockgen/examples/pkg2.Service1")

	assert.Equal(t, want, got)
}

func TestTypeDeclWithoutImport(t *testing.T) {
	want := &TypeDecl{
		pkgImport: "",
		typeName:  "int64",
	}

	got := getTypeDeclarationFromString("int64")

	assert.Equal(t, want, got)
}
