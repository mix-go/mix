package bean

import (
    "fmt"
    "reflect"
)

const (
    PROTOTYPE = "prototype"
    SINGLETON = "singleton"
)

// 创建反射
func NewReflect(i interface{}) func() reflect.Value {
    switch reflect.TypeOf(i).Kind() {
    case reflect.Func:
        return func() reflect.Value {
            return reflect.ValueOf(i)
        }
    case reflect.Struct:
        return func() reflect.Value {
            return reflect.New(reflect.TypeOf(i))
        }
    default:
        panic("Invalid type, use only func and struct")
    }
    return nil
}

// 构造器参数
type ConstructorArgs []interface{}

// 字段
type Fields map[string]interface{}

// 引用
type Reference struct {
    Name string
}

// 创建引用
func NewReference(name string) Reference {
    return Reference{Name: name}
}

type ReturnError struct {
    error
}

func NewReturnError(err error) *ReturnError {
    return &ReturnError{err}
}

// 定义
type BeanDefinition struct {
    Name            string
    Reflect         func() reflect.Value
    Scope           string
    InitMethod      string
    ConstructorArgs ConstructorArgs
    Fields          Fields
    context         *ApplicationContext
}

// 刷新
func (t *BeanDefinition) Refresh() {
    t.context.instances.Store(t.Name, t.instance())
}

// 实例化
func (t *BeanDefinition) instance() interface{} {
    v := t.Reflect()

    // 构造器注入
    if v.Kind() == reflect.Func {
        in := []reflect.Value{}
        for _, a := range t.ConstructorArgs {
            // 引用字段处理
            if _, ok := a.(Reference); ok {
                a = t.context.GetBean(a.(Reference).Name, Fields{}, ConstructorArgs{})
            }
            in = append(in, reflect.ValueOf(a))
        }
        func() {
            defer func() {
                if err := recover(); err != nil {
                    if err, ok := err.(*ReturnError); ok {
                        panic(err)
                    }
                    panic(fmt.Sprintf("Bean name '%s' reflect %s construct failed, %s", t.Name, v.Type().String(), err))
                }
            }()
            res := v.Call(in)
            if len(res) >= 2 {
                if err, ok := res[1].Interface().(error); ok {
                    panic(NewReturnError(err))
                }
            }
            v = res[0]
        }()
        if v.Kind() == reflect.Struct {
            panic(fmt.Sprintf("Bean name '%s' reflect %s return value is not a pointer type", t.Name, v.Type().String()))
        }
    }

    // 字段注入
    for k, p := range t.Fields {
        // 字段检测
        if !v.Elem().FieldByName(k).CanSet() {
            panic(fmt.Sprintf("Bean name '%s' type %s field %s cannot be found or cannot be set", t.Name, reflect.TypeOf(v.Interface()), k))
        }
        // 引用字段处理
        if _, ok := p.(Reference); ok {
            p = t.context.GetBean(p.(Reference).Name, Fields{}, ConstructorArgs{})
        }
        // 类型检测
        if v.Elem().FieldByName(k).Type().String() != reflect.TypeOf(p).String() {
            panic(fmt.Sprintf("Bean name '%s' type %s field %s value of type %s is not assignable to type %s",
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
        if !m.IsValid() {
            panic(fmt.Sprintf("Bean name '%s' type %s init method %s not found",
                t.Name,
                reflect.TypeOf(v.Interface()),
                t.InitMethod,
            ))
        }
        m.Call([]reflect.Value{})
    }

    return v.Interface()
}
