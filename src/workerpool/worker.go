package workerpool

import (
    "fmt"
    "sync"
)

type Handler func(data interface{})

type Worker interface {
    Init(workerPool chan JobQueue, wg sync.WaitGroup, handler Handler)
    Start()
    Stop()
    Handle(data interface{})
}

type WorkerTrait struct {
    workerPool chan JobQueue
    wg         sync.WaitGroup
    handler    Handler
    jobQueue   JobQueue
    quit       chan bool
}

func (t *WorkerTrait) Init(workerPool chan JobQueue, wg sync.WaitGroup, handler Handler) {
    t.workerPool = workerPool
    t.wg = wg
    t.handler = handler
    t.jobQueue = make(chan interface{})
    t.quit = make(chan bool)
}

func (t *WorkerTrait) Start() {
    go func() {
        fmt.Println("add 1")
        t.wg.Add(1)
        defer func() {
            fmt.Println("done 1")
            t.wg.Done()
        }()
        for {
            select {
            case t.workerPool <- t.jobQueue:
                // none
            case data := <-t.jobQueue:
                if data == nil {
                    return
                }
                t.handler(data)
            case <-t.quit:
                close(t.jobQueue)
            }
        }
    }()
}

func (t *WorkerTrait) Stop() {
    go func() {
        t.quit <- true
    }()
}
