package main

import (
	"github.com/mix-go/xcli"
	"github.com/mix-go/dotenv"
	"github.com/mix-go/xstart/commands"
)

func main() {
	xcli.SetName("mixcli").
		SetVersion(commands.FrameworkVersion).
		SetDebug(dotenv.Getenv("APP_DEBUG").Bool(false))
	xcli.AddCommand(commands.Cmds...).Run()
}
