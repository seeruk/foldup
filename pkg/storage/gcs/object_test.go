package gcs

import (
	"context"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/SeerUK/assert"
	xcontext "golang.org/x/net/context"
)

type TestStorageObject struct {
	newWriter bool
}

func (o *TestStorageObject) NewWriter(ctx xcontext.Context) *storage.Writer {
	o.newWriter = true

	return &storage.Writer{}
}

func TestGoogleObject_NewWriteCloser(t *testing.T) {
	t.Run("should create an io.WriteCloser", func(t *testing.T) {
		sob := &TestStorageObject{}
		gob := NewGoogleObject(sob)

		_ = gob.NewWriteCloser(context.Background())

		assert.True(t, sob.newWriter, "Expected newWriter to have been called")
	})
}
