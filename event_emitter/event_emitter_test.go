//go:build challenge

package event_emitter_test

import (
	"testing"

	eventemitter "github.com/N8Brooks/golang-concurrency/event_emitter"
	"github.com/N8Brooks/golang-concurrency/event_emitter/testsuite"
)

func TestEventEmitter(t *testing.T) {
	testsuite.Run(t, func() eventemitter.Emitter {
		return eventemitter.NewEventEmitter()
	})
}
