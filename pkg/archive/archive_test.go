package archive

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/SeerUK/assert"
)

const (
	// Test data directories.
	testDir1 = "testdata/test1"
	testDir2 = "testdata/test2"
	testDir3 = "testdata/test3doesntexist"

	// Test formats.
	testFmtValid   = "test-%s-%d"
	testFmtInvalid = "test-%s-%d-%d"

	testPatternValid = `test\-\w+\-\d+`
)

func TestDirf(t *testing.T) {
	t.Run("should return an archive filename", func(t *testing.T) {
		filename, err := Dirf(testDir1, testFmtValid, TarGz)
		assert.OK(t, err)

		matched, err := regexp.MatchString(testPatternValid, filename)

		assert.OK(t, err)
		assert.True(t, strings.HasPrefix(filename, "testdata/test-test1-"), "Unexpected filename")
		assert.True(t, matched, "Unexpected filename")

		err = os.Remove(filename)
		assert.OK(t, err)
	})

	t.Run("should not error when given an invalid format", func(t *testing.T) {
		// This might seem counter-intuitive, but it's the same behaviour as the fmt package.
		filename, err := Dirf(testDir1, testFmtInvalid, TarGz)
		assert.OK(t, err)

		err = os.Remove(filename)
		assert.OK(t, err)
	})

	t.Run("should create an archive with the returned filename", func(t *testing.T) {
		filename, err := Dirf(testDir1, testFmtValid, TarGz)

		assert.OK(t, err)

		exists := true
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			exists = false
		}

		assert.True(t, exists, "Expected file to exist")

		err = os.Remove(filename)
		assert.OK(t, err)
	})

	t.Run("should error if a non-existent directory is given", func(t *testing.T) {
		filename, err := Dirf(testDir3, testFmtValid, TarGz)
		assert.NotOK(t, err)

		err = os.Remove(filename)
		assert.OK(t, err)
	})

	t.Run("should error if the archive can't be produced", func(t *testing.T) {
		create = func(name string) (*os.File, error) {
			return nil, errors.New("create error")
		}

		filename, err := Dirf(testDir1, testFmtInvalid, TarGz)

		defer revertStubs()
		defer func() {
			if err == nil {
				os.Remove(filename)
			}
		}()

		assert.NotOK(t, err)
	})

	// @todo: This should also test that an actual archive is produced.
}

func TestDirsf(t *testing.T) {
	// @todo: Write these, similar to above, but for multiple directories.
}
