package di

import (
	"fmt"
	"reflect"
	"sync"
)

var container *iContainer

func init() {
	container = New()
}

// New
func New() *iContainer {
	container := &iContainer{}
	return container
}

// Container
func Container() *iContainer {
	return container
}

// Provide
func Provide(objects ...*Object) error {
	return container.Provide(objects...)
}

// Find
func Find(name string, ptr interface{}) error {
	return container.Find(name, ptr)
}

// iContainer
type iContainer struct {
	Objects     []*Object
	tidyObjects sync.Map
	instances   sync.Map
}

// Object
type Object struct {
	Name      string
	New       func() (interface{}, error)
	Singleton bool
	refresher refresher
}

type refresher struct {
	mux    sync.Mutex
	status bool
}

func (t *refresher) On() {
	t.status = true
}

func (t *refresher) Off() {
	t.status = false
}

func (t *refresher) Status() bool {
	return t.status
}

func (t *refresher) Invoke(f func()) {
	if !t.Status() {
		return
	}
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.Status() {
		f()
		t.Off()
	}
}

func (t *Object) Refresh() error {
	if !t.Singleton {
		return fmt.Errorf("error: ")
	}
	t.refresher.On()
	return nil
}

// Provide
func (t *iContainer) Provide(objects ...*Object) error {
	t.tidyObjects = sync.Map{}
	for _, o := range objects {
		if _, ok := t.tidyObjects.Load(o.Name); ok {
			return fmt.Errorf("error: ")
		}
		t.tidyObjects.Store(o.Name, o)
	}
	return nil
}

func (t *iContainer) Object(name string) (*Object, error) {
	v, ok := t.tidyObjects.Load(name)
	if !ok {
		return nil, fmt.Errorf("error: ")
	}
	obj := v.(*Object)
	return obj, nil
}

// Find
func (t *iContainer) Find(name string, ptr interface{}) error {
	obj, err := t.Object(name)
	if err != nil {
		return err
	}
	ptrCopy := func(to, from interface{}) {
		reflect.ValueOf(to).Elem().Set(reflect.ValueOf(from))
	}
	if obj.Singleton {
		refresher := &obj.refresher
		if p, ok := t.instances.Load(name); ok && !refresher.Status() {
			ptrCopy(ptr, p)
			return nil
		}
		v, err := obj.New()
		if err != nil {
			return err
		}
		refresher.Invoke(func() {
			t.instances.Store(name, v)
		})
		p, _ := t.instances.LoadOrStore(name, v) // LoadOrStore 处理并发穿透
		ptrCopy(ptr, p)
		return nil
	} else {
		v, err := obj.New()
		if err != nil {
			return err
		}
		ptrCopy(ptr, v)
	}
	return nil
}
