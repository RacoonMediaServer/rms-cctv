package service

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/reactor"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
)

func (s Service) makeEventReactions(mode rms_cctv.RecordingMode, schedule string) []reactor.Reaction {
	result := []reactor.Reaction{
		s.ReactFactory.NewErrorReaction(s.Notifier),
		s.ReactFactory.NewNotifyReaction(s.Notifier, schedule),
	}
	if mode == rms_cctv.RecordingMode_ByEvents || mode == rms_cctv.RecordingMode_Optimal {
		result = append(result, s.ReactFactory.NewRecordingReaction(mode == rms_cctv.RecordingMode_Optimal))
	}
	return result
}
