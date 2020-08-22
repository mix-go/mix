package dotenv

import (
    "os"
    "strconv"
)

type envValue struct {
    v string
}

func (t *envValue) String(val ...string) string {
    d := ""
    if len(val) >= 1 {
        d = val[0]
    }

    if t.v == "" {
        return d
    }

    return t.v
}

func (t *envValue) Bool(val ...bool) bool {
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

func (t *envValue) Int64(val ...int64) int64 {
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

func (t *envValue) Float64(val ...float64) float64 {

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

func Getenv(key string) *envValue {
    return &envValue{os.Getenv(key)}
}
