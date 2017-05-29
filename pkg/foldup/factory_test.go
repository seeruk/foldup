package foldup

import (
	"errors"
	"testing"

	gstorage "cloud.google.com/go/storage"
	"github.com/SeerUK/assert"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

func TestNewCLIFactory(t *testing.T) {
	t.Run("should not return nil", func(t *testing.T) {
		factory := NewCLIFactory()

		assert.NotEqual(t, nil, factory)
	})
}

func TestCliFactory_CreateGCSGateway(t *testing.T) {
	t.Run("should not error under normal circumstances", func(t *testing.T) {
		defer revertStubs()

		newGCSClient = func(ctx context.Context, opts ...option.ClientOption) (*gstorage.Client, error) {
			return &gstorage.Client{}, nil
		}

		factory := NewCLIFactory()

		_, err := factory.CreateGCSGateway("test-bucket")

		assert.OK(t, err)
	})

	t.Run("should propagate errors creating the GCS client", func(t *testing.T) {
		defer revertStubs()

		newGCSClient = func(ctx context.Context, opts ...option.ClientOption) (*gstorage.Client, error) {
			return nil, errors.New("uh oh")
		}

		factory := NewCLIFactory()

		_, err := factory.CreateGCSGateway("test-bucket")

		assert.NotOK(t, err)
	})
}
