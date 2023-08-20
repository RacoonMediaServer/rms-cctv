package reactions

import (
	"context"
	"github.com/RacoonMediaServer/rms-cctv/internal/accessor"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactor"
	"github.com/RacoonMediaServer/rms-cctv/internal/timeline"
	"go-micro.dev/v4/logger"
	"sync/atomic"
	"time"
)

// TODO: в настройки
const onceEventInterval = 30 * time.Second
const goodQuality uint = 0
const badQuality = 1

type controlFunc func() error

type recordingReaction struct {
	archive   accessor.Archive
	startFunc controlFunc
	stopFunc  controlFunc
	tm        timeline.Defer
	counter   atomic.Int32
}

func newRecordingReaction(archive accessor.Archive, tm timeline.Defer, qualityControl bool) reactor.Reaction {
	r := recordingReaction{
		archive: archive,
		tm:      tm,
	}

	if qualityControl {
		r.startFunc = r.setGoodQuality
		r.stopFunc = r.setBadQuality
	} else {
		r.startFunc = r.startRecording
		r.stopFunc = r.stopRecording
	}

	return &r
}

func (r *recordingReaction) React(l logger.Logger, event iva.PackedEvent) {
	if !event.IsEvent() {
		return
	}

	e := event.AsEvent()
	switch e.Interval {
	case iva.Begin:
		r.start(l)
	case iva.End:
		r.stop(l)
	case iva.Once:
		r.start(l)
		r.tm.Defer(context.TODO(), func() { r.stop(l) }, onceEventInterval)
	}
}

func (r *recordingReaction) start(l logger.Logger) {
	if r.counter.Add(1) != 1 {
		return
	}
	if err := r.startFunc(); err != nil {
		l.Logf(logger.ErrorLevel, "Start recording failed: %s", err)
	} else {
		l.Log(logger.InfoLevel, "Start recording")
	}
}

func (r *recordingReaction) stop(l logger.Logger) {
	if r.counter.Add(-1) != 0 {
		return
	}
	if err := r.stopFunc(); err != nil {
		l.Logf(logger.ErrorLevel, "Stop recording failed: %s", err)
	} else {
		l.Log(logger.InfoLevel, "Stop recording")
	}
}

func (r *recordingReaction) startRecording() error {
	return r.archive.StartRecording()
}

func (r *recordingReaction) stopRecording() error {
	return r.archive.StopRecording()
}

func (r *recordingReaction) setGoodQuality() error {
	return r.archive.SetQuality(goodQuality)
}

func (r *recordingReaction) setBadQuality() error {
	return r.archive.SetQuality(badQuality)
}
