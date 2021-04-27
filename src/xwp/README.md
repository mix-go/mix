> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XWP

通用动态工作池、协程池

A common worker pool

## Installation

```
go get github.com/mix-go/xwp
```

## 单次调度

> 适合处理数据计算、转换等场景

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
- 如果不想阻塞执行，可以使用 `p.Start()` 启动

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

## 常驻调度

> 适合处理 MQ 队列的异步消费

以 Redis 作为 MQ 为例：

~~~go
jobQueue := make(chan interface{}, 200)
p := &xwp.WorkerPool{
    JobQueue:       jobQueue,
    MaxWorkers:     1000,
    InitWorkers:    100,
    MaxIdleWorkers: 100,
    RunI:           &Foo{},
}

ch := make(chan os.Signal)
signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
go func() {
    <-ch
    p.Stop()
}()

go func() {
    for {
        res, err := redis.BRPop(context.Background(), 3*time.Second, "foo").Result()
        if err != nil {
            if strings.Contains(err.Error(), "redis: nil") {
                continue
            }
            fmt.Println(fmt.Sprintf("Redis Error: %s", err))
            p.Stop()
            return
        }
        // brPop命令最后一个键才是值
        jobQueue <- res[1]
    }
}()

p.Run() // 阻塞等待
~~~

## 异常处理

`Do` 方法中执行的代码，可能会出现 `panic` 异常，我们可以通过 `recover` 获取异常信息记录到日志或者执行其他处理

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

## 执行状态

`Stat()` 可以返回 `Workers` 实时执行状态，通常可以使用一个定时器，定时打印或者告警处理

```go
go func() {
    ticker := time.NewTicker(1000 * time.Millisecond)
    for {
        <-ticker.C
        log.Printf("%+v", p.Stat()) // 2021/04/26 14:32:53 &{Active:5 Idle:95 Total:100}
    }
}()
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
