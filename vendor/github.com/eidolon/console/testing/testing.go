package testing

import (
	. "github.com/eidolon/console"
)

// RunCommand makes it easier to run a command in a test, by providing all inputs and output, and
// preparing a command similarly to how it is prepared when run in an application.
func RunCommand(cmd *Command, def *Definition, in *Input, env []string, out *Output) error {
	if cmd.Configure != nil {
		cmd.Configure(def)
	}

	err := MapInput(def, in, env)
	if err != nil {
		return err
	}

	return cmd.Execute(in, out)
}
