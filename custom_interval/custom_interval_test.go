//go:build challenge

package custom_interval_test

import (
	"testing"

	custominterval "github.com/N8Brooks/golang-concurrency/custom_interval"
	"github.com/N8Brooks/golang-concurrency/custom_interval/testsuite"
)

func TestCustomInterval(t *testing.T) {
	testsuite.Run(t, custominterval.CustomInterval, custominterval.CustomClearInterval)
}
