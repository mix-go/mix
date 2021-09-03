package commands

import (
	"fmt"
	"github.com/mix-go/xcli/flag"
)

type HelloCommand struct {
}

func (t *HelloCommand) Main() {
	name := flag.Match("n", "name").String("OpenMix")
	say := flag.Match("say").String("Hello, World!")
	fmt.Printf("%s: %s\n", name, say)
}
