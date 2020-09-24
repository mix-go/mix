package workerpool

import (
    "sync"
)

type Handler func(data interface{})

type Worker interface {
    Init(WorkerID int, workerPool chan JobQueue, wg *sync.WaitGroup, handler Handler)
    Run()
    Stop()
    Do(data interface{})
}

type WorkerTrait struct {
    WorkerID   int
    workerPool chan JobQueue
    wg         *sync.WaitGroup
    handler    Handler
    jobChan    JobQueue
    quit       chan bool
}

func (t *WorkerTrait) Init(WorkerID int, workerPool chan JobQueue, wg *sync.WaitGroup, handler Handler) {
    t.WorkerID = WorkerID
    t.workerPool = workerPool
    t.wg = wg
    t.handler = handler
    t.jobChan = make(chan interface{})
    t.quit = make(chan bool)
}

func (t *WorkerTrait) Run() {
    t.wg.Add(1)
    go func() {
        defer func() {
            t.wg.Done()
            if err := recover(); err != nil {
                t.Run()
                panic(err)
            }
        }()
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

func (t *WorkerTrait) Stop() {
    go func() {
        t.quit <- true
    }()
}
