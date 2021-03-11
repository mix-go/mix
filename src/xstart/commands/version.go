package commands

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	FrameworkVersion = "1.0.24"
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
	fmt.Println(fmt.Sprintf("System      Name:      %s", runtime.GOOS))
	fmt.Println(fmt.Sprintf("Go          Version:   %s", runtime.Version()[2:]))
	fmt.Println(fmt.Sprintf("Framework   Version:   %s", FrameworkVersion))
}
