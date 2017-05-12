package gcs

import (
	"testing"

	"cloud.google.com/go/storage"
	"github.com/SeerUK/assert"
)

type TestStorageBucket struct {
	object string
}

func (b *TestStorageBucket) Object(name string) *storage.ObjectHandle {
	b.object = name

	return &storage.ObjectHandle{}
}

func TestGoogleBucket_Object(t *testing.T) {
	t.Run("should create an object handle", func(t *testing.T) {
		name := "test-object"

		sb := &TestStorageBucket{}
		gb := NewGoogleBucket(sb)

		_ = gb.Object(name)

		assert.Equal(t, sb.object, name)
	})
}
