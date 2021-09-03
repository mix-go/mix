package commands

import (
    "fmt"
	"github.com/mix-go/api-skeleton/di"
	"github.com/mix-go/xcli"
    "runtime"
    "strings"
)

const logo = `             ___         
 ______ ___  _ /__ ___ _____ ______ 
  / __ *__ \/ /\ \/ /__  __ */  __ \
 / / / / / / / /\ \/ _  /_/ // /_/ /
/_/ /_/ /_/_/ /_/\_\  \__, / \____/ 
                     /____/
`

func welcome() {
    fmt.Println(strings.Replace(logo, "*", "`", -1))
    fmt.Println("")
    fmt.Println(fmt.Sprintf("Server      Name:      %s", "mix-api"))
    fmt.Println(fmt.Sprintf("Listen      Addr:      %s", di.Server().Addr))
    fmt.Println(fmt.Sprintf("System      Name:      %s", runtime.GOOS))
    fmt.Println(fmt.Sprintf("Go          Version:   %s", runtime.Version()[2:]))
    fmt.Println(fmt.Sprintf("Framework   Version:   %s", xcli.Version))
}
