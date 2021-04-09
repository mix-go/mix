package xdi

import (
	"fmt"
	"sync"
)

// Object
type Object struct {
	Name string
	// 创建对象的闭包
	New func() (interface{}, error)
	// 每次都创建新的对象
	NewEveryTime bool

	refresher refresher
}

type refresher struct {
	mux sync.Mutex
	val bool
}

func (t *refresher) on() {
	t.val = true
}

func (t *refresher) off() {
	t.val = false
}

func (t *refresher) status() bool {
	return t.val
}

func (t *refresher) invoke(f func()) {
	if !t.status() {
		return
	}
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.status() {
		f()
		t.off()
	}
}

func (t *Object) Refresh() error {
	if t.NewEveryTime {
		return fmt.Errorf("error: '%s' is NewEveryTime, unable to refresh", t.Name)
	}
	t.refresher.on()
	return nil
}
