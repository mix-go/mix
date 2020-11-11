package manifest

import (
	"github.com/mix-go/console"
	"github.com/mix-go/dotenv"
	commands2 "github.com/mix-go/mix/devtool/commands"
	"github.com/mix-go/mix/devtool/manifest/beans"
	"github.com/mix-go/mix/devtool/manifest/commands"
)

var (
	ApplicationDefinition console.ApplicationDefinition
)

func Init() {
	beans.Init()

	ApplicationDefinition = console.ApplicationDefinition{
		Name:     "Mix",
		Version:  commands2.FrameworkVersion,
		Debug:    dotenv.Getenv("APP_DEBUG").Bool(false),
		Beans:    beans.Beans,
		Commands: commands.Commands,
	}
}
