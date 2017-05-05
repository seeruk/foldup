package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
)

func init() {
	// Register the built-in TarGz format.
	RegisterFormat(TarGz, tarGzProducer)
}

// TarGz is a format for creating gzipped tarballs.
const TarGz FormatName = "TarGz"

// tarWriteCloser is an interface that provides functionality for writing data, writing tar headers,
// and closing a tar.
//
// WriteHeader would ideally accept an interface, as there are similar implementations for other
// archive types, but unfortunately that's not how the stdlib was implemented.
type tarWriteCloser interface {
	io.WriteCloser

	WriteHeader(hdr *tar.Header) error
}

// tarGzProducer creates a tarGzArtifact, creating the archive file in the process.
func tarGzProducer(pathname, filename string) (Artifact, error) {
	filename = fmt.Sprintf("%s.tar.gz", filename)

	// To create the file, we need to create the path to the file.
	file, err := create(path.Join(pathname, filename))
	if err != nil {
		return nil, err
	}

	return newTarGzArtifact(file), nil
}

type tarGzArtifact struct {
	fw namedWriteCloser
	gw io.WriteCloser
	tw tarWriteCloser
}

func newTarGzArtifact(fw namedWriteCloser) Artifact {
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
