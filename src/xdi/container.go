package xdi

import (
	"fmt"
	"reflect"
	"sync"
)

var DefaultContainer *Container

func init() {
	DefaultContainer = New()
}

func New() *Container {
	return &Container{}
}

func Provide(objects ...*Object) error {
	return DefaultContainer.Provide(objects...)
}

func Populate(name string, ptr interface{}) error {
	return DefaultContainer.Populate(name, ptr)
}

type Container struct {
	Objects     []*Object
	tidyObjects sync.Map
	instances   sync.Map
}

func (t *Container) Provide(objects ...*Object) error {
	for _, o := range objects {
		if _, ok := t.tidyObjects.Load(o.Name); ok {
			return fmt.Errorf("error: object '%s' existing", o.Name)
		}
		t.tidyObjects.Store(o.Name, o)
	}
	return nil
}

func (t *Container) Object(name string) (*Object, error) {
	v, ok := t.tidyObjects.Load(name)
	if !ok {
		return nil, fmt.Errorf("error: object '%s' not found", name)
	}
	obj := v.(*Object)
	return obj, nil
}

func (t *Container) Populate(name string, ptr interface{}) error {
	obj, err := t.Object(name)
	if err != nil {
		return err
	}
	ptrCopy := func(to, from interface{}) {
		reflect.ValueOf(to).Elem().Set(reflect.ValueOf(from))
	}
	if !obj.NewEverytime {
		refresher := &obj.refresher
		if p, ok := t.instances.Load(name); ok && !refresher.status() {
			ptrCopy(ptr, p)
			return nil
		}
		// 处理并发穿透
		obj.mutex.Lock()
		defer obj.mutex.Unlock()
		p, ok := t.instances.Load(name)
		if ok {
			ptrCopy(ptr, p)
			return nil
		}
		v, err := obj.New()
		if err != nil {
			return err
		}
		t.instances.Store(name, v)
		refresher.off()
		ptrCopy(ptr, v)
		return nil
	} else {
		v, err := obj.New()
		if err != nil {
			return err
		}
		ptrCopy(ptr, v)
		return nil
	}
}
