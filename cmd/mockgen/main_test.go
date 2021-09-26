package main

import (
	"testing"
)

func TestExecute(t *testing.T) {
	const wd = ``
	mocks := []string{
		"impl1:Impl1Service",
	}

	execute(wd, mocks, 0)
}
