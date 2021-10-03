package main

import (
	"testing"

	"github.com/MyNihongo/mockgen/internal/test"
	"github.com/stretchr/testify/assert"
)

func generateMocks(t *testing.T, mocks []string) string {
	wd := test.GetWd("cmd")
	res, err := execute(wd, mocks, 0)

	assert.Nil(t, err)
	return test.FormatFile(res.file)
}

func TestExecute(t *testing.T) {
	const wd = ``
	mocks := []string{}

	res, _ := execute(wd, mocks, 0)
	save(res.file, wd)
}

func TestImplementationOne(t *testing.T) {
	const want = ``
	mocks := []string{
		"impl1:Impl1Service",
	}

	got := generateMocks(t, mocks)
	assert.Equal(t, want, got)
}
