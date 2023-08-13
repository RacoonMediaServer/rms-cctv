package reactor

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"go-micro.dev/v4/logger"
)

type Reaction interface {
	React(l logger.Logger, event iva.PackedEvent)
}
