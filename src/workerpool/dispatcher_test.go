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

func (t *worker) Do(data interface{}) {
    atomic.AddInt64(&count, 1)
}

func (t *worker) Error(err interface{}) {
    // handle err
}

func newWorker() Worker {
    return &worker{}
}

func TestOnce(t *testing.T) {
    a := assert.New(t)

    jobQueue := make(chan interface{}, 200)
    d := NewDispatcher(jobQueue, 15, newWorker)

    go func() {
        for i := 0; i < 10000; i++ {
            jobQueue <- i
        }
        d.Stop()
    }()

    d.Run()

    a.Equal(count, int64(10000))
}
