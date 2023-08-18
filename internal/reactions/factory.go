package reactions

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/accessor"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactor"
	"github.com/teambition/rrule-go"
	"go-micro.dev/v4"
)

type Factory interface {
	NewErrorReaction(pub micro.Publisher) reactor.Reaction
	NewRecordingReaction(archive accessor.Archive, qualityControl bool) reactor.Reaction
	NewNotifyReaction(pub micro.Publisher, camera accessor.Camera, cameraName string, schedule *rrule.Set) reactor.Reaction
}

type factory struct {
}

func NewFactory() Factory {
	return &factory{}
}

func (f factory) NewErrorReaction(pub micro.Publisher) reactor.Reaction {
	return &errorReaction{pub: pub}
}

func (f factory) NewRecordingReaction(archive accessor.Archive, qualityControl bool) reactor.Reaction {
	return &recordingReaction{archive: archive, qualityControl: qualityControl}
}

func (f factory) NewNotifyReaction(pub micro.Publisher, camera accessor.Camera, cameraName string, schedule *rrule.Set) reactor.Reaction {
	return &notifyReaction{
		pub:        pub,
		cam:        camera,
		cameraName: cameraName,
		schedule:   schedule,
	}
}
