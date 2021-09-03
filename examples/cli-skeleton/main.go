package main

import (
	"github.com/mix-go/cli-skeleton/commands"
	_ "github.com/mix-go/cli-skeleton/config/configor"
	_ "github.com/mix-go/cli-skeleton/config/dotenv"
	_ "github.com/mix-go/cli-skeleton/di"
	"github.com/mix-go/dotenv"
	"github.com/mix-go/xcli"
)

func main() {
	xcli.SetName("app").
		SetVersion("0.0.0-alpha").
		SetDebug(dotenv.Getenv("APP_DEBUG").Bool(false))
	xcli.AddCommand(commands.Commands...).Run()
}
