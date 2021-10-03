package test

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"

	gen "github.com/MyNihongo/codegen"
)

const PkgName = `examples`

func FormatFile(file *gen.File) string {
	res, _ := format.Source([]byte(file.GoString()))
	return string(res)
}

func GetWd(rootDir string) string {
	wd, _ := os.Getwd()
	index := strings.LastIndex(wd, rootDir)

	return filepath.Join(wd[:index], PkgName)
}
