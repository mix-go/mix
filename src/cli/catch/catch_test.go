package catch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	run = false
)

func foo() {
	run = true
}

func TestCatchCall(t *testing.T) {
	a := assert.New(t)

	run = false
	Call(foo)
	a.True(run)

	run = false
	Call(func() {
		run = true
	})
	a.True(run)

	testCall1(a)
	testCall2(a)
}

func testCall1(a *assert.Assertions) {
	defer func() {
		if err := recover(); err != nil {
			a.EqualError(err.(error), "Invalid type: 'fn' is not func")
		}
	}()

	Call(nil)
}

func testCall2(a *assert.Assertions) {
	defer func() {
		if err := recover(); err != nil {
			a.EqualError(err.(error), "Invalid type: 'fn' is not func")
		}
	}()

	Call(struct{}{})
}
