package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
)

// A Format represents an archive artifact format that can be produced.
type Format int

const (
	// TarGz is a format that will produce a gzipped tarball with the extension `.tar.gz`.
	TarGz Format = iota
	// Stub is a format used for testing that can be forced to error if necessary.
	Stub
)

var producers = [...]producer{
	tarGzProducer,
	tarGzProducer,
}

// producer is a function that will create a new artifact, taking a file name nad path to produce an
// artifact with an appropriate file name in the request location.
type producer func(pathname, filename string) (artifact, error)

// tarGzProducer creates a tarGzArtifact, creating the archive file in the process.
func tarGzProducer(pathname, filename string) (artifact, error) {
	filename = fmt.Sprintf("%s.tar.gz", filename)

	// To create the file, we need to create the path to the file.
	file, err := create(path.Join(pathname, filename))
	if err != nil {
		return nil, err
	}

	return newTarGzArchive(file), nil
}

// artifact represents an archive to be interacted with.
type artifact interface {
	io.Closer

	// AddFile should take the file at the given path with the given os.FileInfo to the artifact.
	AddFile(path string, info os.FileInfo) error
	// Filename returns the artifact's filename.
	Filename() string
}

type tarGzArtifact struct {
	// @todo: Could we make an interface for this that will let us get the name, and close?
	fw *os.File
	gw *gzip.Writer
	tw *tar.Writer
}

func newTarGzArchive(fw *os.File) artifact {
	gw := gzip.NewWriter(fw)
	tw := tar.NewWriter(gw)

	return &tarGzArtifact{
		fw: fw,
		gw: gw,
		tw: tw,
	}
}

func (a *tarGzArtifact) Close() error {
	if err := a.tw.Close(); err != nil {
		return err
	}

	if err := a.gw.Close(); err != nil {
		return err
	}

	if err := a.fw.Close(); err != nil {
		return err
	}

	return nil
}

func (a *tarGzArtifact) AddFile(path string, info os.FileInfo) error {
	source, err := open(path)
	if err != nil {
		return err
	}

	defer source.Close()

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// @todo: remove leading ./ and ../
	header.Name = path

	err = a.tw.WriteHeader(header)
	if err != nil {
		return err
	}

	if _, err := io.Copy(a.tw, source); err != nil {
		return err
	}

	return nil
}

func (a *tarGzArtifact) Filename() string {
	return a.fw.Name()
}
