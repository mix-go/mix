package xwp

import (
	"fmt"
	"sync"
	"time"
)

// JobQueue 任务队列
type JobQueue chan interface{}

// RunI
type RunI interface {
	Do(data interface{})
}

// WorkerPool 调度器
type WorkerPool struct {
	JobQueue JobQueue

	MaxWorkers int
	// default == MaxWorkers
	InitWorkers int
	// default == MaxWorkers
	MaxIdleWorkers int

	WorkerRun  func(data interface{})
	WorkerRunI RunI

	workers    *sync.Map
	workerPool chan JobQueue
	wg         *sync.WaitGroup
	quit       chan bool
}

// Run 执行
func (t *WorkerPool) Run() {
	t.Start()
	t.Wait()
}

// init
func (t *WorkerPool) init() {
	t.workerPool = make(chan JobQueue, t.MaxIdleWorkers)
	t.wg = &sync.WaitGroup{}
	t.quit = make(chan bool)

	if t.WorkerRun == nil && t.WorkerRunI == nil {
		panic(fmt.Errorf("xwp.WorkerPool WorkerRun & WorkerRunI field is empty"))
	}

	if t.InitWorkers == 0 {
		t.InitWorkers = t.MaxWorkers
	}
	if t.MaxIdleWorkers == 0 {
		t.MaxIdleWorkers = t.MaxWorkers
	}
}

func (t *WorkerPool) addWorker() {
	w := NewWorker(t)
	w.Run()
}

// Start 启动
func (t *WorkerPool) Start() {
	t.init()

	for i := 0; i < t.InitWorkers; i++ {
		NewWorker(t).Run()
	}

	go func() {
		timer := time.NewTimer(time.Millisecond)
		timer.Stop()
		for {
			select {
			case data := <-t.JobQueue:
				if data == nil {
					t.workers.Range(func(key, value interface{}) bool {
						w := value.(Worker)
						w.Stop()
						return true
					})
					return
				}
				for {
					select {
					case ch := <-t.workerPool:
						ch <- data
					default:
						if len(t.workerPool) < t.MaxWorkers {
							NewWorker(t).Run()
						} else {
							timer.Reset(10 * time.Millisecond)
							select {
							case ch := <-t.workerPool:
								timer.Stop()
								ch <- data
							case <-timer.C:
							}
						}
					}
				}
			case <-t.quit:
				close(t.JobQueue)
			}
		}
	}()
}

// Stop 停止
func (t *WorkerPool) Stop() {
	go func() {
		t.quit <- true
	}()
}

// Wait 等待执行完成
func (t *WorkerPool) Wait() {
	t.wg.Wait()
}

type PoolStat struct {
	Active int `json:"active"`
	Idle   int `json:"idle"`
	Total  int `json:"total"`
}

// Wait 等待执行完成
func (t *WorkerPool) Stat() *PoolStat {
	return &PoolStat{
		Active: 0,
		Idle:   0,
		Total:  0,
	}
}
