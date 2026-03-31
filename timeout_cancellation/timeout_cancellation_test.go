//go:build challenge

package timeout_cancellation_test

import (
	"testing"

	timeoutcancellation "github.com/N8Brooks/golang-concurrency/timeout_cancellation"
	"github.com/N8Brooks/golang-concurrency/timeout_cancellation/testsuite"
)

func TestCancellable(t *testing.T) {
	testsuite.Run(t, timeoutcancellation.Cancellable)
}
