package main

import (
	"github.com/mix-go/web-skeleton/commands"
	_ "github.com/mix-go/web-skeleton/config/configor"
	_ "github.com/mix-go/web-skeleton/config/dotenv"
	_ "github.com/mix-go/web-skeleton/di"
	"github.com/mix-go/xcli"
	"github.com/mix-go/xutil/xenv"
)

func main() {
	xcli.SetName("app").
		SetVersion("0.0.0-alpha").
		SetDebug(xenv.Getenv("APP_DEBUG").Bool(false))
	xcli.AddCommand(commands.Commands...).Run()
}
