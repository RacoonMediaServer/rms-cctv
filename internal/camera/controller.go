package camera

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"net/url"
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
	GetProfiles() ([]string, error)
	GetSnapshot(profileToken string) ([]byte, error)
	GetStreamUri(profileToken string) (*url.URL, error)
}

// Controller can send commands to the specified camera
type Controller interface {
	DeviceService
	EventsService
	StreamService
}
