package event

// Event 事件接口
type Event interface {
	IsPropagationStopped() bool
}

// EventTrait 普通事件特性
type EventTrait struct {
}

// IsPropagationStopped 是否停止传播
func (t *EventTrait) IsPropagationStopped() bool {
	return false
}

// StoppableEventTrait 停止传播事件特性
type StoppableEventTrait struct {
}

// IsPropagationStopped 是否停止传播
func (t *StoppableEventTrait) IsPropagationStopped() bool {
	return true
}
