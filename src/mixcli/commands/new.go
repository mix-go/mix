package commands

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/manifoldco/promptui"
	"github.com/mix-go/mixcli/logic"
	"github.com/mix-go/xcli/flag"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	None     = "none"
	CLI      = "cli"
	API      = "api"
	Web      = "web"
	gRPC     = "grpc"
	Gorm     = "gorm"
	Xorm     = "xorm"
	Zap      = "zap"
	Logrus   = "logrus"
	GoRedis  = "go-redis"
	DotEnv   = "dotenv"
	Configor = "configor"
	Viper    = "viper"
)

type NewCommand struct {
}

func (t *NewCommand) Main() {
	name := flag.Arguments().First().String("hello")
	name = strings.ReplaceAll(name, " ", "")

	promp := func(label string, items []string) string {
		prompt := promptui.Select{
			Label: label,
			Items: items,
		}
		prompt.HideSelected = true
		_, result, err := prompt.Run()
		if err != nil {
			return ""
		}
		return result
	}

	var selectType string
	switch promp("Select project type", []string{"CLI", "API", "Web (contains the websocket)", "gRPC"}) {
	case "CLI":
		selectType = CLI
		break
	case "API":
		selectType = API
		break
	case "Web (contains the websocket)":
		selectType = Web
		break
	case "gRPC":
		selectType = gRPC
		break
	default:
		return
	}

	selectEnv := promp("Select .env configuration file library", []string{DotEnv, None})
	selectConf := promp("Select .yml, .json, .toml configuration files library", []string{Configor, Viper, None})

	var selectLog string
	var selectLogItems []string
	if selectType == CLI {
		selectLogItems = []string{Zap, Logrus, None}
	} else {
		selectLogItems = []string{Zap, Logrus}
	}
	selectLog = promp("Select logger library", selectLogItems)

	var selectDb string
	var selectDbItems []string
	if selectType == API || selectType == Web {
		selectDbItems = []string{Gorm, Xorm}
	} else {
		selectDbItems = []string{Gorm, Xorm, None}
	}
	selectDb = promp("Select database library", selectDbItems)

	var selectRedis string
	selectRedis = promp("Select redis library", []string{GoRedis, None})

	t.NewProject(name, selectType, selectEnv, selectConf, selectLog, selectDb, selectRedis)
}

func (t *NewCommand) NewProject(name, selectType, selectEnv, selectConf, selectLog, selectDb, selectRedis string) {
	ver := ""
	switch selectType {
	case CLI, API, Web, gRPC:
		ver = fmt.Sprintf("v%s", SkeletonVersion)
		break
	default:
		fmt.Println("Type error, only be console, api, web, grpc")
		return
	}

	envCmd := "go env GOPATH"
	cmd := exec.Command("go", "env", "GOPATH")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Exec error: %v\n", err)
		return
	}
	goPath := string(out[:len(out)-1])
	if goPath == "" {
		fmt.Printf("$GOPATH is not configured, see '%s'\n", envCmd)
		return
	}

	dr := ":"
	if runtime.GOOS == "windows" {
		dr = ";"
	}
	if strings.Contains(goPath, dr) {
		fmt.Printf("$GOPATH cannot have multiple directories, see '%s'\n", envCmd)
		return
	}

	srcDir := fmt.Sprintf("%s/pkg/mod/github.com/mix-go/%s-skeleton@%s", goPath, selectType, ver)
	if _, err := os.Stat(srcDir); err != nil {
		cmd := exec.Command("go", "get", fmt.Sprintf("github.com/mix-go/%s-skeleton@%s", selectType, ver))
		fmt.Printf("Skeleton local not found, exec 'go get github.com/mix-go/%s-skeleton@%s'\n", selectType, ver)
		total := 0
		switch selectType {
		case CLI:
			total = 7695
		case API:
			total = 13834
		case Web:
			total = 17705
		case gRPC:
			total = 15659
		}
		current := int64(0)
		bar := pb.StartNew(total)
		go func() {
			path := fmt.Sprintf("%s/pkg/mod/cache/download/github.com/mix-go/%s-skeleton/@v/%s.zip", goPath, selectType, ver)
			for {
				f, err := os.Open(path)
				if err != nil {
					continue
				}
				fi, err := f.Stat()
				if err != nil {
					_ = f.Close()
					continue
				}
				current = fi.Size()
				bar.SetCurrent(current)
				_ = f.Close()
				time.Sleep(time.Millisecond * 100)
			}
		}()
		err = cmd.Run()
		if err == nil {
			bar.SetTotal(current)
			bar.SetCurrent(current)
		} else {
			bar.SetTotal(0)
			bar.SetCurrent(0)
		}
		bar.Finish()
		if err != nil {
			fmt.Println(fmt.Sprintf("Exec failed: %s", err.Error()))
			fmt.Println("Please try again, or manually execute 'go get ***'")
			return
		}
		time.Sleep(2 * time.Second) // 等待一会，让 gomod 完成解压
		_ = os.Remove(fmt.Sprintf("%s/bin/%s-skeleton", goPath, selectType))
		fmt.Println(fmt.Sprintf("Skeleton 'github.com/mix-go/%s-skeleton@%s' download is completed", selectType, ver))
	} else {
		fmt.Println(fmt.Sprintf("Skeleton 'github.com/mix-go/%s-skeleton@%s' local found", selectType, ver))
	}

	fmt.Print(" - Generate code")
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dest := fmt.Sprintf("%s/%s", pwd, name)
	if !logic.CopyPath(srcDir, dest) {
		panic(errors.New("Copy dir failed"))
	}
	fmt.Println(" > ok")

	fmt.Print(" - Processing .env")
	if selectEnv == None {
		if err := logic.ReplaceMain(dest, fmt.Sprintf(`_ "github.com/mix-go/%s-skeleton/config/dotenv"`, selectType), ""); err != nil {
			panic(errors.New("Replace failed"))
		}
		_ = os.RemoveAll(fmt.Sprintf("%s/config/dotenv", dest))
		_ = os.RemoveAll(fmt.Sprintf("%s/.env", dest))
	}
	fmt.Println(" > ok")

	fmt.Print(" - Processing conf")
	switch selectConf {
	case Configor:
		_ = os.RemoveAll(fmt.Sprintf("%s/config/viper", dest))
		break
	case Viper:
		if err := logic.ReplaceMain(dest, fmt.Sprintf(`_ "github.com/mix-go/%s-skeleton/config/configor"`, selectType), fmt.Sprintf(`_ "github.com/mix-go/%s-skeleton/config/viper"`, selectType)); err != nil {
			panic(errors.New("Replace failed"))
		}
		_ = os.RemoveAll(fmt.Sprintf("%s/config/configor", dest))
		break
	case None:
		if err := logic.ReplaceMain(dest, fmt.Sprintf(`_ "github.com/mix-go/%s-skeleton/config/configor"`, selectType), ""); err != nil {
			panic(errors.New("Replace failed"))
		}
		_ = os.RemoveAll(fmt.Sprintf("%s/config/viper", dest))
		_ = os.RemoveAll(fmt.Sprintf("%s/config/configor", dest))
		_ = os.RemoveAll(fmt.Sprintf("%s/config/main.go", dest))
		_ = os.RemoveAll(fmt.Sprintf("%s/conf", dest))
	}
	fmt.Println(" > ok")

	fmt.Print(" - Processing logger")
	switch selectLog {
	case Zap:
		if err := logic.ReplaceAll(dest, `logger := di.Logrus`, `logger := di.Zap`); err != nil {
			panic(errors.New("Replace failed"))
		}
		if err := logic.ReplaceAll(dest, `Output: logger.Writer\(\)`, `Output: &di.ZapOutput{Logger: logger}`); err != nil {
			panic(errors.New("Replace failed"))
		}
		_ = os.Remove(fmt.Sprintf("%s/di/logrus.go", dest))
		break
	case Logrus:
		_ = os.Remove(fmt.Sprintf("%s/di/zap.go", dest))
		break
	case None:
		_ = os.Remove(fmt.Sprintf("%s/di/logrus.go", dest))
		_ = os.Remove(fmt.Sprintf("%s/di/zap.go", dest))
		break
	}
	fmt.Println(" > ok")

	fmt.Print(" - Processing database")
	switch selectDb {
	case Gorm:
		_ = os.Remove(fmt.Sprintf("%s/di/xorm.go", dest))
		break
	case Xorm:
		_ = os.Remove(fmt.Sprintf("%s/di/gorm.go", dest))
		break
	case None:
		_ = os.Remove(fmt.Sprintf("%s/di/gorm.go", dest))
		_ = os.Remove(fmt.Sprintf("%s/di/xorm.go", dest))
		break
	}
	fmt.Println(" > ok")

	fmt.Print(" - Processing redis")
	switch selectRedis {
	case GoRedis:
		break
	case None:
		_ = os.Remove(fmt.Sprintf("%s/di/goredis.go", dest))
		break
	}
	fmt.Println(" > ok")

	// 都没有选
	if selectLog == None && selectDb == None && selectRedis == None {
		if err := logic.ReplaceMain(dest, fmt.Sprintf(`_ "github.com/mix-go/%s-skeleton/di"`, selectType), ""); err != nil {
			panic(errors.New("Replace failed"))
		}
	}

	fmt.Print(" - Processing package name")
	if err := logic.ReplaceAll(dest, fmt.Sprintf("github.com/mix-go/%s-skeleton", selectType), name); err != nil {
		panic(errors.New("Replace failed"))
	}
	if err := logic.ReplaceMod(dest); err != nil {
		panic(errors.New("Replace go.mod failed"))
	}
	fmt.Println(" > ok")

	fmt.Println(fmt.Sprintf("Project '%s' is generated", name))
}
