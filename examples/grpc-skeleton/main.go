package main

import (
	"github.com/mix-go/dotenv"
	"github.com/mix-go/grpc-skeleton/commands"
	_ "github.com/mix-go/grpc-skeleton/configor"
	_ "github.com/mix-go/grpc-skeleton/di"
	_ "github.com/mix-go/grpc-skeleton/dotenv"
	"github.com/mix-go/xcli"
)

func main() {
	xcli.SetName("app").
		SetVersion("0.0.0-alpha").
		SetDebug(dotenv.Getenv("APP_DEBUG").Bool(false))
	xcli.AddCommand(commands.Commands...).Run()
}
