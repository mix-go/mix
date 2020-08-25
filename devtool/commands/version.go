package commands

import (
    "fmt"
    "strings"
)

var (
    SkeletonVersion  = "v1.0.0"
    FrameworkVersion = "v1.0.0"
    ConsoleVersion   = "v0.0.0-20200822120924-3cb721b6d16f"
    APIVersion       = ""
    WebVersion       = ""
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
    fmt.Println(fmt.Sprintf("Skeleton     Version: %s", SkeletonVersion))
    fmt.Println(fmt.Sprintf("Framework    Version: %s", FrameworkVersion))
}
