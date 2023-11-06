package xdi

import (
	"errors"
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
		{
			Name: "bar",
			New: func() (i interface{}, e error) {
				return nil, errors.New("error")
			},
		},
	}
	_ = c.Provide(objs...)

	// 测试单例
	var f1 *foo
	_ = c.Populate("foo", &f1)
	var f2 *foo
	_ = c.Populate("foo", &f2)
	a.Equal(fmt.Sprintf("%p", f1), fmt.Sprintf("%p", f2))

	// 错误使用测试
	var f3 foo
	err := c.Populate("foo", f3) // 非指针
	a.Contains(err.Error(), "can only be pointer type")
	var f4 foo
	err = c.Populate("foo", &f4) // New函数返回的指针，但是f3为引用 [panic: reflect.Set: value of type *xdi.foo is not assignable to type xdi.foo]
	a.Contains(err.Error(), "is not assignable to")

	// 测试嵌套依赖
	var f *foo
	_ = c.Populate("foo", &f)
	text := fmt.Sprintf("%#v \n", f.Client)
	a.Contains(text, "Timeout:10000000000")

	// 测试多次失败场景
	var i interface{}
	err = c.Populate("bar", &i)
	a.Equal(err, errors.New("error"))
	err = c.Populate("bar", &i)
	a.Equal(err, errors.New("error"))
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

func TestSingletonConcurrencyError(t *testing.T) {
	a := assert.New(t)

	c := New()
	objs := []*Object{
		{
			Name: "foo",
			New: func() (i interface{}, e error) {
				return nil, errors.New("new error")
			},
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

func TestSingletonConcurrencyRefresh(t *testing.T) {
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
