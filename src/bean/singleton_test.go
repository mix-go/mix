package bean

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "sync"
    "testing"
)

var definitions1 = []BeanDefinition{
    {
        Name:    "foo",
        Scope:   SINGLETON,
        Reflect: NewReflect(foo{}),
    },
}

func TestSingleton(t *testing.T) {
    a := assert.New(t)

    context := NewApplicationContext(definitions1)
    f1 := context.Get("foo").(*foo)
    f2 := context.Get("foo").(*foo)
    f3 := context.Get("foo").(*foo)

    a.Equal(fmt.Sprintf("%p", f1), fmt.Sprintf("%p", f2), fmt.Sprintf("%p", f3))
}

func TestRefresh(t *testing.T) {
    a := assert.New(t)

    context := NewApplicationContext(definitions1)
    f1 := context.Get("foo").(*foo)
    d := context.GetBeanDefinition("foo")
    d.Refresh()
    f2 := context.Get("foo").(*foo)

    a.NotEqual(fmt.Sprintf("%p", f1), fmt.Sprintf("%p", f2))
}

func TestConcurrencyGet(t *testing.T) {
    for i := 0; i < 1000; i++ {
        testConcurrencyGet(t)
    }
}

func testConcurrencyGet(t *testing.T) {
    a := assert.New(t)

    context := NewApplicationContext(definitions1)
    var mp sync.Map
    wg := &sync.WaitGroup{}
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(wg *sync.WaitGroup, i int) {
            defer wg.Done()
            f1 := context.Get("foo").(*foo)
            mp.Store(i, f1)
        }(wg, i)
    }
    wg.Wait()

    ptrs := []interface{}{}
    mp.Range(func(key, value interface{}) bool {
        ptrs = append(ptrs, fmt.Sprintf("%p", value))
        return true
    })

    a.Equal(ptrs[0], ptrs[1], ptrs[2:]...)
}
