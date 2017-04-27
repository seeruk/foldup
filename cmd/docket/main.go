package main

import (
	"os"

	"github.com/SeerUK/docket/pkg/docket/cli"
)

func main() {
	os.Exit(cli.CreateApplication().Run(os.Args[1:], os.Environ()))
}
