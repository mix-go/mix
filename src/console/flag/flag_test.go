package flag

import (
    "github.com/mix-go/console/argv"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestSingle(t *testing.T) {
    a := assert.New(t)

    os.Args = []string{os.Args[0], "foo", "-a=a1", "-b", "--cd", "--ab=ab1", "--de", "de1", "-c", "c1", "--sw", "false"}
    argv.Parse()
    Parse()

    v1 := String("a", "")
    a.Equal(v1, "a1")

    v2 := Bool("b", false)
    a.Equal(v2, true)

    v3 := Bool("cd", false)
    a.Equal(v3, true)

    v4 := Bool("sw", false)
    a.Equal(v4, false)

    v5 := String("ab", "")
    a.Equal(v5, "ab1")

    v6 := String("de", "")
    a.Equal(v6, "de1")

    v7 := String("c", "")
    a.Equal(v7, "c1")
}

func TestMatch(t *testing.T) {
    a := assert.New(t)

    os.Args = []string{os.Args[0], "foo", "-a=a1", "-b", "--bc", "--ab=ab1", "--de", "de1", "-c", "c1", "--sw", "false"}
    argv.Parse()
    Parse()

    v1 := BoolMatch([]string{"b", "bc"}, false)
    a.Equal(v1, true)

    v2 := StringMatch([]string{"a", "ab"}, "")
    a.Equal(v2, "a1")
}

func TestNotFound(t *testing.T) {
    a := assert.New(t)

    os.Args = []string{os.Args[0]}
    argv.Parse()
    Parse()

    v1 := Bool("cde", false)
    a.Equal(v1, false)

    v2 := String("x", "")
    a.Equal(v2, "")

    v3 := BoolMatch([]string{"b", "bc"}, false)
    a.Equal(v3, false)

    v4 := StringMatch([]string{"a", "ab"}, "")
    a.Equal(v4, "")
}
