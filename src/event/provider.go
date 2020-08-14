package event

import "fmt"

type ListenerProvider struct {
    EventListeners map[string][]Listener
}

func (t *ListenerProvider) getListenersForEvent(event Event) []func(event Event) {
    typ := fmt.Sprintf("%T", event)
    iterable := []func(event Event){}
    if listeners, ok := t.EventListeners[typ]; ok {
        for _, listener := range listeners {
            iterable = append(iterable, func(event Event) {
                listener.Process(event)
            })
        }
    }
    return iterable
}

func newListenerProvider(listeners ...Listener) ListenerProvider {
    eventListeners := map[string][]Listener{}
    for _, listener := range listeners {
        for _, event := range listener.Events() {
            typ := fmt.Sprintf("%T", event)
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
