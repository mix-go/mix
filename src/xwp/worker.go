package xwp

import (
	"fmt"
	"sync/atomic"
)

// Handler 处理器
type Handler func(data interface{})

// Worker 工作者
type Worker struct {
	pool    *WorkerPool
	handler Handler
	jobChan JobQueue
	quit    chan bool
}

// NewWorker
func NewWorker(p *WorkerPool) *Worker {
	return &Worker{
		pool: p,
		handler: func(data interface{}) {
			if p.RunF != nil {
				p.RunF(data)
			} else if p.RunI != nil {
				i := p.RunI
				i.Do(data)
			}
		},
		jobChan: make(chan interface{}),
		quit:    make(chan bool),
	}
}

// Run 执行
func (t *Worker) Run() {
	// 先注册
	select {
	case t.pool.workerQueuePool <- t.jobChan:
		atomic.AddInt64(&t.pool.workerCount, 1)
		t.pool.workers.Store(fmt.Sprintf("%p", t), t)
	default:
		return
	}

	t.pool.wg.Add(1)
	go func() {
		defer func() {
			atomic.AddInt64(&t.pool.workerCount, -1)
			t.pool.workers.Delete(fmt.Sprintf("%p", t))
			t.pool.wg.Done()
		}()

		for {
			select {
			case data := <-t.jobChan:
				if data == nil {
					return
				}
				t.handler(data)
				select {
				case t.pool.workerQueuePool <- t.jobChan:
				default:
					return
				}
			case <-t.quit:
				close(t.jobChan)
			}
		}
	}()
}

// Stop 停止
func (t *Worker) Stop() {
	go func() {
		t.quit <- true
	}()
}
