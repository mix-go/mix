package catch

import (
    "errors"
    "github.com/mix-go/console"
    "reflect"
)

// 执行方法并捕获 panic
// 错误会统一交给 error 组件处理
func Call(f interface{}, args ...interface{}) {
    defer func() {
        if err := recover(); err != nil {
            console.LastError = err

            if console.App() == nil {
                panic(err)
            }
            console.App().Error.Handle(err)
        }
    }()

    switch reflect.TypeOf(f).Kind() {
    case reflect.Func:
        v := reflect.ValueOf(f)
        vargs := []reflect.Value{}
        for _, arg := range args {
            vargs = append(vargs, reflect.ValueOf(arg))
        }
        v.Call(vargs)
        break
    default:
        panic(errors.New("Invalid type: 'f' is not func"))
    }
}
