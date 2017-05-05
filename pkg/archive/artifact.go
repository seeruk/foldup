package archive

import (
	"fmt"
	"io"
	"os"
)

// Artifact represents an archive to be interacted with.
type Artifact interface {
	io.Closer

	// AddFile should take the file at the given path with the given os.FileInfo to the artifact.
	AddFile(path string, info os.FileInfo) error
	// Name returns the artifact's name. In most cases it will be the file name.
	Name() string
}

// namedWriteCloser is an interface that provides functionality for writing data, closing, and
// fetching a name to identify what is being written.
//
// This interface is used primarily for being like *os.File for testing.
type namedWriteCloser interface {
	io.WriteCloser

	Name() string
}

// FormatName is a slightly less than magical string that is used to identify an archive artifact
// format. If you're creating your own format, you'll also need to declare a FormatName.
type FormatName string

// A producerFunc is a function that produces an archive artifact in a specific format.
type producerFunc func(pathname, filename string) (Artifact, error)

// A format represents an archive artifact format that can be produced.
type format struct {
	name     FormatName
	producer producerFunc
}

// The formats slice contains all registered archive artifact formats.
var formats []format

// RegisterFormat registers an archive artifact format for use by functions that accept a
// FormatName.
//
// Name is a FormatName which is a string. Any format added should have a corresponding FormatName
// constant available.
// Producer is a producerFunc that will create an archive artifact of the appropriate type.
func RegisterFormat(name FormatName, producer producerFunc) {
	formats = append(formats, format{name, producer})
}

// The findFormatByName function attempts to find a archive artifact format that has been registered
// with the given FormatName. If one cannot be found, an error will be returned.
func findFormatByName(name FormatName) (format, error) {
	for _, format := range formats {
		if format.name == name {
			return format, nil
		}
	}

	return format{}, fmt.Errorf("archive: unable to find format '%v'", name)
}
