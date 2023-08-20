package service

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/reactions"
	"github.com/RacoonMediaServer/rms-cctv/internal/timeline"
	"go-micro.dev/v4"
)

type Service struct {
	CameraManager DeviceManager
	Reactor       Reactor
	Notifier      micro.Publisher
	ReactFactory  reactions.Factory
	Timeline      timeline.Timeline
}
