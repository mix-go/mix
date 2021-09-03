package commands

import (
	"fmt"
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
	fmt.Println(fmt.Sprintf("Server      Name:      %s", "mix-grpc"))
	fmt.Println(fmt.Sprintf("System      Name:      %s", runtime.GOOS))
	fmt.Println(fmt.Sprintf("Go          Version:   %s", runtime.Version()[2:]))
	fmt.Println(fmt.Sprintf("Listen      Addr:      %s", netListener.Addr()))
}
