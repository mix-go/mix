package main

import (
	"github.com/mix-go/grpc-skeleton/commands"
	_ "github.com/mix-go/grpc-skeleton/config/configor"
	_ "github.com/mix-go/grpc-skeleton/config/dotenv"
	_ "github.com/mix-go/grpc-skeleton/di"
	"github.com/mix-go/xcli"
	"github.com/mix-go/xutil/xenv"
)

func main() {
	xcli.SetName("app").
		SetVersion("0.0.0-alpha").
		SetDebug(xenv.Getenv("APP_DEBUG").Bool(false))
	xcli.AddCommand(commands.Commands...).Run()
}
