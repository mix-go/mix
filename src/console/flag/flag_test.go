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

	v1 := Match("a").String()
	a.Equal(v1, "a1")

	v2 := Match("b").Bool()
	a.Equal(v2, true)

	v3 := Match("cd").Bool()
	a.Equal(v3, true)

	v4 := Match("sw").Bool()
	a.Equal(v4, false)

	v5 := Match("ab", "").String()
	a.Equal(v5, "ab1")

	v6 := Match("de", "").String()
	a.Equal(v6, "de1")

	v7 := Match("c", "").String()
	a.Equal(v7, "c1")
}

func TestMatch(t *testing.T) {
	a := assert.New(t)

	os.Args = []string{os.Args[0], "foo", "-a=a1", "-b", "--bc", "--ab=ab1", "--de", "de1", "-c", "c1", "--sw", "false"}
	argv.Parse()
	Parse()

	v1 := Match("b", "bc").Bool()
	a.Equal(v1, true)

	v2 := Match("a", "ab").String()
	a.Equal(v2, "a1")
}

func TestNotFound(t *testing.T) {
	a := assert.New(t)

	os.Args = []string{os.Args[0]}
	argv.Parse()
	Parse()

	v1 := Match("cde").Bool()
	a.Equal(v1, false)

	v2 := Match("x").String()
	a.Equal(v2, "")

	v3 := Match("b", "bc").Bool()
	a.Equal(v3, false)

	v4 := Match("a", "ab").String()
	a.Equal(v4, "")
}

func TestOptions(t *testing.T) {
	a := assert.New(t)

	os.Args = []string{os.Args[0], "foo", "-a=a1", "-b", "--cd", "--ab=ab1", "arg0", "--de", "de1", "-c", "c1", "--sw", "false", "arg1", "arg2"}
	argv.Parse()
	Parse()

	a.Equal(Options().Map(), map[string]string{"--ab": "ab1", "--cd": "", "--de": "de1", "--sw": "false", "-a": "a1", "-b": "", "-c": "c1"})
}

func TestArguments(t *testing.T) {
	a := assert.New(t)

	os.Args = []string{os.Args[0], "foo", "-a=a1", "-b", "--cd", "--ab=ab1", "arg0", "--de", "de1", "-c", "c1", "--sw", "false", "arg1", "arg2"}
	argv.Parse()
	Parse()

	a.Equal(Arguments().Array(), []string{"arg0", "arg1", "arg2"})

	a.Equal(Arguments().First().String(), "arg0")
	a.Equal(Arguments().First().Bool(), true)
}
