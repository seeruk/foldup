package cli

import "github.com/eidolon/console"

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

	application.AddCommands([]*console.Command{})

	return application
}
