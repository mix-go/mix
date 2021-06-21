package varwatch

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"time"
)

type Watcher struct {
	Ptr      interface{}
	Interval time.Duration
	Last     map[string]string
	Nodes    map[string]func()
	Ticker   *time.Ticker
}

func NewWatcher(ptr interface{}, interval time.Duration) (*Watcher, error) {
	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("node are not pointer type")
	}
	return &Watcher{
		Ptr:      ptr,
		Interval: interval,
		Last:     make(map[string]string),
		Nodes:    make(map[string]func()),
	}, nil
}

func (t *Watcher) Watch(tag string, f func()) error {
	t.Nodes[tag] = f
	return nil
}

func (t *Watcher) Run() error {
	if t.Ticker != nil {
		return fmt.Errorf("cannot repeat execution")
	}
	t.Ticker = time.NewTicker(t.Interval)
	go func() {
		for {
			<-t.Ticker.C
			t.do()
		}
	}()
	return nil
}

func (t *Watcher) Stop() {
	if t.Ticker == nil {
		return
	}
	t.Ticker.Stop()
}

func (t *Watcher) trigger(tag string) {
	if f, ok := t.Nodes[tag]; ok {
		f()
	}
}

func (t *Watcher) do() {
	data := make(map[string]string)
	extract(reflect.ValueOf(t.Ptr), data)
	for k, v := range data {
		if _, ok := t.Last[k]; !ok {
			continue
		}
		if v != t.Last[k] {
			t.trigger(k)
		}
	}
	t.Last = data
}

func extract(val reflect.Value, data map[string]string) {
	switch val.Kind() {
	case reflect.Ptr:
		elem := val.Elem()
		if !elem.CanAddr() {
			return
		}
		extract(reflect.ValueOf(elem.Interface()), data)
		break
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if !val.Field(i).CanInterface() {
				continue
			}
			tag := val.Type().Field(i).Tag.Get("varwatch")
			if tag != "" && tag != "-" && tag != "_" {
				data[tag] = hash(val.Field(i).Interface())
			}
			extract(reflect.ValueOf(val.Field(i).Interface()), data)
		}
		break
	case reflect.Map:
		iter := val.MapRange()
		for iter.Next() {
			extract(reflect.ValueOf(iter.Value().Interface()), data)
		}
		break
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			extract(reflect.ValueOf(val.Index(i).Interface()), data)
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
