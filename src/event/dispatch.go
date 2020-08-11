package event

type EventDispatcher struct {
    provider ListenerProvider
}

func (t *EventDispatcher) dispatch(event interface{}) interface{} {
    isStoppableEvent := func(event interface{}) bool {
        switch event.(type) {
        case StoppableEvent, *StoppableEvent:
            return true
        }
        return false
    }
    for _, callback := range t.provider.getListenersForEvent(event) {
        callback(event)
        if isStoppableEvent(event) { // 停止传播
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
