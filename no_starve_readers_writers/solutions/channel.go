package solutions

import "context"

type channelLightswitch struct {
	readers int
	mutex   chan struct{}
}

func newChannelLightswitch() *channelLightswitch {
	mutex := make(chan struct{}, 1)
	mutex <- struct{}{}
	return &channelLightswitch{mutex: mutex}
}

func (ls *channelLightswitch) lock(ctx context.Context, roomEmpty chan struct{}) error {
	if err := acquire(ctx, ls.mutex); err != nil {
		return err
	}

	ls.readers++
	if ls.readers == 1 {
		if err := acquire(ctx, roomEmpty); err != nil {
			ls.readers--
			release(ls.mutex)
			return err
		}
	}

	release(ls.mutex)
	return nil
}

func (ls *channelLightswitch) unlock(roomEmpty chan struct{}) {
	acquireUninterruptibly(ls.mutex)
	ls.readers--
	if ls.readers == 0 {
		release(roomEmpty)
	}
	release(ls.mutex)
}

type Channel struct {
	readSwitch *channelLightswitch
	roomEmpty  chan struct{}
	turnstile  chan struct{}
}

func NewChannel() *Channel {
	roomEmpty := make(chan struct{}, 1)
	turnstile := make(chan struct{}, 1)
	roomEmpty <- struct{}{}
	turnstile <- struct{}{}
	return &Channel{
		readSwitch: newChannelLightswitch(),
		roomEmpty:  roomEmpty,
		turnstile:  turnstile,
	}
}

func (rw *Channel) Reader(ctx context.Context, read func()) error {
	if err := acquire(ctx, rw.turnstile); err != nil {
		return err
	}
	release(rw.turnstile)

	if err := rw.readSwitch.lock(ctx, rw.roomEmpty); err != nil {
		return err
	}
	read()
	rw.readSwitch.unlock(rw.roomEmpty)
	return nil
}

func (rw *Channel) Writer(ctx context.Context, write func()) error {
	if err := acquire(ctx, rw.turnstile); err != nil {
		return err
	}
	if err := acquire(ctx, rw.roomEmpty); err != nil {
		release(rw.turnstile)
		return err
	}

	write()

	release(rw.turnstile)
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
