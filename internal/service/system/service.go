package system

import (
	"context"
	"github.com/RacoonMediaServer/rms-cctv/internal/settings"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	SettingsProvider settings.Provider
}

func (s Service) GetSettings(ctx context.Context, empty *emptypb.Empty, settings *rms_cctv.CctvSettings) error {
	localSettings := s.SettingsProvider.Load()
	settings.OneEventDefaultDurationSec = localSettings.OneEventDefaultDurationSec
	settings.EventNotifyThresholdIntervalSec = localSettings.EventNotifyThresholdIntervalSec
	return nil
}

func (s Service) SetSettings(ctx context.Context, settings *rms_cctv.CctvSettings, empty *emptypb.Empty) error {
	s.SettingsProvider.Save(settings)
	return nil
}
