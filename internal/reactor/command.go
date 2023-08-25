package reactor

import "github.com/RacoonMediaServer/rms-cctv/internal/model"

type setReactionsCommand struct {
	cameraId  model.CameraID
	reactions []Reaction
}

type dropReactionsCommand struct {
	cameraId model.CameraID
}
