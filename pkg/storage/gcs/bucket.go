package gcs

import "cloud.google.com/go/storage"

type StorageBucket interface {
	Object(name string) *storage.ObjectHandle
}

type Bucket interface {
	Object(name string) Object
}

type GoogleBucket struct {
	bucket StorageBucket
}

func NewGoogleBucket(bucket StorageBucket) Bucket {
	return &GoogleBucket{
		bucket: bucket,
	}
}

func (b *GoogleBucket) Object(name string) Object {
	obj := b.bucket.Object(name)

	return NewGoogleObject(obj)
}
