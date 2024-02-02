package service

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/accessor"
	"github.com/RacoonMediaServer/rms-cctv/internal/camera"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactor"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/RacoonMediaServer/rms-packages/pkg/video"
	"time"
)

type DeviceManager interface {
	Register(cam *model.Camera) error
	Unregister(cam *model.Camera) error

	ListCameras() []*rms_cctv.Camera
	Add(cam *model.Camera, consumer camera.EventConsumer) error
	Remove(id model.CameraID) error

	GetCamera(id model.CameraID) (accessor.Camera, error)
	GetArchive(id model.CameraID) (accessor.Archive, error)

	GetStreamUri(id model.CameraID, profile model.Profile, transport video.Transport) (string, error)
	GetReplayUri(id model.CameraID, transport video.Transport, timestamp time.Time) (string, error)
}

type Reactor interface {
	PushEvent(event *iva.PackedEvent)
	SetReactions(cameraId model.CameraID, reactions []reactor.Reaction)
	DropReactions(cameraId model.CameraID)
}

type Database interface {
	AddCamera(camera *model.Camera) error
	LoadCameras() ([]*model.Camera, error)
	RemoveCamera(id model.CameraID) error
}
