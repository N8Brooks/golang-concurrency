package solutions_test

import (
	"testing"

	eventemitter "github.com/N8Brooks/golang-concurrency/event_emitter"
	"github.com/N8Brooks/golang-concurrency/event_emitter/solutions"
	"github.com/N8Brooks/golang-concurrency/event_emitter/testsuite"
)

func TestSyncMutex(t *testing.T) {
	testsuite.Run(t, func() eventemitter.Emitter {
		return solutions.NewEventEmitter()
	})
}
