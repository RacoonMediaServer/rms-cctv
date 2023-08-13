package camera

import (
	"context"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"sync"
	"time"
)

const retryInterval = 10 * time.Second

type EventConsumer func(event *iva.PackedEvent)

type Listener struct {
	cam      EventsService
	consumer EventConsumer

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewListener(cam EventsService, consumer EventConsumer) *Listener {
	l := Listener{
		cam:      cam,
		consumer: consumer,
	}

	l.ctx, l.cancel = context.WithCancel(context.Background())
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		l.listen()
	}()

	return &l
}

func (l *Listener) isDone() bool {
	select {
	case <-l.ctx.Done():
		return true
	default:
		return false
	}
}

func (l *Listener) wait() {
	select {
	case <-l.ctx.Done():
	case <-time.After(retryInterval):
	}
}

func (l *Listener) listen() {
	for !l.isDone() {
		events, err := l.cam.GetEvents()
		if err != nil {
			l.consumer(iva.PackEvent(iva.NewMalfunction(err)))
			l.wait()
			continue
		}

		for _, event := range events {
			l.consumer(iva.PackEvent(event))
		}
	}
}

func (l *Listener) Stop() {
	l.cancel()
	l.wg.Wait()
}
