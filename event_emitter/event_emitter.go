//go:build challenge

package event_emitter

type EventEmitter struct{}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{}
}

func (e *EventEmitter) Subscribe(eventName string, callback Callback) Subscription {
	panic("unimplemented")
}

func (e *EventEmitter) Emit(eventName string, args ...int) []any {
	panic("unimplemented")
}
