package cli

import (
	"testing"

	"github.com/SeerUK/assert"
)

func TestCreateApplication(t *testing.T) {
	t.Run("should create an application, with the right name", func(t *testing.T) {
		app := CreateApplication()

		assert.Equal(t, "foldup", app.Name)
	})
}
