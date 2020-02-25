package flag

import "strconv"

func String(name string, value string) string {
    if v, ok := Options[name]; ok {
        return v
    }
    return value
}

func Bool(name string, value bool) bool {
    v := String(name, "")
    if v == "" {
        return value
    }
    return true
}

func Int64(name string, value int64) int64 {
    v := String(name, "")
    if v == "" {
        return value
    }
    i, err := strconv.ParseInt(v, 10, 64)
    if err != nil {
        return value
    }
    return i
}

func Float64(name string, value float64) float64 {
    v := String(name, "")
    if v == "" {
        return value
    }
    f, err := strconv.ParseFloat(v, 64)
    if err != nil {
        return value
    }
    return f
}
