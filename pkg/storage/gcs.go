package storage

import (
	"context"
	"io"

	"github.com/SeerUK/foldup/pkg/storage/gcs"
)

type GCSGateway struct {
	bucket string
	client gcs.Client
}

func NewGCSGateway(client gcs.Client, bucket string) Gateway {
	return &GCSGateway{
		bucket: bucket,
		client: client,
	}
}

func (g *GCSGateway) Store(ctx context.Context, filename string, reader io.Reader) error {
	writer := g.client.Bucket(g.bucket).Object(filename).NewWriteCloser(ctx)
	_, err := io.Copy(writer, reader)

	writer.Close()

	return err
}
