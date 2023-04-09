package argv

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// 命令行信息
var (
	prog program
	cmd  string
)

// Program 返回命令行程序信息
func Program() *program {
	return &prog
}

// Command 返回当前命令信息
func Command() string {
	return cmd
}

// 命令行程序信息
type program struct {
	Path      string `json:"path"`
	AbsPath   string `json:"absPath"`
	Dir       string `json:"dir"`
	ParentDir string `json:"parentDir"`
	File      string `json:"file"`
}

// 创建命令行程序信息
func newProgram() program {
	p := program{
		Path:      os.Args[0],
		AbsPath:   "",
		Dir:       "",
		ParentDir: "",
		File:      "",
	}

	absPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return p
	}
	p.AbsPath = absPath

	dirRaw, file := filepath.Split(absPath)
	dir := dirRaw[:len(dirRaw)-1]
	p.Dir = dir
	p.File = file

	parentDir, err := filepath.Abs(fmt.Sprintf("%s/../", dir))
	if err != nil {
		return p
	}
	p.ParentDir = parentDir
	return p
}

// 创建当前命令信息
func newCommand(singleton bool) string {
	if len(os.Args) <= 1 || singleton {
		return ""
	}
	cmd := ""
	if ok, _ := regexp.MatchString(`^[a-zA-Z0-9_\-:]+$`, os.Args[1]); ok {
		cmd = os.Args[1]
		if cmd[:1] == "-" {
			cmd = ""
		}
	}
	return cmd
}
