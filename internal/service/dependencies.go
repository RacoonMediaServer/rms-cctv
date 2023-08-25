package service

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/accessor"
	"github.com/RacoonMediaServer/rms-cctv/internal/camera"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactor"
)

type DeviceManager interface {
	Register(cam *model.Camera) error
	Unregister(cam *model.Camera) error

	Add(cam *model.Camera, consumer camera.EventConsumer) error
	Remove(id model.CameraID) error

	GetCamera(id model.CameraID) (accessor.Camera, error)
	GetArchive(id model.CameraID) (accessor.Archive, error)
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
