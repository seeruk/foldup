package archive

import (
	"testing"

	"github.com/SeerUK/assert"
)

func TestRegisterFormat(t *testing.T) {
	t.Run("should add the given format", func(t *testing.T) {
		expected := len(formats) + 1

		RegisterFormat("test", func(pathname, filename string) (Artifact, error) {
			return nil, nil
		})

		assert.Equal(t, expected, len(formats))

		_, err := findFormatByName("test")
		assert.OK(t, err)
	})
}
