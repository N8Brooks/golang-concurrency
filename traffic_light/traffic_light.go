//go:build challenge

// Package trafficlight contains the challenge version of the traffic light
// problem.
//
// Cars arrive at a two-road intersection controlled by one traffic light. Road
// 1 starts green. When a car arrives, it may cross immediately if its road is
// already green; otherwise the light must be switched before that car crosses.
package trafficlight

type TrafficLight struct{}

func NewTrafficLight() *TrafficLight {
	return &TrafficLight{}
}

func (tl *TrafficLight) CarArrived(carID, roadID, direction int, turnGreen, crossCar func()) {
	panic("unimplemented")
}
