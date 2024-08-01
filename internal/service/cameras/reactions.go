package cameras

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactor"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
)

func (s Service) makeEventReactions(c *rms_cctv.Camera) []reactor.Reaction {
	camera, _ := s.CameraManager.GetCamera(model.CameraID(c.Id))
	archive, _ := s.CameraManager.GetArchive(model.CameraID(c.Id))
	result := []reactor.Reaction{
		s.ReactFactory.NewErrorReaction(),
		s.ReactFactory.NewNotifyReaction(camera, c.Name, c.Schedule, s.StateStorage),
	}
	if c.Mode == rms_cctv.RecordingMode_ByEvents || c.Mode == rms_cctv.RecordingMode_Optimal {
		result = append(result, s.ReactFactory.NewRecordingReaction(archive, s.Timeline, c.Mode == rms_cctv.RecordingMode_Optimal))
	}
	return result
}
