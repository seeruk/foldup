package command

import "github.com/eidolon/console"

// StartCommand creates a command to trigger periodic backups.
func StartCommand() *console.Command {
	configure := func(def *console.Definition) {

	}

	execute := func(input *console.Input, output *console.Output) error {
		output.Println("Hello, World!")

		return nil
	}

	return &console.Command{
		Name:        "start",
		Description: "Begin periodically backing up.",
		Configure:   configure,
		Execute:     execute,
	}
}
