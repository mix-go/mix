package xwp

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// JobQueue 任务队列
type JobQueue chan interface{}

// Dispatcher 调度器
type Dispatcher struct {
	WorkerFunc     reflect.Value
	WorkerFuncArgs []reflect.Value
	MaxWorkers     int
	JobQueue       JobQueue
	workers        []Worker
	workerPool     chan JobQueue
	wg             *sync.WaitGroup
	quit           chan bool
}

// Run 执行
func (t *Dispatcher) Run() {
	for i := 0; i < t.MaxWorkers; i++ {
		w := t.WorkerFunc.Call(t.WorkerFuncArgs)[0].Interface().(Worker)
		w.Init(i, t.workerPool, t.wg, w.Do)
		w.Run()
		t.workers = append(t.workers, w)
	}
	t.dispatch()
	t.wait()
}

func (t *Dispatcher) dispatch() {
	go func() {
		for {
			select {
			case data := <-t.JobQueue:
				if data == nil {
					for _, w := range t.workers {
						w.Stop()
					}
					return
				}
				ch := <-t.workerPool
				ch <- data
			case <-t.quit:
				close(t.JobQueue)
			}
		}
	}()
}

// Stop 停止
func (t *Dispatcher) Stop() {
	go func() {
		t.quit <- true
	}()
}

func (t *Dispatcher) wait() {
	t.wg.Wait()
}

// NewDispatcher 创建调度器
func NewDispatcher(jobQueue chan interface{}, maxWorkers int, workerFunc interface{}, args ...interface{}) *Dispatcher {
	if reflect.TypeOf(workerFunc).Kind() != reflect.Func {
		panic(errors.New("WorkerFunc is not a Func type"))
	}
	value := reflect.ValueOf(workerFunc)
	valueArgs := []reflect.Value{}
	for _, arg := range args {
		valueArgs = append(valueArgs, reflect.ValueOf(arg))
	}
	func() {
		defer func() {
			if err := recover(); err != nil {
				panic(fmt.Errorf("WorkerFunc %s", err))
			}
		}()
		values := value.Call(valueArgs)
		if len(values) == 0 {
			panic(errors.New("WorkerFunc did not return"))
		}
		if _, ok := values[0].Interface().(Worker); !ok {
			panic(errors.New("WorkerFunc return type can only be workerPool.Worker"))
		}
	}()

	return &Dispatcher{
		WorkerFunc:     value,
		WorkerFuncArgs: valueArgs,
		MaxWorkers:     maxWorkers,
		JobQueue:       jobQueue,
		workerPool:     make(chan JobQueue, maxWorkers),
		wg:             &sync.WaitGroup{},
		quit:           make(chan bool),
	}
}
