package commands

import (
	"errors"
	"fmt"
	"github.com/mix-go/console/flag"
	"github.com/mix-go/mix/devtool/logic"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	Console = "console"
	API     = "api"
	Web     = "web"
	gRPC    = "grpc"
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
	case Console, API, Web, gRPC:
		ver = fmt.Sprintf("v%s", FrameworkVersion)
		break
	default:
		fmt.Println("Type error, only be console, api, web, grpc")
		return
	}

	if os.Getenv("GOPATH") == "" {
		fmt.Println("$GOPATH is not configured, see 'echo $GOPATH'")
		return
	}

	dr := ":"
	if runtime.GOOS == "windows" {
		dr = ";"
	}
	if strings.Contains(os.Getenv("GOPATH"), dr) {
		fmt.Println("$GOPATH cannot have multiple directories, see 'echo $GOPATH'")
		return
	}

	sdir := fmt.Sprintf("%s/pkg/mod/github.com/mix-go/%s-skeleton@%s", os.Getenv("GOPATH"), typ, ver)
	if _, err := os.Stat(sdir); err != nil {
		fmt.Println(fmt.Sprintf("Skeleton '%s' local not found, exec 'go get github.com/mix-go/%s-skeleton@%s', please wait ...", typ, typ, ver))
		cmd := exec.Command("go", "get", fmt.Sprintf("github.com/mix-go/%s-skeleton@%s", typ, ver))
		err = cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprintf("Exec failed: %s", err.Error()))
			fmt.Println("Please try again, or manually execute 'go get ***'")
			return
		}
		_ = os.Remove(fmt.Sprintf("%s/bin/%s-skeleton", os.Getenv("GOPATH"), typ))
	} else {
		fmt.Println(fmt.Sprintf("Skeleton '%s' local found", typ))
	}

	fmt.Print(" - Generate code")
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dest := fmt.Sprintf("%s/%s", pwd, name)
	if !logic.CopyPath(sdir, dest) {
		panic(errors.New("Copy dir failed"))
	}
	fmt.Println(" > ok")

	fmt.Print(" - Rewrite package name")
	if err := logic.ReplaceName(dest, fmt.Sprintf("github.com/mix-go/%s-skeleton", typ), name); err != nil {
		panic(errors.New("Replace name failed"))
	}
	if err := logic.ReplaceMod(dest); err != nil {
		panic(errors.New("Replace go.mod failed"))
	}
	fmt.Println(" > ok")

	fmt.Println(fmt.Sprintf("Application '%s' is generated", name))
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

type GrpcCommand struct {
	NewCommand
}

func (t *GrpcCommand) Main() {
	name := flag.Match("n", "name").String("hello")
	t.NewProject(name, gRPC)
}
