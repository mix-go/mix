package manifest

import (
    "github.com/mix-go/console"
    "github.com/mix-go/dotenv"
    "github.com/mix-go/mix/devtool/manifest/beans"
    "github.com/mix-go/mix/devtool/manifest/commands"
)

var (
    ApplicationDefinition console.ApplicationDefinition
)

func Init() {
    beans.Init()

    ApplicationDefinition = console.ApplicationDefinition{
        AppName:    "mix",
        AppVersion: "1.0.3",
        AppDebug:   dotenv.Getenv("APP_DEBUG").Bool(false),
        Beans:      beans.Beans,
        Commands:   commands.Commands,
    }
}
