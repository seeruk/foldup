package cli

import (
	"bytes"
	"testing"

	"github.com/SeerUK/assert"
)

func TestCreateApplication(t *testing.T) {
	t.Run("should create an application, with the right name", func(t *testing.T) {
		app := CreateApplication(&bytes.Buffer{})

		assert.Equal(t, "foldup", app.Name)
	})
}
