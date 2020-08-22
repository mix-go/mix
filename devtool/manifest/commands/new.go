package commands

import (
    "github.com/mix-go/bean"
    "github.com/mix-go/console"
    "github.com/mix-go/mix/devtool/commands"
)

func init() {
    Commands = append(Commands,
        console.CommandDefinition{
            Name:  "new",
            Usage: "\tNew project, console,api,web",
            Options: []console.OptionDefinition{
                {
                    Names: []string{"n", "name"},
                    Usage: "Project name",
                },
                {
                    Names: []string{"t", "type"},
                    Usage: "Project type, console,api,web",
                },
            },
            Reflect: bean.NewReflect(commands.NewCommand{}),
        },
    )
}
