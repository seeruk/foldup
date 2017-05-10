package gcs

import "cloud.google.com/go/storage"

// StorageBucket is the interface that lets us mock a *storage.BucketHandle instance. We can
// construct a Bucket with a StorageBucket.
type StorageBucket interface {
	Object(name string) *storage.ObjectHandle
}

// Bucket is used by our Client interface for interacting with buckets in GCS.
type Bucket interface {
	Object(name string) Object
}

// GoogleBucket is an implementation of Bucket that can use the real Google Cloud Storage client
// library (but doesn't have to).
type GoogleBucket struct {
	bucket StorageBucket
}

// NewGoogleBucket produces a new Bucket instance, using GoogleBucket.
func NewGoogleBucket(bucket StorageBucket) Bucket {
	return &GoogleBucket{
		bucket: bucket,
	}
}

// Object wraps a call to the underlying StorageBucket, creating an Object, which is like a
// *storage.ObjectHandle. This should be idempotent.
func (b *GoogleBucket) Object(name string) Object {
	obj := b.bucket.Object(name)

	return NewGoogleObject(obj)
}
