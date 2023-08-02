package iva

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"time"
)

type Interval int

const (
	Begin Interval = iota
	End
	Once
)

type Event struct {
	CameraId  uint32
	Kind      events.Alert_Kind
	Interval  Interval
	Timestamp time.Time
}
