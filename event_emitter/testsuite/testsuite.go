// Package testsuite contains reusable behavioral tests for event emitter implementations.
package testsuite

import (
	"reflect"
	"testing"

	eventemitter "github.com/N8Brooks/golang-concurrency/event_emitter"
)

func Run(t *testing.T, newImpl func() eventemitter.Emitter) {
	t.Helper()

	t.Run("EmitWithoutSubscribers", func(t *testing.T) {
		emitter := newImpl()
		if got := emitter.Emit("firstEvent"); len(got) != 0 {
			t.Fatalf("Emit() = %v, want []", got)
		}
	})

	t.Run("MultipleSubscribersInOrder", func(t *testing.T) {
		emitter := newImpl()
		emitter.Subscribe("firstEvent", func(args ...int) any { return 5 })
		emitter.Subscribe("firstEvent", func(args ...int) any { return 6 })

		if got := emitter.Emit("firstEvent"); !reflect.DeepEqual(got, []any{5, 6}) {
			t.Fatalf("Emit() = %v, want [5 6]", got)
		}
	})

	t.Run("EmitWithArguments", func(t *testing.T) {
		emitter := newImpl()
		emitter.Subscribe("firstEvent", func(args ...int) any {
			return args[0] + args[1] + args[2]
		})

		if got := emitter.Emit("firstEvent", 1, 2, 3); !reflect.DeepEqual(got, []any{6}) {
			t.Fatalf("Emit() = %v, want [6]", got)
		}
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		emitter := newImpl()
		sub1 := emitter.Subscribe("firstEvent", func(args ...int) any { return args[0] + 1 })
		emitter.Subscribe("firstEvent", func(args ...int) any { return args[0] + 2 })
		sub1.Unsubscribe()

		if got := emitter.Emit("firstEvent", 5); !reflect.DeepEqual(got, []any{7}) {
			t.Fatalf("Emit() = %v, want [7]", got)
		}
	})
}
