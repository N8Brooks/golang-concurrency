package trafficlight

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"testing"
)

const (
	// North-South road
	RoadA = iota + 1
	// East-West road
	RoadB
)

const (
	// From North to South on Road A
	South = iota + 1
	// From South to North on Road A
	North
	// From East to West on Road B
	West
	// From West to East on Road B
	East
)

func TestTrafficLight(t *testing.T) {
	directionCounts := [][4]int{
		{1, 1, 1, 1},
		{2, 2, 2, 2},
		{3, 3, 3, 3},
		{10, 10, 10, 10},
	}

	for _, directionCount := range directionCounts {
		t.Logf("testing directionCount: %v", directionCount)
		calls := makeCalls(directionCount)
		tl := NewTrafficLight()
		results := tl.run(calls)
		n := len(calls)
		if err := everyCarCrossedOnce(n, results); err != nil {
			t.Error(err)
		}
		if err := carsCrossedCorrectly(calls, results); err != nil {
			t.Error(err)
		}
	}
}

func makeCalls(directionCount [4]int) [][3]int {
	n := directionCount[0] + directionCount[1] + directionCount[2] + directionCount[3]
	calls := make([][3]int, 0, n)

	for i, count := range directionCount {
		roadID := i/2 + 1
		direction := i + 1
		call := [][3]int{{0, roadID, direction}}
		calls = append(calls, slices.Repeat(call, count)...)
	}

	rand.Shuffle(len(calls), func(i, j int) {
		calls[i], calls[j] = calls[j], calls[i]
	})

	for i := range calls {
		cardID := i + 1
		calls[i][0] = cardID
	}

	return calls
}

func (tl *TrafficLight) run(calls [][3]int) [][2]int {
	n := len(calls)

	var wg sync.WaitGroup
	wg.Add(n)
	// Discriminated union of {carID for turnGreen, carID for crossCar}
	ch := make(chan [2]int, 2*n)

	for _, call := range calls {
		go func() {
			defer wg.Done()
			carID, roadID, direction := call[0], call[1], call[2]
			turnGreen := func() {
				ch <- [2]int{carID, 0}
			}
			crossCar := func() {
				ch <- [2]int{0, carID}
			}
			tl.CarArrived(carID, roadID, direction, turnGreen, crossCar)
		}()
	}

	wg.Wait()
	close(ch)

	results := make([][2]int, 0, n)
	for cal := range ch {
		results = append(results, cal)
	}
	return results
}

func everyCarCrossedOnce(n int, results [][2]int) error {
	counts := make([]int, n)
	for _, result := range results {
		carID := result[1]
		if carID <= 0 {
			continue
		}
		i := carID - 1
		counts[i]++
	}

	errs := []error{}
	for i, count := range counts {
		if count != 1 {
			carID := i + 1
			errs = append(errs, fmt.Errorf("car %d crossed %d times, want 1", carID, count))
		}
	}
	return errors.Join(errs...)
}

func carsCrossedCorrectly(calls [][3]int, results [][2]int) error {
	greenRoadID := RoadA
	errs := []error{}
	for _, result := range results {
		turnGreenVariant, crossedVariant := result[0], result[1]
		if turnGreenVariant > 0 {
			carID := turnGreenVariant
			carIndex := carID - 1
			roadID := calls[carIndex][1]
			if greenRoadID == roadID {
				errs = append(errs, fmt.Errorf("car %d turned green on road %d which was already green", carID, roadID))
			}
			greenRoadID = roadID
		} else if crossedVariant > 0 {
			carID := crossedVariant
			carIndex := carID - 1
			roadID := calls[carIndex][1]
			if greenRoadID != roadID {
				errs = append(errs, fmt.Errorf("car %d crossed on road %d which was not green", carID, roadID))
			}
		}
	}
	return errors.Join(errs...)
}
