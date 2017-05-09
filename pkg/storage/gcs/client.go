package gcs

import (
	"cloud.google.com/go/storage"
)

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
func NewGoogleClient(storageClient StorageClient) Client {
	return &GoogleClient{
		storage: storageClient,
	}
}

// Bucket wraps a call to the underlying StorageClient, fetching a Bucket, which is like a
// *storage.BucketHandle.
func (c *GoogleClient) Bucket(name string) Bucket {
	gb := c.storage.Bucket(name)

	return NewGoogleBucket(gb)
}
