## Mix Worker Pool

通用的工作池类库

A common workpool class library

> 该库还有 php 版本：https://github.com/mix-php/worker-pool

## Usage

- 安装

```
go get -u github.com/mix-go/workerpool
```

- 使用

先创建一个 Worker 对象

~~~
type FooWorker struct {
    workerpool.WorkerTrait
}

func (t *FooWorker) Handle(data interface{}) {
    // do something
}

func NewFooWorker() workerpool.Worker {
    return &FooWorker{}
}
~~~

调度任务

~~~
jobQueue := make(chan interface{}, 200)
d := workerpool.NewDispatcher(NewFooWorker, 15, jobQueue)

go func() {
    // 投放任务
    for i := 0; i < 10000; i++ {
        jobQueue <- i
    }

    // 投放完停止调度
    d.Stop()
}()

d.Start()
d.Wait() // 等待任务全部执行完成并停止全部 Worker
~~~

## License

Apache License Version 2.0, http://www.apache.org/licenses/
