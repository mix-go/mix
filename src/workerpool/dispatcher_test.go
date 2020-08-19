package workerpool

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

var (
    count = 0
)

type worker struct {
    WorkerTrait
}

func (t *worker) Handle(data interface{}) {
    count++
}

func newWorker() Worker {
    return &worker{}
}

func TestOnce(t *testing.T) {
    a := assert.New(t)

    jobQ := make(chan interface{}, 200)
    d := NewDispatcher(jobQ, 5)

    go func() {
        for i := 0; i < 10000; i++ {
            jobQ <- i
        }
        d.Stop()
    }()

    d.Start(newWorker)
    d.Wait()

    fmt.Println("Equal")
    a.Equal(count, 10000)
}
