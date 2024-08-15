package h2o

type H2O struct {
	// Hydrogen semaphore
	hydrogen chan struct{}
	// Oxygen semaphore
	oxygen chan struct{}
	// Molecule barrier
	barrier chan struct{}
}

func NewH2O() *H2O {
	h2o := &H2O{
		hydrogen: make(chan struct{}, 2),
		oxygen:   make(chan struct{}, 1),
		barrier:  make(chan struct{}),
	}
	return h2o
}

func (h2o *H2O) Hydrogen(releaseHydrogen func()) {
	h2o.hydrogen <- struct{}{}
	<-h2o.barrier
	releaseHydrogen()
	<-h2o.hydrogen
}

func (h2o *H2O) Oxygen(releaseOxygen func()) {
	h2o.oxygen <- struct{}{}
	h2o.barrier <- struct{}{}
	h2o.barrier <- struct{}{}
	releaseOxygen()
	<-h2o.oxygen
}
