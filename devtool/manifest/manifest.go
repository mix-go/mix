package manifest

import (
    "github.com/mix-go/console"
    "github.com/mix-go/mix/devtool/manifest/beans"
    "github.com/mix-go/mix/devtool/manifest/commands"
)

var (
    ApplicationDefinition console.ApplicationDefinition
)

func Init() {
    beans.Init()

    ApplicationDefinition = console.ApplicationDefinition{
        AppName:    "app",
        AppVersion: "1.0.2",
        AppDebug:   false,
        Beans:      beans.Beans,
        Commands:   commands.Commands,
    }
}
