package flag

import (
    "fmt"
    "strconv"
)

func String(name string, value string) string {
    key := ""
    if len(name) == 1 {
        key = fmt.Sprintf("-%s", name)
    } else {
        key = fmt.Sprintf("--%s", name)
    }
    if v, ok := Options[key]; ok {
        return v
    }
    return value
}

func StringMultiple(names []string, value string) string {
    for _, name := range names {
        v := String(name, value)
        if v != value {
            return v
        }
    }
    return value
}

func Bool(name string, value bool) bool {
    key := ""
    if len(name) == 1 {
        key = fmt.Sprintf("-%s", name)
    } else {
        key = fmt.Sprintf("--%s", name)
    }
    if v, ok := Options[key]; ok {
        if v == "false" {
            return false
        }
        return true
    }
    return value
}

func BoolMultiple(names []string, value bool) bool {
    for _, name := range names {
        v := Bool(name, value)
        if v != value {
            return v
        }
    }
    return value
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

func Int64Multiple(names []string, value int64) int64 {
    for _, name := range names {
        v := Int64(name, value)
        if v != value {
            return v
        }
    }
    return value
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

func Float64Multiple(names []string, value float64) float64 {
    for _, name := range names {
        v := Float64(name, value)
        if v != value {
            return v
        }
    }
    return value
}
