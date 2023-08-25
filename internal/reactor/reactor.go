package reactor

import (
	"context"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"go-micro.dev/v4/logger"
	"sync"
)

const maxEvents = 100000

type Reactor struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	ch     chan interface{}
	l      logger.Logger

	r map[model.CameraID][]Reaction
}

func New() *Reactor {
	e := &Reactor{
		ch: make(chan interface{}, maxEvents),
		l:  logger.Fields(map[string]interface{}{"from": "reactor"}),
		r:  make(map[model.CameraID][]Reaction),
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		e.process()
	}()
	return e
}

func (r *Reactor) PushEvent(event *iva.PackedEvent) {
	r.l.Logf(logger.DebugLevel, "Got event: %+v", event)
	r.ch <- event
}

func (r *Reactor) SetReactions(cameraId model.CameraID, reactions []Reaction) {
	r.ch <- &setReactionsCommand{cameraId: cameraId, reactions: reactions}
}

func (r *Reactor) DropReactions(cameraId model.CameraID) {
	r.ch <- &dropReactionsCommand{cameraId: cameraId}
}

func (r *Reactor) Stop() {
	r.cancel()
	r.wg.Wait()
}
