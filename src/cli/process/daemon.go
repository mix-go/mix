package process

import (
    "bytes"
    "fmt"
    "github.com/mix-go/cli"
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
        panic(cli.NewUnsupportedError(fmt.Errorf("error: the current operating system does not support daemon execution")))
    }

    if getgid() != 1 {
        panic(fmt.Errorf("error: Daemon() can only be used in the main goroutine"))
    }

    // Getppid 父进程ID: 当父进程已经结束，在Unix中返回的ID是初始进程(1)，在Windows中仍然是同一个进程ID，该进程ID有可能已经被其他进程占用
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
