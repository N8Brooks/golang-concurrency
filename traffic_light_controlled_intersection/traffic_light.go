package trafficlight

type TrafficLight struct {
	greenRoadID chan int
}

func NewTrafficLight() TrafficLight {
	tl := TrafficLight{}
	tl.greenRoadID = make(chan int, 1)
	tl.greenRoadID <- 1
	return tl
}

func (tl *TrafficLight) CarArrived(
	carID int, roadID int, direction int, turnGreen func(), crossCar func(),
) {
	if <-tl.greenRoadID != roadID {
		turnGreen()
	}
	crossCar()
	tl.greenRoadID <- roadID
}
