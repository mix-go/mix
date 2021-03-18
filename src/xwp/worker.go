package xwp

import (
	"sync"
)

// Handler 处理器
type Handler func(data interface{})

// Worker 工作者接口
type Worker interface {
	Init(workerID int, workerPool chan JobQueue, wg *sync.WaitGroup, handler Handler)
	Run()
	Stop()
	Do(data interface{})
}

// WorkerTrait 工作者特征
type WorkerTrait struct {
	WorkerID   int
	workerPool chan JobQueue
	wg         *sync.WaitGroup
	handler    Handler
	jobChan    JobQueue
	quit       chan bool
}

// Init 初始化
func (t *WorkerTrait) Init(workerID int, workerPool chan JobQueue, wg *sync.WaitGroup, handler Handler) {
	t.WorkerID = workerID
	t.workerPool = workerPool
	t.wg = wg
	t.handler = handler
	t.jobChan = make(chan interface{})
	t.quit = make(chan bool)
}

// Run 执行
func (t *WorkerTrait) Run() {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.workerPool <- t.jobChan
		for {
			select {
			case data := <-t.jobChan:
				if data == nil {
					return
				}
				t.handler(data)
				t.workerPool <- t.jobChan
			case <-t.quit:
				close(t.jobChan)
			}
		}
	}()
}

// Stop 停止
func (t *WorkerTrait) Stop() {
	go func() {
		t.quit <- true
	}()
}
