package event

// 事件接口
type Event interface {
    IsPropagationStopped() bool
}

// 普通事件特性
type EventTrait struct {
}

// 是否停止传播
func (t *EventTrait) IsPropagationStopped() bool {
    return false
}

// 停止传播事件特性
type StoppableEventTrait struct {
}

// 是否停止传播
func (t *StoppableEventTrait) IsPropagationStopped() bool {
    return true
}
