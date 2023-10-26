package service

import (
	"context"
	"github.com/RacoonMediaServer/rms-cctv/internal/model"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/teambition/rrule-go"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

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

func (s Service) GetCameras(ctx context.Context, empty *emptypb.Empty, response *rms_cctv.GetCamerasResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) AddCamera(ctx context.Context, c *rms_cctv.Camera, response *rms_cctv.AddCameraResponse) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "service",
		"camera": c.Name,
		"method": "AddCamera",
	})

	_, err := rrule.StrToRRuleSet(c.Schedule)
	if err != nil {
		return makeError(l, "parse camera schedule failed: %w", err)
	}

	cam := model.Camera{Info: c}
	if err := s.CameraManager.Register(&cam); err != nil {
		return makeError(l, "register camera on the external CCTV system failed: %w", err)
	}

	if err = s.Database.AddCamera(&cam); err != nil {
		if err := s.CameraManager.Unregister(&cam); err != nil {
			l.Logf(logger.ErrorLevel, "Unregister camera failed: %s", err)
		}
		return makeError(l, "add camera to database failed: %w", err)
	}

	if err = s.registerCamera(&cam); err != nil {
		if err := s.Database.RemoveCamera(cam.ID); err != nil {
			l.Logf(logger.ErrorLevel, "Remove camera failed: %s", err)
		}
		if err := s.CameraManager.Unregister(&cam); err != nil {
			l.Logf(logger.ErrorLevel, "Unregister camera failed: %s", err)
		}
		return makeError(l, "initialize camera failed: %w", err)
	}

	l.Log(logger.InfoLevel, "Camera added")
	response.CameraId = uint32(cam.ID)
	return nil
}

func (s Service) ModifyCamera(ctx context.Context, c *rms_cctv.Camera, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteCamera(ctx context.Context, request *rms_cctv.DeleteCameraRequest, empty *emptypb.Empty) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "service",
		"camera": request.CameraId,
		"method": "DeleteCamera",
	})

	id := model.CameraID(request.CameraId)
	if err := s.Database.RemoveCamera(id); err != nil {
		return makeError(l, "cannot remove camera from db: %w", err)
	}
	if err := s.CameraManager.Remove(id); err != nil {
		l.Logf(logger.ErrorLevel, "stop camera failed: %w", err)
	}
	return nil
}

func (s Service) GetLiveUri(ctx context.Context, request *rms_cctv.GetLiveUriRequest, response *rms_cctv.GetLiveUriResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetReplayUri(ctx context.Context, request *rms_cctv.GetReplayUriRequest, request2 *rms_cctv.GetReplayUriRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetSnapshot(ctx context.Context, request *rms_cctv.GetSnapshotRequest, response *rms_cctv.GetSnapshotResponse) error {
	l := logger.Fields(map[string]interface{}{
		"from":   "service",
		"camera": request.CameraId,
		"method": "GetSnapshot",
	})
	cam, err := s.CameraManager.GetCamera(model.CameraID(request.CameraId))
	if err != nil {
		return makeError(l, "operation failed: %w", err)
	}
	response.Snapshot, err = cam.TakeSnapshot(model.PrimaryProfile)
	if err != nil {
		return makeError(l, "get snapshot failed: %w", err)
	}
	return nil
}
