package accessor

import "github.com/RacoonMediaServer/rms-cctv/internal/model"

type Camera interface {
	TakeSnapshot(profile model.Profile) ([]byte, error)
}
