package archive

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/SeerUK/assert"
)

const (
	// Test data directories.
	testData = "testdata"
	testDir1 = "testdata/test1"
	testDir2 = "testdata/test2"
	testDir3 = "testdata/test3doesntexist"

	// Test formats.
	testFmtValid     = "test-%s-%d"
	testFmtInvalid   = "test-%s-%d-%d"
	testPatternValid = `test\-\w+\-\d+`
)

func TestDirf(t *testing.T) {
	t.Run("should return an archive filename", func(t *testing.T) {
		filename, err := Dirf(testDir1, testFmtValid, TarGz)
		assert.OK(t, err)

		defer os.Remove(filename)

		matched, err := regexp.MatchString(testPatternValid, filename)

		assert.OK(t, err)
		assert.True(t, strings.HasPrefix(filename, "testdata/test-test1-"), "Unexpected filename")
		assert.True(t, matched, "Unexpected filename")
	})

	t.Run("should not error when given an invalid name format", func(t *testing.T) {
		// This might seem counter-intuitive, but it's the same behaviour as the fmt package.
		filename, err := Dirf(testDir2, testFmtInvalid, TarGz)
		assert.OK(t, err)

		err = os.Remove(filename)
		assert.OK(t, err)
	})

	t.Run("should create an archive file with the returned filename", func(t *testing.T) {
		filename, err := Dirf(testDir1, testFmtValid, TarGz)
		assert.OK(t, err)

		defer os.Remove(filename)

		exists := true
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			exists = false
		}

		assert.True(t, exists, "Expected file to exist")
	})

	t.Run("should error if a non-existent directory is given", func(t *testing.T) {
		filename, err := Dirf(testDir3, testFmtValid, TarGz)
		assert.NotOK(t, err)

		err = os.Remove(filename)
		assert.OK(t, err)
	})

	t.Run("should error if the archive artifact can't be produced", func(t *testing.T) {
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

	t.Run("should error if files can't be added to the archive artifact", func(t *testing.T) {
		stubArtifactRef.addFile = func(path string, info os.FileInfo) error {
			return errors.New("addFile error")
		}

		defer revertStubs()

		_, err := Dirf(testDir2, testFmtValid, "stub")
		assert.NotOK(t, err)
	})

	t.Run("should error if a directory cannot be read", func(t *testing.T) {
		readDir = func(dirname string) ([]os.FileInfo, error) {
			return []os.FileInfo{}, errors.New("readDir error")
		}

		defer revertStubs()

		_, err := Dirf(testDir2, testFmtValid, "stub")
		assert.NotOK(t, err)
	})

	t.Run("should error if a subdirectory cannot be walked", func(t *testing.T) {
		calls := 0

		readDir = func(dirname string) ([]os.FileInfo, error) {
			if calls > 0 {
				return []os.FileInfo{}, errors.New("readDir error")
			}

			calls++

			return ioutil.ReadDir(dirname)
		}

		defer revertStubs()

		_, err := Dirf(testData, testFmtValid, "stub")
		assert.NotOK(t, err)
	})

	t.Run("should error if an invalid archive format is given", func(t *testing.T) {
		_, err := Dirf(testData, testFmtValid, "star-wars_the-force-awakens")
		assert.NotOK(t, err)
	})
}

func TestDirsf(t *testing.T) {
	t.Run("should return a sorted list of archive names", func(t *testing.T) {
		filenames, err := Dirsf([]string{testDir1, testDir2}, testFmtValid, TarGz)

		defer func() {
			for _, filename := range filenames {
				os.Remove(filename)
			}
		}()

		assert.OK(t, err)

		// To preserve the original result, we make a new slice...
		sorted := make([]string, len(filenames))

		// ... copy the data into it ...
		copy(sorted, filenames)

		// ... and then sort it.
		sort.Strings(sorted)

		assert.Equal(t, sorted, filenames)
	})

	t.Run("should create archive files with the returned filenames", func(t *testing.T) {
		filenames, err := Dirsf([]string{testDir1, testDir2}, testFmtValid, TarGz)

		defer func() {
			for _, filename := range filenames {
				os.Remove(filename)
			}
		}()

		assert.OK(t, err)

		exists := true

		for _, filename := range filenames {
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				exists = false
			}
		}

		assert.True(t, exists, "Expected all files to exist")
	})

	t.Run("should error if there is an error archiving a directory", func(t *testing.T) {
		filenames, err := Dirsf([]string{testDir1, testDir2}, testFmtValid, "memento")

		defer func() {
			for _, filename := range filenames {
				os.Remove(filename)
			}
		}()

		assert.NotOK(t, err)
	})
}
