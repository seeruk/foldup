package xioutil_test

import (
	"testing"

	"github.com/SeerUK/assert"
	"github.com/SeerUK/foldup/pkg/xioutil"
)

func TestReadDirsInDir(t *testing.T) {
	t.Run("should fail looking to read non-existent directory", func(t *testing.T) {
		_, err := xioutil.ReadDirsInDir("rumpelstilzchen", true)

		assert.NotOK(t, err)
	})

	t.Run("should no fail reading existing directory", func(t *testing.T) {
		_, err := xioutil.ReadDirsInDir("..", true)

		assert.OK(t, err)
	})

	t.Run("should only find directories", func(t *testing.T) {
		dirs, err := xioutil.ReadDirsInDir("./testdata", false)

		assert.OK(t, err)

		foundFile := false
		foundDir := false

		for _, dir := range dirs {
			switch {
			case !dir.IsDir() && dir.Name() == "normal_file":
				foundFile = true
			case dir.IsDir() && dir.Name() == "normal_dir":
				foundDir = true
			}
		}

		assert.Equal(t, false, foundFile)
		assert.Equal(t, true, foundDir)
	})

	t.Run("should not list hidden directories is hidden is false", func(t *testing.T) {
		dirs, err := xioutil.ReadDirsInDir("./testdata", false)

		assert.OK(t, err)

		foundHidden := false

		for _, dir := range dirs {
			// Having to rely on this at the moment...
			if dir.IsDir() && dir.Name() == ".hidden_dir" {
				foundHidden = true
			}
		}

		assert.Equal(t, false, foundHidden)
	})

	t.Run("should list hidden directories is hidden is true", func(t *testing.T) {
		dirs, err := xioutil.ReadDirsInDir("./testdata", true)

		assert.OK(t, err)

		foundHidden := false

		for _, dir := range dirs {
			// Having to rely on this at the moment...
			if dir.IsDir() && dir.Name() == ".hidden_dir" {
				foundHidden = true
			}
		}

		assert.Equal(t, true, foundHidden)
	})
}
