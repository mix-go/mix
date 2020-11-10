package argv

import (
	"os"
	"path/filepath"
	"regexp"
)

// 命令行信息
var (
	prog program
	cmd  string
)

// 初始化
func init() {
	Parse()
}

// 解析命令行参数
func Parse() {
	prog = newProgram()
	cmd = newCommand()
}

// 返回命令行程序信息
func Program() *program {
	return &prog
}

// 返回当前执行的命令信息
func Command() string {
	return cmd
}

// 命令行程序信息
type program struct {
	Path    string
	AbsPath string
	Dir     string
	File    string
}

// 创建命令行程序信息
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

// 创建当前执行的命令信息
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
