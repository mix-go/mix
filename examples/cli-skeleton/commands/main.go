package commands

import (
	"github.com/mix-go/xcli"
)

var Commands = []*xcli.Command{
	{
		Name:   "hello",
		Short:  "\tEcho demo",
		Options: []*xcli.Option{
			{
				Names: []string{"n", "name"},
				Usage: "Your name",
			},
			{
				Names: []string{"say"},
				Usage: "\tSay ...",
			},
		},
		RunI: &HelloCommand{},
	},
}
