// Package solutions contains implementations of the Promise.allSettled problem.
package solutions

import (
	"sync"

	"github.com/N8Brooks/golang-concurrency/internal/promise"
	promiseallsettled "github.com/N8Brooks/golang-concurrency/promise_all_settled"
)

func PromiseAllSettled(functions []func() promise.Promiser[int]) promise.Promiser[[]promiseallsettled.Obj] {
	var wg sync.WaitGroup
	wg.Add(len(functions))

	results := make([]promiseallsettled.Obj, len(functions))

	for i, function := range functions {
		go func(i int, function func() promise.Promiser[int]) {
			defer wg.Done()
			val, err := function().Result()
			if err != nil {
				results[i] = promiseallsettled.Obj{Status: "rejected", Reason: err.Error()}
			} else {
				results[i] = promiseallsettled.Obj{Status: "fulfilled", Value: val}
			}
		}(i, function)
	}

	wg.Wait()
	return promise.Resolve(results)
}
