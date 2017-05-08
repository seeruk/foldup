package storage

import (
	"context"
	"io"
)

// Gateway provides an interface for interacting with some kind of storage system. This could be
// filesystem-based, in-memory, in some remote storage bucket, etc.
type Gateway interface {
	io.Closer

	Store(ctx context.Context, filename string, in io.Reader) error
}
