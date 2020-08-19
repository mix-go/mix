package workerpool

import (
    "github.com/stretchr/testify/assert"
    "sync/atomic"
    "testing"
)

var count int64

type worker struct {
    WorkerTrait
}

func (t *worker) Handle(data interface{}) {
    atomic.AddInt64(&count, 1)
}

func newWorker() Worker {
    return &worker{}
}

func TestOnce(t *testing.T) {
    a := assert.New(t)

    jobQ := make(chan interface{}, 200)
    d := NewDispatcher(jobQ, 15)

    go func() {
        for i := 0; i < 10000; i++ {
            jobQ <- i
        }
        d.Stop()
    }()

    d.Start(newWorker)
    d.Wait()

    a.Equal(count, int64(10000))
}
