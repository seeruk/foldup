package foldup

import (
	"context"

	gstorage "cloud.google.com/go/storage"
	"github.com/SeerUK/foldup/pkg/storage"
	"github.com/SeerUK/foldup/pkg/storage/gcs"
)

// Factory is used to create dependencies used elsewhere in the application. It's primary reason for
// existence is to facilitate testing, by allowing us to pass a factory into a command (or a subset
// of a factory) to create it's dependencies "dynamically" in a test. The factory should produce
// interfaces, meaning the actual implementations of anything it creates could be fake.
type Factory interface {
	// CreateGCSGateway is used to create a storage gateway. In normal use it should be a GCS
	// gateway, that uses the given bucket to store files.
	CreateGCSGateway(bucket string) (storage.Gateway, error)
}

// For testing
var newGCSClient = gstorage.NewClient

// cliFactory is the default factory for CLI use, creating real implementations of dependencies.
type cliFactory struct{}

// NewCLIFactory produces a new instance of cliFactory.
func NewCLIFactory() Factory {
	return &cliFactory{}
}

func (f *cliFactory) CreateGCSGateway(bucket string) (storage.Gateway, error) {
	storageClient, err := newGCSClient(context.Background())
	if err != nil {
		return nil, err
	}

	client := gcs.NewGoogleClient(storageClient)
	gateway := storage.NewGCSGateway(client, bucket)

	return gateway, nil
}
