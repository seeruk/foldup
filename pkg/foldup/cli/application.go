package cli

import (
	"github.com/SeerUK/foldup/pkg/foldup/cli/command"
	"github.com/eidolon/console"
)

// CreateApplication builds the console application instance. Providing it with some basic
// information like the name and version.
func CreateApplication() *console.Application {
	application := console.NewApplication("foldup", "0.1.0")
	application.Logo = `
███████╗ ██████╗ ██╗     ██████╗ ██╗   ██╗██████╗
██╔════╝██╔═══██╗██║     ██╔══██╗██║   ██║██╔══██╗
█████╗  ██║   ██║██║     ██║  ██║██║   ██║██████╔╝
██╔══╝  ██║   ██║██║     ██║  ██║██║   ██║██╔═══╝
██║     ╚██████╔╝███████╗██████╔╝╚██████╔╝██║
╚═╝      ╚═════╝ ╚══════╝╚═════╝  ╚═════╝ ╚═╝
`

	application.AddCommands(buildCommands())

	return application
}

// buildCommands instantiates all of the commands registered in the application.
func buildCommands() []*console.Command {
	return []*console.Command{
		command.StartCommand(),
	}
}
