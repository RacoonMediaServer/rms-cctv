package reactions

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-cctv/internal/mocks"
	"github.com/RacoonMediaServer/rms-cctv/internal/settings"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"github.com/golang/mock/gomock"
	"go-micro.dev/v4/logger"
	"testing"
	"time"
)

const cameraID = 133

func makeEvent(kind events.Alert_Kind, interval iva.Interval) *iva.PackedEvent {
	e := iva.Event{
		CameraId: cameraID,
		Kind:     kind,
		Interval: interval,
	}
	return iva.PackEvent(&e)
}

func TestRecordingReaction_React(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tm := mocks.NewMockTimeline(ctrl)
	archive := mocks.NewMockArchive(ctrl)
	conf := settings.New()

	r := newRecordingReaction(archive, conf, tm, false)

	archive.EXPECT().StartRecording().Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_TamperDetected, iva.Begin))

	archive.EXPECT().StopRecording().Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_CrossLineDetected, iva.End))

	archive.EXPECT().StartRecording().Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_MotionDetected, iva.Begin))
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_CrossLineDetected, iva.Begin))
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_TamperDetected, iva.Begin))

	r.React(logger.DefaultLogger, *makeEvent(events.Alert_CrossLineDetected, iva.End))
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_MotionDetected, iva.End))
	archive.EXPECT().StopRecording().Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_TamperDetected, iva.End))

	// test recovery
	archive.EXPECT().StopRecording().Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_TamperDetected, iva.End))
	archive.EXPECT().StopRecording().Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_TamperDetected, iva.End))

	archive.EXPECT().StartRecording().Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_MotionDetected, iva.Begin))
	archive.EXPECT().StopRecording().Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_MotionDetected, iva.End))

	archive.EXPECT().StartRecording().Return(nil).Times(1)
	tm.EXPECT().Defer(gomock.Any(), gomock.Any(), gomock.Eq(time.Duration(conf.Load().OneEventDefaultDurationSec)*time.Second))
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_MotionDetected, iva.Once))
}

func TestRecordingReaction_React2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tm := mocks.NewMockTimeline(ctrl)
	archive := mocks.NewMockArchive(ctrl)
	conf := settings.New()

	r := newRecordingReaction(archive, conf, tm, true)

	archive.EXPECT().SetQuality(gomock.Eq(goodQuality)).Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_TamperDetected, iva.Begin))

	archive.EXPECT().SetQuality(gomock.Eq(badQuality)).Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_CrossLineDetected, iva.End))

	archive.EXPECT().SetQuality(gomock.Eq(goodQuality)).Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_MotionDetected, iva.Begin))
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_CrossLineDetected, iva.Begin))
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_TamperDetected, iva.Begin))

	r.React(logger.DefaultLogger, *makeEvent(events.Alert_CrossLineDetected, iva.End))
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_MotionDetected, iva.End))
	archive.EXPECT().SetQuality(gomock.Eq(badQuality)).Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_TamperDetected, iva.End))

	// test recovery
	archive.EXPECT().SetQuality(gomock.Eq(badQuality)).Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_TamperDetected, iva.End))
	archive.EXPECT().SetQuality(gomock.Eq(badQuality)).Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_TamperDetected, iva.End))

	archive.EXPECT().SetQuality(gomock.Eq(goodQuality)).Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_MotionDetected, iva.Begin))
	archive.EXPECT().SetQuality(gomock.Eq(badQuality)).Return(nil).Times(1)
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_MotionDetected, iva.End))

	archive.EXPECT().SetQuality(gomock.Eq(goodQuality)).Return(nil).Times(1)
	tm.EXPECT().Defer(gomock.Any(), gomock.Any(), gomock.Eq(time.Duration(conf.Load().OneEventDefaultDurationSec)*time.Second))
	r.React(logger.DefaultLogger, *makeEvent(events.Alert_MotionDetected, iva.Once))
}
