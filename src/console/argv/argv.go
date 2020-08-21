package argv

import (
    "os"
    "path/filepath"
    "regexp"
)

var (
    program ProgramMeta
    command string
)

func init() {
    Parse()
}

func Parse() {
    program = newProgram()
    command = newCommand()
}

func Program() ProgramMeta {
    return program
}

func Command() string {
    return command
}

type ProgramMeta struct {
    Path    string
    AbsPath string
    Dir     string
    File    string
}

func newProgram() ProgramMeta {
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

func newCommand() string {
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
