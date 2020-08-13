package event

type Event interface {
    isPropagationStopped() bool
}

type EventTrait struct {
}

func (t *EventTrait) isPropagationStopped() bool {
    return false
}

type StoppableEventTrait struct {
}

func (t *StoppableEventTrait) isPropagationStopped() bool {
    return true
}
