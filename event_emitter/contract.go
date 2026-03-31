// Package event_emitter contains the Event Emitter problem and shared contract.
package event_emitter

type Callback func(args ...int) any

type Subscription interface {
	Unsubscribe()
}

type Emitter interface {
	Subscribe(eventName string, callback Callback) Subscription
	Emit(eventName string, args ...int) []any
}
