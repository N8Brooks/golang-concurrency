package solutions

import "context"

type channelLightswitch struct {
	count int
	mutex chan struct{}
}

func newChannelLightswitch() *channelLightswitch {
	mutex := make(chan struct{}, 1)
	mutex <- struct{}{}
	return &channelLightswitch{mutex: mutex}
}

func (ls *channelLightswitch) lock(ctx context.Context, sem chan struct{}) error {
	if err := acquire(ctx, ls.mutex); err != nil {
		return err
	}

	ls.count++
	if ls.count == 1 {
		if err := acquire(ctx, sem); err != nil {
			ls.count--
			release(ls.mutex)
			return err
		}
	}

	release(ls.mutex)
	return nil
}

func (ls *channelLightswitch) unlock(sem chan struct{}) {
	acquireUninterruptibly(ls.mutex)
	ls.count--
	if ls.count == 0 {
		release(sem)
	}
	release(ls.mutex)
}

type Channel struct {
	readSwitch  *channelLightswitch
	writeSwitch *channelLightswitch
	noReaders   chan struct{}
	noWriters   chan struct{}
}

func NewChannel() *Channel {
	noReaders := make(chan struct{}, 1)
	noWriters := make(chan struct{}, 1)
	noReaders <- struct{}{}
	noWriters <- struct{}{}
	return &Channel{
		readSwitch:  newChannelLightswitch(),
		writeSwitch: newChannelLightswitch(),
		noReaders:   noReaders,
		noWriters:   noWriters,
	}
}

func (rw *Channel) Reader(ctx context.Context, read func()) error {
	if err := acquire(ctx, rw.noReaders); err != nil {
		return err
	}
	if err := rw.readSwitch.lock(ctx, rw.noWriters); err != nil {
		release(rw.noReaders)
		return err
	}
	release(rw.noReaders)

	read()

	rw.readSwitch.unlock(rw.noWriters)
	return nil
}

func (rw *Channel) Writer(ctx context.Context, write func()) error {
	if err := rw.writeSwitch.lock(ctx, rw.noReaders); err != nil {
		return err
	}
	if err := acquire(ctx, rw.noWriters); err != nil {
		rw.writeSwitch.unlock(rw.noReaders)
		return err
	}

	write()

	release(rw.noWriters)
	rw.writeSwitch.unlock(rw.noReaders)
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
