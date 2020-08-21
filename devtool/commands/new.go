package commands

import (
    "fmt"
    "github.com/mix-go/console/flag"
)

type NewCommand struct {
}

func (t *NewCommand) Main() {
    name := flag.StringMatch([]string{"n", "name"}, "hello")
    typ := flag.StringMatch([]string{"t", "type"}, "console")
    fmt.Println(name, typ)
}
