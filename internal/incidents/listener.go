package incidents

import (
	"context"
	"github.com/RacoonMediaServer/rms-cctv/internal/camera"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"sync"
	"time"
)

const retryInterval = 10 * time.Second

type EventConsumer interface {
	Event(event *iva.Event)
	Error(err error)
}

type Listener struct {
	cam      camera.Events
	consumer EventConsumer

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewListener(cam camera.Events, consumer EventConsumer) *Listener {
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

	return nil
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
			l.consumer.Error(err)
			l.wait()
			continue
		}

		for _, event := range events {
			l.consumer.Event(event)
		}
	}
}

func (l *Listener) Stop() {
	l.cancel()
	l.wg.Wait()
}
