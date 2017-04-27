package main

import (
	"os"

	"github.com/SeerUK/foldup/pkg/foldup/cli"
)

func main() {
	os.Exit(cli.CreateApplication().Run(os.Args[1:], os.Environ()))
}
