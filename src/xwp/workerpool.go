package xwp

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
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

	// default == runtime.NumCPU()
	MaxWorkers int
	// default == MaxWorkers
	InitWorkers int
	// default == InitWorkers
	MaxIdleWorkers int

	RunF func(data interface{})
	RunI RunI

	workers         *sync.Map
	workerCount     int64
	workerQueuePool chan JobQueue
	wg              *sync.WaitGroup
	quit            chan bool
}

// Run 执行
func (t *WorkerPool) Run() {
	t.Start()
	t.Wait()
}

// init
func (t *WorkerPool) init() {
	if t.MaxWorkers == 0 {
		t.MaxWorkers = runtime.NumCPU()
	}
	if t.InitWorkers == 0 {
		t.InitWorkers = t.MaxWorkers
	}
	if t.MaxIdleWorkers == 0 {
		t.MaxIdleWorkers = t.InitWorkers
	}

	t.workers = &sync.Map{}
	t.workerQueuePool = make(chan JobQueue, t.MaxIdleWorkers)
	t.wg = &sync.WaitGroup{}
	t.quit = make(chan bool)

	if t.RunF == nil && t.RunI == nil {
		panic(fmt.Errorf("xwp.WorkerPool RunF & RunI field is empty"))
	}
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
						w := value.(*Worker)
						w.Stop()
						return true
					})
					return
				}
				func() {
					for {
						select {
						case ch := <-t.workerQueuePool:
							ch <- data
						default:
							if atomic.LoadInt64(&t.workerCount) < int64(t.MaxWorkers) {
								NewWorker(t).Run()
								continue
							} else {
								// 设定时间的监听
								timer.Reset(10 * time.Millisecond)
								select {
								case ch := <-t.workerQueuePool:
									timer.Stop()
									ch <- data
								case <-timer.C:
									continue
								}
							}
						}
						return
					}
				}()
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

type Statistic struct {
	Active int `json:"active"`
	Idle   int `json:"idle"`
	Total  int `json:"total"`
}

// Stats 统计
func (t *WorkerPool) Stats() *Statistic {
	total := int(t.workerCount)
	idle := len(t.workerQueuePool)
	return &Statistic{
		Active: total - idle,
		Idle:   idle,
		Total:  total,
	}
}
