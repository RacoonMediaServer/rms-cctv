package manager

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/camera"
	"net/url"
)

type Device struct {
	Id   camera.ID
	Name string
	Url  *url.URL

	l *camera.Listener
}
