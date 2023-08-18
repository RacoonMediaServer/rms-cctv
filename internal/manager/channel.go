package manager

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/camera"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"net/url"
)

type channel struct {
	cameraUrl *url.URL
	camera    *model.Camera
	l         *camera.Listener
}
