package commands

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/mix-go/cli/flag"
	"github.com/mix-go/xstart/github/loading"
	"github.com/mix-go/xstart/logic"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	CLI  = "console"
	API  = "api"
	Web  = "web"
	gRPC = "grpc"
)

type NewCommand struct {
}

func (t *NewCommand) Main() {
	name := flag.Arguments().First().String("hello")

	prompt := promptui.Select{
		Label: "Select project type:",
		Items: []string{"CLI", "API", "Web (contains the websocket)", "gRPC"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	typ := CLI
	switch result {
	case "CLI":
		typ = CLI
		break
	case "API":
		typ = API
		break
	case "Web (contains the websocket)":
		typ = Web
		break
	case "gRPC":
		typ = gRPC
		break
	}

	t.NewProject(name, typ)
}

func (t *NewCommand) NewProject(name, typ string) {
	ver := ""
	switch typ {
	case CLI, API, Web, gRPC:
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
		cmd := exec.Command("go", "get", fmt.Sprintf("github.com/mix-go/%s-skeleton@%s", typ, ver))
		loader := loading.StartNew(fmt.Sprintf("Skeleton local not found, exec 'go get github.com/mix-go/%s-skeleton@%s'", typ, ver))
		err = cmd.Run()
		loader.Stop()
		if err != nil {
			fmt.Println(fmt.Sprintf("Exec failed: %s", err.Error()))
			fmt.Println("Please try again, or manually execute 'go get ***'")
			return
		}
		_ = os.Remove(fmt.Sprintf("%s/bin/%s-skeleton", os.Getenv("GOPATH"), typ))
		fmt.Println(fmt.Sprintf("Skeleton 'github.com/mix-go/%s-skeleton@%s' download is completed", typ, ver))
	} else {
		fmt.Println(fmt.Sprintf("Skeleton 'github.com/mix-go/%s-skeleton@%s' local found", typ, ver))
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

	fmt.Println(fmt.Sprintf("Project '%s' is generated", name))
}
