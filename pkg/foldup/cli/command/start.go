package command

import (
	"fmt"

	"github.com/SeerUK/foldup/pkg/xioutil"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
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
			panic(err)
		}

		for _, d := range dirs {
			fmt.Println(d.Name())
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
