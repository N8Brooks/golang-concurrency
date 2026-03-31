// Package solutions contains implementations of the Event Emitter problem.
package solutions

import (
	"sync"

	eventemitter "github.com/N8Brooks/golang-concurrency/event_emitter"
)

type EventEmitter struct {
	mu     sync.RWMutex
	nextID int
	events map[string][]subscriber
}

type subscriber struct {
	id int
	cb eventemitter.Callback
}

type Subscription struct {
	emitter *EventEmitter
	event   string
	id      int
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		events: make(map[string][]subscriber),
	}
}

func (e *EventEmitter) Subscribe(eventName string, callback eventemitter.Callback) eventemitter.Subscription {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.nextID++
	id := e.nextID
	e.events[eventName] = append(e.events[eventName], subscriber{id: id, cb: callback})
	return &Subscription{emitter: e, event: eventName, id: id}
}

func (e *EventEmitter) Emit(eventName string, args ...int) []any {
	e.mu.RLock()
	listeners := append([]subscriber(nil), e.events[eventName]...)
	e.mu.RUnlock()

	results := make([]any, len(listeners))
	for i, listener := range listeners {
		results[i] = listener.cb(args...)
	}
	return results
}

func (s *Subscription) Unsubscribe() {
	s.emitter.mu.Lock()
	defer s.emitter.mu.Unlock()

	listeners := s.emitter.events[s.event]
	for i, listener := range listeners {
		if listener.id != s.id {
			continue
		}
		copy(listeners[i:], listeners[i+1:])
		listeners[len(listeners)-1] = subscriber{}
		s.emitter.events[s.event] = listeners[:len(listeners)-1]
		return
	}
}
