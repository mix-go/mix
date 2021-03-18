package main

import (
	"github.com/mix-go/xcli"
	"github.com/mix-go/dotenv"
	"github.com/mix-go/xstart/commands"
)

func main() {
	xcli.SetName("xstart").
		SetVersion("1.1.1").
		SetDebug(dotenv.Getenv("APP_DEBUG").Bool(false))
	xcli.AddCommand(commands.Cmds...).Run()
}
