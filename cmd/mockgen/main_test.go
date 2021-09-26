package main

import (
	"testing"
)

func TestExecute(t *testing.T) {
	const wd = ``
	mocks := []string{
		"impl1",
	}

	execute(wd, mocks, 0)
}
