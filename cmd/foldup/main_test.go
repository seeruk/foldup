package main

import (
	"bytes"
	"testing"

	"github.com/SeerUK/assert"
)

func TestFoldup(t *testing.T) {
	var result int

	writer = &bytes.Buffer{}
	exit = func(code int) {
		result = code
	}

	main()

	assert.Equal(t, 100, result)
}
