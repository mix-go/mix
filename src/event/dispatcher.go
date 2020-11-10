package event

// 事件调度器
type EventDispatcher struct {
	provider ListenerProvider
}

// 调度事件
func (t *EventDispatcher) Dispatch(event Event) Event {
	for _, callback := range t.provider.getListenersForEvent(event) {
		callback(event)
		if event.IsPropagationStopped() {
			break
		}
	}
	return event
}

// 创建调度器
func NewDispatcher(listeners ...Listener) Dispatcher {
	return &EventDispatcher{
		provider: newListenerProvider(listeners...),
	}
}

// 调度器接口
type Dispatcher interface {
	Dispatch(event Event) Event
}
