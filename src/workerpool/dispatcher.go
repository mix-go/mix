package workerpool

import (
    "sync"
)

type JobQueue chan interface{}

type Dispatcher struct {
    JobQueue   JobQueue
    MaxWorkers int
    workers    []Worker
    workerPool chan JobQueue
    wg         *sync.WaitGroup
    quit       chan bool
}

func (t *Dispatcher) Start(fn func() Worker) {
    for i := 0; i < t.MaxWorkers; i++ {
        w := fn()
        w.Init(t.workerPool, t.wg, w.Handle)
        w.Start()

        t.workers = append(t.workers, w)
    }
    t.dispatch()
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

func (t *Dispatcher) Stop() {
    go func() {
        t.quit <- true
    }()
}

func (t *Dispatcher) Wait() {
    t.wg.Wait()
}

func NewDispatcher(jobQueue chan interface{}, maxWorkers int) *Dispatcher {
    return &Dispatcher{
        JobQueue:   jobQueue,
        MaxWorkers: maxWorkers,
        workerPool: make(chan JobQueue, maxWorkers),
        wg:         &sync.WaitGroup{},
        quit:       make(chan bool),
    }
}
