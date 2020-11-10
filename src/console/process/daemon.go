package process

import (
	"bytes"
	"errors"
	"github.com/mix-go/console"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

// Daemon 使当前进程蜕变为一个守护进程
func Daemon() {
	var ok bool
	switch runtime.GOOS {
	case "darwin", "linux":
		ok = true
	case "windows":
		ok = false
	default:
		ok = true
	}

	if !ok {
		panic(console.NewUnsupportError(errors.New("Error: the current OS does not support background execution")))
	}

	if getgid() != 1 {
		panic(errors.New("Daemon() can only be used in Main Goroutine"))
	}

	if os.Getppid() != 1 {
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		if err := cmd.Start(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}
}

// 获取协程id
func getgid() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
