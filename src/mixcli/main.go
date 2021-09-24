package main

import (
	"github.com/mix-go/xcli"
	"github.com/mix-go/dotenv"
	"github.com/mix-go/mixcli/commands"
)

func main() {
	xcli.SetName("mixcli").
		SetVersion(commands.CLIVersion).
		SetDebug(dotenv.Getenv("APP_DEBUG").Bool(false))
	xcli.AddCommand(commands.Cmds...).Run()
}
