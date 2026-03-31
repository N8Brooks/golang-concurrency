package solutions_test

import (
	"testing"

	"github.com/N8Brooks/golang-concurrency/promisify/solutions"
	"github.com/N8Brooks/golang-concurrency/promisify/testsuite"
)

func TestPromise(t *testing.T) {
	testsuite.Run(t, solutions.Promisify)
}
