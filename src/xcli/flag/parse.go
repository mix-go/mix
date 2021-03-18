package flag

import (
	"github.com/mix-go/xcli/argv"
	"os"
	"regexp"
	"strings"
)

// 初始化
func init() {
	Parse()
}

// Parse 解析参数
func Parse() {
	var o = make(map[string]string, 0)
	var a []string
	s := 1
	if argv.Command() == "" {
		s = 0
	}
	ignore := ""
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
			if name[:1] == "-" && value == "" && len(os.Args)-1 >= k+1 && os.Args[k+1][:1] != "-" {
				next := os.Args[k+1]
				ok, _ := regexp.MatchString(`^[\S\s]+$`, next)
				if ok {
					value = next
					ignore = next
				}
			}
		} else {
			name = ""
			if v != ignore {
				a = append(a, v)
			}
		}
		if name != "" {
			o[name] = value
		}
	}
	opts = options(o)
	args = arguments(a)
}
