package process

import (
    "errors"
    "github.com/mix-go/console"
    "os"
    "os/exec"
    "runtime"
)

// 使当前进程蜕变为一个守护进程
func Daemon() {
    ok := false
    switch runtime.GOOS {
    case "darwin", "linux":
        ok = true
    case "windows":
        ok = false
    default:
        ok = true
    }

    if !ok {
        panic(console.UnsupportError(errors.New("The current OS does not support background execution")))
    }

    if os.Getppid() != 1 {
        cmd := exec.Command(os.Args[0], os.Args[1:]...)
        if err := cmd.Start(); err != nil {
            panic(err)
        }
        os.Exit(0)
    }
}
