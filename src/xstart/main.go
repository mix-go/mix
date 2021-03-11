package main

import (
	"github.com/mix-go/cli"
	"github.com/mix-go/dotenv"
	"github.com/mix-go/xstart/commands"
)

func main() {
	cli.SetName("xstart").
		SetVersion("1.1.0").
		SetDebug(dotenv.Getenv("APP_DEBUG").Bool(false))
	cli.AddCommand(commands.Cmds...).Run()
}
