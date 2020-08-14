package workerpool

type Dispatcher struct {
    JobQueue   chan<- interface{}
    MaxWorkers int
    workerPool chan chan<- interface{}
    quit       chan bool
}


