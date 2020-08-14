package event

type Listener interface {
    Events() []Event
    Process(Event)
}

