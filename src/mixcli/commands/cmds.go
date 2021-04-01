package commands

import "github.com/mix-go/xcli"

var Cmds = []*xcli.Command{
	{
		Name:  "new",
		Usage: "\tCreate a project",
		RunI:  &NewCommand{},
	},
}
