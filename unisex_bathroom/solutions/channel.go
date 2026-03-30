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
	empty           chan struct{}
	maleSwitch      *channelLightswitch
	femaleSwitch    *channelLightswitch
	maleMultiplex   chan struct{}
	femaleMultiplex chan struct{}
}

func NewChannel() *Channel {
	empty := make(chan struct{}, 1)
	empty <- struct{}{}
	maleMultiplex := make(chan struct{}, 3)
	femaleMultiplex := make(chan struct{}, 3)
	for range 3 {
		maleMultiplex <- struct{}{}
		femaleMultiplex <- struct{}{}
	}
	return &Channel{
		empty:           empty,
		maleSwitch:      newChannelLightswitch(),
		femaleSwitch:    newChannelLightswitch(),
		maleMultiplex:   maleMultiplex,
		femaleMultiplex: femaleMultiplex,
	}
}

func (b *Channel) Male(ctx context.Context, bathroom func()) error {
	if err := b.maleSwitch.lock(ctx, b.empty); err != nil {
		return err
	}
	if err := acquire(ctx, b.maleMultiplex); err != nil {
		b.maleSwitch.unlock(b.empty)
		return err
	}

	bathroom()

	release(b.maleMultiplex)
	b.maleSwitch.unlock(b.empty)
	return nil
}

func (b *Channel) Female(ctx context.Context, bathroom func()) error {
	if err := b.femaleSwitch.lock(ctx, b.empty); err != nil {
		return err
	}
	if err := acquire(ctx, b.femaleMultiplex); err != nil {
		b.femaleSwitch.unlock(b.empty)
		return err
	}

	bathroom()

	release(b.femaleMultiplex)
	b.femaleSwitch.unlock(b.empty)
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
