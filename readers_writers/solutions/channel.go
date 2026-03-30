package solutions

import "context"

type Channel struct {
	readers   int
	mutex     chan struct{}
	roomEmpty chan struct{}
}

func NewChannel() *Channel {
	mutex := make(chan struct{}, 1)
	roomEmpty := make(chan struct{}, 1)
	mutex <- struct{}{}
	roomEmpty <- struct{}{}
	return &Channel{
		mutex:     mutex,
		roomEmpty: roomEmpty,
	}
}

func (rw *Channel) Reader(ctx context.Context, read func()) error {
	if err := acquire(ctx, rw.mutex); err != nil {
		return err
	}

	rw.readers++
	if rw.readers == 1 {
		if err := acquire(ctx, rw.roomEmpty); err != nil {
			rw.readers--
			release(rw.mutex)
			return err
		}
	}
	release(rw.mutex)

	read()

	acquireUninterruptibly(rw.mutex)
	rw.readers--
	if rw.readers == 0 {
		release(rw.roomEmpty)
	}
	release(rw.mutex)
	return nil
}

func (rw *Channel) Writer(ctx context.Context, write func()) error {
	if err := acquire(ctx, rw.roomEmpty); err != nil {
		return err
	}
	write()
	release(rw.roomEmpty)
	return nil
}

func acquire(ctx context.Context, sem chan struct{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-sem:
		return nil
	}
}

func acquireUninterruptibly(sem chan struct{}) {
	<-sem
}

func release(sem chan struct{}) {
	sem <- struct{}{}
}
