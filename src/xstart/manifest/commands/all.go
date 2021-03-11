package commands

import (
	"github.com/mix-go/console"
	"github.com/mix-go/mix/devtool/commands"
)

var (
	Commands []console.CommandDefinition
)

func init() {
	Commands = append(Commands,
		console.CommandDefinition{
			Name:    "version",
			Usage:   "Prints the current Mix Go version",
			Command: &commands.VersionCommand{},
		},
		console.CommandDefinition{
			Name:  "new",
			Usage: "\tCreate a console application",
			Options: []console.OptionDefinition{
				{
					Names: []string{"n", "name"},
					Usage: "Project name",
				},
			},
			Command: &commands.NewCommand{},
		},
		console.CommandDefinition{
			Name:  "api",
			Usage: "\tCreate a api application",
			Options: []console.OptionDefinition{
				{
					Names: []string{"n", "name"},
					Usage: "Project name",
				},
			},
			Command: &commands.APICommand{},
		},
		console.CommandDefinition{
			Name:  "web",
			Usage: "\tCreate a web application",
			Options: []console.OptionDefinition{
				{
					Names: []string{"n", "name"},
					Usage: "Project name",
				},
			},
			Command: &commands.WebCommand{},
		},
		console.CommandDefinition{
			Name:  "grpc",
			Usage: "\tCreate a gRPC application",
			Options: []console.OptionDefinition{
				{
					Names: []string{"n", "name"},
					Usage: "Project name",
				},
			},
			Command: &commands.GrpcCommand{},
		},
	)
}
