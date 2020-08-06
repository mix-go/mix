package argv

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
    Path    string
    AbsPath string
    Dir     string
    File    string
}

func NewProgram() ProgramMeta {
    abspath, err := filepath.Abs(os.Args[0])
    if err != nil {
        panic(err)
    }
    dir, file := filepath.Split(abspath)
    return ProgramMeta{
        Path:    os.Args[0],
        AbsPath: abspath,
        Dir:     dir[:len(dir)-1],
        File:    file,
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
