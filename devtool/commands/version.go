package commands

import (
    "fmt"
    "strings"
)

var (
    ConsoleSkeletonVersion = "v1.0.1"
    APISkeletonVersion     = "none"
    WebSkeletonVersion     = "none"
    FrameworkVersion       = "v1.0.1"
)

const logo = `             ___         
 ______ ___  _ /__ ___ _____ ______ 
  / __ *__ \/ /\ \/ /__  __ */  __ \
 / / / / / / / /\ \/ _  /_/ // /_/ /
/_/ /_/ /_/_/ /_/\_\  \__, / \____/ 
                     /____/
`

type VersionCommand struct {
}

func (t *VersionCommand) Main() {
    fmt.Println(strings.Replace(logo, "*", "`", -1))
    fmt.Println("")
    fmt.Println(fmt.Sprintf("Console      Version: %s", ConsoleSkeletonVersion))
    fmt.Println(fmt.Sprintf("API          Version: %s", APISkeletonVersion))
    fmt.Println(fmt.Sprintf("Web          Version: %s", WebSkeletonVersion))
    fmt.Println(fmt.Sprintf("Framework    Version: %s", FrameworkVersion))
}
