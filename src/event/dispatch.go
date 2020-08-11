package event

type EventDispatcher struct {
    provider ListenerProvider
}

func (t *EventDispatcher) dispatch(event Event) interface{} {
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
