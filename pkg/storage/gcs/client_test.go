package gcs

import (
	"testing"

	"cloud.google.com/go/storage"
	"github.com/SeerUK/assert"
)

type TestStorageClient struct {
	bucket string
}

func (c *TestStorageClient) Bucket(name string) *storage.BucketHandle {
	c.bucket = name

	return &storage.BucketHandle{}
}

func TestGoogleClient_Bucket(t *testing.T) {
	t.Run("should create a bucket handle", func(t *testing.T) {
		name := "test-bucket"

		sc := &TestStorageClient{}
		gc := NewGoogleClient(sc)

		_ = gc.Bucket(name)

		assert.Equal(t, sc.bucket, name)
	})
}
