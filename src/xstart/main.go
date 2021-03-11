package main

import (
	"github.com/mix-go/console"
	"github.com/mix-go/mix/devtool/manifest"
)

func init() {
	// Manifest
	manifest.Init()
}

func main() {
	// App
	console.NewApplication(manifest.ApplicationDefinition, "eventDispatcher", "error").Run()
}
