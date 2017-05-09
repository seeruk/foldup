package storage

import (
	"context"
	"io"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	gstorage "google.golang.org/api/storage/v1"
)

const gcsScopeReadWrite = "https://www.googleapis.com/auth/devstorage.read_write"

type GCSGateway struct {
	bucket  string
	service *gstorage.Service
}

var gcsDefaultClientFn = google.DefaultClient

func NewGCSClient(bucket string) (Gateway, error) {
	httpClient, err := gcsDefaultClientFn(context.Background(), gcsScopeReadWrite)
	if err != nil {
		return nil, err
	}

	service, err := gstorage.New(httpClient)
	if err != nil {
		return nil, err
	}

	client := &GCSGateway{
		bucket:  bucket,
		service: service,
	}

	return client, nil
}

type objectInsertCallDoer interface {
	Do(opts ...googleapi.CallOption) (*gstorage.Object, error)
}

func (c *GCSGateway) Store(ctx context.Context, filename string, reader io.Reader) error {
	call := c.service.Objects.
		Insert(c.bucket, &gstorage.Object{Name: filename}).
		Media(reader)

	_ = prepObjectInsertCall()

	_, err := doObjectInsertCall(call)

	return err
}

func prepObjectInsertCall() objectInsertCallDoer {
	return nil
}

func doObjectInsertCall(call objectInsertCallDoer) (*gstorage.Object, error) {
	return call.Do()
}
