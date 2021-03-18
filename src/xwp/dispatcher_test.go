package xwp

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	count     int64
	workerMap sync.Map
)

type worker struct {
	WorkerTrait
}

func (t *worker) Do(data interface{}) {
	atomic.AddInt64(&count, 1)
	workerMap.Store(t.WorkerID, 1)
}

func newWorker() Worker {
	return &worker{}
}

func TestOnceRun(t *testing.T) {
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
	workerNum := 0
	workerMap.Range(func(key, value interface{}) bool {
		workerNum++
		return true
	})
	a.Equal(workerNum, 15)
}

func TestStop(t *testing.T) {
	a := assert.New(t)
	jobQueue := make(chan interface{}, 200)
	d := NewDispatcher(jobQueue, 15, newWorker)
	go func() {
		defer func() {
			err := recover()
			a.EqualError(err.(error), "send on closed channel")
		}()
		for {
			jobQueue <- struct {
			}{}
		}
	}()
	go func() {
		time.Sleep(100 * time.Millisecond)
		d.Stop()
	}()
	d.Run()
}
