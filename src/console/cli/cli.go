package cli

import (
    "os"
    "path/filepath"
    "regexp"
)

var (
    Program ProgramMeta = NewProgram()
    Command string      = NewCommand()
)

func Parse() {
    Program = NewProgram()
    Command = NewCommand()
}

type ProgramMeta struct {
    Name string
    Dir  string
    Path string
}

func NewProgram() ProgramMeta {
    dir, file := filepath.Split(os.Args[0])
    return ProgramMeta{
        Name: file,
        Dir:  dir[:len(dir)-1],
        Path: os.Args[0],
    }
}

func NewCommand() string {
    cmd := ""
    if len(os.Args) <= 1 {
        return cmd
    }
    ok, _ := regexp.MatchString(`^[a-zA-Z0-9_\-:]+$`, os.Args[1])
    if ok {
        cmd = os.Args[1]
        if cmd[:1] == "-" {
            cmd = ""
        }
    }
    return cmd
}
