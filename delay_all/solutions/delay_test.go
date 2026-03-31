package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/delay_all/solutions"
	"github.com/N8Brooks/golang-concurrency/delay_all/testsuite"
)

func TestDelay(t *testing.T) {
	testsuite.Run(t, solutions.DelayAll[int])
}
