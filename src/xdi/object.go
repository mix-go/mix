package xdi

import (
	"fmt"
	"sync"
)

type Object struct {
	Name string
	// 创建对象的闭包
	New func() (interface{}, error)
	// 每次都创建新的对象
	NewEverytime bool

	refresher refresher
	mutex     sync.Mutex
}

func (t *Object) Refresh() error {
	if t.NewEverytime {
		return fmt.Errorf("xdi: '%s' is NewEverytime, unable to refresh", t.Name)
	}
	t.refresher.on()
	return nil
}

type refresher struct {
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
