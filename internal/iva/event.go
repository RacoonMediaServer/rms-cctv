package iva

import (
	"time"

	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
)

type Interval int

const (
	End Interval = iota
	Begin
	Once
)

type Event struct {
	CameraId  model.CameraID
	Kind      events.Alert_Kind
	Interval  Interval
	Timestamp time.Time
}

func (i Interval) String() string {
	switch i {
	case End:
		return "End"
	case Begin:
		return "Begin"
	case Once:
		return "Once"
	default:
		return "und"
	}
}
