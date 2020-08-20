package dotenv

import (
    "os"
    "strconv"
)

type Value struct {
    v string
}

func (t *Value) String() string {
    return t.v
}

func (t *Value) Bool() bool {
    switch t.v {
    case "", "0", "false":
        return false
    default:
        return true
    }
}

func (t *Value) Int64() int64 {
    v, _ := strconv.ParseInt(t.v, 10, 64)
    return v
}

func (t *Value) Float64() float64 {
    v, _ := strconv.ParseFloat(t.v, 64)
    return v
}

func Getenv(key string) *Value {
    return &Value{os.Getenv(key)}
}
