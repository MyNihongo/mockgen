package main

import (
	"go/format"
	"testing"

	gen "github.com/MyNihongo/codegen"
	"github.com/stretchr/testify/assert"
)

const pkgName = `example`

func formatFile(file *gen.File) string {
	res, _ := format.Source([]byte(file.GoString()))
	return string(res)
}

func TestGenerateOneService(t *testing.T) {
	const want = ``
	fixture := []*mockDecl{
		{
			mockNameDecl: &mockNameDecl{
				typeName:      "impl1",
				interfaceName: "Impl1Service",
			},
			fields: []*fieldDecl{
				{
					name: "ser1",
					typeDecl: &typeDecl{
						pkgImport: "github.com/MyNihongo/mockgen/examples/pkg1",
						typeName:  "Service1",
					},
				},
			},
		},
	}

	file, err := generateMocks(pkgName, fixture)
	got := formatFile(file)

	assert.NotNil(t, err)
	assert.Equal(t, want, got)
}
