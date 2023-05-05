package main

import (
	"github.com/mix-go/cli-skeleton/commands"
	_ "github.com/mix-go/cli-skeleton/config/configor"
	_ "github.com/mix-go/cli-skeleton/config/dotenv"
	_ "github.com/mix-go/cli-skeleton/di"
	"github.com/mix-go/xcli"
	"github.com/mix-go/xutil/xenv"
)

func main() {
	xcli.SetName("app").
		SetVersion("0.0.0-alpha").
		SetDebug(xenv.Getenv("APP_DEBUG").Bool(false))
	xcli.AddCommand(commands.Commands...).Run()
}
