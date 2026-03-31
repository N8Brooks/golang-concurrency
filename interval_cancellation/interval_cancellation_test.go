//go:build challenge

package interval_cancellation_test

import (
	"testing"

	intervalcancellation "github.com/N8Brooks/golang-concurrency/interval_cancellation"
	"github.com/N8Brooks/golang-concurrency/interval_cancellation/testsuite"
)

func TestCancellable(t *testing.T) {
	testsuite.Run(t, intervalcancellation.Cancellable)
}
