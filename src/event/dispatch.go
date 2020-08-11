package event

type EventDispatcher struct {
    provider ListenerProvider
}

func (t *EventDispatcher) Dispatch(event Event) Event {
    for _, callback := range t.provider.getListenersForEvent(event) {
        callback(event)
        if event.isPropagationStopped() {
            break
        }
    }
    return event;
}

func NewEventDispatcher(listeners ...Listener) *EventDispatcher {
    return &EventDispatcher{
        provider: newListenerProvider(listeners...),
    }
}
