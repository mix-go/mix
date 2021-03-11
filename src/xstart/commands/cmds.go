package commands

import "github.com/mix-go/cli"

var Cmds = []*cli.Command{
	{
		Name:  "new",
		Usage: "\tCreate a project",
		RunI:  &NewCommand{},
	},
}
