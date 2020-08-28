package workerpool

import (
    "sync"
)

type JobQueue chan interface{}

type Dispatcher struct {
    WorkerFunc func() Worker
    MaxWorkers int
    JobQueue   JobQueue
    workers    []Worker
    workerPool chan JobQueue
    wg         *sync.WaitGroup
    quit       chan bool
}

func (t *Dispatcher) Run() {
    for i := 0; i < t.MaxWorkers; i++ {
        w := t.WorkerFunc()
        w.Init(t.workerPool, t.wg, w.Handle)
        w.Start()

        t.workers = append(t.workers, w)
    }
    t.dispatch()
    t.wait()
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

func (t *Dispatcher) wait() {
    t.wg.Wait()
}

func NewDispatcher(workerFunc func() Worker, maxWorkers int, jobQueue chan interface{}) *Dispatcher {
    return &Dispatcher{
        WorkerFunc: workerFunc,
        MaxWorkers: maxWorkers,
        JobQueue:   jobQueue,
        workerPool: make(chan JobQueue, maxWorkers),
        wg:         &sync.WaitGroup{},
        quit:       make(chan bool),
    }
}
