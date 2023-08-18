package manager

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/camera"
	"github.com/RacoonMediaServer/rms-cctv/internal/cctv"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
)

type cameraAccessor struct {
	cam      camera.StreamService
	profiles map[model.Profile]string
}

func (a cameraAccessor) TakeSnapshot(profile model.Profile) ([]byte, error) {
	return a.cam.GetSnapshot(a.profiles[profile])
}

type archiveAccessor struct {
	id       cctv.ID
	recorder cctv.RecordingController
}

func (a archiveAccessor) StartRecording() error {
	return a.recorder.StartRecording(a.id)
}

func (a archiveAccessor) StopRecording() error {
	return a.recorder.StopRecording(a.id)
}

func (a archiveAccessor) SetQuality(quality uint) error {
	return a.recorder.SetQuality(a.id, quality)
}
