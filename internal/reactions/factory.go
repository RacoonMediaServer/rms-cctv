package reactions

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/accessor"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactor"
	"github.com/RacoonMediaServer/rms-cctv/internal/settings"
	"github.com/RacoonMediaServer/rms-cctv/internal/timeline"
	"github.com/teambition/rrule-go"
	"go-micro.dev/v4"
)

type Factory interface {
	NewErrorReaction() reactor.Reaction
	NewRecordingReaction(archive accessor.Archive, tm timeline.Defer, qualityControl bool) reactor.Reaction
	NewNotifyReaction(camera accessor.Camera, cameraName string, schedule *rrule.Set) reactor.Reaction
}

type factory struct {
	pub      micro.Publisher
	settings settings.Loader
}

func NewFactory(pub micro.Publisher, settings settings.Loader) Factory {
	return &factory{
		pub:      pub,
		settings: settings,
	}
}

func (f factory) NewErrorReaction() reactor.Reaction {
	return &errorReaction{pub: f.pub}
}

func (f factory) NewRecordingReaction(archive accessor.Archive, tm timeline.Defer, qualityControl bool) reactor.Reaction {
	return newRecordingReaction(archive, f.settings, tm, qualityControl)
}

func (f factory) NewNotifyReaction(camera accessor.Camera, cameraName string, schedule *rrule.Set) reactor.Reaction {
	return &notifyReaction{
		pub:        f.pub,
		cam:        camera,
		settings:   f.settings,
		cameraName: cameraName,
		schedule:   schedule,
	}
}
