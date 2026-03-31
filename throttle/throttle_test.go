//go:build challenge

package throttle_test

import (
	"testing"

	throttlepkg "github.com/N8Brooks/golang-concurrency/throttle"
	"github.com/N8Brooks/golang-concurrency/throttle/testsuite"
)

func TestThrottle(t *testing.T) {
	testsuite.Run(t, throttlepkg.Throttle)
}
