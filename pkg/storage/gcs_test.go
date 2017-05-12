package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/SeerUK/assert"
	"github.com/SeerUK/foldup/pkg/storage/gcs"
)

type discardWriteCloser struct {
	writer io.Writer
	closed bool
}

func newDiscardWriteCloser(writer io.Writer) io.WriteCloser {
	return &discardWriteCloser{
		writer: writer,
	}
}

func (w *discardWriteCloser) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *discardWriteCloser) Close() error {
	if w.closed {
		return errors.New("storage: already closed")
	}

	w.closed = true

	return nil
}

type testGCSClient struct {
	bucket     gcs.Bucket
	bucketName string
}

func (c *testGCSClient) Bucket(name string) gcs.Bucket {
	c.bucketName = name

	return c.bucket
}

type testGCSBucket struct {
	object     gcs.Object
	objectName string
}

func (b *testGCSBucket) Object(name string) gcs.Object {
	b.objectName = name

	return b.object
}

type testGCSObject struct {
	writeCloser io.WriteCloser
}

func (o *testGCSObject) NewWriteCloser(ctx context.Context) io.WriteCloser {
	return o.writeCloser
}

func newGCSClient(writeCloser io.WriteCloser) gcs.Client {
	object := &testGCSObject{}
	object.writeCloser = writeCloser

	bucket := &testGCSBucket{}
	bucket.object = object

	client := &testGCSClient{}
	client.bucket = bucket

	return client
}

func TestGCSGateway_Store(t *testing.T) {
	t.Run("should not error", func(t *testing.T) {
		bucketName := "test-bucket"
		fileName := "test-file"
		reader := bytes.NewBuffer([]byte("test-data"))

		client := newGCSClient(newDiscardWriteCloser(ioutil.Discard))
		gateway := NewGCSGateway(client, bucketName)

		err := gateway.Store(context.Background(), fileName, reader)

		assert.OK(t, err)
	})

	t.Run("should error if the writer is somehow closed", func(t *testing.T) {
		bucketName := "test-bucket"
		fileName := "test-file"
		reader := bytes.NewBuffer([]byte("test-data"))

		writeCloser := newDiscardWriteCloser(ioutil.Discard)
		err := writeCloser.Close()

		assert.OK(t, err)

		client := newGCSClient(writeCloser)
		gateway := NewGCSGateway(client, bucketName)

		err = gateway.Store(context.Background(), fileName, reader)

		assert.NotOK(t, err)
	})
}
