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

// New
func NewReference(name string) Reference {
    return Reference{Name: name}
}

// 实例化
func (t *Definition) Instance() interface{} {
    s := t.Reflect()
    v := s
    // 构造器注入
    if s.Kind() == reflect.Func {
        in := []reflect.Value{}
        for _, a := range t.ConstructorArgs {
            in = append(in, reflect.ValueOf(a))
        }
        v = s.Call(in)[0]
        if v.Kind() == reflect.Struct {
            panic(fmt.Sprintf("Bean name %s reflect %s return value is not a pointer type", t.Name, s.Type().String()))
        }
    }
    // 字段注入
    for k, p := range t.Fields {
        // 字段检测
        if !v.Elem().FieldByName(k).CanSet() {
            panic(fmt.Sprintf("Bean name %s type %s field %s cannot be found or cannot be set", t.Name, reflect.TypeOf(v.Interface()), k))
        }
        // 引用字段处理
        if _, ok := p.(Reference); ok {
            p = t.Context.GetBean(p.(Reference).Name, Fields{}, ConstructorArgs{})
        }
        // 类型检测
        if v.Elem().FieldByName(k).Type().String() != reflect.TypeOf(p).String() {
            panic(fmt.Sprintf("Bean name %s type %s field %s value of type %s is not assignable to type %s",
                t.Name,
                reflect.TypeOf(v.Interface()),
                k,
                reflect.TypeOf(p).String(), v.Elem().FieldByName(k).Type().String()),
            )
        }
        // 常规字段处理
        v.Elem().FieldByName(k).Set(reflect.ValueOf(p))
    }
    // 执行初始化方法
    if t.InitMethod != "" {
        m := v.MethodByName(t.InitMethod)
        m.Call([]reflect.Value{})
    }
    return v.Interface()
}
