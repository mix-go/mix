package workerpool

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

var (
    num = 0
)

type worker struct {
    WorkerTrait
}

func (t *worker) Handle(data interface{}) {
    fmt.Println(data)
    num++
}

func newWorker() Worker {
    return &worker{}
}

func Test(t *testing.T) {
    a := assert.New(t)

    jobQ := make(chan interface{}, 10)
    d := NewDispatcher(jobQ, 5)

    go func() {
        for i := 0; i < 10; i++ {
            jobQ <- i
        }
        fmt.Println("stop")
        d.Stop()
    }()

    d.Start(newWorker)
    d.Wait()
    a.Equal(num, 10)
}
