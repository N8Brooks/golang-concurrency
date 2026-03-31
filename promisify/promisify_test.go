//go:build challenge

package promisify_test

import (
	"testing"

	promisifypkg "github.com/N8Brooks/golang-concurrency/promisify"
	"github.com/N8Brooks/golang-concurrency/promisify/testsuite"
)

func TestPromisify(t *testing.T) {
	testsuite.Run(t, promisifypkg.Promisify)
}
