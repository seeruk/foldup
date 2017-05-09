package gcs

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	xcontext "golang.org/x/net/context"
)

type StorageObject interface {
	NewWriter(ctx xcontext.Context) *storage.Writer
}

type Object interface {
	NewWriteCloser(ctx context.Context) io.WriteCloser
}

type GoogleObject struct {
	object StorageObject
}

func NewGoogleObject(object StorageObject) Object {
	return &GoogleObject{
		object: object,
	}
}

func (o *GoogleObject) NewWriteCloser(ctx context.Context) io.WriteCloser {
	return o.object.NewWriter(ctx)
}
