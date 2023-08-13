package reactions

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

type notifyReaction struct {
	pub      micro.Publisher
	schedule string
}

func (n notifyReaction) React(l logger.Logger, event iva.PackedEvent) {
	//TODO implement me
	panic("implement me")
}
