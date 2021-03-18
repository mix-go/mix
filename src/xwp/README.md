> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XWP

通用的工作池

A common worker pool

## Installation

- 安装

```
go get -u github.com/mix-go/xwp
```

## Usage

先创建一个 Worker 结构体

~~~
type FooWorker struct {
    xwp.WorkerTrait
}

func (t *FooWorker) Do(data interface{}) {
    // do something
}

func NewFooWorker() xwp.Worker {
    return &FooWorker{}
}
~~~

调度任务

~~~
jobQueue := make(chan interface{}, 200)
d := xwp.NewDispatcher(jobQueue, 15, NewFooWorker)

go func() {
    // 投放任务
    for i := 0; i < 10000; i++ {
        jobQueue <- i
    }

    // 投放完停止调度
    d.Stop()
}()

d.Run() // 阻塞代码，直到任务全部执行完成并且全部 Worker 停止
~~~

异常处理：`Do` 方法中执行的代码，可能会出现 `panic` 异常，我们可以通过 `recover` 获取异常信息记录到日志或者执行其他处理

~~~
func (t *FooWorker) Do(data interface{}) {
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
