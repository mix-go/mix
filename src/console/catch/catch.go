package catch

import (
	"errors"
	"github.com/mix-go/console"
	"reflect"
)

// Use example: `go catch.Call`
// 执行方法
// 捕获 panic，错误会统一交给 error 组件处理
func Call(fn interface{}, args ...interface{}) {
	if fn == nil {
		panic(errors.New("Invalid type: 'fn' is not func"))
	}
	switch reflect.TypeOf(fn).Kind() {
	case reflect.Func:
		defer func() {
			if err := recover(); err != nil {
				Error(err)
			}
		}()
		v := reflect.ValueOf(fn)
		vargs := []reflect.Value{}
		for _, arg := range args {
			vargs = append(vargs, reflect.ValueOf(arg))
		}
		v.Call(vargs)
		break
	default:
		panic(errors.New("Invalid type: 'fn' is not func"))
	}
}

// 捕获错误
func Error(err interface{}) {
	if console.App == nil {
		panic(err)
	}
	console.App.Error.Handle(err)
}
