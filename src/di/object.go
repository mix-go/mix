package di

import (
    "fmt"
    "sync"
)

// Object
type Object struct {
    Name      string
    New       func() (interface{}, error)
    Singleton bool

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
    if !t.Singleton {
        return fmt.Errorf("error: '%s' not a singleton, unable to refresh", t.Name)
    }
    t.refresher.on()
    return nil
}
