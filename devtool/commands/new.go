package commands

import (
    "errors"
    "fmt"
    "github.com/mix-go/console/flag"
    "github.com/mix-go/mix/devtool/logic"
    "os"
    "os/exec"
)

var (
    Console = "console"
    API     = "api"
    Web     = "web"
)

type NewCommand struct {
}

func (t *NewCommand) Main() {
    name := flag.Match("n", "name").String("hello")
    t.NewProject(name, Console)
}

func (t *NewCommand) NewProject(name, typ string) {
    ver := ""
    switch typ {
    case "console":
        ver = ConsoleSkeletonVersion
        break
    case "api":
        ver = APISkeletonVersion
        break
    case "web":
        ver = WebSkeletonVersion
        break
    default:
        fmt.Println("Type error, only be console, api, web")
        return
    }

    sDir := fmt.Sprintf("%s/pkg/mod/github.com/mix-go/mix-%s-skeleton@%s", os.Getenv("GOPATH"), typ, ver)
    if _, err := os.Stat(sDir); err != nil {
        fmt.Println(fmt.Sprintf("Skeleton '%s' local not found, exec 'go get -u github.com/mix-go/mix-%s-skeleton@%s', please wait ...", typ, typ, ver))
        cmd := exec.Command("go", "get", "-u", fmt.Sprintf("github.com/mix-go/mix-%s-skeleton@%s", typ, ver))
        err = cmd.Run()
        if err != nil {
            fmt.Println(fmt.Sprintf("Exec failure: %s", err.Error()))
            fmt.Println("Please try again, or manually execute 'go get -u ***'")
            return
        }
    }

    pwd, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    dest := fmt.Sprintf("%s/%s", pwd, name)
    if !logic.CopyPath(sDir, dest) {
        panic(errors.New("Copy dir failed"))
    }

    if err := logic.ReplaceName(dest, name); err != nil {
        panic(errors.New("Replace name failed"))
    }

    if err := logic.ReplaceMod(dest); err != nil {
        panic(errors.New("Replace go.mod failed"))
    }

    fmt.Println(fmt.Sprintf("Skeleton '%s' is generated", name))
}

type APICommand struct {
    NewCommand
}

func (t *APICommand) Main() {
    name := flag.Match("n", "name").String("hello")
    t.NewProject(name, API)
}

type WebCommand struct {
    NewCommand
}

func (t *WebCommand) Main() {
    name := flag.Match("n", "name").String("hello")
    t.NewProject(name, Web)
}
