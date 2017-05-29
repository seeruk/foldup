package main

import (
	"io"
	"os"

	"github.com/SeerUK/foldup/pkg/foldup/cli"
)

// For testing
var args = os.Args
var exit = os.Exit
var writer io.Writer = os.Stdout

func main() {
	code := cli.CreateApplication(writer).Run(args[1:], os.Environ())

	exit(code)
}
