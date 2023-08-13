package service

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/camera"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-cctv/internal/manager"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactor"
)

type DeviceManager interface {
	Add(device manager.Device, consumer camera.EventConsumer) error
}

type Reactor interface {
	PushEvent(event *iva.PackedEvent)
	SetReactions(cameraId uint32, reactions []reactor.Reaction)
	DropReactions(cameraId uint32)
}
