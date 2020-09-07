package workerpool

import (
    "sync"
)

type Handler func(data interface{})
type ErrorHandler func(err interface{})

type Worker interface {
    Init(WorkerID int, workerPool chan JobQueue, wg *sync.WaitGroup, handler Handler, errorHandler ErrorHandler)
    Run()
    Stop()
    Do(data interface{})
    Error(err interface{})
}

type WorkerTrait struct {
    WorkerID     int
    workerPool   chan JobQueue
    wg           *sync.WaitGroup
    handler      Handler
    errorHandler ErrorHandler
    jobChan      JobQueue
    quit         chan bool
}

func (t *WorkerTrait) Init(WorkerID int, workerPool chan JobQueue, wg *sync.WaitGroup, handler Handler, errorHandler ErrorHandler) {
    t.WorkerID = WorkerID
    t.workerPool = workerPool
    t.wg = wg
    t.handler = handler
    t.errorHandler = errorHandler
    t.jobChan = make(chan interface{})
    t.quit = make(chan bool)
}

func (t *WorkerTrait) Run() {
    t.wg.Add(1)
    go func() {
        defer t.wg.Done()
        defer func() {
            if err := recover(); err != nil {
                t.errorHandler(err)
            }
        }()
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
