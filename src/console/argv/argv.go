package argv

import (
    "os"
    "path/filepath"
    "regexp"
)

var (
    prog program
    cmd  string
)

func init() {
    Parse()
}

func Parse() {
    prog = newProgram()
    cmd = newCommand()
}

func Program() *program {
    return &prog
}

func Command() string {
    return cmd
}

type program struct {
    Path    string
    AbsPath string
    Dir     string
    File    string
}

func newProgram() program {
    abspath, err := filepath.Abs(os.Args[0])
    if err != nil {
        panic(err)
    }
    dir, file := filepath.Split(abspath)
    return program{
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
