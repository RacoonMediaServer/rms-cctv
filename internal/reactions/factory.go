package reactions

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/reactor"
	"go-micro.dev/v4"
)

type Factory interface {
	NewErrorReaction(pub micro.Publisher) reactor.Reaction
	NewRecordingReaction(qualityControl bool) reactor.Reaction
	NewNotifyReaction(pub micro.Publisher, schedule string) reactor.Reaction
}

type factory struct {
}

func NewFactory() Factory {
	return &factory{}
}

func (f factory) NewErrorReaction(pub micro.Publisher) reactor.Reaction {
	return &errorReaction{pub: pub}
}

func (f factory) NewRecordingReaction(qualityControl bool) reactor.Reaction {
	return &recordingReaction{qualityControl: qualityControl}
}

func (f factory) NewNotifyReaction(pub micro.Publisher, schedule string) reactor.Reaction {
	return &notifyReaction{
		pub:      pub,
		schedule: schedule,
	}
}
