package service

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactions"
	"github.com/RacoonMediaServer/rms-cctv/internal/settings"
	"github.com/RacoonMediaServer/rms-cctv/internal/timeline"
	"github.com/teambition/rrule-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

type Service struct {
	Database         Database
	CameraManager    DeviceManager
	Reactor          Reactor
	Notifier         micro.Publisher
	ReactFactory     reactions.Factory
	SettingsProvider settings.Provider
	Timeline         timeline.Timeline
}

func (s Service) Initialize() error {
	cameras, err := s.Database.LoadCameras()
	if err != nil {
		return fmt.Errorf("load cameras failed: %w", err)
	}

	for _, cam := range cameras {
		if err = s.registerCamera(cam); err != nil {
			logger.Errorf("Load camera %d failed: %s", cam.ID, err)
			// TODO: notify malfunction
		}
	}
	logger.Infof("Loaded %d cameras", len(cameras))
	return nil
}

func (s Service) registerCamera(cam *model.Camera) error {
	schedule, err := rrule.StrToRRuleSet(cam.Info.Schedule)
	if err != nil {
		return fmt.Errorf("parse schedule failed: %w", err)
	}

	if err := s.CameraManager.Add(cam, s.Reactor.PushEvent); err != nil {
		return fmt.Errorf("manager: %w", err)
	}

	s.Reactor.SetReactions(cam.ID, s.makeEventReactions(cam.Info, schedule))
	return nil
}
