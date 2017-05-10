package gcs

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	xcontext "golang.org/x/net/context"
)

// StorageObject is the interface that lets use mock a *storage.ObjectHandle instance. We can
// construct on Object with a StorageObject.
type StorageObject interface {
	NewWriter(ctx xcontext.Context) *storage.Writer
}

// Object is used by our Bucket interface for interacting with objects in GCS.
type Object interface {
	NewWriteCloser(ctx context.Context) io.WriteCloser
}

// GoogleObject is an implementation of Object that can use the real Google Cloud Storage client
// library (but doesn't have to).
type GoogleObject struct {
	object StorageObject
}

// NewGoogleObject produces a new Object instance, using GoogleObject.
func NewGoogleObject(object StorageObject) Object {
	return &GoogleObject{
		object: object,
	}
}

// NewWriteCloser wraps a call to the underlying StorageObject, creating an io.WriteCloser, which is
// like a *storage.Writer. This should be idempotent (but the returned writer may write to GCS).
func (o *GoogleObject) NewWriteCloser(ctx context.Context) io.WriteCloser {
	return o.object.NewWriter(ctx)
}
