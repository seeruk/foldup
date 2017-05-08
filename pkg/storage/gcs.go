package storage

import (
	"context"
	"io"

	gstorage "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var gcsNewClient = gstorage.NewClient // For testing

type GCSGateway struct {
	bucketName string
	bucketRef  *gstorage.BucketHandle
	clientRef  *gstorage.Client
	options    []option.ClientOption
}

func NewGCSGateway(ctx context.Context, bktName string, opts []option.ClientOption) (Gateway, error) {
	client, err := gcsNewClient(ctx, opts...)
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

func (g *GCSGateway) Store(ctx context.Context, filename string, in io.Reader) error {
	object := g.bucketRef.Object(filename)

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

func (g *GCSGateway) Close() error {
	return g.clientRef.Close()
}
