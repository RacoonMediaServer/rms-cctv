package reactor

import (
	"fmt"

	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"go-micro.dev/v4/logger"
)

func (r *Reactor) process() {
	for {
		select {
		case cmd := <-r.ch:
			r.handleCommand(cmd)
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *Reactor) handleCommand(cmd interface{}) {
	switch content := cmd.(type) {
	case *iva.PackedEvent:
		r.handleEvent(content)
	case *setReactionsCommand:
		r.setReactions(content)
	case *dropReactionsCommand:
		r.dropReactions(content)
	default:
		panic(fmt.Errorf("unknown event: %T", content))
	}
}

func (r *Reactor) handleEvent(event *iva.PackedEvent) {
	reactions := r.r[event.CameraId()]
	for _, reaction := range reactions {
		reaction.React(r.l, *event)
	}
}

func (r *Reactor) setReactions(cmd *setReactionsCommand) {
	r.r[cmd.cameraId] = cmd.reactions
	r.l.Logf(logger.DebugLevel, "Reactions set for camera %d", cmd.cameraId)
}

func (r *Reactor) dropReactions(cmd *dropReactionsCommand) {
	delete(r.r, cmd.cameraId)
}
