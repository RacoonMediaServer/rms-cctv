package iva

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"time"
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
