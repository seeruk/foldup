package storage

import (
	"context"
	"io"

	"github.com/SeerUK/foldup/pkg/storage/gcs"
)

// GCSGateway implements the Gateway interface for interacting with Google Cloud Storage.
type GCSGateway struct {
	bucket string
	client gcs.Client
}

// NewGCSGateway creates a new Gateway instance, using GCSGateway.
func NewGCSGateway(client gcs.Client, bucket string) Gateway {
	return &GCSGateway{
		bucket: bucket,
		client: client,
	}
}

// Store attempts to write a file via the Gateway.
func (g *GCSGateway) Store(ctx context.Context, filename string, reader io.Reader) error {
	writer := g.client.Bucket(g.bucket).Object(filename).NewWriteCloser(ctx)
	_, err := io.Copy(writer, reader)

	cerr := writer.Close()
	if cerr != nil {
		return cerr
	}

	return err
}
