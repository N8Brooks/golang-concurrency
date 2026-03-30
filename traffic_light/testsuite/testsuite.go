// Package testsuite contains reusable behavioral tests for traffic-light
// implementations.
package testsuite

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"testing"
)

const (
	// RoadA is the north-south road and starts green.
	RoadA = iota + 1
	// RoadB is the east-west road.
	RoadB
)

const (
	// South goes from north to south on RoadA.
	South = iota + 1
	// North goes from south to north on RoadA.
	North
	// West goes from east to west on RoadB.
	West
	// East goes from west to east on RoadB.
	East
)

type TrafficLight interface {
	CarArrived(carID, roadID, direction int, turnGreen, crossCar func())
}

type event struct {
	turnGreenCarID int
	crossCarID     int
}

type call struct {
	carID     int
	roadID    int
	direction int
}

func makeCalls(directionCounts [4]int) []call {
	n := directionCounts[0] + directionCounts[1] + directionCounts[2] + directionCounts[3]
	calls := make([]call, 0, n)

	for i, count := range directionCounts {
		roadID := i/2 + 1
		direction := i + 1
		base := []call{{roadID: roadID, direction: direction}}
		calls = append(calls, slices.Repeat(base, count)...)
	}

	rng := rand.New(rand.NewSource(1))
	rng.Shuffle(len(calls), func(i, j int) {
		calls[i], calls[j] = calls[j], calls[i]
	})

	for i := range calls {
		calls[i].carID = i + 1
	}

	return calls
}

func runCalls(tl TrafficLight, calls []call) []event {
	var wg sync.WaitGroup
	wg.Add(len(calls))

	ch := make(chan event, 2*len(calls))
	for _, call := range calls {
		call := call
		go func() {
			defer wg.Done()
			tl.CarArrived(call.carID, call.roadID, call.direction, func() {
				ch <- event{turnGreenCarID: call.carID}
			}, func() {
				ch <- event{crossCarID: call.carID}
			})
		}()
	}

	wg.Wait()
	close(ch)

	events := make([]event, 0, len(calls))
	for event := range ch {
		events = append(events, event)
	}
	return events
}

func everyCarCrossedOnce(calls []call, events []event) error {
	counts := make([]int, len(calls))
	for _, event := range events {
		if event.crossCarID == 0 {
			continue
		}
		counts[event.crossCarID-1]++
	}

	errs := []error{}
	for i, count := range counts {
		if count != 1 {
			errs = append(errs, fmt.Errorf("car %d crossed %d times, want 1", i+1, count))
		}
	}
	return errors.Join(errs...)
}

func eventSequenceIsValid(calls []call, events []event, initialGreenRoadID int) (int, error) {
	greenRoadID := initialGreenRoadID
	errs := []error{}

	for _, event := range events {
		switch {
		case event.turnGreenCarID > 0:
			roadID := calls[event.turnGreenCarID-1].roadID
			if roadID == greenRoadID {
				errs = append(errs, fmt.Errorf("car %d turned green for road %d which was already green", event.turnGreenCarID, roadID))
			}
			greenRoadID = roadID
		case event.crossCarID > 0:
			roadID := calls[event.crossCarID-1].roadID
			if roadID != greenRoadID {
				errs = append(errs, fmt.Errorf("car %d crossed on road %d while road %d was green", event.crossCarID, roadID, greenRoadID))
			}
		}
	}

	return greenRoadID, errors.Join(errs...)
}

func Run(t *testing.T, newImpl func() TrafficLight) {
	t.Helper()

	t.Run("RoadAStartsGreen", func(t *testing.T) {
		tl := newImpl()
		turned := false
		crossed := false

		tl.CarArrived(1, RoadA, South, func() {
			turned = true
		}, func() {
			crossed = true
		})

		if turned {
			t.Fatal("turnGreen ran even though RoadA starts green")
		}
		if !crossed {
			t.Fatal("car did not cross")
		}
	})

	t.Run("RoadBSwitchesOnce", func(t *testing.T) {
		tl := newImpl()
		turns := 0
		crosses := 0

		tl.CarArrived(1, RoadB, East, func() {
			turns++
		}, func() {
			crosses++
		})

		if turns != 1 {
			t.Fatalf("got %d turnGreen calls, want 1", turns)
		}
		if crosses != 1 {
			t.Fatalf("got %d crosses, want 1", crosses)
		}
	})

	t.Run("OnlyRoadA", func(t *testing.T) {
		tl := newImpl()
		calls := makeCalls([4]int{5, 5, 0, 0})
		events := runCalls(tl, calls)

		if err := everyCarCrossedOnce(calls, events); err != nil {
			t.Fatal(err)
		}
		if _, err := eventSequenceIsValid(calls, events, RoadA); err != nil {
			t.Fatal(err)
		}

		for _, event := range events {
			if event.turnGreenCarID != 0 {
				t.Fatalf("unexpected turnGreen call for car %d while only RoadA cars arrived", event.turnGreenCarID)
			}
		}
	})

	t.Run("MixedTraffic", func(t *testing.T) {
		for _, directionCounts := range [][4]int{
			{1, 1, 1, 1},
			{2, 2, 2, 2},
			{3, 3, 3, 3},
			{10, 10, 10, 10},
		} {
			t.Run(fmt.Sprintf("Counts%v", directionCounts), func(t *testing.T) {
				tl := newImpl()
				calls := makeCalls(directionCounts)
				events := runCalls(tl, calls)

				if err := everyCarCrossedOnce(calls, events); err != nil {
					t.Fatal(err)
				}
				if _, err := eventSequenceIsValid(calls, events, RoadA); err != nil {
					t.Fatal(err)
				}
			})
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		tl := newImpl()
		greenRoadID := RoadA
		for range 5 {
			calls := makeCalls([4]int{2, 2, 2, 2})
			events := runCalls(tl, calls)

			if err := everyCarCrossedOnce(calls, events); err != nil {
				t.Fatal(err)
			}
			var err error
			greenRoadID, err = eventSequenceIsValid(calls, events, greenRoadID)
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}

func Benchmark(b *testing.B, newImpl func() TrafficLight) {
	b.Helper()
	b.ReportAllocs()

	for _, directionCounts := range [][4]int{
		{1, 1, 1, 1},
		{10, 10, 10, 10},
	} {
		b.Run(fmt.Sprintf("Counts%v", directionCounts), func(b *testing.B) {
			for b.Loop() {
				tl := newImpl()
				calls := makeCalls(directionCounts)
				runCalls(tl, calls)
			}
		})
	}
}
