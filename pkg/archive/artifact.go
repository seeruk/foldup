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
)

// producers is an array of producer functions that are ordered to match up with the Format
// constants.
var producers = [...]producer{
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

	return newTarGzArtifact(file), nil
}

// namedWriteCloser is an interface that provides functionality for writing data, closing, and
// fetching a name to identify what is being written.
//
// This interface is used primarily for being like *os.File for testing.
type namedWriteCloser interface {
	io.WriteCloser

	Name() string
}

// tarWriteCloser is an interface that provides functionality for writing data, writing tar headers,
// and closing a tar.
//
// WriteHeader would ideally accept an interface, as there are similar implementations for other
// archive types, but unfortunately that's not how the stdlib was implemented.
type tarWriteCloser interface {
	io.WriteCloser

	WriteHeader(hdr *tar.Header) error
}

// artifact represents an archive to be interacted with.
type artifact interface {
	io.Closer

	// AddFile should take the file at the given path with the given os.FileInfo to the artifact.
	AddFile(path string, info os.FileInfo) error
	// Name returns the artifact's name. In most cases it will be the file name.
	Name() string
}

type tarGzArtifact struct {
	fw namedWriteCloser
	gw io.WriteCloser
	tw tarWriteCloser
}

func newTarGzArtifact(fw namedWriteCloser) artifact {
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

	header, err := tar.FileInfoHeader(info, path)
	if err != nil {
		return err
	}

	// @todo: We need to handle these still... this must be things like symlinks?
	if !info.Mode().IsRegular() {
		return nil
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

func (a *tarGzArtifact) Name() string {
	return a.fw.Name()
}
