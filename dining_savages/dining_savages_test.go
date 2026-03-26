package diningsavages_test

import (
	"context"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"

	diningsavages "github.com/N8Brooks/golang-concurrency/dining_savages"
)

func TestDiningSavages(t *testing.T) {
	const numSavages = 3
	const numEaten = 100
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		ds := diningsavages.NewDiningSavages()

		var servings atomic.Int64

		putServingsInPot := func() int {
			newServings := rand.Int64N(9) + 1
			oldServings := servings.Swap(newServings)
			if oldServings > 0 {
				t.Errorf("cook filled the pot with %d servings, but there were still %d servings left in the pot", newServings, oldServings)
			} else if oldServings == 0 {
				t.Logf("cook filled the pot with %d servings", newServings)
			} else {
				panic("servings should never be negative")
			}
			return int(newServings)
		}

		go ds.Cook(ctx, putServingsInPot)

		var eaten atomic.Int64
		done := make(chan struct{})
		stop := sync.OnceFunc(func() {
			cancel()
			close(done)
		})

		for i := range numSavages {
			var hasServing atomic.Bool

			getServingFromPot := func() {
				if servings.Add(-1) < 0 {
					t.Errorf("savage %d got a serving from the pot, but there were no servings left", i)
				} else {
					t.Logf("savage %d got a serving from the pot, %d servings left", i, servings.Load())
				}
				if hasServing.Swap(true) {
					t.Errorf("savage %d got a serving from the pot, but they already had a serving", i)
				} else {
					t.Logf("savage %d got a serving from the pot", i)
				}
			}

			eat := func() {
				if !hasServing.Swap(false) {
					t.Errorf("savage %d ate a serving, but they didn't have a serving", i)
				} else {
					t.Logf("savage %d ate a serving", i)
				}
				if eaten.Add(1) >= numEaten {
					stop()
				}
			}

			go ds.Savage(ctx, getServingFromPot, eat)
		}

		<-done
		synctest.Wait()
	})
}
