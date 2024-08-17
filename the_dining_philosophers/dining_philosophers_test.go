package diningphilosophers

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

const (
	None = iota
	Left
	Right
)

const (
	Pick = iota + 1
	Put
	Eat
)

const (
	// NumPhilosophers is the number of philosophers and the number of forks
	NumPhilosophers = 5
	// NumCallsPerPhilosopher is the number of function calls per philosopher
	NumCallsPerPhilosopher = 5
	// NumCallsPerN is the number of function calls for each n
	NumCallsPerN = NumPhilosophers * NumCallsPerPhilosopher
)

func TestDiningPhilosophers(t *testing.T) {
	ns := []int{1, 13}

	for _, n := range ns {
		dp := NewDiningPhilosophers()
		actual, err := dp.run(n)
		if err != nil {
			t.Fatalf("dp.run(%d) = %v, want nil", n, err)
		}
		if err := validateCallCounts(actual); err != nil {
			t.Errorf("dp.run(%d) = %v, want nil", n, err)
		}
		if err := validateCallOrder(actual); err != nil {
			t.Errorf("dp.run(%d) = %v, want nil", n, err)
		}
	}
}

func (dp *DiningPhilosophers) run(n int) ([][3]int, error) {
	if err := validateN(n); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(NumCallsPerPhilosopher * n)
	ch := make(chan [3]int, NumCallsPerN*n)

	for _, i := range makeRandSeq(n) {
		go func() {
			defer wg.Done()
			dp.WantsToEat(i,
				func() { ch <- [3]int{i, Left, Pick} },
				func() { ch <- [3]int{i, Right, Pick} },
				func() { ch <- [3]int{i, None, Eat} },
				func() { ch <- [3]int{i, Left, Put} },
				func() { ch <- [3]int{i, Right, Put} },
			)
		}()
	}

	wg.Wait()
	close(ch)

	actual := make([][3]int, 0, NumCallsPerN*n)
	for call := range ch {
		actual = append(actual, call)
	}
	return actual, nil
}

func makeRandSeq(n int) []int {
	sequence := make([]int, 0, NumPhilosophers*n)
	for i := 0; i < NumPhilosophers; i++ {
		for j := 0; j < n; j++ {
			sequence = append(sequence, i)
		}
	}
	rand.Shuffle(len(sequence), func(i, j int) {
		sequence[i], sequence[j] = sequence[j], sequence[i]
	})
	return sequence
}

func validateN(n int) error {
	if n < 1 {
		return errors.New("n must be greater than 0")
	}
	if n > 60 {
		return errors.New("n must be less than or equal to 60")
	}
	return nil
}

func validateCallCounts(calls [][3]int) error {
	n := len(calls) / NumCallsPerN
	counts := [NumCallsPerN]int{}

	for _, call := range calls {
		philosopher := call[0]
		fork := call[1]
		operation := call[2]

		var funcID int
		if fork == Left && operation == Pick {
			funcID = 0
		} else if fork == Right && operation == Pick {
			funcID = 1
		} else if fork == None && operation == Eat {
			funcID = 2
		} else if fork == Left && operation == Put {
			funcID = 3
		} else if fork == Right && operation == Put {
			funcID = 4
		} else {
			return fmt.Errorf("invalid call: %v", call)
		}

		counts[philosopher*NumCallsPerPhilosopher+funcID]++
	}

	for _, count := range counts {
		if count != n {
			return errors.New("invalid call count")
		}
	}

	return nil
}

func validateCallOrder(calls [][3]int) error {
	forks := [NumPhilosophers]int{-1, -1, -1, -1, -1}
	for _, call := range calls {
		philosopher := call[0]

		var i int
		fork := call[1]
		switch fork {
		case None:
			i = -1
		case Left:
			i = philosopher
		case Right:
			i = (philosopher + 1) % NumPhilosophers
		}

		operation := call[2]
		switch operation {
		case Pick:
			if forks[i] != -1 {
				return errors.New("fork is already picked")
			}
			forks[i] = philosopher
		case Put:
			if forks[i] != philosopher {
				return errors.New("fork is not picked")
			}
			forks[i] = -1
		case Eat:
			if forks[philosopher] != philosopher || forks[(philosopher+1)%NumPhilosophers] != philosopher {
				return errors.New("forks are not picked")
			}
		}
	}

	return nil
}
