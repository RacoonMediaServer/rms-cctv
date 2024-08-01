package reactions

import (
	"time"

	"github.com/RacoonMediaServer/rms-cctv/internal/accessor"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-cctv/internal/settings"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

// TODO: в конфиг
type notifyReaction struct {
	pub        micro.Publisher
	cam        accessor.Camera
	cameraName string
	schedule   string
	registry   Registry
	settings   settings.Loader
	lastTime   time.Time
	state      StateStorage
}

func (n *notifyReaction) thresholdInterval() time.Duration {
	return time.Duration(n.settings.Load().EventNotifyThresholdIntervalSec) * time.Second
}

func (n *notifyReaction) React(l logger.Logger, event iva.PackedEvent) {
	if !event.IsEvent() {
		return
	}
	e := event.AsEvent()
	if e.Interval == iva.End {
		return
	}

	sched := n.registry.Find(n.schedule, true)
	if !sched.Empty() && !sched.IsActiveNow() && !n.state.IsNobodyAtHome() {
		l.Logf(logger.DebugLevel, "Skip event, schedule deny notifications")
		return
	}

	if e.Kind == events.Alert_TamperDetected {
		n.processEvent(l, e)
		return
	}

	now := time.Now()
	if now.Sub(n.lastTime) < n.thresholdInterval() {
		l.Logf(logger.DebugLevel, "Event threshold reached, ignore")
		return
	}
	n.lastTime = now

	n.processEvent(l, e)
}

func (n *notifyReaction) processEvent(l logger.Logger, event *iva.Event) {
	l.Logf(logger.InfoLevel, "Event on camera %s: %s (interval = %d)", n.cameraName, event.Kind.String(), event.Interval)
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
