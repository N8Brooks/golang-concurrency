// Package solutions contains implementations of the traffic light problem.
package solutions

type Channel struct {
	greenRoadID chan int
}

func NewChannel() *Channel {
	tl := &Channel{
		greenRoadID: make(chan int, 1),
	}
	tl.greenRoadID <- 1
	return tl
}

func (tl *Channel) CarArrived(carID, roadID, direction int, turnGreen, crossCar func()) {
	if <-tl.greenRoadID != roadID {
		turnGreen()
	}
	crossCar()
	tl.greenRoadID <- roadID
}
