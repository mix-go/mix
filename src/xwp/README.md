> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XWP

通用的工作池

A common worker pool

## Installation

```
go get github.com/mix-go/xwp
```

## Usage

先创建一个结构体用来处理任务，使用类型断言转换任务数据类型，例如：`i := data.(int)` 

~~~go
type Foo struct {
}

func (t *Foo) Do(data interface{}) {
    // do something
}
~~~

调度任务

- 也可以使用 `RunF` 采用闭包来处理任务

~~~go
jobQueue := make(chan interface{}, 200)

p := &xwp.WorkerPool{
    JobQueue:       jobQueue,
    MaxWorkers:     1000,
    InitWorkers:    100,
    MaxIdleWorkers: 100,
    RunI:           &Foo{},
}

go func() {
    // 投放任务
    for i := 0; i < 10000; i++ {
        jobQueue <- i
    }

    // 投放完停止调度
    p.Stop()
}()

p.Run() // 阻塞等待
~~~

异常处理：`Do` 方法中执行的代码，可能会出现 `panic` 异常，我们可以通过 `recover` 获取异常信息记录到日志或者执行其他处理

~~~go
func (t *Foo) Do(data interface{}) {
    defer func() {
        if err := recover(); err != nil {
            // handle error
        }
    }()
    // do something
}
~~~

## License

Apache License Version 2.0, http://www.apache.org/licenses/
