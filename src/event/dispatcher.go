package event

type EventDispatcher struct {
    provider ListenerProvider
}

func (t *EventDispatcher) Dispatch(event Event) Event {
    for _, callback := range t.provider.getListenersForEvent(event) {
        callback(event)
        if event.IsPropagationStopped() {
            break
        }
    }
    return event;
}

func NewDispatcher(listeners ...Listener) Dispatcher {
    return &EventDispatcher{
        provider: newListenerProvider(listeners...),
    }
}

type Dispatcher interface {
    Dispatch(event Event) Event
}
