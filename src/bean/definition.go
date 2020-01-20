package bean

import (
    "fmt"
    "reflect"
)

const (
    PROTOTYPE = "prototype";
    SINGLETON = "singleton";
)

// 定义
type Definition struct {
    Name            string
    Reflect         func() reflect.Value
    Scope           string
    InitMethod      string
    ConstructorArgs ConstructorArgs
    Fields          Fields
    Context         *ApplicationContext
}

// 构造器参数
type ConstructorArgs []interface{}

// 字段
type Fields map[string]interface{}

// 引用
type Reference struct {
    Name string
}

// 实例化
func (c *Definition) Instance() interface{} {
    src := c.Reflect()
    v := src
    // 构造器注入
    if src.Kind() == reflect.Func {
        in := []reflect.Value{}
        for _, a := range c.ConstructorArgs {
            in = append(in, reflect.ValueOf(a))
        }
        v = src.Call(in)[0]
        if v.Kind() == reflect.Struct {
            panic(fmt.Sprintf("Bean name '%s' reflect %s return value is not a pointer type", c.Name, src.Type().String()))
        }
    }
    // 字段注入
    for k, p := range c.Fields {
        // 字段检测
        if !v.Elem().FieldByName(k).CanSet() {
            panic(fmt.Sprintf("Bean name '%s' type %s field %s cannot be found or cannot be set", c.Name, reflect.TypeOf(v.Interface()), k))
        }
        // 引用字段处理
        if _, ok := p.(Reference); ok {
            p = c.Context.GetBean(p.(Reference).Name, Fields{}, ConstructorArgs{})
        }
        // 类型检测
        if v.Elem().FieldByName(k).Type().String() != reflect.TypeOf(p).String() {
            panic(fmt.Sprintf("Bean name '%s' type %s field %s value of type %s is not assignable to type %s",
                c.Name,
                reflect.TypeOf(v.Interface()),
                k,
                reflect.TypeOf(p).String(), v.Elem().FieldByName(k).Type().String()),
            )
        }
        // 常规字段处理
        v.Elem().FieldByName(k).Set(reflect.ValueOf(p))
    }
    // 执行初始化方法
    if c.InitMethod != "" {
        m := v.MethodByName(c.InitMethod)
        m.Call([]reflect.Value{})
    }
    return v.Interface()
}
