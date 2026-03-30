package solutions

import "context"

type Channel struct {
	firstDone  chan struct{}
	secondDone chan struct{}
}

func NewChannel() *Channel {
	return &Channel{
		firstDone:  make(chan struct{}, 1),
		secondDone: make(chan struct{}, 1),
	}
}

func (p *Channel) First(ctx context.Context, first func()) error {
	first()

	select {
	case p.firstDone <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Channel) Second(ctx context.Context, second func()) error {
	select {
	case <-p.firstDone:
	case <-ctx.Done():
		return ctx.Err()
	}

	second()

	select {
	case p.secondDone <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Channel) Third(ctx context.Context, third func()) error {
	select {
	case <-p.secondDone:
	case <-ctx.Done():
		return ctx.Err()
	}

	third()
	return nil
}
