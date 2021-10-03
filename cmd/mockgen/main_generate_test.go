package main

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"

	gen "github.com/MyNihongo/codegen"
	"github.com/MyNihongo/mockgen/internal/loader"
	"github.com/stretchr/testify/assert"
)

const pkgName = `example`

func formatFile(file *gen.File) string {
	res, _ := format.Source([]byte(file.GoString()))
	return string(res)
}

func getWd() string {
	wd, _ := os.Getwd()
	index := strings.LastIndex(wd, "cmd")

	return filepath.Join(wd[:index], "examples")
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
					TypeDecl: loader.NewTypeDecl(
						"github.com/MyNihongo/mockgen/examples/pkg1",
						"Service1",
					),
				},
			},
		},
	}

	wd := getWd()
	file, err := generateMocks(wd, pkgName, fixture)
	got := formatFile(file)

	assert.NotNil(t, err)
	assert.Equal(t, want, got)
}
