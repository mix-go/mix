package event

type Event interface {
    IsPropagationStopped() bool
}

type EventTrait struct {
}

func (t *EventTrait) IsPropagationStopped() bool {
    return false
}

type StoppableEventTrait struct {
}

func (t *StoppableEventTrait) IsPropagationStopped() bool {
    return true
}
