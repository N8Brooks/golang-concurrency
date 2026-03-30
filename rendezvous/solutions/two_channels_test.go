package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/rendezvous/solutions"
	"github.com/N8Brooks/golang-concurrency/rendezvous/testsuite"
)

func TestTwoChannel(t *testing.T) {
	testsuite.Run(t, func() testsuite.Rendezvous {
		return solutions.NewTwoChannel()
	})
}
