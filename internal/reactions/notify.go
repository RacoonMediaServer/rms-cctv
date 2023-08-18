package reactions

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/accessor"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"github.com/teambition/rrule-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

type notifyReaction struct {
	pub        micro.Publisher
	cam        accessor.Camera
	cameraName string
	schedule   *rrule.Set
}

func (n notifyReaction) React(l logger.Logger, event iva.PackedEvent) {
	if !event.IsEvent() {
		return
	}
	e := event.AsEvent()
	if e.Interval == iva.End {
		return
	}

	if e.Kind == events.Alert_TamperDetected {
		n.processEvent(l, e)
		return
	}
	// TODO
}

func (n notifyReaction) processEvent(l logger.Logger, event *iva.Event) {
	snapshot, err := n.cam.TakeSnapshot(model.PrimaryProfile)
	if err != nil {
		l.Logf(logger.ErrorLevel, "Get snapshot failed: %s", err)
	}
	alert := &events.Alert{
		Sender:         "rms-cctv",
		Timestamp:      event.Timestamp.Unix(),
		Camera:         n.cameraName,
		Kind:           event.Kind,
		Image:          snapshot,
		ImageMimeType:  "image/jpeg",
		DurationSec:    0,
		ExtraOffsetSec: 0,
	}
	notify(l, n.pub, alert)
}
