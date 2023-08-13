package reactions

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"go-micro.dev/v4/logger"
)

type recordingReaction struct {
	qualityControl bool
}

func (r recordingReaction) React(l logger.Logger, event iva.PackedEvent) {
	if !event.IsEvent() {
		return
	}
	//TODO implement me
	panic("implement me")
}
