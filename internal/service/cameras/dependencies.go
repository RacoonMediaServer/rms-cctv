package cameras

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/accessor"
	"github.com/RacoonMediaServer/rms-cctv/internal/camera"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactor"
	"github.com/RacoonMediaServer/rms-packages/pkg/media"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"time"
)

type DeviceManager interface {
	Register(cam *model.Camera) error
	Unregister(cam *model.Camera) error

	Add(cam *model.Camera, consumer camera.EventConsumer) error
	Modify(id model.CameraID, keepDays uint32, mode rms_cctv.RecordingMode) error
	Remove(id model.CameraID) error

	GetCamera(id model.CameraID) (accessor.Camera, error)
	GetArchive(id model.CameraID) (accessor.Archive, error)

	GetStreamUri(id model.CameraID, profile model.Profile, transport media.Transport) (string, error)
	GetReplayUri(id model.CameraID, transport media.Transport, timestamp time.Time) (string, error)
}

type Reactor interface {
	PushEvent(event *iva.PackedEvent)
	SetReactions(cameraId model.CameraID, reactions []reactor.Reaction)
	DropReactions(cameraId model.CameraID)
}

type Database interface {
	AddCamera(camera *model.Camera) error
	LoadCameras() ([]*model.Camera, error)
	GetCamera(id model.CameraID) (*model.Camera, error)
	UpdateCamera(camera *model.Camera) error
	RemoveCamera(id model.CameraID) error
}

type Schedules interface {
	GetSchedule(id string, defaultIfNotExists bool) *model.Schedule
}
