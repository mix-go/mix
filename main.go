package main

import (
    console2 "github.com/mix-go/console"
    "github.com/mix-go/mix/devtool/manifest"
)

func init() {
    // Manifest
    manifest.Init()
}

func main() {
    // App
    console2.NewApplication(manifest.ApplicationDefinition, "eventDispatcher", "error").Run()
}
