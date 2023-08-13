package camera

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"time"
)

const maxNetworkTimeout = 20 * time.Second

type DeviceService interface {
	GetDeviceName() (string, error)
}

type EventsService interface {
	// GetEvents returns actual events
	GetEvents() ([]*iva.Event, error)
}

type StreamService interface {
	GetSnapshot(profileToken string) ([]byte, error)
	GetStreamUri(profileToken string) (string, error)
}

// Controller can send commands to the specified camera
type Controller interface {
	DeviceService
	EventsService
	StreamService
}
