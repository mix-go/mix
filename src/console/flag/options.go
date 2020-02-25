package flag

import (
    "mix/src/console/cli"
    "os"
    "regexp"
    "strings"
)

var (
    Options = NewOptions()
)

func NewOptions() map[string]string {
    ops := make(map[string]string, 0)
    s := 1
    if cli.Command == "" {
        s = 0
    }
    for k, v := range os.Args {
        if k <= s {
            continue
        }
        name := v
        value := ""
        if strings.Contains(v, "=") {
            name = strings.Split(v, "=")[0]
            value = v[strings.Index(v, "=")+1:]
        }
        if (len(name) >= 1 && name[:1] == "-") || (len(name) >= 2 && name[:2] == "--") {
            // 无值参数处理
            if name[:1] == "-" && value == "" && len(os.Args)-1 >= k+1 && os.Args[k+1][:1] != "-" {
                next := os.Args[k+1]
                ok, _ := regexp.MatchString(`^[\S\s]+$`, next)
                if ok {
                    value = next
                }
            }
        } else {
            name = ""
        }
        if name != "" {
            ops[name] = value
        }
    }
    return ops
}
