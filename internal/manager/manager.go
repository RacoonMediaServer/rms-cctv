package manager

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-cctv/internal/camera"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"go-micro.dev/v4/logger"
	"sync"
)

type Manager struct {
	mu      sync.RWMutex
	devices map[camera.ID]*Device
	l       logger.Logger
	f       camera.Factory
}

func New(factory camera.Factory) *Manager {
	return &Manager{
		devices: make(map[camera.ID]*Device),
		l:       logger.Fields(map[string]interface{}{"from": "manager"}),
		f:       factory,
	}
}

func (m *Manager) Add(device Device, consumer camera.EventConsumer) error {
	wrappedConsumer := func(event *iva.PackedEvent) {
		event.SetCameraId(uint32(device.Id))
		consumer(event)
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.devices[device.Id]; ok {
		return fmt.Errorf("camera %d already exists", device.Id)
	}
	controller := m.f.New(context.TODO(), device.Url, camera.AutoDetect)
	device.l = camera.NewListener(controller, wrappedConsumer)

	m.devices[device.Id] = &device
	return nil
}
