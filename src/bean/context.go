package bean

import (
    "fmt"
    "sync"
)

// 创建实例
func NewApplicationContext(definitions []BeanDefinition) *ApplicationContext {
    context := &ApplicationContext{Definitions: definitions}
    context.Init()
    return context
}

// 应用上下文
type ApplicationContext struct {
    Definitions     []BeanDefinition
    tidyDefinitions sync.Map
    instances       sync.Map
}

// 初始化
func (t *ApplicationContext) Init() {
    t.tidyDefinitions = sync.Map{}
    for _, d := range t.Definitions {
        ptr := &BeanDefinition{
            Name:            d.Name,
            Reflect:         d.Reflect,
            Scope:           d.Scope,
            InitMethod:      d.InitMethod,
            ConstructorArgs: d.ConstructorArgs,
            Fields:          d.Fields,
            context:         t,
        }
        t.tidyDefinitions.Store(d.Name, ptr)
    }
}

// 获取定义
func (t *ApplicationContext) GetBeanDefinition(name string) *BeanDefinition {
    var (
        inf interface{}
        ok  bool
    )
    if inf, ok = t.tidyDefinitions.Load(name); !ok {
        panic(fmt.Sprintf("Bean not found: %s", name))
    }
    return inf.(*BeanDefinition)
}

// 获取实例
func (t *ApplicationContext) GetBean(name string, fields Fields, args ConstructorArgs) interface{} {
    bd := merge(t.GetBeanDefinition(name), fields, args)
    if bd.Scope == SINGLETON {
        if ins, ok := t.instances.Load(name); ok {
            return ins
        }
        val := bd.instance()
        ins, _ := t.instances.LoadOrStore(name, val) // LoadOrStore 处理并发穿透
        return ins
    }
    return bd.instance()
}

// 快速获取实例
func (c *ApplicationContext) Get(name string) interface{} {
    return c.GetBean(name, Fields{}, ConstructorArgs{})
}

// 判断组件是否存在
func (c *ApplicationContext) Has(name string) (ok bool) {
    ok = true
    defer func() {
        if err := recover(); err != nil {
            ok = false
        }
    }()
    c.GetBeanDefinition(name)
    return ok
}

// 合并
// args | fields 内的字段会替换之前定义的值
// args 内的 nil 值将会忽略，不会替换处理
func merge(bd *BeanDefinition, fields Fields, args ConstructorArgs) *BeanDefinition {
    f := len(fields) > 0
    a := len(args) > 0
    if !(f || a) {
        return bd
    }
    nbd := &BeanDefinition{
        Name:            bd.Name,
        Scope:           bd.Scope,
        Reflect:         bd.Reflect,
        InitMethod:      bd.InitMethod,
        ConstructorArgs: nil,
        Fields:          nil,
        context:         bd.context,
    }
    if f {
        // 合并替换字段
        tmp := Fields{}
        for k, v := range bd.Fields {
            tmp[k] = v
        }
        for k, v := range fields {
            tmp[k] = v
        }
        nbd.Fields = tmp
    }
    if a {
        // 合并替换参数，nil 忽略
        tmp := ConstructorArgs{}
        tmp = append(tmp, bd.ConstructorArgs...)
        for k, v := range args {
            if v == nil {
                continue
            }
            ok := false
            for sk, _ := range bd.ConstructorArgs {
                if sk == k {
                    ok = true
                }
            }
            if ok {
                tmp[k] = v
            } else {
                tmp = append(tmp, v)
            }
        }
        nbd.ConstructorArgs = tmp
    }
    return nbd
}
