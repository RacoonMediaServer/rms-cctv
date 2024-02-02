package cctv

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/config"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
)

func New(conf config.Backend, f servicemgr.ClientFactory) Backend {
	l := logger.Fields(map[string]interface{}{
		"from":    "backend",
		"backend": conf.Type,
	})
	return newExternalBackend(l, f)
}
