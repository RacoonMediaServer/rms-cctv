package iva

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/misc"
	"time"
)

type Malfunction struct {
	CameraId  model.CameraID
	Timestamp time.Time
	Backtrace string
	Err       error
}

func NewMalfunction(err error) *Malfunction {
	return &Malfunction{
		Timestamp: time.Now(),
		Backtrace: misc.GetStackTrace(),
		Err:       err,
	}
}
