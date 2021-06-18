package vwatch

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"time"
)

type Watcher struct {
	Ptr    interface{}
	Last   map[uintptr]string
	Nodes  map[uintptr]func()
	Ticker *time.Ticker
}

func NewWatcher(ptr interface{}) (*Watcher, error) {
	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("node are not pointer type")
	}
	return &Watcher{
		Ptr:   ptr,
		Last:  make(map[uintptr]string),
		Nodes: make(map[uintptr]func()),
	}, nil
}

func (t *Watcher) Watch(node interface{}, f func()) error {
	val := reflect.ValueOf(node)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("node are not pointer type")
	}
	t.Nodes[val.Pointer()] = f
	return nil
}

func (t *Watcher) Run(interval time.Duration) {
	if t.Ticker != nil {
		return
	}
	t.Ticker = time.NewTicker(interval)
	go func() {
		for {
			<-t.Ticker.C
			t.do()
		}
	}()
}

func (t *Watcher) Stop() {
	if t.Ticker == nil {
		return
	}
	t.Ticker.Stop()
}

func (t *Watcher) trigger(ptr uintptr) {
	if f, ok := t.Nodes[ptr]; ok {
		f()
	}
}

func (t *Watcher) do() {
	data := make(map[uintptr]string)
	extract(reflect.ValueOf(t.Ptr), data)
	fmt.Println(data)
	for k, v := range data {
		if v != t.Last[k] {
			t.trigger(k)
		}
	}
}

func extract(val reflect.Value, data map[uintptr]string) {
	switch val.Kind() {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.String,
		reflect.Map,
		reflect.Slice,
		reflect.Array:
		if val.CanAddr() {
			data[val.Pointer()] = hash(val.Interface())
		}
		break
	case reflect.Ptr:
		elem := val.Elem()
		if !elem.CanAddr() {
			return
		}
		data[elem.Addr().Pointer()] = hash(elem.Interface())
		extract(reflect.ValueOf(elem.Interface()), data)
		break
	case reflect.Struct:
		if val.CanAddr() {
			data[val.Pointer()] = hash(val.Interface())
		}
		for i := 0; i < val.NumField(); i++ {
			if !val.Field(i).CanInterface() {
				continue
			}
			extract(reflect.ValueOf(val.Field(i).Interface()), data)
		}
		break
	}
}

func hash(i interface{}) string {
	s := fmt.Sprintf("%#v", i)
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
