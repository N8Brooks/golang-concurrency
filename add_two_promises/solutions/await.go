// Package solutions contains implementations of the Add Two Promises problem.
package solutions

import "github.com/N8Brooks/golang-concurrency/internal/promise"

func AddTwoPromises(a, b promise.Promiser[int]) promise.Promiser[int] {
	return promise.Resolve(a.Await() + b.Await())
}
