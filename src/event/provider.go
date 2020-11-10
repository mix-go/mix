package event

import "fmt"

// 监听供应器
type ListenerProvider struct {
	EventListeners map[string][]Listener
}

// 获取监听器列表
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

// 创建监听供应器
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
