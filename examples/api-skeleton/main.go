package main

import (
	"github.com/mix-go/api-skeleton/commands"
	_ "github.com/mix-go/api-skeleton/config/configor"
	_ "github.com/mix-go/api-skeleton/config/dotenv"
	_ "github.com/mix-go/api-skeleton/di"
	"github.com/mix-go/xcli"
	"github.com/mix-go/xutil/xenv"
)

func main() {
	xcli.SetName("app").
		SetVersion("0.0.0-alpha").
		SetDebug(xenv.Getenv("APP_DEBUG").Bool(false))
	xcli.AddCommand(commands.Commands...).Run()
}
