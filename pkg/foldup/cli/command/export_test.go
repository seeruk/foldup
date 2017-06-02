package command

import (
	"context"
	"io"
	"os"

	"github.com/SeerUK/foldup/pkg/archive"
	"github.com/SeerUK/foldup/pkg/scheduling"
	"github.com/SeerUK/foldup/pkg/storage"
)

// backupTestStorageGateway is used as a no-op storage gateway for the backup command during
// testing.
type backupTestStorageGateway struct {
	storeError error
}

func (f *backupTestStorageGateway) Store(ctx context.Context, filename string, in io.Reader) error {
	return f.storeError
}

// backupTestFactory is used to create dependencies for the backup command during testing.
type backupTestFactory struct {
	createGCSGatewayGateway storage.Gateway
	createGCSGatewayError   error
}

func (f *backupTestFactory) CreateGCSGateway(bucket string) (storage.Gateway, error) {
	if f.createGCSGatewayGateway == nil {
		f.createGCSGatewayGateway = &backupTestStorageGateway{}
	}

	return f.createGCSGatewayGateway, f.createGCSGatewayError
}

func revertStubs() {
	archiveDirsf = archive.Dirsf
	osOpen = os.Open
	osRemove = os.Remove
	scheduleFunc = scheduling.ScheduleFunc
}
