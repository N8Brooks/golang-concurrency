//go:build challenge

package delay_all_test

import (
	"testing"

	delayall "github.com/N8Brooks/golang-concurrency/delay_all"
	"github.com/N8Brooks/golang-concurrency/delay_all/testsuite"
)

func TestDelayAll(t *testing.T) {
	testsuite.Run(t, delayall.DelayAll[int])
}
