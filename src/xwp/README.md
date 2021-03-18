> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix Worker Pool

通用的工作池

A common worker pool

> 该库还有 php 版本：https://github.com/mix-php/worker-pool

## Installation

- 安装

```
go get -u github.com/mix-go/xwp
```

## Usage

先创建一个 Worker 结构体

~~~
type FooWorker struct {
    workerpool.WorkerTrait
}

func (t *FooWorker) Do(data interface{}) {
    // do something
}

func NewFooWorker() workerpool.Worker {
    return &FooWorker{}
}
~~~

调度任务

~~~
jobQueue := make(chan interface{}, 200)
d := workerpool.NewDispatcher(jobQueue, 15, NewFooWorker)

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
