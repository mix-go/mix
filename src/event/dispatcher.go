package event

type Dispatcher struct {
    provider ListenerProvider
}

func (t *Dispatcher) Dispatch(event Event) Event {
    for _, callback := range t.provider.getListenersForEvent(event) {
        callback(event)
        if event.isPropagationStopped() {
            break
        }
    }
    return event;
}

func NewDispatcher(listeners ...Listener) *Dispatcher {
    return &Dispatcher{
        provider: newListenerProvider(listeners...),
    }
}
