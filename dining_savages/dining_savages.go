package diningsavages

import (
	"context"
)

type DiningSavages struct {
	servings chan chan struct{}
}

func NewDiningSavages() *DiningSavages {
	return &DiningSavages{
		servings: make(chan chan struct{}),
	}
}

func (ds *DiningSavages) Savage(ctx context.Context, getServingFromPot, eat func()) {
	for {
		select {
		case <-ctx.Done():
			return
		case consumed := <-ds.servings:
			getServingFromPot()
			eat()
			select {
			case <-ctx.Done():
				return
			case consumed <- struct{}{}:
			}
		}
	}
}

func (ds *DiningSavages) Cook(ctx context.Context, fillPot func() int) {
	for {
		numServings := fillPot()
		consumed := make(chan struct{}, numServings)
		for range numServings {
			select {
			case <-ctx.Done():
				return
			case ds.servings <- consumed:
			}
		}
		for range numServings {
			select {
			case <-ctx.Done():
				return
			case <-consumed:
			}
		}
	}
}
