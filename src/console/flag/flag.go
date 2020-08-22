package flag

import (
    "fmt"
    "strconv"
)

type flagValue struct {
    v string
}

func (t *flagValue) String(val ...string) string {
    d := ""
    if len(val) >= 1 {
        d = val[0]
    }

    if t.v == "" {
        return d
    }

    return t.v
}

func (t *flagValue) Bool(val ...bool) bool {
    d := false
    if len(val) >= 1 {
        d = val[0]
    }

    switch t.v {
    case "":
        return d
    case "0", "false":
        return false
    default:
        return true
    }
}

func (t *flagValue) Int64(val ...int64) int64 {
    d := int64(0)
    if len(val) >= 1 {
        d = val[0]
    }

    if t.v == "" {
        return d
    }

    v, _ := strconv.ParseInt(t.v, 10, 64)
    return v
}

func (t *flagValue) Float64(val ...float64) float64 {

    d := float64(0)
    if len(val) >= 1 {
        d = val[0]
    }

    if t.v == "" {
        return d
    }

    v, _ := strconv.ParseFloat(t.v, 64)
    return v
}

func Match(names ...string) *flagValue {
    for _, name := range names {
        v := value(name)
        if v != "" {
            return &flagValue{v}
        }
    }
    return &flagValue{}
}

func value(name string) string {
    key := ""
    if len(name) == 1 {
        key = fmt.Sprintf("-%s", name)
    } else {
        key = fmt.Sprintf("--%s", name)
    }
    if v, ok := Options()[key]; ok {
        return v
    }
    return key
}
