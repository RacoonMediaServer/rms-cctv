package cctv

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/config"
	"go-micro.dev/v4/logger"
)

func New(conf config.Backend) Backend {
	l := logger.Fields(map[string]interface{}{
		"from":    "backend",
		"backend": conf.Type,
	})
	switch conf.Type {
	case "debug":
		return newDebugBackend(l)
	case "external":
		return newExternalBackend(l, conf)
	default:
		panic("unknown backend type")
	}
}
