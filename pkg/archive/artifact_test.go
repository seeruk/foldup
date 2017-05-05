package archive

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/SeerUK/assert"
)

func TestTarGzProducer(t *testing.T) {
	t.Run("should produce an artifact with '.tar.gz' extension", func(t *testing.T) {
		artifact, err := tarGzProducer("testdata", "ext")
		defer os.Remove(artifact.Name())

		assert.OK(t, err)
		assert.True(t, strings.HasSuffix(artifact.Name(), ".tar.gz"), "Expected .tar.gz suffix")
	})

	t.Run("should error the artifact file can't be created", func(t *testing.T) {
		artifact, err := tarGzProducer("idontexist", "andshouldfail")

		defer func() {
			if err == nil {
				os.Remove(artifact.Name())
			}
		}()

		assert.NotOK(t, err)
	})
}

func TestNewTarGzArtifact(t *testing.T) {
	t.Run("should create an artifact using the given namedWriteCloser", func(t *testing.T) {
		nwc := stubArchiveWriter{}
		nwc.name = func() string {
			return "constructor"
		}

		artifact := newTarGzArtifact(&nwc)

		assert.Equal(t, "constructor", artifact.Name())
	})
}

func TestTarGzArtifact_Close(t *testing.T) {
	t.Run("should error if any of the writers are already closed", func(t *testing.T) {
		fw := &stubArchiveWriter{}
		gw := &stubArchiveWriter{}
		tw := &stubArchiveWriter{}

		artifact := tarGzArtifact{
			fw: fw,
			gw: gw,
			tw: tw,
		}

		assert.OK(t, artifact.Close())

		fw.close = func() error {
			return errors.New("fw closed")
		}

		err := artifact.Close()
		assert.NotOK(t, err)
		assert.Equal(t, "fw closed", err.Error())

		gw.close = func() error {
			return errors.New("gw closed")
		}

		err = artifact.Close()
		assert.NotOK(t, err)
		assert.Equal(t, "gw closed", err.Error())

		tw.close = func() error {
			return errors.New("tw closed")
		}

		err = artifact.Close()
		assert.NotOK(t, err)
		assert.Equal(t, "tw closed", err.Error())
	})
}

func TestTarGzArtifact_AddFile(t *testing.T) {
	t.Run("should actually add files to the resulting archive", func(t *testing.T) {
		// Create the archive artifact
		artifact, err := tarGzProducer("testdata", "adds_files")
		defer os.Remove(artifact.Name())

		// Add some files
		filename1 := "testdata/test2/test2_1.txt"
		filename2 := "testdata/test2/test2_2.txt"

		info1, err := stat(filename1)
		assert.OK(t, err)

		info2, err := stat(filename2)
		assert.OK(t, err)

		assert.OK(t, err)
		assert.OK(t, artifact.AddFile("testdata/test1/test.txt", info1))
		assert.OK(t, artifact.AddFile("testdata/test1/test.txt", info2))
		assert.OK(t, artifact.Close())

		// Read the archive
		fr, err := os.Open(artifact.Name())
		assert.OK(t, err)

		gr, err := gzip.NewReader(fr)
		assert.OK(t, err)

		tr := tar.NewReader(gr)

		actual := 0
		expected := 2

		for {
			_, err := tr.Next()
			if err == io.EOF {
				break
			}

			assert.OK(t, err)

			actual++
		}

		// Finally, check that the amount of files we found was equal to the files we put in.
		assert.Equal(t, expected, actual)
	})

	// @todo: Test that it handles symlinks

	t.Run("should error if the file can't be opened", func(t *testing.T) {
		artifact, err := tarGzProducer("testdata", "invalid_path")
		defer os.Remove(artifact.Name())

		assert.OK(t, err)
		assert.NotOK(t, artifact.AddFile("this/path/doesnt/exist", nil))
	})

	t.Run("should error if the info passed is bad", func(t *testing.T) {
		artifact, err := tarGzProducer("testdata", "invalid_info")
		defer os.Remove(artifact.Name())

		assert.OK(t, err)
		assert.NotOK(t, artifact.AddFile("testdata/test1/test.txt", nil))
	})

	t.Run("should error if the header fails to write to the tar", func(t *testing.T) {
		fw := &stubArchiveWriter{}
		gw := &stubArchiveWriter{}
		tw := &stubArchiveWriter{}

		artifact := tarGzArtifact{
			fw: fw,
			gw: gw,
			tw: tw,
		}

		tw.writeHeader = func(*tar.Header) error {
			return errors.New("tw write header error")
		}

		info, err := stat("testdata/test1/test.txt")
		assert.OK(t, err)

		err = artifact.AddFile("testdata/test1/test.txt", info)

		assert.NotOK(t, err)
	})

	t.Run("should error if the file fails to write to the tar", func(t *testing.T) {
		fw := &stubArchiveWriter{}
		gw := &stubArchiveWriter{}
		tw := &stubArchiveWriter{}

		artifact := tarGzArtifact{
			fw: fw,
			gw: gw,
			tw: tw,
		}

		tw.close = func() error {
			return errors.New("tw closed")
		}

		info, err := stat("testdata/test1/test.txt")
		assert.OK(t, err)

		err = artifact.AddFile("testdata/test1/test.txt", info)

		assert.NotOK(t, err)
	})
}

func TestTarGzArtifact_Name(t *testing.T) {
	t.Run("should create an artifact using the given namedWriteCloser", func(t *testing.T) {
		artifact, err := tarGzProducer("testdata", "name")
		defer os.Remove(artifact.Name())

		assert.OK(t, err)
		assert.Equal(t, "testdata/name.tar.gz", artifact.Name())
	})
}
