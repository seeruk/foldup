package storage

import (
	"context"
	"io"

	gstorage "cloud.google.com/go/storage"
	xcontext "golang.org/x/net/context"
	"google.golang.org/api/option"
)

// gcsBucketRef is an interface that allows us to mock a Google Cloud Storage bucket handle.
type gcsBucketRef interface {
	Attrs(ctx xcontext.Context) (*gstorage.BucketAttrs, error)
	Object(name string) *gstorage.ObjectHandle
}

// gcsClientRef is an interface that allows us to mock a Google Cloud Storage client.
type gcsClientRef interface {
	io.Closer

	Bucket(name string) *gstorage.BucketHandle
}

var newGcsClientRefFn = gstorage.NewClient // For testing

// newGcsClientRef is used to force the use of the new context package, and our interface upon the
// Google Cloud Storage client library. It'd be nice if they provided some interfaces of their own!
func newGcsClientRef(ctx context.Context, opts ...option.ClientOption) (gcsClientRef, error) {
	return newGcsClientRefFn(ctx, opts...)
}

// GCSGateway is a gateway for interacting with Google Cloud Storage.
type GCSGateway struct {
	bucketName string
	bucketRef  gcsBucketRef
	clientRef  gcsClientRef
	options    []option.ClientOption
}

// NewGCSGateway sets up a new Google Cloud Storage gateway, including acquiring a client reference
// and a bucket reference. If the client fails to set up, or the bucket can't be found then this
// function will return an error.
func NewGCSGateway(ctx context.Context, bktName string, opts []option.ClientOption) (Gateway, error) {
	client, err := newGcsClientRef(ctx, opts...)
	if err != nil {
		return nil, err
	}

	bkt := client.Bucket(bktName)

	// Probe for the bucket's existence. We won't be creating it if it doesn't exist.
	_, err = bkt.Attrs(ctx)
	if err != nil {
		return nil, err
	}

	gateway := &GCSGateway{
		bucketName: bktName,
		bucketRef:  bkt,
		clientRef:  client,
	}

	return gateway, nil
}

// Store takes a name (used as the object name) and a reader (to pull data from) and writes the data
// into the storage. We use `io.Copy` to safely handle large files.
func (g *GCSGateway) Store(ctx context.Context, name string, in io.Reader) error {
	object := g.bucketRef.Object(name)

	out := object.NewWriter(ctx)

	// We use io.Copy here to avoid reading giant files into memory all at once easily. Passing in
	// something like a byte array wouldn't work well.
	_, err := io.Copy(out, in)
	if err != nil {
		return err
	}

	if err := out.Close(); err != nil {
		return err
	}

	return nil
}

// Close closes the underlying Google Cloud Storage client.
func (g *GCSGateway) Close() error {
	return g.clientRef.Close()
}
