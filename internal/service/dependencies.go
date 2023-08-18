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
	GetCamera(id uint32) (accessor.Camera, error)
	GetArchive(id uint32) (accessor.Archive, error)
}

type Reactor interface {
	PushEvent(event *iva.PackedEvent)
	SetReactions(cameraId uint32, reactions []reactor.Reaction)
	DropReactions(cameraId uint32)
}
