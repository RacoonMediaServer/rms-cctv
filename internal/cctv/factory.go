package cctv

import "go-micro.dev/v4/logger"

type BackendType int

const (
	DebugBackend BackendType = iota
	RestBackend
)

func New(t BackendType) Backend {
	switch t {
	case DebugBackend:
		return &debugBackend{
			channels: map[ID]*channel{},
			archives: map[ID]*archive{},
			l:        logger.Fields(map[string]interface{}{"from": "debug-backend"}),
		}
	default:
		panic("unknown backend type")
	}
}
