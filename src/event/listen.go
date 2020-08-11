package event

import "fmt"

type ListenerProvider struct {
    EventListeners map[string][]Listener
}

func (t *ListenerProvider) getListenersForEvent(event interface{}) []func(event interface{}) {
    typ := fmt.Sprintf("%T", event)
    iterable := []func(i interface{}){}
    if listeners, ok := t.EventListeners[typ]; ok {
        for _, listener := range listeners {
            iterable = append(iterable, func(event interface{}) {
                listener.Process(event)
            })
        }
    }
    return iterable
}

func newListenerProvider(listeners ...Listener) ListenerProvider {
    eventListeners := map[string][]Listener{}
    for _, listener := range listeners {
        for _, e := range listener.Events() {
            typ := fmt.Sprintf("%T", e)
            if _, ok := eventListeners[typ]; ok {
                eventListeners[typ] = append(eventListeners[typ], listener)
            } else {
                eventListeners[typ] = []Listener{listener}
            }
        }
    }
    return ListenerProvider{
        EventListeners: eventListeners,
    }
}

type Listener interface {
    Events() []interface{}
    Process(i interface{})
}

type StoppableEvent interface {
    isPropagationStopped() bool
}
