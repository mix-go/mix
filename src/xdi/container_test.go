package xdi

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sync"
	"testing"
	"time"
)

type foo struct {
	Bar    string
	Client *http.Client
}

func TestPopulate(t *testing.T) {
	a := assert.New(t)

	c := New()
	objs := []*Object{
		{
			Name: "client",
			New: func() (i interface{}, e error) {
				timeout := time.Second * 10
				return &http.Client{
					Timeout: timeout,
				}, nil
			},
		},
		{
			Name: "foo",
			New: func() (i interface{}, e error) {
				var hc *http.Client
				if err := c.Populate("client", &hc); err != nil {
					panic(err)
				}
				return &foo{
					Bar:    "",
					Client: hc,
				}, nil
			},
		},
	}
	_ = c.Provide(objs...)

	var f *foo
	_ = c.Populate("foo", &f)
	text := fmt.Sprintf("%#v \n", f.Client)

	a.Contains(text, "Timeout:10000000000")
}

func TestSingletonConcurrency(t *testing.T) {
	a := assert.New(t)

	c := New()
	objs := []*Object{
		{
			Name: "foo",
			New: func() (i interface{}, e error) {
				return &foo{
					Bar: "",
				}, nil
			},
			Singleton: true,
		},
	}
	_ = c.Provide(objs...)

	for i := 0; i < 1000; i++ {
		var mp sync.Map
		wg := &sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, i int) {
				defer wg.Done()

				var f *foo
				_ = c.Populate("foo", &f)

				mp.Store(i, f)
			}(wg, i)
		}
		wg.Wait()

		ptrs := []interface{}{}
		mp.Range(func(key, value interface{}) bool {
			ptrs = append(ptrs, fmt.Sprintf("%p", value))
			return true
		})

		//fmt.Println(ptrs...)

		a.Equal(ptrs[0], ptrs[1], ptrs[2:]...)
	}
}

func TestSingletonRefreshConcurrency(t *testing.T) {
	a := assert.New(t)

	c := New()
	objs := []*Object{
		{
			Name: "foo",
			New: func() (i interface{}, e error) {
				return &foo{
					Bar: "",
				}, nil
			},
			Singleton: true,
		},
	}
	_ = c.Provide(objs...)

	for i := 0; i < 1000; i++ {
		var mp sync.Map
		wg := &sync.WaitGroup{}

		var f *foo
		_ = c.Populate("foo", &f)
		o, _ := c.Object("foo")
		_ = o.Refresh()

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, i int) {
				defer wg.Done()

				var f *foo
				_ = c.Populate("foo", &f)

				mp.Store(i, f)
			}(wg, i)
		}
		wg.Wait()

		ptrs := []interface{}{}
		mp.Range(func(key, value interface{}) bool {
			ptrs = append(ptrs, fmt.Sprintf("%p", value))
			return true
		})

		//fmt.Println(ptrs...)

		a.Equal(ptrs[0], ptrs[1], ptrs[2:]...)
	}
}
