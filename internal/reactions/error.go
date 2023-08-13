package reactions

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

type errorReaction struct {
	pub micro.Publisher
}

func (e errorReaction) React(l logger.Logger, event iva.PackedEvent) {
	if !event.IsMalfunction() {
		return
	}

	l = l.Fields(map[string]interface{}{"camera": event.CameraId(), "reaction": "error"})
	m := event.AsMalfunction()
	n := &events.Malfunction{
		Sender:     "rms-cctv",
		Timestamp:  m.Timestamp.Unix(),
		Error:      m.Err.Error(),
		System:     events.Malfunction_Cameras,
		Code:       events.Malfunction_CannotAccess,
		StackTrace: m.Backtrace,
	}

	l.Logf(logger.InfoLevel, "Error reaction triggered on camera")

	notify(l, e.pub, n)
}
