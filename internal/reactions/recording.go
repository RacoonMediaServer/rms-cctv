package reactions

import (
	"context"
	"github.com/RacoonMediaServer/rms-cctv/internal/accessor"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactor"
	"github.com/RacoonMediaServer/rms-cctv/internal/settings"
	"github.com/RacoonMediaServer/rms-cctv/internal/timeline"
	"go-micro.dev/v4/logger"
	"sync/atomic"
	"time"
)

// TODO: в настройки
const onceEventInterval = 30 * time.Second
const goodQuality uint = 0
const badQuality uint = 1

type controlFunc func() error

type recordingReaction struct {
	archive   accessor.Archive
	startFunc controlFunc
	stopFunc  controlFunc
	tm        timeline.Defer
	settings  settings.Loader
	counter   atomic.Int32
}

func newRecordingReaction(archive accessor.Archive, settings settings.Loader, tm timeline.Defer, qualityControl bool) reactor.Reaction {
	r := recordingReaction{
		archive:  archive,
		tm:       tm,
		settings: settings,
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
	newVal := r.counter.Add(-1)
	if newVal < 0 {
		l.Logf(logger.WarnLevel, "Events order malfunction, reset state")
		newVal = 0
		r.counter.Store(newVal)
	}
	if newVal != 0 {
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
