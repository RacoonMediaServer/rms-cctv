package timeline

import (
	"context"
	"sort"
	"sync"
	"time"
)

const tickResolution = 100 * time.Millisecond

type timeline struct {
	ch     chan *point
	q      []*point
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	t      *time.Ticker
}

type point struct {
	when    time.Time
	handler Handler
}

func New() Timeline {
	t := timeline{
		t: time.NewTicker(tickResolution),
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.process()
	}()
	return &t
}

func (t *timeline) process() {
	for {
		select {
		case p := <-t.ch:
			t.add(p)
		case now := <-t.t.C:
			t.runExceeded(now)
		case <-t.ctx.Done():
			return
		}
	}
}

func (t *timeline) add(p *point) {
	t.q = append(t.q, p)
	sort.SliceStable(t.q, func(i, j int) bool {
		return t.q[i].when.Before(t.q[j].when)
	})
}

func (t *timeline) runExceeded(now time.Time) {
	for len(t.q) != 0 {
		p := t.q[0]
		if now.Before(p.when) {
			return
		}
		p.handler()
		t.q = t.q[1:]
	}
}

func (t *timeline) Defer(ctx context.Context, handler Handler, after time.Duration) {
	p := point{
		when:    time.Now().Add(after),
		handler: handler,
	}

	t.ch <- &p
}

func (t *timeline) Stop() {
	t.cancel()
	t.wg.Wait()
	close(t.ch)
}
