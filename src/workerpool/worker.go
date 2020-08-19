package workerpool

import (
    "sync"
)

type Handler func(data interface{})

type Worker interface {
    Init(workerPool chan JobQueue, wg *sync.WaitGroup, handler Handler)
    Start()
    Stop()
    Handle(data interface{})
}

type WorkerTrait struct {
    workerPool chan JobQueue
    wg         *sync.WaitGroup
    handler    Handler
    jobChan    JobQueue
    quit       chan bool
}

func (t *WorkerTrait) Init(workerPool chan JobQueue, wg *sync.WaitGroup, handler Handler) {
    t.workerPool = workerPool
    t.wg = wg
    t.handler = handler
    t.jobChan = make(chan interface{})
    t.quit = make(chan bool)
}

func (t *WorkerTrait) Start() {
    t.wg.Add(1)
    go func() {
        defer t.wg.Done()
        for {
            select {
            case t.workerPool <- t.jobChan:
                // none
            case data := <-t.jobChan:
                if data == nil {
                    return
                }
                t.handler(data)
            case <-t.quit:
                close(t.jobChan)
            }
        }
    }()
}

func (t *WorkerTrait) Stop() {
    go func() {
        t.quit <- true
    }()
}
