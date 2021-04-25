package xwp

import (
	"fmt"
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
			if p.WorkerRun != nil {
				p.WorkerRun(data)
			} else if p.WorkerRunI != nil {
				i := p.WorkerRunI
				i.Do(data)
			}
		},
		jobChan: make(chan interface{}),
		quit:    make(chan bool),
	}
}

// Run 执行
func (t *Worker) Run() {
	t.pool.wg.Add(1)
	go func() {
		defer func() {
			t.pool.workers.Delete(fmt.Sprintf("%p", t))
			t.pool.wg.Done()
		}()

		select {
		case t.pool.workerPool <- t.jobChan:
			t.pool.workers.Store(fmt.Sprintf("%p", t), t)
		default:
			return
		}
		for {
			select {
			case data := <-t.jobChan:
				if data == nil {
					return
				}
				t.handler(data)
				select {
				case t.pool.workerPool <- t.jobChan:
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
