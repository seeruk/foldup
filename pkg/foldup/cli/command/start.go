package command

import (
	"os"
	"path"

	"github.com/SeerUK/foldup/pkg/archive"
	"github.com/SeerUK/foldup/pkg/storage"
	"github.com/SeerUK/foldup/pkg/xioutil"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
	gstorage "google.golang.org/api/storage/v1"
)

// StartCommand creates a command to trigger periodic backups.
func StartCommand() *console.Command {
	var dirname string

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&dirname),
			Spec:  "DIRNAME",
			Desc:  "The directory to archive folders from.",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		dirs, err := xioutil.ReadDirsInDir(dirname, false)
		if err != nil {
			return err
		}

		relativePaths := []string{}

		for _, d := range dirs {
			relativePaths = append(relativePaths, path.Join(dirname, d.Name()))
		}

		archives, err := archive.Dirsf(relativePaths, "backup-%s-%d", archive.TarGz)
		if err != nil {
			return err
		}

		//gateway, err := storage.NewGCSGateway(context.Background(), "backups-sierra", nil)
		service, err := storage.NewGCSClient()
		if err != nil {
			return err
		}

		for _, a := range archives {
			in, err := os.Open(a)
			if err != nil {
				return err
			}

			_, err = service.Objects.Insert("backups-sierra", &gstorage.Object{Name: a}).Media(in).Do()
			if err != nil {
				return err
			}

			err = os.Remove(a)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return &console.Command{
		Name:        "start",
		Description: "Begin periodically backing up.",
		Configure:   configure,
		Execute:     execute,
	}
}
