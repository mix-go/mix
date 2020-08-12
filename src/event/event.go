package event

type Event interface {
    isPropagationStopped() bool
}

type EventTrait struct {
}

func (t *EventTrait) isPropagationStopped() bool {
    return false
}
