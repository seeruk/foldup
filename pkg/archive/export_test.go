package archive

import (
	"archive/tar"
	"io/ioutil"
	"os"
)

func revertStubs() {
	create = os.Create
	open = os.Open
	readDir = ioutil.ReadDir
	stat = os.Stat
}

type stubArchiveWriter struct {
	name        func() string
	write       func([]byte) (n int, err error)
	writeHeader func(*tar.Header) error
	close       func() error
}

func (b *stubArchiveWriter) Name() string {
	if b.name != nil {
		return b.name()
	}

	return ""
}

func (b *stubArchiveWriter) WriteHeader(hdr *tar.Header) error {
	if b.writeHeader != nil {
		return b.writeHeader(hdr)
	}

	return nil
}

func (b *stubArchiveWriter) Write(bs []byte) (n int, err error) {
	if b.write != nil {
		return b.write(bs)
	}

	return 0, nil
}

func (b *stubArchiveWriter) Close() error {
	if b.close != nil {
		return b.close()
	}

	return nil
}
