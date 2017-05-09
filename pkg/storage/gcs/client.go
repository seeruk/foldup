package gcs

import (
	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// @todo: How do we create StorageClient in a way we can test? Bearing in mind storage.NewClient
// actually makes a request to GCS, we need to avoid using that in tests.

var newStorageClientFn = storage.NewClient

// @todo: Move to export_test.go / an _test.go file.
func newFakeStorageClient(ctx context.Context, opts ...option.ClientOption) (*storage.Client, error) {
	return &storage.Client{}, nil
}

// newStorageClient uses the (maybe patched) newStorageClientFn to produce a new StorageClient
// instance. This is used to wrap the *storage.Client struct in an interface we can re-use.
func newStorageClient() (StorageClient, error) {
	return newStorageClientFn(context.Background())
}

// StorageClient is the interface that lets us mock a *storage.Client instance, we can construct a
// Client with a StorageClient.
type StorageClient interface {
	Bucket(name string) *storage.BucketHandle
}

// Client is the client interface we'll be using in our code that intends to use GCS.
type Client interface {
	Bucket(name string) Bucket
}

// GoogleClient is an implementation of Client that can use the real Google Cloud Storage client
// library (but doesn't have to).
type GoogleClient struct {
	storage StorageClient
}

// NewGoogleClient produces a new Client instance, using GoogleClient.
func NewGoogleClient() (Client, error) {
	storageClient, err := newStorageClient()
	if err != nil {
		return nil, err
	}

	client := &GoogleClient{
		storage: storageClient,
	}

	return client, nil
}

// NewGoogleClientWithStorageClient produces a new Client instance, using GoogleClient, given a
// specific StorageClient implementation.
func NewGoogleClientWithStorageClient(client StorageClient) Client {
	return &GoogleClient{
		storage: client,
	}
}

// Bucket wraps a call to the underlying StorageClient, fetching a Bucket, which is like a
// *storage.BucketHandle.
func (c *GoogleClient) Bucket(name string) Bucket {
	gb := c.storage.Bucket(name)

	return NewGoogleBucket(gb)
}
